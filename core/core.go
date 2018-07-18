package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/wallet"
	"gopkg.in/natefinch/lumberjack.v2"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/fsrepo"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var fileLogFormat = logging.MustStringFormatter(
	`%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
)
var log = logging.MustGetLogger("core")

const Version = "0.0.5"

// Node is the single TextileNode instance
var Node *TextileNode

// TextileNode is the main node interface for textile functionality
type TextileNode struct {
	Wallet *wallet.Wallet
	server *http.Server
	mux    sync.Mutex
}

// NodeConfig is used to configure the node
type NodeConfig struct {
	LogLevel     logging.Level
	LogFiles     bool
	WalletConfig wallet.Config
}

// NewNode creates a new TextileNode
func NewNode(config NodeConfig) (*TextileNode, string, error) {
	// TODO: shouldn't need to manually remove these
	repoLockFile := filepath.Join(config.WalletConfig.RepoPath, fsrepo.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(config.WalletConfig.RepoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// log handling
	var backendFile *logging.LogBackend
	if config.LogFiles {
		w := &lumberjack.Logger{
			Filename:   path.Join(config.WalletConfig.RepoPath, "logs", "textile.log"),
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     30, // days
		}
		backendFile = logging.NewLogBackend(w, "", 0)
	} else {
		backendFile = logging.NewLogBackend(os.Stdout, "", 0)
	}
	backendFileFormatter := logging.NewBackendFormatter(backendFile, fileLogFormat)
	logging.SetBackend(backendFileFormatter)
	logging.SetLevel(config.LogLevel, "")

	// create a wallet
	config.WalletConfig.Version = Version
	wall, mnemonic, err := wallet.NewWallet(config.WalletConfig)
	if err != nil {
		return nil, "", err
	}

	// finally, construct our node
	node := &TextileNode{Wallet: wall}

	return node, mnemonic, nil
}

// StopWallet starts the wallet
func (t *TextileNode) StartWallet() (online chan struct{}, err error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.Wallet.Start()
}

// StopWallet stops the wallet
func (t *TextileNode) StopWallet() error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.Wallet.Stop()
}

// StartServer starts the server
func (t *TextileNode) StartServer() {
	router := gin.Default()
	router.GET("/ipfs/:cid/:path", gateway)
	v1 := router.Group("/api/v1")
	{
		v1.POST("/pin", pin)
	}
	t.server = &http.Server{
		Addr:    t.Wallet.GetServerAddress(),
		Handler: router,
	}

	// start the server
	errc := make(chan error)
	go func() {
		errc <- t.server.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err != http.ErrServerClosed {
					log.Errorf("server error: %s", err)
				}
				if !ok {
					log.Info("server was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("server listening at %s\n", t.server.Addr)
}

// StopServer stops the server
func (t *TextileNode) StopServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := t.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down server: %s", err)
		return err
	}
	cancel()
	return nil
}

// GetServerAddress returns the server's address
func (t *TextileNode) GetServerAddress() string {
	return t.server.Addr
}

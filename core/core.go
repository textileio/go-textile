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
	"github.com/textileio/textile-go/cafe"
)

var fileLogFormat = logging.MustStringFormatter(
	`%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
)
var log = logging.MustGetLogger("core")

const Version = "0.0.7"

// Node is the single TextileNode instance
var Node *TextileNode

// TextileNode is the main node interface for textile functionality
type TextileNode struct {
	Wallet  *wallet.Wallet
	gateway *http.Server
	api     *http.Server
	mux     sync.Mutex
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

// StartGateway starts the gateway
func (t *TextileNode) StartGateway() {
	router := gin.Default()
	router.GET("/ipfs/:cid/:path", gateway)
	v1 := router.Group("/api/v1")
	{
		v1.POST("/pin", pin)
	}
	t.gateway = &http.Server{
		Addr:    t.Wallet.GetGatewayAddress(),
		Handler: router,
	}

	// start the gateway
	errc := make(chan error)
	go func() {
		errc <- t.gateway.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err != http.ErrServerClosed {
					log.Errorf("gateway error: %s", err)
				}
				if !ok {
					log.Info("gateway was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("gateway listening at %s\n", t.gateway.Addr)
}

// StopGateway stops the gateway
func (t *TextileNode) StopGateway() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := t.gateway.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down gateway: %s", err)
		return err
	}
	cancel()
	return nil
}

// GetGatewayAddress returns the gateway's address
func (t *TextileNode) GetGatewayAddress() string {
	return t.gateway.Addr
}

// StartAPI starts the api
func (t *TextileNode) StartAPI() {
	router := cafe.Router()
	t.api = &http.Server{
		Addr:    t.Wallet.GetAPIAddress(),
		Handler: router,
	}

	// start the api
	errc := make(chan error)
	go func() {
		errc <- t.api.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err != http.ErrServerClosed {
					log.Errorf("api error: %s", err)
				}
				if !ok {
					log.Info("api was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("api listening at %s\n", t.api.Addr)
}

// StopAPI stops the api
func (t *TextileNode) StopAPI() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := t.api.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down api: %s", err)
		return err
	}
	cancel()
	return nil
}

// GetAPIAddress returns the gateway's address
func (t *TextileNode) GetAPIAddress() string {
	return t.api.Addr
}

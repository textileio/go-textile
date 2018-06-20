package core

import (
	"context"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/wallet"
	"gopkg.in/natefinch/lumberjack.v2"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/fsrepo"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var fileLogFormat = logging.MustStringFormatter(
	`%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
)
var log = logging.MustGetLogger("core")

const Version = "0.0.2"
const threadPublishInterval = time.Minute * 1

// Node is the single TextileNode instance
var Node *TextileNode

// TextileNode is the main node interface for textile functionality
type TextileNode struct {
	Wallet  *wallet.Wallet
	gateway *http.Server
	mux     sync.Mutex
}

// NodeConfig is used to configure the node
type NodeConfig struct {
	LogLevel     logging.Level
	LogFiles     bool
	WalletConfig wallet.Config
}

// NewNode creates a new TextileNode
func NewNode(config NodeConfig) (*TextileNode, error) {
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
	wall, err := wallet.NewWallet(config.WalletConfig)
	if err != nil {
		return nil, err
	}

	// setup gateway
	gateway := &http.Server{Addr: wall.GetGatewayAddress()}

	// finally, construct our node
	node := &TextileNode{
		Wallet:  wall,
		gateway: gateway,
	}

	return node, nil
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

// StopGateway starts the gateway
func (t *TextileNode) StartGateway() error {
	// try to register our handler
	t.registerGatewayHandler()

	// start the server
	errc := make(chan error)
	go func() {
		errc <- t.gateway.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err.Error() != "http: Server closed" {
					log.Errorf("gateway error: %s", err)
				}
				if !ok {
					log.Info("decrypting gateway was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("decrypting gateway (readonly) server listening at %s\n", t.gateway.Addr)
	return nil
}

// StopGateway stops the gateway
func (t *TextileNode) StopGateway() error {
	cgCtx, cancelCGW := context.WithCancel(context.Background())
	if err := t.gateway.Shutdown(cgCtx); err != nil {
		log.Errorf("error shutting down gateway: %s", err)
		return err
	}
	cancelCGW()
	return nil
}

// GetGatewayAddress returns the gateway's address
func (t *TextileNode) GetGatewayAddress() string {
	return t.gateway.Addr
}

// StartPublishing continuously publishes the latest update in each thread
func (t *TextileNode) StartPublishing() {
	t.Wallet.PublishThreads() // start now
	ticker := time.NewTicker(threadPublishInterval)
	defer func() {
		ticker.Stop()
		defer func() {
			if recover() != nil {
				log.Error("publishing ticker already stopped")
			}
		}()
	}()
	go func() {
		for range ticker.C {
			t.Wallet.PublishThreads()
		}
	}()

	// we can stop when the node stops
	for {
		if !t.Wallet.Started() {
			return
		}
		select {
		case <-t.Wallet.Done():
			log.Info("publishing stopped")
			return
		}
	}
}

// registerGatewayHandler registers a handler for the gateway
func (t *TextileNode) registerGatewayHandler() {
	defer func() {
		if recover() != nil {
			log.Debug("gateway handler already registered")
		}
	}()
	// NOTE: always returning 404 in the event of an error seems most secure as it doesn't reveal existence
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("gateway request: %s", r.URL.RequestURI())
		parsed, contentType := parsePath(r.URL.Path)

		// look for block id
		blockId := r.URL.Query()["block"]
		if blockId != nil {
			block, err := t.Wallet.GetBlock(blockId[0])
			if err != nil {
				log.Errorf("error finding block %s: %s", blockId[0], err)
				return
			}
			thrd := t.Wallet.GetThread(block.Id)
			if thrd == nil {
				log.Errorf("could not find thread for block: %s", block.Id)
				return
			}
			file, err := thrd.GetFileData(parsed, block)
			if err != nil {
				log.Errorf("error decrypting path %s: %s", parsed, err)
				w.WriteHeader(404)
				return
			}
			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
			w.Write(file)
			return
		}

		// get raw file
		file, err := t.Wallet.GetDataAtPath(parsed)
		if err != nil {
			log.Errorf("error getting raw path %s: %s", parsed, err)
			w.WriteHeader(404)
			return
		}

		// if key is provided, try to decrypt the file with it
		key := r.URL.Query()["key"]
		if key != nil {
			plain, err := crypto.DecryptAES(file, []byte(key[0]))
			if err != nil {
				log.Errorf("error decrypting %s: %s", parsed, err)
				w.WriteHeader(404)
				return
			}
			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
			w.Write(plain)
			return
		}

		// lastly, just return the raw bytes (standard gateway)
		w.Write(file)
	})
}

func parsePath(path string) (parsed string, contentType string) {
	parts := strings.Split(path, ".")
	parsed = parts[0]
	if len(parts) == 1 {
		return parsed, ""
	}
	switch parts[1] {
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "gif":
		contentType = "image/gif"
	}
	return parsed, contentType
}

package core

import (
	"github.com/op/go-logging"
	w "github.com/textileio/textile-go/wallet"
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

const Version = "0.1.2"

// Node is the single TextileNode instance
var Node *TextileNode

// TextileNode is the main node interface for textile functionality
type TextileNode struct {
	Wallet  *w.Wallet
	gateway *http.Server
	mux     sync.Mutex
}

// NodeConfig is used to configure the node
type NodeConfig struct {
	WalletConfig w.Config
	LogLevel     logging.Level
	LogFiles     bool
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
		logger := &lumberjack.Logger{
			Filename:   path.Join(config.WalletConfig.RepoPath, "logs", "textile.log"),
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     30, // days
		}
		backendFile = logging.NewLogBackend(logger, "", 0)
	} else {
		backendFile = logging.NewLogBackend(os.Stdout, "", 0)
	}
	backendFileFormatter := logging.NewBackendFormatter(backendFile, fileLogFormat)
	logging.SetBackend(backendFileFormatter)
	logging.SetLevel(config.LogLevel, "")

	// create a wallet
	config.WalletConfig.Version = Version
	wallet, mnemonic, err := w.NewWallet(config.WalletConfig)
	if err != nil {
		return nil, "", err
	}

	// construct our node
	node := &TextileNode{Wallet: wallet}

	return node, mnemonic, nil
}

// StopWallet starts the wallet
func (t *TextileNode) StartWallet() error {
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

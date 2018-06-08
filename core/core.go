package core

import (
	"context"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/crypto"
	trepo "github.com/textileio/textile-go/repo"
	tconfig "github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/wallet"
	"gopkg.in/natefinch/lumberjack.v2"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/fsrepo"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	RepoPath      string
	CentralApiURL string
	IsMobile      bool
	IsServer      bool
	LogLevel      logging.Level
	LogFiles      bool
	SwarmPort     string
}

// NewNode creates a new TextileNode
func NewNode(config NodeConfig) (*TextileNode, error) {
	// TODO: shouldn't need to manually remove these
	repoLockFile := filepath.Join(config.RepoPath, fsrepo.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(config.RepoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// log handling
	var backendFile *logging.LogBackend
	if config.LogFiles {
		w := &lumberjack.Logger{
			Filename:   path.Join(config.RepoPath, "logs", "textile.log"),
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

	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, "")
	if err != nil {
		return nil, err
	}

	// we may be running in an uninitialized state.
	err = trepo.DoInit(config.RepoPath, config.IsMobile, Version, sqliteDB.Config().Init, sqliteDB.Config().Configure)
	if err != nil && err != trepo.ErrRepoExists {
		return nil, err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(config.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return nil, err
	}

	// setup gateway
	gwAddr, err := repo.GetConfigKey("Addresses.Gateway")
	if err != nil {
		log.Errorf("error getting ipfs config: %s", err)
		return nil, err
	}
	gateway := &http.Server{Addr: gwAddr.(string)}

	// if a specific swarm port was selected, set it in the config
	if config.SwarmPort != "" {
		log.Infof("using specified swarm port: %s", config.SwarmPort)
		if err := tconfig.Update(repo, "Addresses.Swarm", []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", config.SwarmPort),
			fmt.Sprintf("/ip6/::/tcp/%s", config.SwarmPort),
		}); err != nil {
			return nil, err
		}
	}

	// if this is a server node, apply the ipfs server profile
	if config.IsServer {
		if err := tconfig.Update(repo, "Addresses.NoAnnounce", tconfig.DefaultServerFilters); err != nil {
			return nil, err
		}
		if err := tconfig.Update(repo, "Swarm.AddrFilters", tconfig.DefaultServerFilters); err != nil {
			return nil, err
		}
		if err := tconfig.Update(repo, "Swarm.EnableRelayHop", true); err != nil {
			return nil, err
		}
		if err := tconfig.Update(repo, "Discovery.MDNS.Enabled", false); err != nil {
			return nil, err
		}
		log.Info("applied server profile")
	}

	// clean central api url
	if len(config.CentralApiURL) > 0 {
		ca := config.CentralApiURL
		if ca[len(ca)-1:] == "/" {
			ca = ca[0 : len(ca)-1]
		}
		config.CentralApiURL = ca
	}

	// finally, construct our node
	node := &TextileNode{
		Wallet: &wallet.Wallet{
			RepoPath:       config.RepoPath,
			Datastore:      sqliteDB,
			CentralUserAPI: fmt.Sprintf("%s/api/v1/users", config.CentralApiURL),
			IsMobile:       config.IsMobile,
		},
		gateway: gateway,
	}

	return node, nil
}

func (t *TextileNode) StartWallet() (online chan struct{}, err error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	online, err = t.Wallet.Start()
	if err != nil {
		return nil, err
	}

	// construct decrypting http gateway
	var gwpErrc <-chan error
	gwpErrc, err = t.startGateway()
	if err != nil {
		log.Errorf("error starting decrypting gateway: %s", err)
		return nil, err
	}
	go func() {
		for {
			select {
			case err, ok := <-gwpErrc:
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

	return online, nil
}

func (t *TextileNode) StopWallet() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	// shutdown the gateway
	cgCtx, cancelCGW := context.WithCancel(context.Background())
	if err := t.gateway.Shutdown(cgCtx); err != nil {
		log.Errorf("error shutting down gateway: %s", err)
		return err
	}

	if err := t.Wallet.Stop(); err != nil {
		return err
	}

	// force the gateway closed if it's not already closed
	cancelCGW()

	return nil
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

		// look for block id
		blockId := r.URL.Query()["block"]
		if blockId != nil {
			file, err := t.Wallet.GetFile(r.URL.Path, blockId[0])
			if err != nil {
				log.Errorf("error decrypting path %s: %s", r.URL.Path, err)
				w.WriteHeader(404)
				return
			}
			w.Write(file)
			return
		}

		// get raw file
		file, err := wallet.GetDataAtPath(t.Wallet.Ipfs, r.URL.Path)
		if err != nil {
			log.Errorf("error getting raw path %s: %s", r.URL.Path, err)
			w.WriteHeader(404)
			return
		}

		// if key is provided, try to decrypt the file with it
		key := r.URL.Query()["key"]
		if key != nil {
			plain, err := crypto.DecryptAES(file, []byte(key[0]))
			if err != nil {
				log.Errorf("error decrypting %s: %s", r.URL.Path, err)
				w.WriteHeader(404)
				return
			}
			w.Write(plain)
			return
		}

		// lastly, just return the raw bytes (standard gateway)
		w.Write(file)
	})
}

// startGateway starts the secure HTTP gatway server
func (t *TextileNode) startGateway() (<-chan error, error) {
	// try to register our handler
	t.registerGatewayHandler()

	// Start the HTTPS server in a goroutine
	errc := make(chan error)
	go func() {
		errc <- t.gateway.ListenAndServe()
		close(errc)
	}()
	log.Infof("decrypting gateway (readonly) server listening at %s\n", t.gateway.Addr)

	return errc, nil
}

// startPublishing continuously publishes the latest update in each thread
func (t *TextileNode) startPublishing() {
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
		case <-t.Wallet.Ipfs.Context().Done():
			log.Info("publishing stopped")
			return
		}
	}
}

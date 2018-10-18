package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/archive"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/net"
	"github.com/textileio/textile-go/net/service"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/thread"
	"gopkg.in/natefinch/lumberjack.v2"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	utilmain "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/commands"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/config"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/fsrepo"
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

// Version is the core version identifier
const Version = "0.2.0"

// Node is the single Textile instance
var Node *Textile

// Update is used to notify UI listeners of changes
type Update struct {
	Id   string     `json:"id"`
	Name string     `json:"name"`
	Type UpdateType `json:"type"`
}

// UpdateType indicates a type of node update
type UpdateType int

const (
	// ThreadAdded is emitted when a thread is added
	ThreadAdded UpdateType = iota
	// ThreadRemoved is emitted when a thread is removed
	ThreadRemoved
	// AccountPeerAdded is emitted when an account peer (device) is added
	AccountPeerAdded
	// AccountPeerRemoved is emitted when an account peer (device) is removed
	AccountPeerRemoved
)

// AddDataResult wraps added data content id and key
type AddDataResult struct {
	Id      string           `json:"id"`
	Key     string           `json:"key"`
	Archive *archive.Archive `json:"archive,omitempty"`
}

// InitConfig is used to setup a textile node
type InitConfig struct {
	Account    keypair.Full
	PinCode    string
	RepoPath   string
	SwarmPorts string
	IsMobile   bool
	IsServer   bool
	LogLevel   logging.Level
	LogFiles   bool
}

// RunConfig is used to define run options for a textile node
type RunConfig struct {
	PinCode  string
	RepoPath string
	LogLevel logging.Level
	LogFiles bool
}

// Textile is the main Textile node structure
type Textile struct {
	version          string
	context          oldcmds.Context
	repoPath         string
	cancel           context.CancelFunc
	ipfs             *core.IpfsNode
	datastore        repo.Datastore
	started          bool
	threads          []*thread.Thread
	online           chan struct{}
	done             chan struct{}
	updates          chan Update
	threadUpdates    chan thread.Update
	notifications    chan repo.Notification
	threadsService   *net.ThreadsService
	cafeService      *net.CafeService
	cafeRequestQueue *net.CafeRequestQueue
	mux              sync.Mutex
}

// common errors
var ErrAccountRequired = errors.New("account required")
var ErrStarted = errors.New("node is started")
var ErrStopped = errors.New("node is stopped")
var ErrOffline = errors.New("node is offline")
var ErrThreadLoaded = errors.New("thread is loaded")

// InitRepo initializes a new node repo
func InitRepo(config InitConfig) error {
	// ensure init has not been run
	if fsrepo.IsInitialized(config.RepoPath) {
		return repo.ErrRepoExists
	}

	// log handling
	setupLogging(config.RepoPath, config.LogLevel, config.LogFiles)

	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, config.PinCode)
	if err != nil {
		return err
	}

	// init ipfs repo
	if err := repo.DoInit(config.RepoPath, func() error {
		if err := sqliteDB.Config().Init(config.PinCode); err != nil {
			return err
		}
		return sqliteDB.Config().Configure(&config.Account, config.IsMobile, time.Now())
	}); err != nil {
		return err
	}

	// open the repo
	rep, err := fsrepo.Open(config.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}

	// if a specific swarm port was selected, set it in the config
	if err := applySwarmPortConfigOption(rep, config.SwarmPorts); err != nil {
		return err
	}

	// if this is a server node, apply the ipfs server profile
	if err := applyServerConfigOption(rep, config.IsServer); err != nil {
		return err
	}

	// add account key to ipfs keystore for resolving ipns profile
	sk, err := config.Account.LibP2PPrivKey()
	if err != nil {
		return err
	}
	return rep.Keystore().Put("account", sk)

	// TODO: discover other devices

	//fmt.Println("Publishing new account peer identity...")
	//
	//// create a tmp node
	//node, err := NewTextile(RunConfig{
	//	PinCode:  config.PinCode,
	//	RepoPath: config.RepoPath,
	//	LogLevel: config.LogLevel,
	//	LogFiles: config.LogFiles,
	//})
	//if err != nil {
	//	return err
	//}
	//
	//// add new peer to account profile
	//if err := node.Start(); err != nil {
	//	return err
	//}
	//<-node.OnlineCh()
	//if _, err := node.PublishAccountProfile(nil); err != nil {
	//	log.Errorf("error publishing profile: %s", err)
	//}
	//return nil
}

// NewTextile runs a node out of an initialized repo
func NewTextile(config RunConfig) (*Textile, error) {
	// ensure init has been run
	if !fsrepo.IsInitialized(config.RepoPath) {
		return nil, repo.ErrRepoDoesNotExist
	}

	// force open the repo and datastore (fixme please)
	removeLocks(config.RepoPath)

	// log handling
	setupLogging(config.RepoPath, config.LogLevel, config.LogFiles)

	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, config.PinCode)
	if err != nil {
		return nil, err
	}

	// run all migrations if needed
	if err := repo.MigrateUp(config.RepoPath, config.PinCode, false); err != nil {
		return nil, err
	}

	// TODO: put cafes into bootstrap?

	// open repo
	//rep, err := fsrepo.Open(config.RepoPath)
	//if err != nil {
	//	log.Errorf("error opening repo: %s", err)
	//	return nil, err
	//}

	// ensure bootstrap addresses are latest in config
	//if err := ensureBootstrapConfig(rep); err != nil {
	//	return nil, err
	//}

	return &Textile{
		version:   Version,
		repoPath:  config.RepoPath,
		datastore: sqliteDB,
	}, nil
}

// Start
func (t *Textile) Start() error {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.started {
		return ErrStarted
	}
	defer func() {
		t.done = make(chan struct{})
		t.started = true

		addr, err := t.Address()
		if err != nil {
			log.Error(err.Error())
			return
		}
		accntId, err := t.Id()
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("node is started")
		log.Infof("account address: %s", addr)
		log.Infof("account id: %s", accntId.Pretty())
	}()
	log.Info("starting node...")

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		log.Errorf("setting file descriptor limit: %s", err)
	}

	// check db
	if err := t.touchDatastore(); err != nil {
		return err
	}

	// load account
	accnt, err := t.Account()
	if err != nil {
		return err
	}

	// build update channels
	t.online = make(chan struct{})
	t.updates = make(chan Update, 10)
	t.threadUpdates = make(chan thread.Update, 10)
	t.notifications = make(chan repo.Notification, 10)

	// start the ipfs node
	log.Debug("creating an ipfs node...")
	if err := t.createIPFS(false); err != nil {
		log.Errorf("error creating offline ipfs node: %s", err)
		return err
	}
	go func() {
		defer close(t.online)
		if err := t.createIPFS(true); err != nil {
			log.Errorf("error creating online ipfs node: %s", err)
			return
		}

		// setup thread service
		t.threadsService = net.NewThreadsService(
			accnt,
			t.ipfs,
			t.datastore,
			t.GetThread,
			t.sendNotification,
		)

		// setup cafe service
		t.cafeService = net.NewCafeService(accnt, t.ipfs, t.datastore)

		// start store queue
		if t.Mobile() {
			go t.cafeRequestQueue.Flush()
		} else {
			go t.cafeRequestQueue.Run()
		}

		// print swarm addresses
		if err := ipfs.PrintSwarmAddrs(t.ipfs); err != nil {
			log.Errorf(err.Error())
		}
		log.Info("node is online")
	}()

	// build a store request queue
	t.cafeRequestQueue = net.NewCafeRequestQueue(
		func() *net.CafeService {
			return t.cafeService
		},
		func() *core.IpfsNode {
			return t.ipfs
		},
		t.datastore,
	)

	// setup threads
	for _, mod := range t.datastore.Threads().List("") {
		_, err := t.loadThread(&mod)
		if err == ErrThreadLoaded {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop the node
func (t *Textile) Stop() error {
	t.mux.Lock()
	defer t.mux.Unlock()
	if !t.started {
		return ErrStopped
	}
	defer func() {
		t.started = false
		close(t.done)
	}()
	log.Info("stopping node...")

	// close ipfs node
	t.context.Close()
	t.cancel()
	if err := t.ipfs.Close(); err != nil {
		log.Errorf("error closing ipfs node: %s", err)
		return err
	}

	// close db connection
	t.datastore.Close()
	dsLockFile := filepath.Join(t.repoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// wipe threads
	t.threads = nil

	// close update channels
	close(t.updates)
	close(t.threadUpdates)
	close(t.notifications)

	log.Info("node is stopped")

	return nil
}

// Started returns whether or not node is started
func (t *Textile) Started() bool {
	return t.started
}

// Online returns whether or not node is online
func (t *Textile) Online() bool {
	if t.ipfs == nil {
		return false
	}
	return t.started && t.ipfs.OnlineMode()
}

// Mobile returns whether or not node is configured for a mobile device
func (t *Textile) Mobile() bool {
	if err := t.touchDatastore(); err != nil {
		log.Errorf("error calling is mobile: %s", err)
		return false
	}
	mobile, err := t.datastore.Config().GetMobile()
	if err != nil {
		log.Errorf("error calling is mobile: %s", err)
		return false
	}
	return mobile
}

// Version return core node version
func (t *Textile) Version() string {
	return t.version
}

// Ipfs returns the underlying ipfs node
func (t *Textile) Ipfs() *core.IpfsNode {
	return t.ipfs
}

// OnlineCh returns the online channel
func (t *Textile) OnlineCh() <-chan struct{} {
	return t.online
}

// DoneCh returns the core node done channel
func (t *Textile) DoneCh() <-chan struct{} {
	return t.done
}

// Ping pings another peer
func (t *Textile) Ping(pid peer.ID) (service.PeerStatus, error) {
	return t.cafeService.Ping(pid)
}

// Update returns the node update channel
func (t *Textile) Updates() <-chan Update {
	return t.updates
}

// ThreadUpdates returns the thread update channel
func (t *Textile) ThreadUpdates() <-chan thread.Update {
	return t.threadUpdates
}

// Notifications returns the notifications channel
func (t *Textile) Notifications() <-chan repo.Notification {
	return t.notifications
}

// GetPeerId returns peer id
func (t *Textile) GetPeerId() (peer.ID, error) {
	if !t.started {
		return "", ErrStopped
	}
	return t.ipfs.Identity, nil
}

// GetPrivKey returns the current peer private key
func (t *Textile) GetPeerPrivKey() (libp2pc.PrivKey, error) {
	if !t.started {
		return nil, ErrStopped
	}
	if t.ipfs.PrivateKey == nil {
		if err := t.ipfs.LoadPrivateKey(); err != nil {
			return nil, err
		}
	}
	return t.ipfs.PrivateKey, nil
}

// GetPeerPubKey returns the current peer public key
func (t *Textile) GetPeerPubKey() (libp2pc.PubKey, error) {
	sk, err := t.GetPeerPrivKey()
	if err != nil {
		return nil, err
	}
	return sk.GetPublic(), nil
}

// GetRepoPath returns the node's repo path
func (t *Textile) GetRepoPath() string {
	return t.repoPath
}

// GetDataAtPath returns raw data behind an ipfs path
func (t *Textile) GetDataAtPath(path string) ([]byte, error) {
	if !t.started {
		return nil, ErrStopped
	}
	return ipfs.GetDataAtPath(t.ipfs, path)
}

// GetLinksAtPath returns ipld links behind an ipfs path
func (t *Textile) GetLinksAtPath(path string) ([]*ipld.Link, error) {
	if !t.started {
		return nil, ErrStopped
	}
	return ipfs.GetLinksAtPath(t.ipfs, path)
}

// createIPFS creates an IPFS node
func (t *Textile) createIPFS(online bool) error {
	// open repo
	rep, err := fsrepo.Open(t.repoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}

	// determine routing
	routing := core.DHTOption
	if t.Mobile() {
		routing = core.DHTClientOption
	}

	// assemble node config
	cfg := &core.BuildCfg{
		Repo:      rep,
		Permanent: true, // temporary way to signify that node is permanent
		Online:    online,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
			"mplex":  true,
		},
		Routing: routing,
	}

	// create the node
	cctx, cancel := context.WithCancel(context.Background())
	nd, err := core.NewNode(cctx, cfg)
	if err != nil {
		return err
	}
	nd.SetLocal(!online)

	// build the context
	ctx := oldcmds.Context{}
	ctx.Online = online
	ctx.ConfigRoot = t.repoPath
	ctx.LoadConfig = func(path string) (*config.Config, error) {
		return fsrepo.ConfigAt(t.repoPath)
	}
	ctx.ConstructNode = func() (*core.IpfsNode, error) {
		return nd, nil
	}

	// attach to textile node
	if t.cancel != nil {
		t.cancel()
	}
	if t.ipfs != nil {
		if err := t.ipfs.Close(); err != nil {
			log.Errorf("error closing prev ipfs node: %s", err)
			return err
		}
	}
	t.context = ctx
	t.cancel = cancel
	t.ipfs = nd

	return nil
}

// getThreadByBlock returns the thread owning the given block
func (t *Textile) getThreadByBlock(block *repo.Block) (*thread.Thread, error) {
	if block == nil {
		return nil, errors.New("block is empty")
	}
	var thrd *thread.Thread
	for _, t := range t.threads {
		if t.Id == block.ThreadId {
			thrd = t
			break
		}
	}
	if thrd == nil {
		return nil, errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadId))
	}
	return thrd, nil
}

// loadThread loads a thread into memory from the given on-disk model
func (t *Textile) loadThread(mod *repo.Thread) (*thread.Thread, error) {
	if _, loaded := t.GetThread(mod.Id); loaded != nil {
		return nil, ErrThreadLoaded
	}
	id := mod.Id // save value locally
	threadConfig := &thread.Config{
		RepoPath: t.repoPath,
		Ipfs: func() *core.IpfsNode {
			return t.ipfs
		},
		Blocks:        t.datastore.Blocks,
		Peers:         t.datastore.ThreadPeers,
		Notifications: t.datastore.Notifications,
		GetHead: func() (string, error) {
			thrd := t.datastore.Threads().Get(id)
			if thrd == nil {
				return "", errors.New(fmt.Sprintf("could not re-load thread: %s", id))
			}
			return thrd.Head, nil
		},
		UpdateHead: func(head string) error {
			if err := t.datastore.Threads().UpdateHead(id, head); err != nil {
				return err
			}
			t.cafeRequestQueue.Add(id, repo.CafeStoreThreadRequest)
			return nil
		},
		NewBlock:       t.threadsService.NewBlock,
		SendMessage:    t.threadsService.SendMessage,
		AddCafeRequest: t.cafeRequestQueue.Add,
		GetUsername:    t.GetUsername,
		SendUpdate:     t.sendThreadUpdate,
	}
	thrd, err := thread.NewThread(mod, threadConfig)
	if err != nil {
		return nil, err
	}
	t.threads = append(t.threads, thrd)
	return thrd, nil
}

// sendUpdate adds an update to the update channel
func (t *Textile) sendUpdate(update Update) {
	defer func() {
		if recover() != nil {
			log.Error("update channel already closed")
		}
	}()
	t.updates <- update
}

// sendThreadUpdate adds a thread update to the update channel
func (t *Textile) sendThreadUpdate(update thread.Update) {
	defer func() {
		if recover() != nil {
			log.Error("thread update channel already closed")
		}
	}()
	t.threadUpdates <- update
}

// sendNotification adds a notification to the notification channel
func (t *Textile) sendNotification(notification *repo.Notification) error {
	// add to db
	if err := t.datastore.Notifications().Add(notification); err != nil {
		return err
	}

	// broadcast
	defer func() {
		if recover() != nil {
			log.Error("notification channel already closed")
		}
	}()
	t.notifications <- *notification

	return nil
}

// touchDatastore ensures that we have a good db connection
func (t *Textile) touchDatastore() error {
	if err := t.datastore.Ping(); err != nil {
		log.Debug("re-opening datastore...")
		sqliteDB, err := db.Create(t.repoPath, "")
		if err != nil {
			log.Errorf("error re-opening datastore: %s", err)
			return err
		}
		t.datastore = sqliteDB
	}
	return nil
}

// setupLogging handles log settings
func setupLogging(repoPath string, level logging.Level, files bool) {
	var backendFile *logging.LogBackend
	if files {
		logger := &lumberjack.Logger{
			Filename:   path.Join(repoPath, "logs", "textile.log"),
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
	logging.SetLevel(level, "")
}

// removeLocks force deletes the IPFS repo and SQLite DB lock files
func removeLocks(repoPath string) {
	repoLockFile := filepath.Join(repoPath, fsrepo.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(repoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)
}

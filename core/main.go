package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ipfsconfig "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	utilmain "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/commands"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"
	logger "gx/ipfs/QmcaSwFc5RBg8yCq54QURwEU4nwjfCpjbpmaAm4VbdGLKv/go-logging"

	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/service"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = logging.Logger("tex-core")

// Version is the core version identifier
const Version = "1.0.0-rc7"

// kQueueFlushFreq how often to flush the message queues
const kQueueFlushFreq = time.Second * 60

// kMobileQueueFlushFreq how often to flush the message queues on mobile
const kMobileQueueFlush = time.Second * 20

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

// InitConfig is used to setup a textile node
type InitConfig struct {
	Account      *keypair.Full
	PinCode      string
	RepoPath     string
	SwarmPorts   string
	ApiAddr      string
	CafeApiAddr  string
	GatewayAddr  string
	IsMobile     bool
	IsServer     bool
	LogLevel     logger.Level
	LogToDisk    bool
	CafeOpen     bool
	CafePublicIP string
}

// MigrateConfig is used to define options during a major migration
type MigrateConfig struct {
	PinCode  string
	RepoPath string
}

// RunConfig is used to define run options for a textile node
type RunConfig struct {
	PinCode  string
	RepoPath string
}

// Textile is the main Textile node structure
type Textile struct {
	context        oldcmds.Context
	repoPath       string
	config         *config.Config
	account        *keypair.Full
	cancel         context.CancelFunc
	node           *core.IpfsNode
	datastore      repo.Datastore
	started        bool
	threads        []*Thread
	online         chan struct{}
	done           chan struct{}
	updates        chan Update
	threadUpdates  *broadcast.Broadcaster
	notifications  chan NotificationInfo
	threadsService *ThreadsService
	threadsOutbox  *ThreadsOutbox
	cafeService    *CafeService
	cafeOutbox     *CafeOutbox
	cafeInbox      *CafeInbox
	mux            sync.Mutex
	writer         io.Writer
}

// common errors
var ErrAccountRequired = errors.New("account required")
var ErrStarted = errors.New("node is started")
var ErrStopped = errors.New("node is stopped")
var ErrOffline = errors.New("node is offline")

// InitRepo initializes a new node repo
func InitRepo(conf InitConfig) error {
	if fsrepo.IsInitialized(conf.RepoPath) {
		return repo.ErrRepoExists
	}

	if conf.Account == nil {
		return ErrAccountRequired
	}

	setupLogging(conf.RepoPath, conf.LogLevel, conf.LogToDisk)

	// init repo
	if err := repo.Init(conf.RepoPath, Version); err != nil {
		return err
	}

	rep, err := fsrepo.Open(conf.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}
	defer rep.Close()

	// apply ipfs config opts
	if err := applySwarmPortConfigOption(rep, conf.SwarmPorts); err != nil {
		return err
	}
	if err := applyServerConfigOption(rep, conf.IsServer); err != nil {
		return err
	}

	sqliteDb, err := db.Create(conf.RepoPath, conf.PinCode)
	if err != nil {
		return err
	}
	if err := sqliteDb.Config().Init(conf.PinCode); err != nil {
		return err
	}
	if err := sqliteDb.Config().Configure(conf.Account, time.Now()); err != nil {
		return err
	}

	return applyTextileConfigOptions(conf)
}

// MigrateRepo runs _all_ repo migrations, including major
func MigrateRepo(conf MigrateConfig) error {
	if !fsrepo.IsInitialized(conf.RepoPath) {
		return repo.ErrRepoDoesNotExist
	}

	// force open the repo and datastore
	removeLocks(conf.RepoPath)

	// run _all_ repo migrations if needed
	return repo.MigrateUp(conf.RepoPath, conf.PinCode, false)
}

// NewTextile runs a node out of an initialized repo
func NewTextile(conf RunConfig) (*Textile, error) {
	if !fsrepo.IsInitialized(conf.RepoPath) {
		return nil, repo.ErrRepoDoesNotExist
	}

	// check if repo needs a major migration
	if err := repo.Stat(conf.RepoPath); err != nil {
		return nil, err
	}

	// force open the repo and datastore
	removeLocks(conf.RepoPath)

	node := &Textile{
		repoPath:      conf.RepoPath,
		updates:       make(chan Update, 10),
		threadUpdates: broadcast.NewBroadcaster(10),
		notifications: make(chan NotificationInfo, 10),
	}

	var err error
	node.config, err = config.Read(conf.RepoPath)
	if err != nil {
		return nil, err
	}

	llevel, err := logger.LogLevel(strings.ToUpper(node.config.Logs.LogLevel))
	if err != nil {
		llevel = logger.ERROR
	}
	node.writer = setupLogging(conf.RepoPath, llevel, node.config.Logs.LogToDisk)

	// run all minor repo migrations if needed
	if err := repo.MigrateUp(conf.RepoPath, conf.PinCode, false); err != nil {
		return nil, err
	}

	sqliteDb, err := db.Create(conf.RepoPath, conf.PinCode)
	if err != nil {
		return nil, err
	}
	node.datastore = sqliteDb

	accnt, err := node.datastore.Config().GetAccount()
	if err != nil {
		return nil, err
	}
	node.account = accnt

	return node, nil
}

// Start creates an ipfs node and starts textile services
func (t *Textile) Start() error {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.started {
		return ErrStarted
	}
	log.Info("starting node...")

	// raise file descriptor limit
	changed, limit, err := utilmain.ManageFdLimit()
	if err != nil {
		log.Errorf("error setting fd limit: %s", err)
	}
	log.Debugf("fd limit: %d (changed %t)", limit, changed)

	if err := t.touchDatastore(); err != nil {
		return err
	}

	swarmPorts, err := loadSwarmPorts(t.repoPath)
	if err != nil {
		return err
	}
	if swarmPorts == nil {
		return errors.New("failed to load swarm ports")
	}

	t.online = make(chan struct{})

	t.cafeInbox = NewCafeInbox(
		func() *CafeService {
			return t.cafeService
		},
		func() *ThreadsService {
			return t.threadsService
		},
		func() *core.IpfsNode {
			return t.node
		},
		t.datastore,
	)
	t.cafeOutbox = NewCafeOutbox(
		func() *CafeService {
			return t.cafeService
		},
		func() *core.IpfsNode {
			return t.node
		},
		t.datastore,
	)
	t.threadsOutbox = NewThreadsOutbox(
		func() *ThreadsService {
			return t.threadsService
		},
		func() *core.IpfsNode {
			return t.node
		},
		t.datastore,
		t.cafeOutbox,
	)

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

		t.threadsService = NewThreadsService(
			t.account,
			t.node,
			t.datastore,
			t.Thread,
			t.AddThread,
			t.sendNotification,
		)

		t.cafeService = NewCafeService(t.account, t.node, t.datastore, t.cafeInbox)
		t.cafeService.setAddrs(t.config.Addresses.CafeAPI, t.config.Cafe.Host.PublicIP, *swarmPorts)
		if t.config.Cafe.Host.Open {
			t.cafeService.open = true
			t.startCafeApi(t.config.Addresses.CafeAPI)
		}

		go t.runQueues()

		if err := ipfs.PrintSwarmAddrs(t.node); err != nil {
			log.Errorf(err.Error())
		}
		log.Info("node is online")
	}()

	for _, mod := range t.datastore.Threads().List() {
		if _, err := t.loadThread(&mod); err == ErrThreadLoaded {
			continue
		}
		if err != nil {
			return err
		}
	}

	t.done = make(chan struct{})
	t.started = true

	log.Info("node is started")
	log.Infof("peer id: %s", t.node.Identity.Pretty())
	log.Infof("account address: %s", t.account.Address())

	return t.addAccountThread()
}

// Stop destroys the ipfs node and shutsdown textile services
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

	// close apis
	if err := t.stopCafeApi(); err != nil {
		return err
	}

	// close ipfs node
	t.context.Close()
	t.cancel()
	if err := t.node.Close(); err != nil {
		log.Errorf("error closing ipfs node: %s", err)
		return err
	}

	// close db connection
	t.datastore.Close()
	dsLockFile := filepath.Join(t.repoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// wipe threads
	t.threads = nil

	log.Info("node is stopped")

	return nil
}

// CloseChns closes update channels
func (t *Textile) CloseChns() {
	close(t.updates)
	t.threadUpdates.Close()
	close(t.notifications)
}

// Started returns node started status
func (t *Textile) Started() bool {
	return t.started
}

// Online returns node online status
func (t *Textile) Online() bool {
	if t.node == nil {
		return false
	}
	return t.started && t.node.OnlineMode()
}

// Mobile returns whether or not node is configured for a mobile device
func (t *Textile) Mobile() bool {
	return t.config.IsMobile
}

// Writer returns the output writer (logger / stdout)
func (t *Textile) Writer() io.Writer {
	return t.writer
}

// Ipfs returns the underlying ipfs node
func (t *Textile) Ipfs() *core.IpfsNode {
	return t.node
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

// UpdateCh returns the node update channel
func (t *Textile) UpdateCh() <-chan Update {
	return t.updates
}

// GetThreadUpdateListener returns the thread update channel
func (t *Textile) GetThreadUpdateListener() *broadcast.Listener {
	return t.threadUpdates.Listen()
}

// NotificationsCh returns the notifications channel
func (t *Textile) NotificationCh() <-chan NotificationInfo {
	return t.notifications
}

// PeerId returns peer id
func (t *Textile) PeerId() (peer.ID, error) {
	return t.node.Identity, nil
}

// RepoPath returns the node's repo path
func (t *Textile) RepoPath() string {
	return t.repoPath
}

// DataAtPath returns raw data behind an ipfs path
func (t *Textile) DataAtPath(path string) ([]byte, error) {
	return ipfs.DataAtPath(t.node, path)
}

// LinksAtPath returns ipld links behind an ipfs path
func (t *Textile) LinksAtPath(path string) ([]*ipld.Link, error) {
	return ipfs.LinksAtPath(t.node, path)
}

// createIPFS creates an IPFS node
func (t *Textile) createIPFS(online bool) error {
	rep, err := fsrepo.Open(t.repoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}

	routing := core.DHTOption
	if t.Mobile() {
		routing = core.DHTClientOption
	}

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

	cctx, cancel := context.WithCancel(context.Background())
	nd, err := core.NewNode(cctx, cfg)
	if err != nil {
		return err
	}
	nd.SetLocal(!online)

	ctx := oldcmds.Context{}
	ctx.Online = online
	ctx.ConfigRoot = t.repoPath
	ctx.LoadConfig = func(path string) (*ipfsconfig.Config, error) {
		return fsrepo.ConfigAt(t.repoPath)
	}
	ctx.ConstructNode = func() (*core.IpfsNode, error) {
		return nd, nil
	}

	// attach to textile node
	if t.cancel != nil {
		t.cancel()
	}
	if t.node != nil {
		if err := t.node.Close(); err != nil {
			log.Errorf("error closing prev ipfs node: %s", err)
			return err
		}
	}
	t.context = ctx
	t.cancel = cancel
	t.node = nd

	return nil
}

// runQueues runs each message queue
func (t *Textile) runQueues() {
	var freq time.Duration
	if t.Mobile() {
		freq = kMobileQueueFlush
	} else {
		freq = kQueueFlushFreq
	}

	tick := time.NewTicker(freq)
	defer tick.Stop()

	t.flushQueues()

	for {
		select {
		case <-tick.C:
			t.flushQueues()
		case <-t.done:
			return
		}
	}
}

// flushQueues flushes each message queue
func (t *Textile) flushQueues() {
	if err := t.touchDatastore(); err != nil {
		log.Error(err)
		return
	}

	go func() {
		t.threadsOutbox.Flush()
		t.cafeInbox.CheckMessages()
		t.cafeOutbox.Flush()
	}()
}

// threadByBlock returns the thread owning the given block
func (t *Textile) threadByBlock(block *repo.Block) (*Thread, error) {
	if block == nil {
		return nil, errors.New("block is empty")
	}

	var thrd *Thread
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
func (t *Textile) loadThread(mod *repo.Thread) (*Thread, error) {
	if loaded := t.Thread(mod.Id); loaded != nil {
		return nil, ErrThreadLoaded
	}

	threadConfig := &ThreadConfig{
		RepoPath: t.repoPath,
		Config:   t.config,
		Node: func() *core.IpfsNode {
			return t.node
		},
		Datastore: t.datastore,
		Service: func() *ThreadsService {
			if t.threadsService == nil {
				return NewDummyThreadsService(t.account, t.node)
			}
			return t.threadsService
		},
		ThreadsOutbox: t.threadsOutbox,
		CafeOutbox:    t.cafeOutbox,
		SendUpdate:    t.sendThreadUpdate,
	}

	thrd, err := NewThread(mod, threadConfig)
	if err != nil {
		return nil, err
	}
	t.threads = append(t.threads, thrd)

	return thrd, nil
}

// sendUpdate adds an update to the update channel
func (t *Textile) sendUpdate(update Update) {
	t.updates <- update
}

// sendThreadUpdate adds a thread update to the update channel
func (t *Textile) sendThreadUpdate(update ThreadUpdate) {
	t.threadUpdates.Send(update)
}

// sendNotification adds a notification to the notification channel
func (t *Textile) sendNotification(notification *repo.Notification) error {
	if err := t.datastore.Notifications().Add(notification); err != nil {
		return err
	}

	t.notifications <- t.NotificationInfo(*notification)
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

// setupLogging hijacks the ipfs logging system, putting output to files
func setupLogging(repoPath string, level logger.Level, files bool) io.Writer {
	var writer io.Writer
	if files {
		writer = &lumberjack.Logger{
			Filename:   path.Join(repoPath, "logs", "textile.log"),
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     30, // days
		}
	} else {
		writer = os.Stdout
	}
	backendFile := logger.NewLogBackend(writer, "", 0)
	logger.SetBackend(backendFile)
	logging.SetAllLoggers(level)

	// tmp until we have a log command to alter subsystems
	logging.SetLogLevel("tex-core", "debug")
	logging.SetLogLevel("tex-service", "debug")
	logging.SetLogLevel("tex-gateway", "debug")
	logging.SetLogLevel("tex-ipfs", "debug")
	logging.SetLogLevel("tex-mill", "debug")
	logging.SetLogLevel("tex-mobile", "debug")

	return writer
}

// removeLocks force deletes the IPFS repo and SQLite DB lock files
func removeLocks(repoPath string) {
	repoLockFile := filepath.Join(repoPath, fsrepo.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(repoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)
}

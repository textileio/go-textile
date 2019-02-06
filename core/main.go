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
	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	utilmain "gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/commands"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/repo/fsrepo"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"
	logger "gx/ipfs/QmcaSwFc5RBg8yCq54QURwEU4nwjfCpjbpmaAm4VbdGLKv/go-logging"

	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/service"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var log = logging.Logger("tex-core")

// Version is the core version identifier
const Version = "1.0.0-rc33"

// kQueueFlushFreq how often to flush the message queues
const kQueueFlushFreq = time.Second * 60

// kMobileQueueFlushFreq how often to flush the message queues on mobile
const kMobileQueueFlush = time.Second * 40

// Update is used to notify UI listeners of changes
type Update struct {
	Id   string     `json:"id"`
	Key  string     `json:"key"`
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
	Account         *keypair.Full
	PinCode         string
	RepoPath        string
	SwarmPorts      string
	ApiAddr         string
	CafeApiAddr     string
	GatewayAddr     string
	IsMobile        bool
	IsServer        bool
	LogToDisk       bool
	Debug           bool
	CafeOpen        bool
	CafePublicIP    string
	CafeURL         string
	CafeNeighborURL string
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
	Debug    bool
}

// Textile is the main Textile node structure
type Textile struct {
	context       oldcmds.Context
	repoPath      string
	config        *config.Config
	account       *keypair.Full
	cancel        context.CancelFunc
	node          *core.IpfsNode
	datastore     repo.Datastore
	started       bool
	loadedThreads []*Thread
	online        chan struct{}
	done          chan struct{}
	updates       chan Update
	threadUpdates *broadcast.Broadcaster
	notifications chan NotificationInfo
	threads       *ThreadsService
	threadsOutbox *ThreadsOutbox
	cafe          *CafeService
	cafeOutbox    *CafeOutbox
	cafeInbox     *CafeInbox
	mux           sync.Mutex
	writer        io.Writer
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

	logLevels := map[string]string{}
	if conf.Debug {
		logLevels = getTextileDebugLevels()
	}
	if _, err := setLogLevels(conf.RepoPath, logLevels, conf.LogToDisk); err != nil {
		return err
	}

	// init repo
	if err := repo.Init(conf.RepoPath, Version); err != nil {
		return err
	}

	rep, err := fsrepo.Open(conf.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}
	defer func() {
		if err := rep.Close(); err != nil {
			log.Error(err)
		}
	}()

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

	// add self as a contact
	ipfsConf, err := rep.Config()
	if err != nil {
		return err
	}
	if err := sqliteDb.Contacts().Add(&repo.Contact{
		Id:      ipfsConf.Identity.PeerID,
		Address: conf.Account.Address(),
	}); err != nil {
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

	logLevels := map[string]string{}
	if conf.Debug {
		logLevels = getTextileDebugLevels()
	}
	node.writer, err = setLogLevels(conf.RepoPath, logLevels, node.config.Logs.LogToDisk)
	if err != nil {
		return nil, err
	}

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

	// open db
	if err := t.touchDatastore(); err != nil {
		return err
	}

	t.online = make(chan struct{})

	t.cafeInbox = NewCafeInbox(t.cafeService, t.threadsService, t.Ipfs, t.datastore)
	t.cafeOutbox = NewCafeOutbox(t.cafeService, t.Ipfs, t.datastore)
	t.threadsOutbox = NewThreadsOutbox(t.threadsService, t.Ipfs, t.datastore, t.cafeOutbox)
	t.threads = NewThreadsService(t.account, t.Ipfs, t.datastore, t.Thread, t.sendNotification)
	t.cafe = NewCafeService(t.account, t.Ipfs, t.datastore, t.cafeInbox)

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

		t.threads.Start()
		t.threads.online = true

		t.cafe.Start()
		t.cafe.online = true

		if t.config.Cafe.Host.Open {
			swarmPorts, err := loadSwarmPorts(t.repoPath)
			if err != nil {
				log.Errorf("error loading swarm ports: %s", err)
			} else {
				t.cafe.setAddrs(t.config, *swarmPorts)
			}

			t.cafe.open = true
			t.startCafeApi(t.config.Addresses.CafeAPI)
		}

		go t.runQueues()

		if err := ipfs.PrintSwarmAddrs(t.node); err != nil {
			log.Errorf(err.Error())
		}
		log.Info("node is online")

		// tmp. publish contact for migrated users.
		// this normally only happens when contact details are changed,
		// will be removed at some point in the future.
		if err := t.PublishContact(); err != nil {
			log.Errorf(err.Error())
		}
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
	if err := os.Remove(dsLockFile); err != nil {
	}

	// wipe threads
	t.loadedThreads = nil

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
	return t.cafe.Ping(pid)
}

// UpdateCh returns the node update channel
func (t *Textile) UpdateCh() <-chan Update {
	return t.updates
}

// ThreadUpdateListener returns the thread update channel
func (t *Textile) ThreadUpdateListener() *broadcast.Listener {
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

// SetLogLevels provides node scoped access to the logging system
func (t *Textile) SetLogLevels(logLevels map[string]string) error {
	if _, err := setLogLevels(t.repoPath, logLevels, t.config.Logs.LogToDisk); err != nil {
		return err
	}
	return nil
}

// threadsService returns the threads service
func (t *Textile) threadsService() *ThreadsService {
	return t.threads
}

// cafeService returns the cafe service
func (t *Textile) cafeService() *CafeService {
	return t.cafe
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
		if err := t.cafeInbox.CheckMessages(); err != nil {
			log.Errorf("error checking messages: %s", err)
		}
		t.cafeOutbox.Flush()
	}()
}

// threadByBlock returns the thread owning the given block
func (t *Textile) threadByBlock(block *repo.Block) (*Thread, error) {
	if block == nil {
		return nil, errors.New("block is empty")
	}

	var thrd *Thread
	for _, t := range t.loadedThreads {
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
		RepoPath:           t.repoPath,
		Config:             t.config,
		Node:               t.Ipfs,
		Datastore:          t.datastore,
		Service:            t.threadsService,
		ThreadsOutbox:      t.threadsOutbox,
		CafeOutbox:         t.cafeOutbox,
		SendUpdate:         t.sendThreadUpdate,
		ContactDisplayInfo: t.ContactDisplayInfo,
	}

	thrd, err := NewThread(mod, threadConfig)
	if err != nil {
		return nil, err
	}
	t.loadedThreads = append(t.loadedThreads, thrd)

	return thrd, nil
}

// sendUpdate adds an update to the update channel
func (t *Textile) sendUpdate(update Update) {
	for _, k := range internalThreadKeys {
		if update.Key == k {
			return
		}
	}
	t.updates <- update
}

// sendThreadUpdate adds a thread update to the update channel
func (t *Textile) sendThreadUpdate(update ThreadUpdate) {
	for _, k := range internalThreadKeys {
		if update.ThreadKey == k {
			return
		}
	}
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

// setLogLevels hijacks the ipfs logging system, putting output to files
func setLogLevels(repoPath string, logLevels map[string]string, disk bool) (io.Writer, error) {
	var writer io.Writer
	if disk {
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
	logging.SetAllLoggers(logger.ERROR)

	for key, value := range logLevels {
		if err := logging.SetLogLevel(key, value); err != nil {
			return nil, err
		}
	}
	return writer, nil
}

// getTextileDebugLevels returns a map of textile's logging subsystems set to debug
func getTextileDebugLevels() map[string]string {
	levels := make(map[string]string)
	for _, system := range logging.GetSubsystems() {
		if strings.HasPrefix(system, "tex") {
			levels[system] = "DEBUG"
		}
	}
	return levels
}

// removeLocks force deletes the IPFS repo and SQLite DB lock files
func removeLocks(repoPath string) {
	repoLockFile := filepath.Join(repoPath, fsrepo.LockFile)
	if err := os.Remove(repoLockFile); err != nil {
	}
	dsLockFile := filepath.Join(repoPath, "datastore", "LOCK")
	if err := os.Remove(dsLockFile); err != nil {
	}
}

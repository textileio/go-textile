package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	utilmain "github.com/ipfs/go-ipfs/cmd/ipfs/util"
	oldcmds "github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/repo/config"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/service"
	logger "github.com/whyrusleeping/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = logging.Logger("tex-core")

// kJobFreq how often to flush the message queues
const kJobFreq = time.Second * 60

// kMobileJobFreq how often to flush the message queues on mobile
const kMobileJobFreq = time.Second * 40

// kSyncAccountFreq how often to run account sync
const kSyncAccountFreq = time.Hour

// InitConfig is used to setup a textile node
type InitConfig struct {
	Account         *keypair.Full
	PinCode         string
	RepoPath        string
	SwarmPorts      string
	ApiAddr         string
	CafeApiAddr     string
	GatewayAddr     string
	ProfilingAddr   string
	IsMobile        bool
	IsServer        bool
	LogToDisk       bool
	Debug           bool
	CafeOpen        bool
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
	PinCode           string
	RepoPath          string
	CafeOutboxHandler CafeOutboxHandler
	Debug             bool
}

// Textile is the main Textile node structure
type Textile struct {
	context           oldcmds.Context
	repoPath          string
	config            *config.Config
	account           *keypair.Full
	cancel            context.CancelFunc
	node              *core.IpfsNode
	datastore         repo.Datastore
	started           bool
	loadedThreads     []*Thread
	online            chan struct{}
	done              chan struct{}
	updates           chan *pb.AccountUpdate
	threadUpdates     *broadcast.Broadcaster
	notifications     chan *pb.Notification
	threads           *ThreadsService
	blockOutbox       *BlockOutbox
	cafe              *CafeService
	cafeOutbox        *CafeOutbox
	cafeOutboxHandler CafeOutboxHandler
	cafeInbox         *CafeInbox
	cancelSync        *broadcast.Broadcaster
	mux               sync.Mutex
	writer            io.Writer
}

// common errors
var ErrAccountRequired = fmt.Errorf("account required")
var ErrStarted = fmt.Errorf("node is started")
var ErrStopped = fmt.Errorf("node is stopped")
var ErrOffline = fmt.Errorf("node is offline")

// InitRepo initializes a new node repo
func InitRepo(conf InitConfig) error {
	if fsrepo.IsInitialized(conf.RepoPath) {
		return repo.ErrRepoExists
	}

	if conf.Account == nil {
		return ErrAccountRequired
	}

	logLevel := &pb.LogLevel{
		Systems: make(map[string]pb.LogLevel_Level),
	}
	if conf.Debug {
		logLevel = getTextileDebugLevels()
	}
	_, err := setLogLevels(conf.RepoPath, logLevel, conf.LogToDisk)
	if err != nil {
		return err
	}

	// init repo
	err = repo.Init(conf.RepoPath, conf.IsMobile, conf.IsServer)
	if err != nil {
		return err
	}

	rep, err := fsrepo.Open(conf.RepoPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := rep.Close(); err != nil {
			log.Error(err.Error())
		}
	}()

	// apply ipfs config opts
	err = applySwarmPortConfigOption(rep, conf.SwarmPorts)
	if err != nil {
		return err
	}

	sqliteDb, err := db.Create(conf.RepoPath, conf.PinCode)
	if err != nil {
		return err
	}
	err = sqliteDb.Config().Init(conf.PinCode)
	if err != nil {
		return err
	}
	err = sqliteDb.Config().Configure(conf.Account, time.Now())
	if err != nil {
		return err
	}

	ipfsConf, err := rep.Config()
	if err != nil {
		return err
	}

	// add self as a contact
	err = sqliteDb.Peers().Add(&pb.Peer{
		Id:      ipfsConf.Identity.PeerID,
		Address: conf.Account.Address(),
	})
	if err != nil {
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
	err := repo.Stat(conf.RepoPath)
	if err != nil {
		return nil, err
	}

	// force open the repo and datastore
	removeLocks(conf.RepoPath)

	node := &Textile{
		repoPath:          conf.RepoPath,
		updates:           make(chan *pb.AccountUpdate, 10),
		threadUpdates:     broadcast.NewBroadcaster(10),
		notifications:     make(chan *pb.Notification, 10),
		cafeOutboxHandler: conf.CafeOutboxHandler,
	}

	node.config, err = config.Read(conf.RepoPath)
	if err != nil {
		return nil, err
	}

	logLevel := &pb.LogLevel{
		Systems: make(map[string]pb.LogLevel_Level),
	}
	if conf.Debug {
		logLevel = getTextileDebugLevels()
	}
	node.writer, err = setLogLevels(conf.RepoPath, logLevel, node.config.Logs.LogToDisk)
	if err != nil {
		return nil, err
	}

	// run all minor repo migrations if needed
	err = repo.MigrateUp(conf.RepoPath, conf.PinCode, false)
	if err != nil {
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

	t.online = make(chan struct{})
	t.done = make(chan struct{})

	plugins, err := repo.LoadPlugins(t.repoPath)
	if err != nil {
		return err
	}

	// ensure older peers get latest profiles
	if t.Mobile() {
		err = ensureMobileConfig(t.repoPath)
		if err != nil {
			return err
		}
	}
	if t.Server() {
		err = ensureServerConfig(t.repoPath)
		if err != nil {
			return err
		}
	}

	// raise file descriptor limit
	changed, limit, err := utilmain.ManageFdLimit()
	if err != nil {
		log.Errorf("error setting fd limit: %s", err)
	}
	log.Debugf("fd limit: %d (changed %t)", limit, changed)

	// open db
	err = t.touchDatastore()
	if err != nil {
		return err
	}

	// create queues
	t.cafeInbox = NewCafeInbox(
		t.cafeService,
		t.threadsService,
		t.Ipfs,
		t.datastore)
	t.cafeOutbox = NewCafeOutbox(
		t.Ipfs,
		t.datastore,
		t.cafeOutboxHandler)
	t.blockOutbox = NewBlockOutbox(
		t.threadsService,
		t.Ipfs,
		t.datastore,
		t.cafeOutbox)

	// create services
	t.threads = NewThreadsService(
		t.account,
		t.Ipfs,
		t.datastore,
		t.Thread,
		t.handleThreadAdd,
		t.RemoveThread,
		t.sendNotification)
	t.cafe = NewCafeService(
		t.account,
		t.Ipfs,
		t.datastore,
		t.cafeInbox)

	if t.cafeOutbox.handler == nil {
		t.cafeOutbox.handler = t.cafe
	}

	// start the ipfs node
	log.Debug("creating an ipfs node...")
	err = t.createIPFS(plugins, false)
	if err != nil {
		return err
	}
	go func() {
		defer close(t.online)
		err := t.createIPFS(plugins, true)
		if err != nil {
			log.Errorf("error creating ipfs node: %s", err)
			return
		}

		t.threads.Start()
		t.threads.online = true

		t.cafe.Start()
		t.cafe.online = true

		if t.config.Cafe.Host.Open {
			go func() {
				t.cafe.setAddrs(t.config)
				t.cafe.open = true
				t.startCafeApi(t.config.Addresses.CafeAPI)
			}()
		}

		go t.runJobs()

		err = ipfs.PrintSwarmAddrs(t.node)
		if err != nil {
			log.Errorf(err.Error())
		}
		log.Info("node is online")

		// tmp. publish contact for migrated users.
		// this normally only happens when peer details are changed,
		// will be removed at some point in the future.
		err = t.publishPeer()
		if err != nil {
			log.Errorf(err.Error())
		}
	}()

	for _, mod := range t.datastore.Threads().List().Items {
		_, err = t.loadThread(mod)
		if err != nil {
			if err == ErrThreadLoaded {
				continue
			} else {
				return err
			}
		}
	}

	go t.loadThreadSchemas()

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

	// stop sync if in progress
	if t.cancelSync != nil {
		t.cancelSync.Close()
		t.cancelSync = nil
	}

	// close apis
	err := t.stopCafeApi()
	if err != nil {
		return err
	}

	// close ipfs node
	err = t.node.Close()
	if err != nil {
		return err
	}

	// close db connection
	t.datastore.Close()
	dsLockFile := filepath.Join(t.repoPath, "datastore", "LOCK")
	_ = os.Remove(dsLockFile)

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
	return t.started && t.node.IsOnline
}

// Mobile returns whether or not node is configured for a mobile device
func (t *Textile) Mobile() bool {
	return t.config.IsMobile
}

// Server returns whether or not node is configured for a server
func (t *Textile) Server() bool {
	return t.config.IsServer
}

// Datastore returns the underlying sqlite datastore interface
func (t *Textile) Datastore() repo.Datastore {
	return t.datastore
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

// UpdateCh returns the account update channel
func (t *Textile) UpdateCh() <-chan *pb.AccountUpdate {
	return t.updates
}

// ThreadUpdateListener returns the thread update channel
func (t *Textile) ThreadUpdateListener() *broadcast.Listener {
	return t.threadUpdates.Listen()
}

// NotificationsCh returns the notifications channel
func (t *Textile) NotificationCh() <-chan *pb.Notification {
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

// SetLogLevel provides node scoped access to the logging system
func (t *Textile) SetLogLevel(level *pb.LogLevel) error {
	_, err := setLogLevels(t.repoPath, level, t.config.Logs.LogToDisk)
	return err
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
func (t *Textile) createIPFS(plugins *loader.PluginLoader, online bool) error {
	rep, err := fsrepo.Open(t.repoPath)
	if err != nil {
		return err
	}

	routing := libp2p.DHTOption
	if t.Mobile() {
		routing = libp2p.DHTClientOption
	}

	cctx, _ := context.WithCancel(context.Background())
	nd, err := core.NewNode(cctx, &core.BuildCfg{
		Repo:      rep,
		Permanent: true, // temporary way to signify that node is permanent
		Online:    online,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
			"mplex":  true,
		},
		Routing: routing,
	})
	if err != nil {
		return err
	}
	nd.IsDaemon = true

	if t.node != nil {
		err = t.node.Close()
		if err != nil {
			return err
		}
	}
	t.node = nd

	return nil
}

// runJobs runs each message queue
func (t *Textile) runJobs() {
	var freq time.Duration
	if t.Mobile() {
		freq = kMobileJobFreq
	} else {
		freq = kJobFreq
	}

	tick := time.NewTicker(freq)
	defer tick.Stop()

	go t.flushQueues()
	t.maybeSyncAccount()
	t.runGC()

	for {
		select {
		case <-tick.C:
			if err := t.touchDatastore(); err != nil {
				log.Error(err.Error())
				return
			}

			go t.flushQueues()
			t.maybeSyncAccount()

		case <-t.done:
			return
		}
	}
}

// flushQueues flushes each message queue
func (t *Textile) flushQueues() {
	t.cafeOutbox.Flush()
	t.blockOutbox.Flush()
	if err := t.cafeInbox.CheckMessages(); err != nil {
		log.Errorf("error checking messages: %s", err)
	}
}

// threadByBlock returns the thread owning the given block
func (t *Textile) threadByBlock(block *pb.Block) (*Thread, error) {
	if block == nil {
		return nil, fmt.Errorf("block is empty")
	}

	var thrd *Thread
	for _, l := range t.loadedThreads {
		if l.Id == block.Thread {
			thrd = l
			break
		}
	}
	if thrd == nil {
		return nil, fmt.Errorf("could not find thread: %s", block.Thread)
	}
	return thrd, nil
}

// loadThread loads a thread into memory from the given on-disk model
func (t *Textile) loadThread(mod *pb.Thread) (*Thread, error) {
	if loaded := t.Thread(mod.Id); loaded != nil {
		return nil, ErrThreadLoaded
	}

	thrd, err := NewThread(mod, &ThreadConfig{
		RepoPath:    t.repoPath,
		Config:      t.config,
		Account:     t.account,
		Node:        t.Ipfs,
		Datastore:   t.datastore,
		Service:     t.threadsService,
		BlockOutbox: t.blockOutbox,
		CafeOutbox:  t.cafeOutbox,
		AddPeer:     t.addPeer,
		PushUpdate:  t.sendThreadUpdate,
	})
	if err != nil {
		return nil, err
	}
	t.loadedThreads = append(t.loadedThreads, thrd)

	return thrd, nil
}

// loadThreadSchemas loads thread schemas that were not found locally during startup
func (t *Textile) loadThreadSchemas() {
	<-t.online
	var err error
	for _, l := range t.loadedThreads {
		err = l.loadSchema()
		if err != nil {
			log.Errorf("unable to load schema %s: %s", l.schemaId, err)
		}
	}
}

// sendUpdate sends an update to the update channel
func (t *Textile) sendUpdate(update *pb.AccountUpdate) {
	if (update.Type == pb.AccountUpdate_THREAD_ADDED ||
		update.Type == pb.AccountUpdate_THREAD_REMOVED) &&
		update.Id == t.config.Account.Thread {
		return
	}
	t.updates <- update
}

// sendThreadUpdate sends a feed item to the update channel
func (t *Textile) sendThreadUpdate(block *pb.Block, key string) {
	if key == t.account.Address() {
		return
	}

	update, err := t.feedItem(block, feedItemOpts{})
	if err != nil {
		log.Errorf("error building thread update: %s", err)
		return
	}

	t.threadUpdates.Send(update)
}

// sendNotification adds a notification to the notification channel
func (t *Textile) sendNotification(note *pb.Notification) error {
	if err := t.datastore.Notifications().Add(note); err != nil {
		return err
	}

	t.notifications <- t.NotificationView(note)
	return nil
}

// touchDatastore ensures that we have a good db connection
func (t *Textile) touchDatastore() error {
	if err := t.datastore.Ping(); err != nil {
		log.Debug("re-opening datastore...")

		sqliteDB, err := db.Create(t.repoPath, "")
		if err != nil {
			return err
		}
		t.datastore = sqliteDB
	}

	return nil
}

// runGC periodically runs repo blockstore GC
func (t *Textile) runGC() {
	errc := make(chan error)
	go func() {
		errc <- corerepo.PeriodicGC(t.node.Context(), t.node)
		close(errc)
	}()
	go func() {
		for {
			select {
			case <-t.node.Context().Done():
				log.Debug("blockstore GC shutdown")
				return
			case err, ok := <-errc:
				if !ok {
					return
				}
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}()
}

// setLogLevels hijacks the ipfs logging system, putting output to files
func setLogLevels(repoPath string, level *pb.LogLevel, disk bool) (io.Writer, error) {
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

	var err error
	for key, value := range level.Systems {
		err = logging.SetLogLevel(key, value.String())
		if err != nil {
			return nil, err
		}
	}
	return writer, nil
}

// getTextileDebugLevels returns a map of textile's logging subsystems set to debug
func getTextileDebugLevels() *pb.LogLevel {
	levels := make(map[string]pb.LogLevel_Level)
	for _, system := range logging.GetSubsystems() {
		if strings.HasPrefix(system, "tex") {
			levels[system] = pb.LogLevel_DEBUG
		}
	}
	return &pb.LogLevel{Systems: levels}
}

// removeLocks force deletes the IPFS repo and SQLite DB lock files
func removeLocks(repoPath string) {
	repoLockFile := filepath.Join(repoPath, fsrepo.LockFile)
	_ = os.Remove(repoLockFile)
	dsLockFile := filepath.Join(repoPath, "datastore", "LOCK")
	_ = os.Remove(dsLockFile)
}

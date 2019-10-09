package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	utilmain "github.com/ipfs/go-ipfs/cmd/ipfs/util"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/bootstrap"
	"github.com/ipfs/go-ipfs/core/corerepo"
	corenode "github.com/ipfs/go-ipfs/core/node"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
	"github.com/ipfs/go-metrics-interface"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/repo/config"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/service"
	"github.com/textileio/go-textile/util"
	logger "github.com/whyrusleeping/go-logging"
	"go.uber.org/fx"
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
	BaseRepoPath    string
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
	CheckMessages     func() error
	Debug             bool
}

// Textile is the main Textile node structure
type Textile struct {
	repoPath          string
	pinCode           string
	config            *config.Config
	account           *keypair.Full
	ctx               context.Context
	stop              func() error
	node              *core.IpfsNode
	started           bool
	datastore         repo.Datastore
	loadedThreads     []*Thread
	online            chan struct{}
	done              chan struct{}
	updates           chan *pb.AccountUpdate
	threadUpdates     *broadcast.Broadcaster
	notifications     chan *pb.Notification
	threads           *ThreadsService
	blockOutbox       *BlockOutbox
	blockDownloads    *BlockDownloads
	cafe              *CafeService
	cafeOutbox        *CafeOutbox
	cafeOutboxHandler CafeOutboxHandler
	cafeInbox         *CafeInbox
	checkMessages     func() error
	cancelSync        *broadcast.Broadcaster
	lock              sync.Mutex
	writer            io.Writer
}

// common errors
var ErrAccountRequired = fmt.Errorf("account required")
var ErrStarted = fmt.Errorf("node is started")
var ErrStopped = fmt.Errorf("node is stopped")
var ErrOffline = fmt.Errorf("node is offline")
var ErrMissingRepoConfig = fmt.Errorf("you must specify InitConfig.RepoPath or InitConfig.BaseRepoPath and InitConfig.Account")

// Repo returns the actual location of the configured repo
func (conf InitConfig) Repo() (string, error) {
	if len(conf.RepoPath) > 0 {
		return conf.RepoPath, nil
	} else if len(conf.BaseRepoPath) > 0 && conf.Account != nil {
		return path.Join(conf.BaseRepoPath, conf.Account.Address()), nil
	} else {
		return "", ErrMissingRepoConfig
	}

}

// RepoExists return whether or not the configured repo already exists
func (conf InitConfig) RepoExists() (bool, error) {
	repoPath, err := conf.Repo()
	if err != nil {
		return false, err
	}
	return RepoExists(repoPath), nil
}

// RepoExists return whether or not the repo at repoPath exists
func RepoExists(repoPath string) bool {
	return fsrepo.IsInitialized(repoPath)
}

// AccountRepoExists return whether or not the repo at repoPath exists
func AccountRepoExists(baseRepoPath string, accountAddress string) bool {
	return fsrepo.IsInitialized(path.Join(baseRepoPath, accountAddress))
}

// InitRepo initializes a new node repo
func InitRepo(conf InitConfig) error {
	exists, err := conf.RepoExists()
	if err != nil {
		return err
	}

	if exists {
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

	repoPath, err := conf.Repo()
	if err != nil {
		return err
	}

	_, err = setLogLevels(repoPath, logLevel, conf.LogToDisk, !conf.IsMobile)
	if err != nil {
		return err
	}

	// init repo
	err = repo.Init(repoPath, conf.IsMobile, conf.IsServer)
	if err != nil {
		return err
	}

	rep, err := fsrepo.Open(repoPath)
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

	sqliteDb, err := db.Create(repoPath, conf.PinCode)
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
		pinCode:           conf.PinCode,
		updates:           make(chan *pb.AccountUpdate, 10),
		threadUpdates:     broadcast.NewBroadcaster(10),
		notifications:     make(chan *pb.Notification, 10),
		cafeOutboxHandler: conf.CafeOutboxHandler,
		checkMessages:     conf.CheckMessages,
	}

	node.config, err = config.Read(node.repoPath)
	if err != nil {
		return nil, err
	}

	logLevel := &pb.LogLevel{
		Systems: make(map[string]pb.LogLevel_Level),
	}
	if conf.Debug {
		logLevel = getTextileDebugLevels()
	}
	node.writer, err = setLogLevels(node.repoPath, logLevel,
		node.config.Logs.LogToDisk, !node.config.IsMobile)
	if err != nil {
		return nil, err
	}

	// run all minor repo migrations if needed
	err = repo.MigrateUp(node.repoPath, node.pinCode, false)
	if err != nil {
		return nil, err
	}

	sqliteDb, err := db.Create(node.repoPath, node.pinCode)
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
	t.lock.Lock()
	if t.started {
		t.lock.Unlock()
		return ErrStarted
	}
	log.Info("starting node...")

	t.online = make(chan struct{})
	t.done = make(chan struct{})

	_, err := repo.LoadPlugins(t.repoPath)
	if err != nil {
		return err
	}

	// ensure older peers get latest profiles
	if t.Mobile() {
		err = ensureProfile(mobileProfile, t.repoPath)
	} else if t.Server() {
		err = ensureProfile(serverProfile, t.repoPath)
	} else {
		err = ensureProfile(desktopProfile, t.repoPath)
	}
	if err != nil {
		return err
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
	t.blockDownloads = NewBlockDownloads(
		t.Ipfs,
		t.datastore,
		t.Thread)
	t.cafeInbox = NewCafeInbox(
		t.cafeService,
		t.threadsService,
		t.Ipfs,
		t.datastore)
	t.cafeOutbox = NewCafeOutbox(
		t.Ipfs,
		t.datastore,
		t.cafeOutboxHandler,
		t.FlushBlocks)
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
	err = t.createNode()
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			close(t.online)
			t.lock.Unlock()
			t.runJobs()
		}()

		err = t.node.Bootstrap(bootstrap.DefaultBootstrapConfig)
		if err != nil {
			log.Errorf("error bootstrapping ipfs node: %s", err)
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

		err = ipfs.PrintSwarmAddrs(t.node)
		if err != nil {
			log.Errorf(err.Error())
		}
		log.Info("node is online")

		// ensure the peer table is not empty by adding our bootstraps
		boots, err := config.TextileBootstrapPeers()
		if err != nil {
			log.Errorf(err.Error())
		}
		for _, p := range boots {
			t.node.Peerstore.AddAddrs(p.ID, p.Addrs, peerstore.PermanentAddrTTL)
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

type loggingWaitGroup struct {
	n  string
	wg sync.WaitGroup
}

func (lwg *loggingWaitGroup) Add(delta int, src string) {
	log.Debugf("%s wait added delta %d (src=%s)", lwg.n, delta, src)
	lwg.wg.Add(delta)
}

func (lwg *loggingWaitGroup) Done(src string) {
	log.Debugf("%s wait done (src=%s)", lwg.n, src)
	lwg.wg.Done()
}

func (lwg *loggingWaitGroup) Wait(src string) {
	log.Debugf("%s waiting (src=%s)", lwg.n, src)
	lwg.wg.Wait()
}

// stopGroup is used to block shutdown. Workers must add to the wait group counter
// before Stop blocks on wait.
var stopGroup = loggingWaitGroup{n: "stop"}

// Stop destroys the ipfs node and shutsdown textile services
func (t *Textile) Stop() error {
	stopGroup.Wait("Stop")

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
	err = t.stop()
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

// WaitAdd add delta wait to the stop wait group
func (t *Textile) WaitAdd(delta int, src string) {
	stopGroup.Add(delta, src)
}

// WaitDone marks a wait as done in the stop wait group
func (t *Textile) WaitDone(src string) {
	stopGroup.Done(src)
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

// Inbox returns the cafe inbox
func (t *Textile) Inbox() *CafeInbox {
	return t.cafeInbox
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

// Publish sends 'data' to 'topic'
func (t *Textile) Publish(payload []byte, topic string) error {
	return ipfs.Publish(t.node, topic, payload)
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
func (t *Textile) SetLogLevel(level *pb.LogLevel, color bool) error {
	_, err := setLogLevels(t.repoPath, level, t.config.Logs.LogToDisk, color)
	return err
}

// FlushBlocks flushes the block message outbox
func (t *Textile) FlushBlocks() {
	query := fmt.Sprintf("status=%d", pb.Block_QUEUED)
	queued := t.datastore.Blocks().List("", -1, query)
	sort.SliceStable(queued.Items, func(i, j int) bool {
		return util.ProtoTime(queued.Items[i].Date).Before(
			util.ProtoTime(queued.Items[j].Date))
	})
	wg := sync.WaitGroup{}
	for _, block := range queued.Items {
		if t.datastore.CafeRequests().SyncGroupComplete(block.Id) {
			wg.Add(1)
			go func(block *pb.Block) {
				var posted bool
				defer func() {
					t.blockOutbox.Flush()
					if posted {
						go t.cafeOutbox.Flush(true)
					} else if t.cafeOutbox.handler != nil {
						t.cafeOutbox.handler.Flush()
					}
					wg.Done()
				}()

				thread := t.Thread(block.Thread)
				if thread == nil {
					return
				}

				// if this is not a join, ensure it will hava at least one parent
				if block.Type != pb.Block_JOIN {
					heads, err := thread.Heads()
					if err != nil {
						log.Warningf("error getting heads: %s", err)
						return
					}
					if len(heads) == 0 {
						return
					}
				}

				err := thread.post(block)
				if err != nil {
					log.Errorf("error posting block %s: %s", block.Id, err)
					if block.Attempts+1 >= maxDownloadAttempts {
						err = t.datastore.Blocks().Delete(block.Id)
					} else {
						err = t.datastore.Blocks().AddAttempt(block.Id)
					}
					if err != nil {
						log.Errorf("error handling post error: %s", err)
					}
					return
				}
				posted = true

				err = t.datastore.CafeRequests().DeleteBySyncGroup(block.Id)
				if err != nil {
					log.Error(err)
				} else {
					log.Debugf("deleted sync group: %s", block.Id)
				}
			}(block)
		}
	}
	wg.Wait()
}

// FlushCafes flushes the cafe request outbox
func (t *Textile) FlushCafes() {
	stopGroup.Add(1, "FlushCafes")
	go func() {
		defer stopGroup.Done("FlushCafes")
		t.cafeOutbox.Flush(false)
	}()
}

// threadsService returns the threads service
func (t *Textile) threadsService() *ThreadsService {
	return t.threads
}

// cafeService returns the cafe service
func (t *Textile) cafeService() *CafeService {
	return t.cafe
}

// createNode constructs an IpfsNode
func (t *Textile) createNode() error {
	rep, err := fsrepo.Open(t.repoPath)
	if err != nil {
		return err
	}

	routing := libp2p.DHTOption
	if t.Mobile() {
		routing = libp2p.DHTClientOption
	}

	cfg := &core.BuildCfg{
		Repo:      rep,
		Permanent: true, // temporary way to signify that node is permanent
		Online:    true,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
			"mplex":  true,
		},
		Routing: routing,
	}

	ctx := context.Background()
	ctx = metrics.CtxScope(ctx, "ipfs")

	n := &core.IpfsNode{}
	t.ctx = ctx

	app := fx.New(
		corenode.IPFS(ctx, cfg),

		fx.NopLogger,
		fx.Extract(n),
	)

	var once sync.Once
	var stopErr error
	t.stop = func() error {
		once.Do(func() {
			stopErr = app.Stop(context.Background())
		})
		return stopErr
	}
	n.IsOnline = cfg.Online
	n.IsDaemon = true

	go func() {
		// Note that some services use contexts to signal shutting down, which is
		// very suboptimal. This needs to be here until that's addressed somehow
		<-ctx.Done()
		err := t.stop()
		if err != nil {
			log.Error("failure on stop: ", err)
		}
	}()

	if app.Err() != nil {
		return app.Err()
	}

	if err := app.Start(ctx); err != nil {
		return err
	}

	t.node = n

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

	if t.Mobile() {
		t.runConditionalGC()
	} else {
		t.runPeriodicGC()
	}

	for {
		select {
		case <-tick.C:
			go t.flushQueues()
			t.maybeSyncAccount()

		case <-t.done:
			return
		}
	}
}

// flushQueues flushes each message queue
func (t *Textile) flushQueues() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.cafeOutbox.Flush(false)
	var err error
	if t.checkMessages != nil {
		err = t.checkMessages()
	} else {
		err = t.cafeInbox.CheckMessages()
	}
	if err != nil {
		log.Errorf("error checking messages: %s", err)
	}
	t.blockDownloads.Flush()
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
		RepoPath:       t.repoPath,
		Config:         t.config,
		Account:        t.account,
		Node:           t.Ipfs,
		Datastore:      t.datastore,
		Service:        t.threadsService,
		BlockOutbox:    t.blockOutbox,
		BlockDownloads: t.blockDownloads,
		CafeOutbox:     t.cafeOutbox,
		AddPeer:        t.AddPeer,
		PushUpdate:     t.sendThreadUpdate,
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

		sqliteDB, err := db.Create(t.repoPath, t.pinCode)
		if err != nil {
			return err
		}
		t.datastore = sqliteDB
	}

	return nil
}

// runPeriodicGC periodically runs repo blockstore GC
func (t *Textile) runPeriodicGC() {
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

// runConditionalGC runs repo blockstore GC once, if needed
func (t *Textile) runConditionalGC() {
	err := corerepo.ConditionalGC(t.node.Context(), t.node, 0)
	if err != nil {
		log.Errorf("error running conditional gc: %s", err)
	}
}

// setLogLevels hijacks the ipfs logging system, putting output to files
func setLogLevels(repoPath string, level *pb.LogLevel, disk bool, color bool) (io.Writer, error) {
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

	var format string
	if color {
		format = logging.LogFormats["color"]
	} else {
		format = logging.LogFormats["nocolor"]
	}
	logger.SetFormatter(logger.MustStringFormatter(format))
	logging.SetAllLoggers(logger.ERROR)

	var err error
	for key, value := range level.Systems {
		if key == "*" {
			for _, s := range logging.GetSubsystems() {
				err = logging.SetLogLevel(s, value.String())
				if err != nil {
					return nil, err
				}
			}
		}
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

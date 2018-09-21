package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/archive"
	"github.com/textileio/textile-go/core/thread"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/net"
	serv "github.com/textileio/textile-go/net/service"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/storage"
	"gopkg.in/natefinch/lumberjack.v2"
	"gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	utilmain "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/commands"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/config"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/fsrepo"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
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

const Version = "0.1.9"

// Node is the single Textile instance
var Node *Textile

type Config struct {
	RepoPath string
	PinCode  string

	SwarmPorts string

	IsMobile bool
	IsServer bool

	LogLevel logging.Level
	LogFiles bool

	CafeAddr string
}

type Update struct {
	Id   string     `json:"id"`
	Name string     `json:"name"`
	Type UpdateType `json:"type"`
}

type UpdateType int

const (
	ThreadAdded UpdateType = iota
	ThreadRemoved
	DeviceAdded
	DeviceRemoved
)

// AddDataResult wraps added data content id and key
type AddDataResult struct {
	Id      string           `json:"id"`
	Key     string           `json:"key"`
	Archive *archive.Archive `json:"archive,omitempty"`
}

// Textile is the main Textile node structure
type Textile struct {
	version            string
	context            oldcmds.Context
	repoPath           string
	cancel             context.CancelFunc
	ipfs               *core.IpfsNode
	datastore          repo.Datastore
	service            *serv.TextileService
	cafeAddr           string
	isMobile           bool
	started            bool
	threads            []*thread.Thread
	online             chan struct{}
	done               chan struct{}
	updates            chan Update
	threadUpdates      chan thread.Update
	notifications      chan repo.Notification
	messageStorage     storage.OfflineMessagingStorage
	messageRetriever   *net.MessageRetriever
	pointerRepublisher *net.PointerRepublisher
	pinner             *net.Pinner
	mux                sync.Mutex
}

const pingTimeout = time.Second * 10

var ErrStarted = errors.New("node is started")
var ErrStopped = errors.New("node is stopped")
var ErrOffline = errors.New("node is offline")
var ErrProfileNotFound = errors.New("profile not found")
var ErrThreadLoaded = errors.New("thread is loaded")
var ErrNoCafeHost = errors.New("cafe host address is not set")

func NewTextile(config Config) (*Textile, error) {
	repoLockFile := filepath.Join(config.RepoPath, fsrepo.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(config.RepoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// log handling
	var backendFile *logging.LogBackend
	if config.LogFiles {
		logger := &lumberjack.Logger{
			Filename:   path.Join(config.RepoPath, "logs", "textile.log"),
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

	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, "")
	if err != nil {
		return nil, err
	}

	// we may be running in an uninitialized state.
	err = repo.DoInit(config.RepoPath, Version, func() error {
		sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return err
		}
		if err := sqliteDB.Config().Init(config.PinCode); err != nil {
			return err
		}
		return sqliteDB.Config().Configure(sk, time.Now())
	})
	if err != nil && err != repo.ErrRepoExists {
		return nil, err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	rep, err := fsrepo.Open(config.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return nil, err
	}

	// if a specific swarm port was selected, set it in the config
	if err := applySwarmPortConfigOption(rep, config.SwarmPorts); err != nil {
		return nil, err
	}

	// ensure bootstrap addresses are latest in config (without wiping repo)
	if err := ensureBootstrapConfig(rep); err != nil {
		return nil, err
	}

	// if this is a server node, apply the ipfs server profile
	if err := applyServerConfigOption(rep, config.IsServer); err != nil {
		return nil, err
	}

	return &Textile{
		version:   Version,
		repoPath:  config.RepoPath,
		datastore: sqliteDB,
		isMobile:  config.IsMobile,
		cafeAddr:  config.CafeAddr,
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

		pk, err := t.GetPeerPubKey()
		if err != nil {
			log.Errorf("error loading peer pk: %s", err)
			return
		}
		pks, err := ipfs.EncodeKey(pk)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Infof("wallet is started, peer pk: %s", pks)
	}()
	log.Info("starting wallet...")
	t.online = make(chan struct{})
	t.updates = make(chan Update, 10)
	t.threadUpdates = make(chan thread.Update, 10)
	t.notifications = make(chan repo.Notification, 10)

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		log.Errorf("setting file descriptor limit: %s", err)
	}

	// check db
	if err := t.touchDatastore(); err != nil {
		return err
	}

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

		// wait for dht to bootstrap
		<-dht.DefaultBootstrapConfig.DoneChan

		// set offline message storage
		t.messageStorage = storage.NewCafeStorage(t.ipfs, t.repoPath, func(id *cid.Cid) error {
			if t.pinner == nil {
				return nil
			}
			tokens, err := t.GetCafeTokens(false)
			if err != nil {
				return err
			}
			hash := id.Hash().B58String()
			if err := net.Pin(t.ipfs, hash, tokens, t.pinner.Url()); err != nil {
				if err == net.ErrTokenExpired {
					tokens, err := t.GetCafeTokens(true)
					if err != nil {
						return err
					}
					return net.Pin(t.ipfs, hash, tokens, t.pinner.Url())
				} else {
					return err
				}
			}
			return nil
		})

		// service is now configurable
		t.service = serv.NewService(t.ipfs, t.datastore, t.GetThread, t.sendNotification)

		// build the message retriever
		mrCfg := net.MRConfig{
			Datastore: t.datastore,
			Ipfs:      t.ipfs,
			Service:   t.service,
			PrefixLen: 14,
			SendAck:   t.sendOfflineAck,
			SendError: t.sendError,
		}
		t.messageRetriever = net.NewMessageRetriever(mrCfg)

		// build the pointer republisher
		t.pointerRepublisher = net.NewPointerRepublisher(t.ipfs, t.datastore)

		// start jobs if not mobile
		if !t.isMobile {
			go t.messageRetriever.Run()
			go t.pointerRepublisher.Run()
		} else {
			go t.pointerRepublisher.Republish()
		}

		// print swarm addresses
		if err := ipfs.PrintSwarmAddrs(t.ipfs); err != nil {
			log.Errorf("failed to read listening addresses: %s", err)
		}
		log.Info("wallet is online")
	}()

	// build a pin requester
	if t.GetCafeApiAddr() != "" {
		pinnerCfg := &net.PinnerConfig{
			Datastore: t.datastore,
			Ipfs: func() *core.IpfsNode {
				return t.ipfs
			},
			Url:       fmt.Sprintf("%s/pin", t.GetCafeApiAddr()),
			GetTokens: t.GetCafeTokens,
		}
		t.pinner = net.NewPinner(pinnerCfg)

		// start ticker job if not mobile
		if !t.isMobile {
			go t.pinner.Run()
		} else {
			go t.pinner.Pin()
		}
	}

	// re-pub profile
	go func() {
		<-t.Online()
		if _, err := t.PublishProfile(nil); err != nil {
			log.Errorf("error publishing profile: %s", err)
		}
	}()

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
	log.Info("stopping wallet...")

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

	// shutdown message retriever
	select {
	case t.messageRetriever.DoneChan <- struct{}{}:
	default:
	}

	// close update channels
	close(t.updates)
	close(t.threadUpdates)
	close(t.notifications)

	log.Info("wallet is stopped")

	return nil
}

func (t *Textile) Started() bool {
	return t.started
}

func (t *Textile) IsOnline() bool {
	if t.ipfs == nil {
		return false
	}
	return t.started && t.ipfs.OnlineMode()
}

func (t *Textile) Version() string {
	return t.version
}

func (t *Textile) Ipfs() *core.IpfsNode {
	return t.ipfs
}

func (t *Textile) FetchMessages() error {
	if !t.IsOnline() {
		return ErrOffline
	}
	if t.messageRetriever.IsFetching() {
		return net.ErrFetching
	}
	go t.messageRetriever.FetchPointers()
	return nil
}

func (t *Textile) Online() <-chan struct{} {
	return t.online
}

func (t *Textile) Done() <-chan struct{} {
	return t.done
}

func (t *Textile) Updates() <-chan Update {
	return t.updates
}

func (t *Textile) ThreadUpdates() <-chan thread.Update {
	return t.threadUpdates
}

func (t *Textile) Notifications() <-chan repo.Notification {
	return t.notifications
}

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
	if t.isMobile {
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
		Peers:         t.datastore.Peers,
		Notifications: t.datastore.Notifications,
		GetHead: func() (string, error) {
			m := t.datastore.Threads().Get(id)
			if m == nil {
				return "", errors.New(fmt.Sprintf("could not re-load thread: %s", id))
			}
			return m.Head, nil
		},
		UpdateHead: func(head string) error {
			if err := t.datastore.Threads().UpdateHead(id, head); err != nil {
				return err
			}
			return nil
		},
		Send:          t.SendMessage,
		NewEnvelope:   t.NewEnvelope,
		PutPinRequest: t.putPinRequest,
		GetUsername:   t.GetUsername,
		SendUpdate:    t.sendThreadUpdate,
	}
	thrd, err := thread.NewThread(mod, threadConfig)
	if err != nil {
		return nil, err
	}
	t.threads = append(t.threads, thrd)
	return thrd, nil
}

// putPinRequest adds a pin request to the pinner
func (t *Textile) putPinRequest(id string) error {
	if t.pinner == nil {
		return nil
	}
	return t.pinner.Put(id)
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

// touchDB ensures that we have a good db connection
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

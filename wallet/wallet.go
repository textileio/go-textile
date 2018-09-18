package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/core/cafe"
	"github.com/textileio/textile-go/net"
	serv "github.com/textileio/textile-go/net/service"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/storage"
	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/wallet/thread"
	"gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	utilmain "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/commands"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/config"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/fsrepo"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"os"
	"path/filepath"
	"time"
)

var log = logging.MustGetLogger("wallet")

type Config struct {
	Version  string
	RepoPath string
	Mnemonic *string

	SwarmPorts string

	IsMobile bool
	IsServer bool

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
	Id      string          `json:"id"`
	Key     string          `json:"key"`
	Archive *client.Archive `json:"archive,omitempty"`
}

// Wallet is the main Textile node structure
type Wallet struct {
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
}

const pingTimeout = time.Second * 10

var ErrStarted = errors.New("node is started")
var ErrStopped = errors.New("node is stopped")
var ErrOffline = errors.New("node is offline")
var ErrProfileNotFound = errors.New("profile not found")
var ErrThreadLoaded = errors.New("thread is loaded")
var ErrNoCafeHost = errors.New("cafe host address is not set")

func NewWallet(config Config) (*Wallet, string, error) {
	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, "")
	if err != nil {
		return nil, "", err
	}

	// we may be running in an uninitialized state.
	mnemonic, err := repo.DoInit(config.RepoPath, config.Version, config.Mnemonic, sqliteDB.Config().Init, sqliteDB.Config().Configure)
	if err != nil && err != repo.ErrRepoExists {
		return nil, "", err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	rep, err := fsrepo.Open(config.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return nil, "", err
	}

	// if a specific swarm port was selected, set it in the config
	if err := applySwarmPortConfigOption(rep, config.SwarmPorts); err != nil {
		return nil, "", err
	}

	// ensure bootstrap addresses are latest in config (without wiping repo)
	if err := ensureBootstrapConfig(rep); err != nil {
		return nil, "", err
	}

	// if this is a server node, apply the ipfs server profile
	if err := applyServerConfigOption(rep, config.IsServer); err != nil {
		return nil, "", err
	}

	return &Wallet{
		version:   config.Version,
		repoPath:  config.RepoPath,
		datastore: sqliteDB,
		isMobile:  config.IsMobile,
		cafeAddr:  config.CafeAddr,
	}, mnemonic, nil
}

// Start
func (w *Wallet) Start() error {
	if w.started {
		return ErrStarted
	}
	defer func() {
		w.done = make(chan struct{})
		w.started = true

		pk, err := w.GetPeerPubKey()
		if err != nil {
			log.Errorf("error loading peer pk: %s", err)
			return
		}
		pks, err := util.EncodeKey(pk)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Infof("wallet is started, peer pk: %s", pks)
	}()
	log.Info("starting wallet...")
	w.online = make(chan struct{})
	w.updates = make(chan Update, 10)
	w.threadUpdates = make(chan thread.Update, 10)
	w.notifications = make(chan repo.Notification, 10)

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		log.Errorf("setting file descriptor limit: %s", err)
	}

	// check db
	if err := w.touchDatastore(); err != nil {
		return err
	}

	// start the ipfs node
	log.Debug("creating an ipfs node...")
	if err := w.createIPFS(false); err != nil {
		log.Errorf("error creating offline ipfs node: %s", err)
		return err
	}
	go func() {
		defer close(w.online)
		if err := w.createIPFS(true); err != nil {
			log.Errorf("error creating online ipfs node: %s", err)
			return
		}

		// wait for dht to bootstrap
		<-dht.DefaultBootstrapConfig.DoneChan

		// set offline message storage
		w.messageStorage = storage.NewCafeStorage(w.ipfs, w.repoPath, func(id *cid.Cid) error {
			if w.pinner == nil {
				return nil
			}
			tokens, err := w.GetCafeTokens(false)
			if err != nil {
				return err
			}
			hash := id.Hash().B58String()
			if err := net.Pin(w.ipfs, hash, tokens, w.pinner.Url()); err != nil {
				if err == net.ErrTokenExpired {
					tokens, err := w.GetCafeTokens(true)
					if err != nil {
						return err
					}
					return net.Pin(w.ipfs, hash, tokens, w.pinner.Url())
				} else {
					return err
				}
			}
			return nil
		})

		// service is now configurable
		w.service = serv.NewService(w.ipfs, w.datastore, w.GetThread, w.sendNotification)

		// build the message retriever
		mrCfg := net.MRConfig{
			Datastore: w.datastore,
			Ipfs:      w.ipfs,
			Service:   w.service,
			PrefixLen: 14,
			SendAck:   w.sendOfflineAck,
			SendError: w.sendError,
		}
		w.messageRetriever = net.NewMessageRetriever(mrCfg)

		// build the pointer republisher
		w.pointerRepublisher = net.NewPointerRepublisher(w.ipfs, w.datastore)

		// start jobs if not mobile
		if !w.isMobile {
			go w.messageRetriever.Run()
			go w.pointerRepublisher.Run()
		} else {
			go w.pointerRepublisher.Republish()
		}

		// print swarm addresses
		if err := util.PrintSwarmAddrs(w.ipfs); err != nil {
			log.Errorf("failed to read listening addresses: %s", err)
		}
		log.Info("wallet is online")
	}()

	// build a pin requester
	if w.GetCafeApiAddr() != "" {
		pinnerCfg := &net.PinnerConfig{
			Datastore: w.datastore,
			Ipfs: func() *core.IpfsNode {
				return w.ipfs
			},
			Url:       fmt.Sprintf("%s/pin", w.GetCafeApiAddr()),
			GetTokens: w.GetCafeTokens,
		}
		w.pinner = net.NewPinner(pinnerCfg)

		// start ticker job if not mobile
		if !w.isMobile {
			go w.pinner.Run()
		} else {
			go w.pinner.Pin()
		}
	}

	// re-pub profile
	go func() {
		<-w.Online()
		if _, err := w.PublishProfile(nil); err != nil {
			log.Errorf("error publishing profile: %s", err)
		}
	}()

	// setup threads
	for _, mod := range w.datastore.Threads().List("") {
		_, err := w.loadThread(&mod)
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
func (w *Wallet) Stop() error {
	if !w.started {
		return ErrStopped
	}
	defer func() {
		w.started = false
		close(w.done)
	}()
	log.Info("stopping wallet...")

	// close ipfs node
	w.context.Close()
	w.cancel()
	if err := w.ipfs.Close(); err != nil {
		log.Errorf("error closing ipfs node: %s", err)
		return err
	}

	// close db connection
	w.datastore.Close()
	dsLockFile := filepath.Join(w.repoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// wipe threads
	w.threads = nil

	// shutdown message retriever
	select {
	case w.messageRetriever.DoneChan <- struct{}{}:
	default:
	}

	// close update channels
	close(w.updates)
	close(w.threadUpdates)
	close(w.notifications)

	log.Info("wallet is stopped")

	return nil
}

func (w *Wallet) Started() bool {
	return w.started
}

func (w *Wallet) IsOnline() bool {
	if w.ipfs == nil {
		return false
	}
	return w.started && w.ipfs.OnlineMode()
}

func (w *Wallet) Version() string {
	return w.version
}

func (w *Wallet) Ipfs() *core.IpfsNode {
	return w.ipfs
}

func (w *Wallet) FetchMessages() error {
	if !w.IsOnline() {
		return ErrOffline
	}
	if w.messageRetriever.IsFetching() {
		return net.ErrFetching
	}
	go w.messageRetriever.FetchPointers()
	return nil
}

func (w *Wallet) Online() <-chan struct{} {
	return w.online
}

func (w *Wallet) Done() <-chan struct{} {
	return w.done
}

func (w *Wallet) Updates() <-chan Update {
	return w.updates
}

func (w *Wallet) ThreadUpdates() <-chan thread.Update {
	return w.threadUpdates
}

func (w *Wallet) Notifications() <-chan repo.Notification {
	return w.notifications
}

func (w *Wallet) GetRepoPath() string {
	return w.repoPath
}

// GetDataAtPath returns raw data behind an ipfs path
func (w *Wallet) GetDataAtPath(path string) ([]byte, error) {
	if !w.started {
		return nil, ErrStopped
	}
	return util.GetDataAtPath(w.ipfs, path)
}

// createIPFS creates an IPFS node
func (w *Wallet) createIPFS(online bool) error {
	// open repo
	rep, err := fsrepo.Open(w.repoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}

	// determine routing
	routing := core.DHTOption
	if w.isMobile {
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
	ctx.ConfigRoot = w.repoPath
	ctx.LoadConfig = func(path string) (*config.Config, error) {
		return fsrepo.ConfigAt(w.repoPath)
	}
	ctx.ConstructNode = func() (*core.IpfsNode, error) {
		return nd, nil
	}

	// attach to textile node
	if w.cancel != nil {
		w.cancel()
	}
	if w.ipfs != nil {
		if err := w.ipfs.Close(); err != nil {
			log.Errorf("error closing prev ipfs node: %s", err)
			return err
		}
	}
	w.context = ctx
	w.cancel = cancel
	w.ipfs = nd

	return nil
}

func (w *Wallet) getThreadByBlock(block *repo.Block) (*thread.Thread, error) {
	if block == nil {
		return nil, errors.New("block is empty")
	}
	var thrd *thread.Thread
	for _, t := range w.threads {
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

func (w *Wallet) loadThread(mod *repo.Thread) (*thread.Thread, error) {
	if _, loaded := w.GetThread(mod.Id); loaded != nil {
		return nil, ErrThreadLoaded
	}
	id := mod.Id // save value locally
	threadConfig := &thread.Config{
		RepoPath: w.repoPath,
		Ipfs: func() *core.IpfsNode {
			return w.ipfs
		},
		Blocks:        w.datastore.Blocks,
		Peers:         w.datastore.Peers,
		Notifications: w.datastore.Notifications,
		GetHead: func() (string, error) {
			m := w.datastore.Threads().Get(id)
			if m == nil {
				return "", errors.New(fmt.Sprintf("could not re-load thread: %s", id))
			}
			return m.Head, nil
		},
		UpdateHead: func(head string) error {
			if err := w.datastore.Threads().UpdateHead(id, head); err != nil {
				return err
			}
			return nil
		},
		Send:          w.SendMessage,
		NewEnvelope:   w.NewEnvelope,
		PutPinRequest: w.putPinRequest,
		GetUsername:   w.GetUsername,
		SendUpdate:    w.sendThreadUpdate,
	}
	thrd, err := thread.NewThread(mod, threadConfig)
	if err != nil {
		return nil, err
	}
	w.threads = append(w.threads, thrd)
	return thrd, nil
}

// putPinRequest adds a pin request to the pinner
func (w *Wallet) putPinRequest(id string) error {
	if w.pinner == nil {
		return nil
	}
	return w.pinner.Put(id)
}

// sendUpdate adds an update to the update channel
func (w *Wallet) sendUpdate(update Update) {
	defer func() {
		if recover() != nil {
			log.Error("update channel already closed")
		}
	}()
	w.updates <- update
}

// sendThreadUpdate adds a thread update to the update channel
func (w *Wallet) sendThreadUpdate(update thread.Update) {
	defer func() {
		if recover() != nil {
			log.Error("thread update channel already closed")
		}
	}()
	w.threadUpdates <- update
}

// sendNotification adds a notification to the notification channel
func (w *Wallet) sendNotification(notification *repo.Notification) error {
	// add to db
	if err := w.datastore.Notifications().Add(notification); err != nil {
		return err
	}

	// broadcast
	defer func() {
		if recover() != nil {
			log.Error("notification channel already closed")
		}
	}()
	w.notifications <- *notification

	return nil
}

// touchDB ensures that we have a good db connection
func (w *Wallet) touchDatastore() error {
	if err := w.datastore.Ping(); err != nil {
		log.Debug("re-opening datastore...")
		sqliteDB, err := db.Create(w.repoPath, "")
		if err != nil {
			log.Errorf("error re-opening datastore: %s", err)
			return err
		}
		w.datastore = sqliteDB
	}
	return nil
}

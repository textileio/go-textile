package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/cafe"
	"github.com/textileio/textile-go/core/cafe"
	"github.com/textileio/textile-go/net"
	serv "github.com/textileio/textile-go/net/service"
	trepo "github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/storage"
	"github.com/textileio/textile-go/wallet/thread"
	"github.com/textileio/textile-go/wallet/util"
	"gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
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

type Wallet struct {
	version            string
	context            oldcmds.Context
	repoPath           string
	cancel             context.CancelFunc
	ipfs               *core.IpfsNode
	datastore          trepo.Datastore
	service            *serv.TextileService
	cafeAddr           string
	isMobile           bool
	started            bool
	threads            []*thread.Thread
	done               chan struct{}
	updates            chan Update
	messageStorage     storage.OfflineMessagingStorage
	messageRetriever   *net.MessageRetriever
	pointerRepublisher *net.PointerRepublisher
	pinner             *net.Pinner
}

const pingTimeout = time.Second * 10

var ErrStarted = errors.New("node is already started")
var ErrStopped = errors.New("node is already stopped")
var ErrOffline = errors.New("node is offline")
var ErrThreadLoaded = errors.New("thread is already loaded")
var ErrNoCafeHost = errors.New("cafe host address is not set")

func NewWallet(config Config) (*Wallet, string, error) {
	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, "")
	if err != nil {
		return nil, "", err
	}

	// we may be running in an uninitialized state.
	mnemonic, err := trepo.DoInit(config.RepoPath, config.Version, config.Mnemonic, sqliteDB.Config().Init, sqliteDB.Config().Configure)
	if err != nil && err != trepo.ErrRepoExists {
		return nil, "", err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(config.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return nil, "", err
	}

	// if a specific swarm port was selected, set it in the config
	if err := applySwarmPortConfigOption(repo, config.SwarmPorts); err != nil {
		return nil, "", err
	}

	// ensure bootstrap addresses are latest in config (without wiping repo)
	if err := ensureBootstrapConfig(repo); err != nil {
		return nil, "", err
	}

	// if this is a server node, apply the ipfs server profile
	if err := applyServerConfigOption(repo, config.IsServer); err != nil {
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
func (w *Wallet) Start() (chan struct{}, error) {
	if w.started {
		return nil, ErrStarted
	}
	defer func() {
		w.done = make(chan struct{})
		w.started = true
	}()
	log.Info("starting wallet...")
	onlineCh := make(chan struct{})
	w.updates = make(chan Update)

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		log.Errorf("setting file descriptor limit: %s", err)
	}

	// check db
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}

	// start the ipfs node
	log.Debug("creating an ipfs node...")
	if err := w.createIPFS(false); err != nil {
		log.Errorf("error creating offline ipfs node: %s", err)
		return nil, err
	}
	go func() {
		defer close(onlineCh)
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
			// get token
			tokens, _ := w.datastore.Profile().GetTokens()
			if tokens == nil {
				return nil
			}
			return net.Pin(w.ipfs, id.Hash().B58String(), tokens, w.pinner.Url())
		})

		// service is now configurable
		w.service = serv.NewService(w.ipfs, w.datastore, w.GetThread, w.AddThread)

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
		}

		// re-pub profile
		go func() {
			if _, err := w.PublishProfile(); err != nil {
				log.Errorf("error publishing profile: %s", err)
			}
		}()

		// print swarm addresses
		if err := util.PrintSwarmAddrs(w.ipfs); err != nil {
			log.Errorf("failed to read listening addresses: %s", err)
		}
		log.Info("wallet is online")
	}()

	// build a pin requester
	if w.GetCafeAddr() != "" {
		pinnerCfg := net.PinnerConfig{
			Datastore: w.datastore,
			Ipfs: func() *core.IpfsNode {
				return w.ipfs
			},
			Api: fmt.Sprintf("%s/pin", w.GetCafeAddr()),
		}
		w.pinner = net.NewPinner(pinnerCfg)

		// start ticker job if not mobile
		if !w.isMobile {
			go w.pinner.Run()
		}
	}

	// setup threads
	for _, mod := range w.datastore.Threads().List("") {
		_, err := w.loadThread(&mod)
		if err == ErrThreadLoaded {
			continue
		}
		if err != nil {
			return nil, err
		}
	}

	log.Info("wallet is started")

	return onlineCh, nil
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
	if err := os.Remove(dsLockFile); err != nil {
		log.Warningf("remove ds lock failed: %s", err)
	}

	// wipe threads
	for _, t := range w.Threads() {
		t.Close()
	}
	w.threads = nil

	// wipe services
	w.messageStorage = nil
	w.service = nil
	w.messageRetriever = nil
	w.pointerRepublisher = nil
	w.pinner = nil

	// close updates
	close(w.updates)

	log.Info("wallet is stopped")

	return nil
}

func (w *Wallet) Started() bool {
	return w.started
}

func (w *Wallet) Online() bool {
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

func (w *Wallet) RefreshMessages() error {
	if !w.Online() {
		return ErrOffline
	}
	w.messageRetriever.Add(1)
	go w.messageRetriever.FetchPointers()
	go w.pointerRepublisher.Republish()
	return nil
}

func (w *Wallet) RunPinner() {
	if w.pinner == nil {
		return
	}
	go w.pinner.Pin()
}

func (w *Wallet) Updates() <-chan Update {
	return w.updates
}

func (w *Wallet) Done() <-chan struct{} {
	return w.done
}

func (w *Wallet) GetRepoPath() string {
	return w.repoPath
}

// GetCafeAddr returns the cafe address is set
func (w *Wallet) GetCafeAddr() string {
	if w.cafeAddr == "" {
		return ""
	}
	return fmt.Sprintf("%s/api/%s", w.cafeAddr, cafe.Version)
}

// GetId returns peer id
func (w *Wallet) GetId() (string, error) {
	if !w.started {
		return "", ErrStopped
	}
	return w.ipfs.Identity.Pretty(), nil
}

// GetPrivKey returns the current user's master secret key
func (w *Wallet) GetPrivKey() (libp2pc.PrivKey, error) {
	if !w.started {
		return nil, ErrStopped
	}
	if w.ipfs.PrivateKey == nil {
		if err := w.ipfs.LoadPrivateKey(); err != nil {
			return nil, err
		}
	}
	return w.ipfs.PrivateKey, nil
}

// GetPubKey returns the current user's master public key
func (w *Wallet) GetPubKey() (libp2pc.PubKey, error) {
	secret, err := w.GetPrivKey()
	if err != nil {
		return nil, err
	}
	return secret.GetPublic(), nil
}

// GetPubKeyString returns the base64 encoded public ipfs peer key
func (w *Wallet) GetPubKeyString() (string, error) {
	pk, err := w.GetPubKey()
	if err != nil {
		return "", err
	}
	pkb, err := pk.Bytes()
	if err != nil {
		return "", err
	}
	return libp2pc.ConfigEncodeKey(pkb), nil
}

func (w *Wallet) Threads() []*thread.Thread {
	return w.threads
}

func (w *Wallet) GetThread(id string) (*int, *thread.Thread) {
	for i, thrd := range w.threads {
		if thrd.Id == id {
			return &i, thrd
		}
	}
	return nil, nil
}

// GetBlock searches for a local block associated with the given target
func (w *Wallet) GetBlock(id string) (*trepo.Block, error) {
	block := w.datastore.Blocks().Get(id)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
}

// GetBlockByDataId searches for a local block associated with the given data id
func (w *Wallet) GetBlockByDataId(dataId string) (*trepo.Block, error) {
	block := w.datastore.Blocks().GetByDataId(dataId)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
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
	repo, err := fsrepo.Open(w.repoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}

	// assemble node config
	cfg := &core.BuildCfg{
		Repo:      repo,
		Permanent: true, // temporary way to signify that node is permanent
		Online:    online,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
			"mplex":  true,
		},
		Routing: core.DHTOption,
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

func (w *Wallet) getThreadByBlock(block *trepo.Block) (*thread.Thread, error) {
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

func (w *Wallet) loadThread(mod *trepo.Thread) (*thread.Thread, error) {
	_, loaded := w.GetThread(mod.Id)
	if loaded != nil {
		return nil, ErrThreadLoaded
	}
	id := mod.Id // save value locally
	threadConfig := &thread.Config{
		RepoPath: w.repoPath,
		Ipfs: func() *core.IpfsNode {
			return w.ipfs
		},
		Blocks: w.datastore.Blocks,
		Peers:  w.datastore.Peers,
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
		Send:        w.SendMessage,
		NewEnvelope: w.NewEnvelope,
		PutPinRequest: func(id string) error {
			if w.pinner == nil {
				return nil
			}
			return w.pinner.Put(id)
		},
	}
	thrd, err := thread.NewThread(mod, threadConfig)
	if err != nil {
		return nil, err
	}
	w.threads = append(w.threads, thrd)
	return thrd, nil
}

func (w *Wallet) sendUpdate(update Update) {
	defer func() {
		if recover() != nil {
			log.Error("update channel already closed")
		}
	}()
	select {
	case w.updates <- update:
	default:
	}
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

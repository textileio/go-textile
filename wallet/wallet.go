package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	cmodels "github.com/textileio/textile-go/central/models"
	"github.com/textileio/textile-go/core/central"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net"
	serv "github.com/textileio/textile-go/net/service"
	"github.com/textileio/textile-go/pb"
	trepo "github.com/textileio/textile-go/repo"
	tconfig "github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/thread"
	"github.com/textileio/textile-go/wallet/util"
	"gx/ipfs/QmSFihvoND3eDaAYRCeLgLPt62yCPgMZs1NSZmKFEtJQQw/go-libp2p-floodsub"
	"gx/ipfs/QmSwZMWwFZSUpe5muU2xgTUwppH24KfMwdPXiwbEp2c6G5/go-libp2p-swarm"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	libp2pn "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	utilmain "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/commands"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/config"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/fsrepo"
	uio "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/unixfs/io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var log = logging.MustGetLogger("wallet")

type Config struct {
	Version        string
	RepoPath       string
	CentralAPI     string
	IsMobile       bool
	IsServer       bool
	SwarmPort      string
	MasterMnemonic *string
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

type Wallet struct {
	context        oldcmds.Context
	repoPath       string
	gatewayAddr    string
	cancel         context.CancelFunc
	ipfs           *core.IpfsNode
	datastore      trepo.Datastore
	service        *serv.TextileService
	centralAPI     string
	isMobile       bool
	started        bool
	threads        []*thread.Thread
	done           chan struct{}
	updates        chan Update
	lastRelayTouch time.Time
}

const (
	pingTimeout        = time.Second * 10
	relayTouchInterval = time.Minute * 2
)

var ErrStarted = errors.New("node is already started")
var ErrStopped = errors.New("node is already stopped")
var ErrOffline = errors.New("node is offline")
var ErrThreadExists = errors.New("thread already exists")
var ErrThreadLoaded = errors.New("thread is already loaded")

func NewWallet(config Config) (*Wallet, string, error) {
	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, "")
	if err != nil {
		return nil, "", err
	}

	// we may be running in an uninitialized state.
	mnemonic, err := trepo.DoInit(config.RepoPath, config.IsMobile, config.Version, config.MasterMnemonic,
		sqliteDB.Config().Init, sqliteDB.Config().Configure)
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

	// save gateway address
	gwAddr, err := repo.GetConfigKey("Addresses.Gateway")
	if err != nil {
		log.Errorf("error getting gateway address: %s", err)
		return nil, "", err
	}

	// if a specific swarm port was selected, set it in the config
	if config.SwarmPort != "" {
		log.Infof("using specified swarm port: %s", config.SwarmPort)
		if err := tconfig.Update(repo, "Addresses.Swarm", []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", config.SwarmPort),
			fmt.Sprintf("/ip6/::/tcp/%s", config.SwarmPort),
		}); err != nil {
			return nil, "", err
		}
	}

	// if this is a server node, apply the ipfs server profile
	if config.IsServer {
		if err := tconfig.Update(repo, "Addresses.NoAnnounce", tconfig.DefaultServerFilters); err != nil {
			return nil, "", err
		}
		if err := tconfig.Update(repo, "Swarm.AddrFilters", tconfig.DefaultServerFilters); err != nil {
			return nil, "", err
		}
		if err := tconfig.Update(repo, "Swarm.EnableRelayHop", true); err != nil {
			return nil, "", err
		}
		if err := tconfig.Update(repo, "Discovery.MDNS.Enabled", false); err != nil {
			return nil, "", err
		}
		log.Info("applied server profile")
	}

	// clean central api url
	if len(config.CentralAPI) > 0 {
		ca := config.CentralAPI
		if ca[len(ca)-1:] == "/" {
			ca = ca[0 : len(ca)-1]
		}
		config.CentralAPI = ca
	}

	return &Wallet{
		repoPath:    config.RepoPath,
		gatewayAddr: gwAddr.(string),
		datastore:   sqliteDB,
		centralAPI:  config.CentralAPI,
		isMobile:    config.IsMobile,
		updates:     make(chan Update),
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
		w.lastRelayTouch = time.Time{}
	}()
	log.Info("starting wallet...")
	onlineCh := make(chan struct{})

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

		// service is now configurable
		w.service = serv.NewService(w.ipfs, w.datastore, w.GetThread, w.AddThread)

		// print swarm addresses
		if err := util.PrintSwarmAddrs(w.ipfs); err != nil {
			log.Errorf("failed to read listening addresses: %s", err)
		}
		log.Info("wallet is online")
	}()

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

	// wipe service
	w.service = nil

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

// Updates returns a read-only channel of updates
func (w *Wallet) Updates() <-chan Update {
	return w.updates
}

func (w *Wallet) Done() <-chan struct{} {
	return w.done
}

func (w *Wallet) GetGatewayAddress() string {
	return w.gatewayAddr
}

func (w *Wallet) GetRepoPath() string {
	return w.repoPath
}

// SignUp requests a new username and token from the central api and saves them locally
func (w *Wallet) SignUp(reg *cmodels.Registration) error {
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debugf("signup: %s %s %s %s %s", reg.Username, "xxxxxx", reg.Identity.Type, reg.Identity.Value, reg.Referral)

	// remote signup
	res, err := central.SignUp(reg, w.GetCentralUserAPI())
	if err != nil {
		log.Errorf("signup error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signup error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.datastore.Profile().SignIn(
		reg.Username,
		res.Session.AccessToken, res.Session.RefreshToken,
	); err != nil {
		log.Errorf("local signin error: %s", err)
		return err
	}
	return nil
}

// SignIn requests a token with a username from the central api and saves them locally
func (w *Wallet) SignIn(creds *cmodels.Credentials) error {
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debugf("signin: %s %s", creds.Username, "xxxxxx")

	// remote signin
	res, err := central.SignIn(creds, w.GetCentralUserAPI())
	if err != nil {
		log.Errorf("signin error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signin error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.datastore.Profile().SignIn(
		creds.Username,
		res.Session.AccessToken, res.Session.RefreshToken,
	); err != nil {
		log.Errorf("local signin error: %s", err)
		return err
	}
	return nil
}

// SignOut deletes the locally saved user info (username and tokens)
func (w *Wallet) SignOut() error {
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debug("signing out...")

	// remote is stateless, so we just ditch the local token
	if err := w.datastore.Profile().SignOut(); err != nil {
		log.Errorf("local signout error: %s", err)
		return err
	}
	return nil
}

// IsSignedIn returns whether or not a user is signed in
func (w *Wallet) IsSignedIn() (bool, error) {
	if err := w.touchDatastore(); err != nil {
		return false, err
	}
	_, err := w.datastore.Profile().GetUsername()
	return err == nil, nil
}

// GetUsername returns the current user's username
func (w *Wallet) GetUsername() (string, error) {
	if err := w.touchDatastore(); err != nil {
		return "", err
	}
	return w.datastore.Profile().GetUsername()
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

// GetAccessToken returns the current access_token (jwt) for central
func (w *Wallet) GetAccessToken() (string, error) {
	if err := w.touchDatastore(); err != nil {
		return "", err
	}
	at, _, err := w.datastore.Profile().GetTokens()
	if err != nil {
		return "", err
	}
	return at, nil
}

func (w *Wallet) GetCentralAPI() string {
	return w.centralAPI
}

func (w *Wallet) GetCentralUserAPI() string {
	return fmt.Sprintf("%s/api/v1/users", w.centralAPI)
}

func (w *Wallet) Threads() []*thread.Thread {
	return w.threads
}

func (w *Wallet) GetThread(id string) *thread.Thread {
	for _, thrd := range w.threads {
		if thrd.Id == id {
			return thrd
		}
	}
	return nil
}

func (w *Wallet) GetThreadByName(name string) (*int, *thread.Thread) {
	for i, thrd := range w.threads {
		if thrd.Name == name {
			return &i, thrd
		}
	}
	return nil, nil
}

// AddThread adds a thread with a given name and secret key
func (w *Wallet) AddThread(name string, secret libp2pc.PrivKey) (*thread.Thread, error) {
	existing, err := w.getThreadModelByName(name)
	if err != nil {
		return nil, err
	}
	// not ideal way to check for existence, but want to skip
	// all the heavy crypto stuff below if we know for sure this thread already exists
	if existing != nil {
		return nil, ErrThreadExists
	}
	log.Debugf("adding a new thread: %s", name)

	// index a new thread
	skb, err := secret.Bytes()
	if err != nil {
		return nil, err
	}
	pkb, err := secret.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	pk := libp2pc.ConfigEncodeKey(pkb)
	threadModel := &trepo.Thread{
		Id:      pk,
		Name:    name,
		PrivKey: skb,
	}
	if err := w.datastore.Threads().Add(threadModel); err != nil {
		return nil, err
	}

	// load as active thread
	thrd, err := w.loadThread(threadModel)
	if err != nil {
		return nil, err
	}

	// invite each device to the new thread
	for _, device := range w.Devices() {
		dpkb, err := libp2pc.ConfigDecodeKey(device.Id)
		if err != nil {
			return nil, err
		}
		dpk, err := libp2pc.UnmarshalPublicKey(dpkb)
		if err != nil {
			return nil, err
		}
		if _, err := thrd.AddInvite(dpk); err != nil {
			return nil, err
		}
	}

	// notify listeners
	w.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadAdded})

	return thrd, nil
}

// AddThreadWithMnemonic adds a thread with a given name and mnemonic phrase
func (w *Wallet) AddThreadWithMnemonic(name string, mnemonic *string) (*thread.Thread, string, error) {
	existing, err := w.getThreadModelByName(name)
	if err != nil {
		return nil, "", err
	}
	// not ideal way to check for existence, but want to skip
	// all the heavy crypto stuff below if we know for sure this thread already exists
	if existing != nil {
		return nil, "", ErrThreadExists
	}
	if mnemonic != nil {
		log.Debugf("regenerating keypair from mnemonic for: %s", name)
	} else {
		log.Debugf("generating keypair for: %s", name)
	}
	secret, mnem, err := util.PrivKeyFromMnemonic(mnemonic)
	if err != nil {
		return nil, "", err
	}
	thrd, err := w.AddThread(name, secret)
	if err != nil {
		return nil, "", err
	}
	return thrd, mnem, nil
}

// RemoveThread removes a thread by name
func (w *Wallet) RemoveThread(name string) error {
	i, thrd := w.GetThreadByName(name) // gets the loaded thread
	if thrd == nil {
		return errors.New("thread not found")
	}

	// clean up
	thrd.Close()
	copy(w.threads[*i:], w.threads[*i+1:])
	w.threads[len(w.threads)-1] = nil
	w.threads = w.threads[:len(w.threads)-1]

	// remove model from db
	if err := w.datastore.Threads().DeleteByName(name); err != nil {
		return err
	}
	log.Infof("removed thread '%s'", name)

	// TODO: create a "left" block?

	// notify listeners
	w.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadRemoved})

	return nil
}

// PublishThreads publishes HEAD for each thread
func (w *Wallet) PublishThreads() {
	for _, t := range w.threads {
		go func(thrd *thread.Thread) {
			thrd.PostHead()
		}(t)
	}
}

// Devices lists all devices
func (w *Wallet) Devices() []trepo.Device {
	return w.datastore.Devices().List("")
}

// AddDevice creates an invite for every current and future thread
func (w *Wallet) AddDevice(name string, pk libp2pc.PubKey) error {
	if !w.Online() {
		return ErrOffline
	}

	// index a new device
	pkb, err := pk.Bytes()
	if err != nil {
		return err
	}
	deviceModel := &trepo.Device{
		Id:   libp2pc.ConfigEncodeKey(pkb),
		Name: name,
	}
	if err := w.datastore.Devices().Add(deviceModel); err != nil {
		return err
	}
	log.Infof("added device '%s'", name)

	// invite device to existing threads
	for _, thrd := range w.threads {
		if _, err := thrd.AddInvite(pk); err != nil {
			return err
		}
	}

	// notify listeners
	w.sendUpdate(Update{Id: deviceModel.Id, Name: deviceModel.Name, Type: DeviceAdded})

	return nil
}

// RemoveDevice removes a device by name
func (w *Wallet) RemoveDevice(name string) error {
	device := w.datastore.Devices().GetByName(name)
	if device == nil {
		return errors.New("device not found")
	}
	if err := w.datastore.Devices().DeleteByName(name); err != nil {
		return err
	}
	log.Infof("removed device '%s'", name)

	// TODO: uninvite?

	// notify listeners
	w.sendUpdate(Update{Id: device.Id, Name: device.Name, Type: DeviceRemoved})

	return nil
}

// AddPhoto add a photo to the local ipfs node
func (w *Wallet) AddPhoto(path string) (*model.AddResult, error) {
	// get a key to encrypt with
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, err
	}

	// read file from disk
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode image
	reader, format, err := util.DecodeImage(file)
	if err != nil {
		return nil, err
	}

	// make a thumbnail
	reader.Seek(0, 0)
	var thumbFormat util.ThumbnailFormat
	if format == "gif" {
		thumbFormat = util.GIF
	} else {
		thumbFormat = util.JPEG
	}
	thumb, err := util.MakeThumbnail(reader, thumbFormat, model.ThumbnailWidth)
	if err != nil {
		return nil, err
	}

	// get some meta data
	username, _ := w.datastore.Profile().GetUsername() // ignore if not present (not signed in)
	mpk, err := w.GetPubKey()
	if err != nil {
		return nil, err
	}
	mpkb, err := mpk.Bytes()
	if err != nil {
		return nil, err
	}

	// path info
	fpath := file.Name()
	ext := strings.ToLower(filepath.Ext(fpath))

	// get metadata
	reader.Seek(0, 0)
	meta, err := util.GetMetadata(reader, fpath, ext, username)
	if err != nil {
		return nil, err
	}
	metab, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	// encrypt files
	reader.Seek(0, 0)
	photocypher, err := util.GetEncryptedReaderBytes(reader, key)
	if err != nil {
		return nil, err
	}
	thumbcypher, err := crypto.EncryptAES(thumb, key)
	if err != nil {
		return nil, err
	}
	metacypher, err := crypto.EncryptAES(metab, key)
	if err != nil {
		return nil, err
	}
	mpkcypher, err := crypto.EncryptAES(mpkb, key)
	if err != nil {
		return nil, err
	}

	// create a virtual directory for the photo
	dirb := uio.NewDirectory(w.ipfs.DAG)
	err = util.AddFileToDirectory(w.ipfs, dirb, photocypher, "photo")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(w.ipfs, dirb, thumbcypher, "thumb")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(w.ipfs, dirb, metacypher, "meta")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(w.ipfs, dirb, mpkcypher, "pk")
	if err != nil {
		return nil, err
	}

	// pin the directory
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := util.PinDirectory(w.ipfs, dir, []string{"photo"}); err != nil {
		return nil, err
	}
	id := dir.Cid().Hash().B58String()

	// create and init a new multipart request
	request := &net.MultipartRequest{}
	request.Init(filepath.Join(w.repoPath, "tmp"), id)

	// add files to request
	if err := request.AddFile(photocypher, "photo"); err != nil {
		return nil, err
	}
	if err := request.AddFile(thumbcypher, "thumb"); err != nil {
		return nil, err
	}
	if err := request.AddFile(metacypher, "meta"); err != nil {
		return nil, err
	}
	if err := request.AddFile(mpkcypher, "pk"); err != nil {
		return nil, err
	}

	// finish request
	if err := request.Finish(); err != nil {
		return nil, err
	}

	// all done
	return &model.AddResult{Id: id, Key: key, RemoteRequest: request}, nil
}

// GetBlock searches for a local block associated with the given target
func (w *Wallet) GetBlock(id string) (*trepo.Block, error) {
	block := w.datastore.Blocks().Get(id)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
}

// GetBlockByTarget searches for a local block associated with the given target
func (w *Wallet) GetBlockByTarget(target string) (*trepo.Block, error) {
	block := w.datastore.Blocks().GetByTarget(target)
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

// ConnectPeer connect to another ipfs peer (i.e., ipfs swarm connect)
func (w *Wallet) ConnectPeer(addrs []string) ([]string, error) {
	if !w.Online() {
		return nil, ErrOffline
	}
	snet, ok := w.ipfs.PeerHost.Network().(*swarm.Network)
	if !ok {
		return nil, errors.New("peerhost network was not swarm")
	}

	swrm := snet.Swarm()
	pis, err := util.PeersWithAddresses(addrs)
	if err != nil {
		return nil, err
	}

	output := make([]string, len(pis))
	for i, pi := range pis {
		swrm.Backoff().Clear(pi.ID)

		output[i] = "connect " + pi.ID.Pretty()

		err := w.ipfs.PeerHost.Connect(w.ipfs.Context(), pi)
		if err != nil {
			return nil, fmt.Errorf("%s failure: %s", output[i], err)
		}
		output[i] += " success"
	}
	return output, nil
}

// PingPeer pings a peer num times, returning the result to out chan
func (w *Wallet) PingPeer(addrs string, num int, out chan string) error {
	if !w.started {
		return ErrStopped
	}
	if !w.Online() {
		return ErrOffline
	}
	addr, pid, err := util.ParsePeerParam(addrs)
	if addr != nil {
		w.ipfs.Peerstore.AddAddr(pid, addr, pstore.TempAddrTTL) // temporary
	}

	if len(w.ipfs.Peerstore.Addrs(pid)) == 0 {
		// Make sure we can find the node in question
		log.Debugf("looking up peer: %s", pid.Pretty())

		ctx, cancel := context.WithTimeout(w.ipfs.Context(), pingTimeout)
		defer cancel()
		p, err := w.ipfs.Routing.FindPeer(ctx, pid)
		if err != nil {
			err = fmt.Errorf("peer lookup error: %s", err)
			log.Errorf(err.Error())
			return err
		}
		w.ipfs.Peerstore.AddAddrs(p.ID, p.Addrs, pstore.TempAddrTTL)
	}

	ctx, cancel := context.WithTimeout(w.ipfs.Context(), pingTimeout*time.Duration(num))
	defer cancel()
	pings, err := w.ipfs.Ping.Ping(ctx, pid)
	if err != nil {
		log.Errorf("error pinging peer %s: %s", pid.Pretty(), err)
		return err
	}

	var done bool
	var total time.Duration
	for i := 0; i < num && !done; i++ {
		select {
		case <-ctx.Done():
			done = true
			close(out)
			break
		case t, ok := <-pings:
			if !ok {
				done = true
				close(out)
				break
			}
			total += t
			msg := fmt.Sprintf("ping %s completed after %f seconds", pid.Pretty(), t.Seconds())
			select {
			case out <- msg:
			default:
			}
			log.Debug(msg)
			time.Sleep(time.Second)
		}
	}
	return nil
}

func (w *Wallet) Peers() ([]libp2pn.Conn, error) {
	if !w.Online() {
		return nil, ErrOffline
	}
	return w.ipfs.PeerHost.Network().Conns(), nil
}

func (w *Wallet) Publish(topic string, payload []byte) error {
	if !w.Online() {
		return ErrOffline
	}
	if w.lastRelayTouch.Add(relayTouchInterval).Before(time.Now()) {
		w.lastRelayTouch = time.Now()
		log.Debug("connecting to relay...")
		out, err := w.ConnectPeer([]string{fmt.Sprintf("/p2p-circuit/ipfs/%s", tconfig.RemoteRelayNode)})
		if err != nil {
			return err
		}
		for _, o := range out {
			log.Debug(o)
		}
	}
	return w.ipfs.Floodsub.Publish(topic, payload)
}

func (w *Wallet) Subscribe(topic string) (*floodsub.Subscription, error) {
	if !w.Online() {
		return nil, ErrOffline
	}
	return w.ipfs.Floodsub.Subscribe(topic)
}

// createIPFS creates an IPFS node
func (w *Wallet) createIPFS(online bool) error {
	// open repo
	repo, err := fsrepo.Open(w.repoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}

	// determine the best routing
	var routingOption core.RoutingOption
	if w.isMobile {
		routingOption = core.DHTClientOption
	} else {
		routingOption = core.DHTOption
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
		Routing: routingOption,
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

func (w *Wallet) getThreadModelByName(name string) (*trepo.Thread, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	return w.datastore.Threads().GetByName(name), nil
}

func (w *Wallet) getThreadByBlock(block *trepo.Block) (*thread.Thread, error) {
	if block == nil {
		return nil, errors.New("block is empty")
	}
	var thrd *thread.Thread
	for _, t := range w.threads {
		if t.Id == block.ThreadPubKey {
			thrd = t
			break
		}
	}
	if thrd == nil {
		return nil, errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadPubKey))
	}
	return thrd, nil
}

func (w *Wallet) loadThread(model *trepo.Thread) (*thread.Thread, error) {
	_, loaded := w.GetThreadByName(model.Name)
	if loaded != nil {
		return nil, ErrThreadLoaded
	}
	id := model.Id // save value locally
	threadConfig := &thread.Config{
		RepoPath: w.repoPath,
		Ipfs: func() *core.IpfsNode {
			return w.ipfs
		},
		Blocks: func() trepo.BlockStore {
			return w.datastore.Blocks()
		},
		Peers: func() trepo.PeerStore {
			return w.datastore.Peers()
		},
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
		Publish: func(payload []byte) error {
			if err := w.Publish(id, payload); err != nil {
				return err
			}
			return nil
		},
		Send: func(message *pb.Message, peerId string) error {
			return w.sendMessage(message, peerId)
		},
	}
	thrd, err := thread.NewThread(model, threadConfig)
	if err != nil {
		return nil, err
	}
	w.threads = append(w.threads, thrd)
	return thrd, nil
}

func (w *Wallet) sendUpdate(update Update) {
	select {
	case w.updates <- update:
	default:
	}
	defer func() {
		if recover() != nil {
			log.Error("update channel already closed")
		}
	}()
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

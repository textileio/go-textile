package wallet

import (
	"time"
	"github.com/op/go-logging"
	"github.com/tyler-smith/go-bip39"
	trepo "github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/crypto"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	uio "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/unixfs/io"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/net"
	"os"
	"strings"
	"path/filepath"
	"encoding/json"
	"github.com/textileio/textile-go/core/central"
	cmodels "github.com/textileio/textile-go/central/models"
	"errors"
	"fmt"
	"encoding/base64"
	oldcmds "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/commands"
	"context"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/fsrepo"
	utilmain "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/cmd/ipfs/util"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/repo/config"
	"gx/ipfs/QmSwZMWwFZSUpe5muU2xgTUwppH24KfMwdPXiwbEp2c6G5/go-libp2p-swarm"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
)

var log = logging.MustGetLogger("wallet")

type Wallet struct {
	Context        oldcmds.Context
	RepoPath       string
	Cancel         context.CancelFunc
	Ipfs           *core.IpfsNode
	Datastore      trepo.Datastore
	CentralUserAPI string
	IsMobile       bool
	started        bool
	threads        []*Thread
}

type Metadata struct {
	Username string    `json:"un,omitempty"`
	Created  time.Time `json:"cts,omitempty"`
	Added    time.Time `json:"ats,omitempty"`
}

type FileMetadata struct {
	Metadata
	Name string `json:"name,omitempty"`
	Ext  string `json:"ext,omitempty"`
}

type AddResult struct {
	Id            string
	Key           []byte
	RemoteRequest *net.MultipartRequest
}

const pingTimeout = time.Second * 10
const pinTimeout = time.Minute * 1
const catTimeout = time.Second * 30

// ErrRunning is an error for when node start is called on a started node
var ErrStarted = errors.New("node is already started")

// ErrStopped is an error for  when node stop is called on a stopped node
var ErrStopped = errors.New("node is already stopped")

// ErrOffline is an error for when online resources are requested on an offline node
var ErrOffline = errors.New("node is offline")

// Start
func (w *Wallet) Start() (chan struct{}, error) {
	if w.started {
		return nil, ErrStarted
	}
	defer func() { w.started = true }()
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

		// print swarm addresses
		if err := printSwarmAddrs(w.Ipfs); err != nil {
			log.Errorf("failed to read listening addresses: %s", err)
		}
		log.Info("wallet is online")
	}()

	// setup threads
	for _, thread := range w.Threads() {
		w.setupThread(thread)
	}

	log.Info("wallet is started")

	return onlineCh, nil
}

// createIPFS creates an IPFS node
func (w *Wallet) createIPFS(online bool) error {
	// open repo
	repo, err := fsrepo.Open(w.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return err
	}

	// determine the best routing
	var routingOption core.RoutingOption
	if w.IsMobile {
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
	ctx.ConfigRoot = w.RepoPath
	ctx.LoadConfig = func(path string) (*config.Config, error) {
		return fsrepo.ConfigAt(w.RepoPath)
	}
	ctx.ConstructNode = func() (*core.IpfsNode, error) {
		return nd, nil
	}

	// attach to textile node
	if w.Cancel != nil {
		w.Cancel()
	}
	if w.Ipfs != nil {
		if err := w.Ipfs.Close(); err != nil {
			log.Errorf("error closing prev ipfs node: %s", err)
			return err
		}
	}
	w.Context = ctx
	w.Cancel = cancel
	w.Ipfs = nd

	return nil
}

// Stop the node
func (w *Wallet) Stop() error {
	if !w.started {
		return ErrStopped
	}
	defer func() { w.started = false }()
	log.Info("stopping wallet...")

	// close ipfs node
	w.Context.Close()
	w.Cancel()
	if err := w.Ipfs.Close(); err != nil {
		log.Errorf("error closing ipfs node: %s", err)
		return err
	}

	// close db connection
	w.Datastore.Close()
	dsLockFile := filepath.Join(w.RepoPath, "datastore", "LOCK")
	if err := os.Remove(dsLockFile); err != nil {
		log.Warningf("remove ds lock failed: %s", err)
	}

	// wipe threads
	w.threads = nil

	log.Info("wallet is stopped")

	return nil
}

func (w *Wallet) Started() bool {
	return w.started
}

func (w *Wallet) Online() bool {
	if w.Ipfs == nil {
		return false
	}
	return w.Ipfs.OnlineMode()
}

// SignUp requests a new username and token from the central api and saves them locally
func (w *Wallet) SignUp(reg *cmodels.Registration) error {
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debugf("signup: %s %s %s %s %s", reg.Username, "xxxxxx", reg.Identity.Type, reg.Identity.Value, reg.Referral)

	// remote signup
	res, err := central.SignUp(reg, w.CentralUserAPI)
	if err != nil {
		log.Errorf("signup error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signup error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.Datastore.Profile().SignIn(reg.Username, res.Session.AccessToken, res.Session.RefreshToken); err != nil {
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
	res, err := central.SignIn(creds, w.CentralUserAPI)
	if err != nil {
		log.Errorf("signin error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signin error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.Datastore.Profile().SignIn(creds.Username, res.Session.AccessToken, res.Session.RefreshToken); err != nil {
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
	if err := w.Datastore.Profile().SignOut(); err != nil {
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
	_, err := w.Datastore.Profile().GetUsername()
	return err == nil, nil
}

// GetUsername returns the current user's username
func (w *Wallet) GetUsername() (string, error) {
	if err := w.touchDatastore(); err != nil {
		return "", err
	}
	un, err := w.Datastore.Profile().GetUsername()
	if err != nil {
		return "", err
	}
	return un, nil
}

// GetAccessToken returns the current access_token (jwt) for central
func (w *Wallet) GetAccessToken() (string, error) {
	if err := w.touchDatastore(); err != nil {
		return "", err
	}
	at, _, err := w.Datastore.Profile().GetTokens()
	if err != nil {
		return "", err
	}
	return at, nil
}

// AddThread add a thread with a given name and mnemonic phrase
func (w *Wallet) AddThread(name string, mnemonic string) (*Thread, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	log.Debugf("adding a new thread: %s", name)

	// use phrase if provided
	if mnemonic == "" {
		var err error
		mnemonic, err = createMnemonic(bip39.NewEntropy, bip39.NewMnemonic)
		if err != nil {
			return nil, err
		}
		log.Debugf("generating Ed25519 keypair for: %s", name)
	} else {
		log.Debugf("regenerating Ed25519 keypair from mnemonic phrase for: %s", name)
	}

	// create the bip39 seed from the phrase
	seed := bip39.NewSeed(mnemonic, "")
	kb, err := identityKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}

	// convert to a libp2p crypto private key
	sk, err := libp2pc.UnmarshalPrivateKey(kb)
	if err != nil {
		return nil, err
	}

	// get public key as b64 string since we need it as a primary table key
	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	pk := libp2pc.ConfigEncodeKey(pkb)

	// finally, index a new thread
	thread := &Thread{
		Id:       pk,
		Name:     name,
		PrivKey:  sk,
	}
	if err := w.Datastore.Threads().Add(thread); err != nil {
		return nil, err
	}
	w.setupThread(thread)

	return thread, nil
}

func (w *Wallet) Threads() []*Thread {
	return w.threads
}

func (w *Wallet) GetThread(id string) *Thread {
	for _, thread := range w.threads {
		if thread.Id == id {
			return thread
		}
	}
	return nil
}

func (w *Wallet) GetThreadByName(name string) *Thread {
	for _, thread := range w.threads {
		if thread.Name == name {
			return thread
		}
	}
	return nil
}

// PublishThreads publishes HEAD for each thread
func (w *Wallet) PublishThreads() {
	for _, t := range w.threads {
		go func(thread *Thread) {
			thread.Publish()
		}(t)
	}
}

// TODO: add node master pk to dir
func (w *Wallet) AddPhoto(path string) (*AddResult, error) {
	// get a key to encrypt with
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, err
	}

	// read file from disk
	photo, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer photo.Close()

	// make a thumbnail
	thumb, err := makeThumbnail(photo, thumbnailWidth)
	if err != nil {
		return nil, err
	}

	// path info
	fpath := photo.Name()
	ext := strings.ToLower(filepath.Ext(fpath))

	// get username, ignoring if not present (not signed in)
	username, _ := w.Datastore.Profile().GetUsername()

	// get metadata
	meta, err := getMetadata(photo, fpath, ext, username)
	if err != nil {
		return nil, err
	}
	metab, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	// encrypt files
	photocypher, err := getEncryptedReaderBytes(photo, key)
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


	// create a virtual directory for the photo
	dirb := uio.NewDirectory(w.Ipfs.DAG)
	err = addFileToDirectory(w.Ipfs, dirb, photocypher, "photo")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(w.Ipfs, dirb, thumbcypher, "thumb")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(w.Ipfs, dirb, metacypher, "meta")
	if err != nil {
		return nil, err
	}

	// pin the directory
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := pinDirectory(w.Ipfs, dir, []string{"photo"}); err != nil {
		return nil, err
	}
	id := dir.Cid().Hash().B58String()

	// create and init a new multipart request
	request := &net.MultipartRequest{}
	request.Init(filepath.Join(w.RepoPath, "tmp"), id)

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

	// finish request
	if err := request.Finish(); err != nil {
		return nil, err
	}

	// all done
	return &AddResult{Id: id, Key: key, RemoteRequest: request}, nil
}

func (w *Wallet) FindBlock(target string) (*Block, error) {
	block := w.Datastore.Blocks().GetByTarget(target)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
}

// GetFile cats data from ipfs and tries to decrypt it with the provided block
// e.g., Qm../thumb, Qm../photo, Qm../meta, Qm../caption
func (w *Wallet) GetFile(path string, blockId string) ([]byte, error) {
	if !w.started {
		return nil, ErrStopped
	}

	// get thread for decryption
	thread, block, err := w.getThreadBlock(blockId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// get bytes
	cypher, err := GetDataAtPath(w.Ipfs, path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// decrypt the file key
	key, err := thread.Decrypt(block.TargetKey)
	if err != nil {
		log.Errorf("error decrypting key: %s", err)
		return nil, err
	}

	// finally, decrypt the file
	return crypto.DecryptAES(cypher, key)
}

// GetFileBase64 returns data encoded as base64 under an ipfs path
func (w *Wallet) GetFileBase64(path string, blockId string) (string, error) {
	file, err := w.GetFile(path, blockId)
	if err != nil {
		return "error", err
	}
	return base64.StdEncoding.EncodeToString(file), nil
}

func (w *Wallet) GetIPFSPeerID() (string, error) {
	if !w.started {
		return "", ErrStopped
	}
	return w.Ipfs.Identity.Pretty(), nil
}

// GetIPFSPubKeyString returns the base64 encoded public ipfs peer key
func (w *Wallet) GetIPFSPubKeyString() (string, error) {
	if !w.started {
		return "", ErrStopped
	}
	pkb, err := w.Ipfs.PrivateKey.GetPublic().Bytes()
	if err != nil {
		log.Errorf("error getting pub key bytes: %s", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pkb), nil
}

// ConnectPeer connect to another ipfs peer (i.e., ipfs swarm connect)
func (w *Wallet) ConnectPeer(addrs []string) ([]string, error) {
	if !w.started {
		return nil, ErrStopped
	}
	if !w.Online() {
		return nil, ErrOffline
	}
	snet, ok := w.Ipfs.PeerHost.Network().(*swarm.Network)
	if !ok {
		return nil, errors.New("peerhost network was not swarm")
	}

	swrm := snet.Swarm()

	pis, err := peersWithAddresses(addrs)
	if err != nil {
		return nil, err
	}

	output := make([]string, len(pis))
	for i, pi := range pis {
		swrm.Backoff().Clear(pi.ID)

		output[i] = "connect " + pi.ID.Pretty()

		err := w.Ipfs.PeerHost.Connect(w.Ipfs.Context(), pi)
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
	addr, pid, err := parsePeerParam(addrs)
	if addr != nil {
		w.Ipfs.Peerstore.AddAddr(pid, addr, pstore.TempAddrTTL) // temporary
	}

	if len(w.Ipfs.Peerstore.Addrs(pid)) == 0 {
		// Make sure we can find the node in question
		log.Debugf("looking up peer: %s", pid.Pretty())

		ctx, cancel := context.WithTimeout(w.Ipfs.Context(), pingTimeout)
		defer cancel()
		p, err := w.Ipfs.Routing.FindPeer(ctx, pid)
		if err != nil {
			err = fmt.Errorf("peer lookup error: %s", err)
			log.Errorf(err.Error())
			return err
		}
		w.Ipfs.Peerstore.AddAddrs(p.ID, p.Addrs, pstore.TempAddrTTL)
	}

	ctx, cancel := context.WithTimeout(w.Ipfs.Context(), pingTimeout*time.Duration(num))
	defer cancel()
	pings, err := w.Ipfs.Ping.Ping(ctx, pid)
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

func (w *Wallet) setupThread(thread *Thread) {
	thread.repoPath = w.RepoPath
	thread.ipfs = w.Ipfs
	thread.blocks = w.Datastore.Blocks()
	thread.update = func(head string) error {
		if err := w.Datastore.Threads().UpdateHead(thread.Id, head); err != nil {
			return err
		}
		thread.Head = head
		return nil
	}
	w.threads = append(w.threads, thread)
}

func (w *Wallet) getThreadBlock(blockId string) (*Thread, *Block, error) {
	block := w.Datastore.Blocks().Get(blockId)
	if block == nil {
		return nil, nil, errors.New(fmt.Sprintf("block %s not found locally", blockId))
	}
	threadId := libp2pc.ConfigEncodeKey(block.ThreadPubKey)
	var thread *Thread
	for _, t := range w.threads {
		if t.Id == threadId {
			thread = t
			break
		}
	}
	if thread == nil {
		return nil, nil, errors.New(fmt.Sprintf("could not find thread: %s", threadId))
	}
	return thread, block, nil
}

// touchDB ensures that we have a good db connection
func (w *Wallet) touchDatastore() error {
	if err := w.Datastore.Ping(); err != nil {
		log.Debug("re-opening datastore...")
		sqliteDB, err := db.Create(w.RepoPath, "")
		if err != nil {
			log.Errorf("error re-opening datastore: %s", err)
			return err
		}
		w.Datastore = sqliteDB
	}
	return nil
}

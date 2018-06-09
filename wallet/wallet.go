package wallet

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/op/go-logging"
	"github.com/rwcarlsen/goexif/exif"
	cmodels "github.com/textileio/textile-go/central/models"
	"github.com/textileio/textile-go/core/central"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net"
	trepo "github.com/textileio/textile-go/repo"
	tconfig "github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/thread"
	"github.com/textileio/textile-go/wallet/util"
	"github.com/tyler-smith/go-bip39"
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
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var log = logging.MustGetLogger("wallet")

type Config struct {
	Version    string
	RepoPath   string
	CentralAPI string
	IsMobile   bool
	IsServer   bool
	SwarmPort  string
}

type Wallet struct {
	context        oldcmds.Context
	repoPath       string
	gatewayAddr    string
	cancel         context.CancelFunc
	ipfs           *core.IpfsNode
	datastore      trepo.Datastore
	centralUserAPI string
	isMobile       bool
	started        bool
	threads        []*thread.Thread
	done           chan struct{}
}

const pingTimeout = time.Second * 10

// ErrRunning is an error for when node start is called on a started node
var ErrStarted = errors.New("node is already started")

// ErrStopped is an error for  when node stop is called on a stopped node
var ErrStopped = errors.New("node is already stopped")

// ErrOffline is an error for when online resources are requested on an offline node
var ErrOffline = errors.New("node is offline")

func NewWallet(config Config) (*Wallet, error) {
	// get database handle
	sqliteDB, err := db.Create(config.RepoPath, "")
	if err != nil {
		return nil, err
	}

	// we may be running in an uninitialized state.
	err = trepo.DoInit(config.RepoPath, config.IsMobile, config.Version, sqliteDB.Config().Init, sqliteDB.Config().Configure)
	if err != nil && err != trepo.ErrRepoExists {
		return nil, err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(config.RepoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return nil, err
	}

	// save gateway address
	gwAddr, err := repo.GetConfigKey("Addresses.Gateway")
	if err != nil {
		log.Errorf("error getting gateway address: %s", err)
		return nil, err
	}

	// if a specific swarm port was selected, set it in the config
	if config.SwarmPort != "" {
		log.Infof("using specified swarm port: %s", config.SwarmPort)
		if err := tconfig.Update(repo, "Addresses.Swarm", []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", config.SwarmPort),
			fmt.Sprintf("/ip6/::/tcp/%s", config.SwarmPort),
		}); err != nil {
			return nil, err
		}
	}

	// if this is a server node, apply the ipfs server profile
	if config.IsServer {
		if err := tconfig.Update(repo, "Addresses.NoAnnounce", tconfig.DefaultServerFilters); err != nil {
			return nil, err
		}
		if err := tconfig.Update(repo, "Swarm.AddrFilters", tconfig.DefaultServerFilters); err != nil {
			return nil, err
		}
		if err := tconfig.Update(repo, "Swarm.EnableRelayHop", true); err != nil {
			return nil, err
		}
		if err := tconfig.Update(repo, "Discovery.MDNS.Enabled", false); err != nil {
			return nil, err
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
		repoPath:       config.RepoPath,
		gatewayAddr:    gwAddr.(string),
		datastore:      sqliteDB,
		centralUserAPI: fmt.Sprintf("%s/api/v1/users", config.CentralAPI),
		isMobile:       config.IsMobile,
	}, nil
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
		if err := util.PrintSwarmAddrs(w.ipfs); err != nil {
			log.Errorf("failed to read listening addresses: %s", err)
		}
		log.Info("wallet is online")
	}()

	// setup threads
	for _, mod := range w.datastore.Threads().List("") {
		_, err := w.loadThread(&mod)
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
	w.threads = nil

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
	return w.ipfs.OnlineMode()
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
	res, err := central.SignUp(reg, w.centralUserAPI)
	if err != nil {
		log.Errorf("signup error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signup error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.datastore.Profile().SignIn(reg.Username, res.Session.AccessToken, res.Session.RefreshToken); err != nil {
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
	res, err := central.SignIn(creds, w.centralUserAPI)
	if err != nil {
		log.Errorf("signin error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signin error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.datastore.Profile().SignIn(creds.Username, res.Session.AccessToken, res.Session.RefreshToken); err != nil {
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
	un, err := w.datastore.Profile().GetUsername()
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
	at, _, err := w.datastore.Profile().GetTokens()
	if err != nil {
		return "", err
	}
	return at, nil
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

func (w *Wallet) GetThreadByName(name string) *thread.Thread {
	for _, thrd := range w.threads {
		if thrd.Name == name {
			return thrd
		}
	}
	return nil
}

// AddThread adds a thread with a given name and secret key
func (w *Wallet) AddThread(name string, secret libp2pc.PrivKey) (*thread.Thread, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	log.Debugf("adding a new thread: %s", name)

	skb, err := secret.Bytes()
	if err != nil {
		return nil, err
	}
	pkb, err := secret.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	pk := libp2pc.ConfigEncodeKey(pkb)

	// index a new thread
	threadModel := &trepo.Thread{
		Id:      pk,
		Name:    name,
		PrivKey: skb,
	}
	if err := w.datastore.Threads().Add(threadModel); err != nil {
		return nil, err
	}
	thrd, err := w.loadThread(threadModel)
	if err != nil {
		return nil, err
	}
	return thrd, nil
}

// AddThreadWithMnemonic adds a thread with a given name and mnemonic phrase
func (w *Wallet) AddThreadWithMnemonic(name string, mnemonic string) (*thread.Thread, error) {
	// use phrase if provided
	if mnemonic == "" {
		var err error
		mnemonic, err = util.CreateMnemonic(bip39.NewEntropy, bip39.NewMnemonic)
		if err != nil {
			return nil, err
		}
		log.Debugf("generating Ed25519 keypair for: %s", name)
	} else {
		log.Debugf("regenerating Ed25519 keypair from mnemonic phrase for: %s", name)
	}

	// create the bip39 seed from the phrase
	seed := bip39.NewSeed(mnemonic, "")
	key, err := util.IdentityKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}

	// convert to a libp2p crypto private key
	secret, err := libp2pc.UnmarshalPrivateKey(key)
	if err != nil {
		return nil, err
	}

	return w.AddThread(name, secret)
}

// PublishThreads publishes HEAD for each thread
func (w *Wallet) PublishThreads() {
	for _, t := range w.threads {
		go func(thrd *thread.Thread) {
			thrd.Publish()
		}(t)
	}
}

// TODO: add node master pk to dir
func (w *Wallet) AddPhoto(path string) (*model.AddResult, error) {
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
	thumb, err := makeThumbnail(photo, model.ThumbnailWidth)
	if err != nil {
		return nil, err
	}

	// path info
	fpath := photo.Name()
	ext := strings.ToLower(filepath.Ext(fpath))

	// get username, ignoring if not present (not signed in)
	username, _ := w.datastore.Profile().GetUsername()

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
	photocypher, err := util.GetEncryptedReaderBytes(photo, key)
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

	// finish request
	if err := request.Finish(); err != nil {
		return nil, err
	}

	// all done
	return &model.AddResult{Id: id, Key: key, RemoteRequest: request}, nil
}

func (w *Wallet) FindBlock(target string) (*trepo.Block, error) {
	block := w.datastore.Blocks().GetByTarget(target)
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
	thrd, block, err := w.getThreadBlock(blockId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// get bytes
	cypher, err := util.GetDataAtPath(w.ipfs, path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// decrypt the file key
	key, err := thrd.Decrypt(block.TargetKey)
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

func (w *Wallet) GetDataAtPath(path string) ([]byte, error) {
	if !w.started {
		return nil, ErrStopped
	}
	return util.GetDataAtPath(w.ipfs, path)
}

func (w *Wallet) GetIPFSPeerID() (string, error) {
	if !w.started {
		return "", ErrStopped
	}
	return w.ipfs.Identity.Pretty(), nil
}

// GetIPFSPubKeyString returns the base64 encoded public ipfs peer key
func (w *Wallet) GetIPFSPubKeyString() (string, error) {
	if !w.started {
		return "", ErrStopped
	}
	pkb, err := w.ipfs.PrivateKey.GetPublic().Bytes()
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

// TODO: connect to relay if needed
func (w *Wallet) Publish(topic string, payload []byte) error {
	if !w.Online() {
		return ErrOffline
	}
	return w.ipfs.Floodsub.Publish(topic, payload)
}

func (w *Wallet) Subscribe(topic string) (*floodsub.Subscription, error) {
	if !w.Online() {
		return nil, ErrOffline
	}
	return w.ipfs.Floodsub.Subscribe(topic)
}

func (w *Wallet) IFPSPeers() ([]libp2pn.Conn, error) {
	if !w.Online() {
		return nil, ErrOffline
	}
	return w.ipfs.PeerHost.Network().Conns(), nil
}

// WaitForInvite to join
// TODO: needs cleanup to handle a generic invite
func (w *Wallet) WaitForInvite() {
	if !w.Online() {
		return
	}
	// we're in a lonesome state here, we can just sub to our own
	// peer id and hope somebody sends us a priv key to join a thread with
	self := w.ipfs.Identity.Pretty()
	sub, err := w.ipfs.Floodsub.Subscribe(self)
	if err != nil {
		log.Errorf("error creating subscription: %s", err)
		return
	}
	log.Infof("waiting for invite at own peer id: %s\n", self)

	ctx, cancel := context.WithCancel(context.Background())
	cancelCh := make(chan struct{})
	go func() {
		for {
			msg, err := sub.Next(ctx)
			if err == io.EOF || err == context.Canceled {
				log.Debugf("wait subscription ended: %s", err)
				return
			} else if err != nil {
				log.Debugf(err.Error())
				return
			}
			from := msg.GetFrom().Pretty()
			log.Infof("got pairing request from: %s\n", from)

			// get private peer key and decrypt the phrase
			skb, err := crypto.Decrypt(w.ipfs.PrivateKey, msg.GetData())
			if err != nil {
				log.Errorf("error decrypting msg data: %s", err)
				return
			}
			secret, err := libp2pc.UnmarshalPrivateKey(skb)
			if err != nil {
				log.Errorf("error unmarshaling mobile private key: %s", err)
				return
			}

			// create a new album for the room
			// TODO: let user name this or take phone's name, e.g., bob's iphone
			// TODO: or auto name it, cause this means only one pairing can happen
			_, err = w.AddThread("mobile", secret)
			if err != nil {
				log.Errorf("error adding mobile thread: %s", err)
				return
			}

			// we're done
			close(cancelCh)
		}
	}()

	for {
		select {
		case <-cancelCh:
			cancel()
			return
		case <-w.ipfs.Context().Done():
			cancel()
			return
		}
	}
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

func (w *Wallet) loadThread(model *trepo.Thread) (*thread.Thread, error) {
	id := model.Id // save value locally
	threadConfig := &thread.Config{
		RepoPath: w.repoPath,
		Ipfs:     w.ipfs,
		Blocks:   w.datastore.Blocks(),
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
	}
	thrd, err := thread.NewThread(model, threadConfig)
	if err != nil {
		return nil, err
	}
	w.threads = append(w.threads, thrd)
	return thrd, nil
}

func (w *Wallet) getThreadBlock(blockId string) (*thread.Thread, *trepo.Block, error) {
	block := w.datastore.Blocks().Get(blockId)
	if block == nil {
		return nil, nil, errors.New(fmt.Sprintf("block %s not found locally", blockId))
	}
	threadId := libp2pc.ConfigEncodeKey(block.ThreadPubKey)
	var thrd *thread.Thread
	for _, t := range w.threads {
		if t.Id == threadId {
			thrd = t
			break
		}
	}
	if thrd == nil {
		return nil, nil, errors.New(fmt.Sprintf("could not find thread: %s", threadId))
	}
	return thrd, block, nil
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

func makeThumbnail(photo *os.File, width int) ([]byte, error) {
	img, _, err := image.Decode(photo)
	if err != nil {
		return nil, err
	}
	thumb := imaging.Resize(img, width, 0, imaging.Lanczos)
	buff := new(bytes.Buffer)
	if err = jpeg.Encode(buff, thumb, nil); err != nil {
		return nil, err
	}
	photo.Seek(0, 0) // be kind, rewind
	return buff.Bytes(), nil
}

// TODO: get image size info
func getMetadata(photo *os.File, path string, ext string, username string) (model.PhotoMetadata, error) {
	var created time.Time
	var lat, lon float64
	x, err := exif.Decode(photo)
	if err == nil {
		// time taken
		createdTmp, err := x.DateTime()
		if err == nil {
			created = createdTmp
		}
		// coords taken
		latTmp, lonTmp, err := x.LatLong()
		if err == nil {
			lat, lon = latTmp, lonTmp
		}
	}
	meta := model.PhotoMetadata{
		FileMetadata: model.FileMetadata{
			Metadata: model.Metadata{
				Username: username,
				Created:  created,
				Added:    time.Now(),
			},
			Name: strings.TrimSuffix(filepath.Base(path), ext),
			Ext:  ext,
		},
		Latitude:  lat,
		Longitude: lon,
	}
	photo.Seek(0, 0) // be kind, rewind
	return meta, nil
}

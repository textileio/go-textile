package core

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/op/go-logging"
	"github.com/tyler-smith/go-bip39"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/textileio/textile-go/net"
	trepo "github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/repo/photos"

	utilmain "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/commands"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreapi"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/config"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo"
	lockfile "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo/lock"

	"database/sql"
	"encoding/json"
	"gx/ipfs/QmSFihvoND3eDaAYRCeLgLPt62yCPgMZs1NSZmKFEtJQQw/go-libp2p-floodsub"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

const VERSION = "0.0.1"

var fileLogFormat = logging.MustStringFormatter(
	`%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
)

var log = logging.MustGetLogger("core")

var Node *TextileNode

type TextileNode struct {
	// Context for issuing IPFS commands
	Context oldcmds.Context

	// IPFS node object
	IpfsNode *core.IpfsNode

	// The path to the openbazaar repo in the file system
	RepoPath string

	// Database for storing node specific data
	Datastore trepo.Datastore

	// Function to call for shutdown
	Cancel context.CancelFunc

	// Indicated if gc and gateways are up
	ServicesUp bool

	// IPFS configuration used to instantiate new ipfs nodes
	ipfsConfig *core.BuildCfg

	// Whether or not we're running on a mobile device
	isMobile bool

	// Signal we're leaving the room
	leaveRoomCh chan struct{}

	// Signal we left the room
	LeftRoomCh chan struct{}
}

type PhotoList struct {
	Hashes []string `json:"hashes"`
}

func NewNode(repoPath string, isMobile bool, logLevel logging.Level) (*TextileNode, error) {
	// shutdown is not clean here yet, so we have to hackily remove
	// the lockfile that should have been removed on shutdown
	// before we start up again
	// TODO: Figure out how to make this work as intended, without doing this
	repoLockFile := filepath.Join(repoPath, lockfile.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(repoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// log handling
	w := &lumberjack.Logger{
		Filename:   path.Join(repoPath, "logs", "textile.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     30, // days
	}
	backendFile := logging.NewLogBackend(w, "", 0)
	backendFileFormatter := logging.NewBackendFormatter(backendFile, fileLogFormat)
	logging.SetBackend(backendFileFormatter)
	logging.SetLevel(logLevel, "")

	// get database handle for wallet indexes
	sqliteDB, err := db.Create(repoPath, "")
	if err != nil {
		return nil, err
	}

	// we may be running in an uninitialized state.
	err = trepo.DoInit(repoPath, isMobile, sqliteDB.Config().Init)
	if err != nil && err != trepo.ErrRepoExists {
		return nil, err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		log.Errorf("error opening repo: %s", err)
		return nil, err
	}

	// determine the best routing
	var routingOption core.RoutingOption
	if isMobile {
		// TODO: Determine best value for this setting on mobile
		// cfg.Swarm.DisableNatPortMap = true
		routingOption = core.DHTClientOption
	} else {
		routingOption = core.DHTOption
	}

	// assemble node config
	ncfg := &core.BuildCfg{
		Repo:      repo,
		Permanent: true, // temporary way to signify that node is permanent
		Online:    true,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
			"mplex":  true,
		},
		Routing: routingOption,
	}

	// finally, construct our node
	return &TextileNode{
		RepoPath:    repoPath,
		Datastore:   sqliteDB,
		ipfsConfig:  ncfg,
		isMobile:    isMobile,
		leaveRoomCh: make(chan struct{}),
		LeftRoomCh:  make(chan struct{}),
	}, nil
}

func (t *TextileNode) ConfigureDatastore(mnemonic string) error {
	log.Info("configuring textile datastore...")
	if mnemonic == "" {
		var err error
		mnemonic, err = createMnemonic(bip39.NewEntropy, bip39.NewMnemonic)
		if err != nil {
			log.Errorf("error creating mnemonic: %s", err)
			return err
		}
		log.Infof("generating %v-bit Ed25519 keypair...", trepo.NBitsForKeypair)
	} else {
		log.Infof("regenerating Ed25519 keypair from mnemonic phrase...")
	}
	seed := bip39.NewSeed(mnemonic, "")
	identityKey, err := identityKeyFromSeed(seed, trepo.NBitsForKeypair)
	if err != nil {
		log.Errorf("error creating identity from seed: %s", err)
		return err
	}

	return t.Datastore.Config().Configure(mnemonic, identityKey, time.Now())
}

func (t *TextileNode) IsDatastoreConfigured() bool {
	_, err := t.Datastore.Config().GetMnemonic()
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Errorf("error checking if datastore is configured: %s", err)
			return false
		}
	}
	return true
}

func (t *TextileNode) Start() error {
	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		log.Errorf("setting file descriptor limit: %s", err)
	}

	cctx, cancel := context.WithCancel(context.Background())
	t.Cancel = cancel

	ctx := oldcmds.Context{}

	if t.IpfsNode != nil {
		return nil
	}

	log.Info("starting node...")
	nd, err := core.NewNode(cctx, t.ipfsConfig)
	if err != nil {
		return err
	}
	nd.SetLocal(false)

	if err := printSwarmAddrs(nd); err != nil {
		log.Errorf("failed to read listening addresses: %s", err)
	}

	ctx.Online = true
	ctx.ConfigRoot = t.RepoPath
	ctx.LoadConfig = func(path string) (*config.Config, error) {
		return fsrepo.ConfigAt(t.RepoPath)
	}
	ctx.ConstructNode = func() (*core.IpfsNode, error) {
		return nd, nil
	}

	t.Context = ctx
	t.IpfsNode = nd

	if t.isMobile {
		log.Info("mobile node is ready")
	} else {
		log.Info("desktop node is ready")
	}

	return nil
}

func (t *TextileNode) StartServices() (<-chan error, error) {
	if t.isMobile {
		return nil, errors.New("services not available on mobile")
	}

	// repo blockstore GC
	var gcErrc <-chan error
	var err error
	gcErrc, err = runGC(t.IpfsNode.Context(), t.IpfsNode)
	if err != nil {
		log.Errorf("error starting gc: %s", err)
		return nil, err
	}

	// construct http gateway
	var gwErrc <-chan error
	gwErrc, err = serveHTTPGateway(&t.Context)
	if err != nil {
		log.Errorf("error starting gateway: %s", err)
		return nil, err
	}

	// construct decrypting http gateway proxy
	var gwpErrc <-chan error
	gwpErrc, err = ServeHTTPGatewayProxy(t)
	if err != nil {
		log.Errorf("error starting decrypting gateway: %s", err)
		return nil, err
	}
	t.ServicesUp = true

	// merge error channels
	return merge(gwErrc, gcErrc, gwpErrc), nil
}

func (t *TextileNode) Stop() error {
	repoLockFile := filepath.Join(t.RepoPath, lockfile.LockFile)
	if err := os.Remove(repoLockFile); err != nil {
		log.Errorf("error removing lock: %s", err)
		return err
	}
	t.Datastore.Close()
	dsLockFile := filepath.Join(t.RepoPath, "datastore", "LOCK")
	if err := os.Remove(dsLockFile); err != nil {
		log.Errorf("error removing ds lock: %s", err)
		return err
	}
	if err := t.IpfsNode.Close(); err != nil {
		log.Errorf("error closing ipfs node: %s", err)
		return err
	}
	t.IpfsNode = nil
	return nil
}

func (t *TextileNode) JoinRoom(datac chan string) {
	rid, err := t.GetRoomID()
	if err != nil {
		return
	}

	// create the subscription
	sub, err := t.IpfsNode.Floodsub.Subscribe(rid.Pretty())
	if err != nil {
		log.Errorf("error creating subscription: %s", err)
		return
	}
	log.Infof("joined room: %s\n", rid.Pretty())

	leave := func() {
		sub.Cancel()
		// not all listeners may be interested in this signal
		// so we do it non-blocking
		select {
		case t.LeftRoomCh <- struct{}{}:
		default:
		}
		log.Infof("left room: %s\n", sub.Topic())
	}

	defer close(datac)
	go func() {
		for {
			// unload new message
			msg, err := sub.Next(t.IpfsNode.Context())
			if err == io.EOF || err == context.Canceled {
				return
			} else if err != nil {
				log.Infof(err.Error())
				return
			}

			// handle the update
			if err = t.handleRoomUpdate(msg); err != nil {
				log.Errorf("error handling room update: %s", err)
				continue
			}

			datac <- string(msg.GetData())
		}
	}()

	// block so we can shutdown with the leave room signal
	for {
		select {
		case <-t.leaveRoomCh:
			leave()
			return
		case <-t.IpfsNode.Context().Done():
			leave()
			return
		}
	}
}

func (t *TextileNode) LeaveRoom() {
	t.leaveRoomCh <- struct{}{}
}

func (t *TextileNode) WaitForRoom() {
	// we're in a lonesome state here, we can just sub to our own
	// peer id and hope somebody sends us a priv key to join a room with
	id := t.IpfsNode.Identity.Pretty()
	sub, err := t.IpfsNode.Floodsub.Subscribe(id)
	if err != nil {
		log.Errorf("error creating subscription: %s", err)
		return
	}
	log.Infof("waiting for room at own peer id: %s\n", id)

	defer sub.Cancel()
	go func() {
		for {
			msg, err := sub.Next(t.IpfsNode.Context())
			if err == io.EOF || err == context.Canceled {
				return
			} else if err != nil {
				log.Infof(err.Error())
				return
			}
			from := msg.GetFrom().Pretty()
			log.Infof("got pairing request from: %s\n", from)

			// get private peer key and decrypt the phrase
			sk, err := t.UnmarshalPrivatePeerKey()
			if err != nil {
				log.Errorf("error unmarshaling priv peer key: %s", err)
				return
			}
			p, err := net.Decrypt(sk, msg.GetData())
			if err != nil {
				log.Errorf("error decrypting msg data: %s", err)
				return
			}
			ps := string(p)
			log.Infof("decrypted mnemonic phrase as: %s\n", ps)

			// setup datastore with phrase and close sub
			_ = t.ConfigureDatastore(ps)
			return
		}
	}()

	for {
		select {
		case <-t.IpfsNode.Context().Done():
			return
		}
	}
}

func (t *TextileNode) AddPhoto(path string, thumb string) (*net.MultipartRequest, error) {
	// read file from disk
	p, err := os.Open(path)
	if err != nil {
		log.Errorf("error opening photo: %s", err)
		return nil, err
	}
	defer p.Close()

	th, err := os.Open(thumb)
	if err != nil {
		log.Errorf("error opening thumb: %s", err)
		return nil, err
	}
	defer th.Close()

	// unmarshal private key
	sk, err := t.UnmarshalPrivateKey()
	if err != nil {
		log.Errorf("error unmarshaling priv key: %s", err)
		return nil, err
	}

	// get last photo
	var lc string
	recent := t.Datastore.Photos().GetPhotos("", 1)
	if len(recent) > 0 {
		lc = recent[0].Cid
	}

	// add it
	mr, md, err := photos.Add(t.IpfsNode, sk.GetPublic(), p, th, lc)
	if err != nil {
		log.Errorf("error adding photo: %s", err)
		return nil, err
	}

	// index
	err = t.Datastore.Photos().Put(mr.Boundary, lc, md)
	if err != nil {
		log.Errorf("error indexing photo: %s", err)
		return nil, err
	}

	// publish
	tp, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		log.Errorf("error getting id from priv key: %s", err)
		return nil, err
	}
	err = t.IpfsNode.Floodsub.Publish(tp.Pretty(), []byte(mr.Boundary))
	if err != nil {
		log.Errorf("error publishing photo update: %s", err)
		return nil, err
	}

	return mr, nil
}

func (t *TextileNode) GetPhotos(offsetId string, limit int) *PhotoList {
	// query for available hashes
	list := t.Datastore.Photos().GetPhotos(offsetId, limit)

	// return json list of hashes
	res := &PhotoList{
		Hashes: make([]string, len(list)),
	}
	for i := range list {
		res.Hashes[i] = list[i].Cid
	}

	return res
}

// pass in Qm../thumb, or Qm../photo for full image
func (t *TextileNode) GetFile(path string) ([]byte, error) {
	// get bytes
	cb, err := t.getDataAtPath(path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// unmarshal private key
	sk, err := t.UnmarshalPrivateKey()
	if err != nil {
		return nil, err
	}
	b, err := net.Decrypt(sk, cb)
	if err != nil {
		log.Errorf("error decrypting file: %s", err)
		return nil, err
	}

	return b, err
}

func (t *TextileNode) GetMetaData(hash string) (*photos.Metadata, error) {
	b, err := t.GetFile(fmt.Sprintf("%s/meta", hash))
	if err != nil {
		log.Errorf("error getting meta file: %s", err)
		return nil, err
	}

	var data *photos.Metadata
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Errorf("error unmarshaling meta file: %s", err)
		return nil, err
	}

	return data, nil
}

func (t *TextileNode) GetLastHash(hash string) (string, error) {
	b, err := t.GetFile(fmt.Sprintf("%s/last", hash))
	if err != nil {
		log.Errorf("error getting last hash file: %s", err)
		return "", err
	}

	return string(b), nil
}

func (t *TextileNode) UnmarshalPrivatePeerKey() (libp2p.PrivKey, error) {
	cfg, err := t.Context.GetConfig()
	if err != nil {
		return nil, err
	}
	skb, err := base64.StdEncoding.DecodeString(cfg.Identity.PrivKey)
	if err != nil {
		return nil, err
	}
	sk, err := libp2p.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}

	// check
	id2, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	if id2 != t.IpfsNode.Identity {
		return nil, fmt.Errorf("private key in config does not match id: %s != %s", t.IpfsNode.Identity, id2)
	}

	return sk, nil
}

func (t *TextileNode) UnmarshalPrivateKey() (libp2p.PrivKey, error) {
	kb, err := t.Datastore.Config().GetIdentityKey()
	if err != nil {
		log.Errorf("error unmarshaling priv key: %s", err)
		return nil, err
	}
	return libp2p.UnmarshalPrivateKey(kb)
}

func (t *TextileNode) GetRoomID() (*peer.ID, error) {
	sk, err := t.UnmarshalPrivateKey()
	if err != nil {
		return nil, err
	}
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		log.Errorf("error getting id from priv key: %s", err)
		return nil, err
	}
	return &id, nil
}

func (t *TextileNode) GetPublicPeerKeyString() (string, error) {
	sk, err := t.UnmarshalPrivatePeerKey()
	if err != nil {
		log.Errorf("error unmarshaling priv peer key: %s", err)
		return "", err
	}
	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		log.Errorf("error getting pub key bytes: %s", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pkb), nil
}

func (t *TextileNode) getDataAtPath(path string) ([]byte, error) {
	// convert string to an ipfs path
	ip, err := coreapi.ParsePath(path)
	if err != nil {
		return nil, err
	}

	api := coreapi.NewCoreAPI(t.IpfsNode)
	r, err := api.Unixfs().Cat(t.IpfsNode.Context(), ip)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

func (t *TextileNode) handleRoomUpdate(msg *floodsub.Message) error {
	// unpack message
	from := msg.GetFrom().Pretty()
	hash := string(msg.GetData())
	api := coreapi.NewCoreAPI(t.IpfsNode)
	log.Infof("got update from: %s", from)

	// recurse back in time starting at this hash
	err := t.handleHash(hash, api)
	if err != nil {
		log.Errorf("error handling hash: %s", err)
	}

	return nil
}

func (t *TextileNode) handleHash(hash string, api iface.CoreAPI) error {
	log.Infof("handling update: %s...", hash)

	// convert string to an ipfs path
	ip, err := coreapi.ParsePath(hash)
	if err != nil {
		return err
	}

	// check if we aleady have this hash
	set := t.Datastore.Photos().GetPhoto(hash)
	if set != nil {
		log.Infof("%s exists, aborting update", hash)
		return nil
	}

	// pin it
	log.Infof("pinning %s recursively...", hash)
	err = api.Pin().Add(t.IpfsNode.Context(), ip, api.Pin().WithRecursive(true))
	if err != nil {
		return err
	}

	// unpack data set
	log.Infof("unpacking %s...", hash)
	md, err := t.GetMetaData(hash)
	if err != nil {
		return err
	}
	last, err := t.GetLastHash(hash)
	if err != nil {
		return err
	}

	// index
	log.Infof("indexing %s...", hash)
	err = t.Datastore.Photos().Put(hash, last, md)
	if err != nil {
		return err
	}

	// check last hash
	return t.handleHash(last, api)
}

func createMnemonic(newEntropy func(int) ([]byte, error), newMnemonic func([]byte) (string, error)) (string, error) {
	entropy, err := newEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := newMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

func identityKeyFromSeed(seed []byte, bits int) ([]byte, error) {
	hm := hmac.New(sha256.New, []byte("scythian horde"))
	hm.Write(seed)
	reader := bytes.NewReader(hm.Sum(nil))
	sk, _, err := libp2p.GenerateKeyPairWithReader(libp2p.Ed25519, bits, reader)
	if err != nil {
		return nil, err
	}
	encodedKey, err := sk.Bytes()
	if err != nil {
		return nil, err
	}
	return encodedKey, nil
}

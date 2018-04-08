package core

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/op/go-logging"
	"github.com/tyler-smith/go-bip39"

	"github.com/textileio/textile-go/net"
	trepo "github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/db"
	"github.com/textileio/textile-go/repo/wallet"

	utilmain "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/commands"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreapi"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/config"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo"
	lockfile "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo/lock"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

var VERSION = "0.0.1"

var stdoutLogFormat = logging.MustStringFormatter(
	`%{color:reset}%{color}%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
)

var logger logging.Backend

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

	// IPFS configuration used to instantiate new ipfs nodes
	ipfsConfig *core.BuildCfg

	// Whether or not we're running on a mobile device
	isMobile bool
}

type PhotoList struct {
	Hashes []string `json:"hashes"`
}

func NewNode(repoPath string, isMobile bool) (*TextileNode, error) {

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		fmt.Errorf("setting file descriptor limit: %s", err)
	}

	// shutdown is not clean here yet, so we have to hackily remove
	// the lockfile that should have been removed on shutdown
	// before we start up again
	// TODO: Figure out how to make this work as intended, without doing this
	repoLockFile := filepath.Join(repoPath, lockfile.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(repoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// get database handle for wallet indexes
	sqliteDB, err := db.Create(repoPath, "")
	if err != nil {
		return nil, err
	}

	// we may be running in an uninitialized state.
	err = trepo.DoInit(os.Stdout, repoPath, sqliteDB.Config().Init)
	if err != nil && err != trepo.ErrRepoExists {
		return nil, err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	var routingOption core.RoutingOption
	if isMobile {
		// some of the below are taken from the not-yet-released "lowpower" profile preset
		cfg, err := repo.Config()
		if err != nil {
			return nil, err
		}
		cfg.Reprovider.Interval = "0"
		cfg.Swarm.ConnMgr.LowWater = 20
		cfg.Swarm.ConnMgr.HighWater = 40
		cfg.Swarm.ConnMgr.GracePeriod = time.Minute.String()
		// TODO: Determine best value for this setting on mobile
		// cfg.Swarm.DisableNatPortMap = true

		routingOption = core.DHTClientOption
	} else {
		routingOption = core.DHTOption
	}

	// Start assembling node config
	ncfg := &core.BuildCfg{
		Repo:      repo,
		Permanent: true, // It is temporary way to signify that node is permanent
		Online:    true,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": false,
			"mplex":  true,
		},
		Routing: routingOption,
	}

	return &TextileNode{RepoPath: repoPath, Datastore: sqliteDB, ipfsConfig: ncfg, isMobile: isMobile}, nil
}

func (t *TextileNode) ConfigureDatastore(mnemonic string) error {
	if mnemonic == "" {
		var err error
		mnemonic, err = createMnemonic(bip39.NewEntropy, bip39.NewMnemonic)
		if err != nil {
			return err
		}
		fmt.Println("Generating Ed25519 keypair...")
	} else {
		fmt.Println("Regenerating Ed25519 keypair from mnemonic phrase...")
	}
	seed := bip39.NewSeed(mnemonic, "")
	identityKey, err := identityKeyFromSeed(seed, trepo.NBitsForKeypair)
	if err != nil {
		return err
	}

	return t.Datastore.Config().Configure(mnemonic, identityKey, time.Now())
}

func (t *TextileNode) Start() error {
	cctx, cancel := context.WithCancel(context.Background())
	t.Cancel = cancel

	ctx := oldcmds.Context{}

	if t.IpfsNode != nil {
		return nil
	}

	fmt.Println("Starting node...")
	fmt.Println("Repo directory: ", t.RepoPath)

	nd, err := core.NewNode(cctx, t.ipfsConfig)
	if err != nil {
		return err
	}
	nd.SetLocal(false)

	if err := PrintSwarmAddrs(nd); err != nil {
		fmt.Errorf("failed to read listening addresses: %s", err)
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

	errc := make(chan error)
	go func() {
		_, err := ctx.ConstructNode()
		errc <- err
		close(errc)
	}()

	fmt.Printf("Node is ready\n")
	for err := range Merge(errc) {
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (t *TextileNode) Stop() error {
	repoLockFile := filepath.Join(t.RepoPath, lockfile.LockFile)
	if err := os.Remove(repoLockFile); err != nil {
		return err
	}
	dsLockFile := filepath.Join(t.RepoPath, "datastore", "LOCK")
	if err := os.Remove(dsLockFile); err != nil {
		return err
	}
	if err := t.IpfsNode.Close(); err != nil {
		return err
	}
	t.IpfsNode = nil
	return nil
}

func (t *TextileNode) AddPhoto(path string, thumb string) (*net.MultipartRequest, error) {
	// read file from disk
	p, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer p.Close()

	th, err := os.Open(thumb)
	if err != nil {
		return nil, err
	}
	defer th.Close()

	// unmarshal private key
	sk, err := t.unmarshalPrivateKey()
	if err != nil {
		return nil, err
	}

	// add it
	mr, err := wallet.AddPhoto(t.IpfsNode, sk, p, th)
	if err != nil {
		return nil, err
	}

	// index
	t.Datastore.Photos().Put(mr.Boundary, time.Now())

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
	// convert string to a ipfs path
	ipath, err := coreapi.ParsePath(path)
	if err != nil {
		return nil, err
	}

	api := coreapi.NewCoreAPI(t.IpfsNode)
	r, err := api.Unixfs().Cat(t.IpfsNode.Context(), ipath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// read bytes and convert to base64 string
	cb, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// unmarshal private key
	sk, err := t.unmarshalPrivateKey()
	if err != nil {
		return nil, err
	}
	b, err := net.Decrypt(sk, cb)
	if err != nil {
		return nil, err
	}

	return b, err
}

func (t *TextileNode) unmarshalPrivateKey() (libp2p.PrivKey, error) {
	kb, err := t.Datastore.Config().GetIdentityKey()
	if err != nil {
		return nil, err
	}
	return libp2p.UnmarshalPrivateKey(kb)
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

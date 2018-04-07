package mobile

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	tcore "github.com/textileio/textile-go/core"
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

type Node struct {
	node       *tcore.TextileNode
	config     MobileConfig
	cancel     context.CancelFunc
	ipfsConfig *core.BuildCfg
}
type Mobile struct{}

type PhotoList struct {
	Hashes []string `json:"hashes"`
}

func NewTextile(repoPath string, apiHost string) *Node {
	nodeconfig := MobileConfig{
		RepoPath: repoPath,
		ApiHost:  apiHost,
	}

	var m Mobile
	node, err := m.NewNode(nodeconfig)
	if err != nil {
		fmt.Println(err)
	}
	return node
}

func (m *Mobile) NewNode(config MobileConfig) (*Node, error) {

	// shutdown is not clean here yet, so we have to hackily remove
	// the lockfile that should have been removed on shutdown
	// before we start up again
	repoLockFile := filepath.Join(config.RepoPath, lockfile.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(config.RepoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		fmt.Errorf("setting file descriptor limit: %s", err)
	}

	// get database handle for wallet indexes
	sqliteDB, err := db.Create(config.RepoPath, "")
	if err != nil {
		return nil, err
	}

	// we may be running in an uninitialized state.
	err = trepo.DoInit(os.Stdout, config.RepoPath, time.Now(), sqliteDB.Config().Init)
	if err != nil && err != trepo.ErrRepoExists {
		return nil, err
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(config.RepoPath)
	if err != nil {
		return nil, err
	}

	// tweak default (textile) config for mobile
	// some of the below are taken from the not-yet-released "lowpower" profile preset
	cfg, err := repo.Config()
	if err != nil {
		return nil, err
	}
	cfg.Reprovider.Interval = "0"
	cfg.Swarm.ConnMgr.LowWater = 20
	cfg.Swarm.ConnMgr.HighWater = 40
	cfg.Swarm.ConnMgr.GracePeriod = time.Minute.String()
	//cfg.Swarm.DisableNatPortMap = true

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
		Routing: core.DHTClientOption,
	}

	// Textile node setup
	tcore.Node = &tcore.TextileNode{
		RepoPath:  config.RepoPath,
		Datastore: sqliteDB,
	}

	return &Node{config: config, node: tcore.Node, ipfsConfig: ncfg}, nil
}

func (n *Node) Start() error {
	cctx, cancel := context.WithCancel(context.Background())
	n.cancel = cancel

	ctx := oldcmds.Context{}

	// TODO: we may need a check to ensure it is running
	if n.node.IpfsNode != nil {
		return nil
	}

	fmt.Println("Starting node...")
	fmt.Println("Repo directory: ", n.config.RepoPath)

	nd, err := core.NewNode(cctx, n.ipfsConfig)
	if err != nil {
		return err
	}
	nd.SetLocal(false)

	if err := tcore.PrintSwarmAddrs(nd); err != nil {
		fmt.Errorf("failed to read listening addresses: %s", err)
	}

	ctx.Online = true
	ctx.ConfigRoot = n.config.RepoPath
	ctx.LoadConfig = func(path string) (*config.Config, error) {
		return fsrepo.ConfigAt(n.config.RepoPath)
	}
	ctx.ConstructNode = func() (*core.IpfsNode, error) {
		return nd, nil
	}

	n.node.Context = ctx
	n.node.IpfsNode = nd

	errc := make(chan error)
	go func() {
		_, err := ctx.ConstructNode()
		errc <- err
		close(errc)
	}()

	fmt.Printf("Node is ready\n")
	for err := range tcore.Merge(errc) {
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (n *Node) Stop() error {
	repoLockFile := filepath.Join(tcore.Node.RepoPath, lockfile.LockFile)
	if err := os.Remove(repoLockFile); err != nil {
		return err
	}
	dsLockFile := filepath.Join(tcore.Node.RepoPath, "datastore", "LOCK")
	if err := os.Remove(dsLockFile); err != nil {
		return err
	}
	return tcore.Node.IpfsNode.Close()
}

func (n *Node) AddPhoto(path string, thumb string) (*net.MultipartRequest, error) {
	// read file from disk
	p, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer p.Close()

	t, err := os.Open(thumb)
	if err != nil {
		return nil, err
	}
	defer t.Close()

	// unmarshal private key
	sk, err := n.unmarshalPrivateKey()
	if err != nil {
		return nil, err
	}

	// add it
	mr, err := wallet.AddPhoto(n.node.IpfsNode, sk, p, t)
	if err != nil {
		return nil, err
	}

	// index
	n.node.Datastore.Photos().Put(mr.Boundary, time.Now())

	return mr, nil
}

func (n *Node) GetPhotos(offsetId string, limit int) (string, error) {
	// query for available hashes
	list := n.node.Datastore.Photos().GetPhotos(offsetId, limit)

	// return json list of hashes
	res := &PhotoList{
		Hashes: make([]string, len(list)),
	}
	for i := range list {
		res.Hashes[i] = list[i].Cid
	}

	// convert to json
	jsonb, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(jsonb), nil
}

// pass in Qm../thumb, or Qm../photo for full image
func (n *Node) GetPhotoBase64String(path string) (string, error) {
	// convert string to a ipfs path
	ipath, err := coreapi.ParsePath(path)
	if err != nil {
		return "", nil
	}

	api := coreapi.NewCoreAPI(n.node.IpfsNode)
	r, err := api.Unixfs().Cat(n.node.IpfsNode.Context(), ipath)
	if err != nil {
		return "", nil
	}
	defer r.Close()

	// read bytes and convert to base64 string
	cb, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	// unmarshal private key
	sk, err := n.unmarshalPrivateKey()
	if err != nil {
		return "", err
	}
	b, err := net.Decrypt(sk, cb)
	if err != nil {
		return "", err
	}

	// do the encoding
	bs64 := base64.StdEncoding.EncodeToString(b)

	return bs64, nil
}

// provides the user's recovery phrase for their private key
func (n *Node) GetRecoveryPhrase() (string, error) {
	mnemonic, err := n.node.Datastore.Config().GetMnemonic()
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func (n *Node) unmarshalPrivateKey() (libp2p.PrivKey, error) {
	kb, err := n.node.Datastore.Config().GetIdentityKey()
	if err != nil {
		return nil, err
	}
	return libp2p.UnmarshalPrivateKey(kb)
}

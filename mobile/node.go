package mobile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"encoding/json"
	"encoding/base64"
	"io/ioutil"

	tcore "github.com/textileio/textile-go/core"
	trepo "github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/wallet"
	"github.com/textileio/textile-go/repo/db"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	oldcmds "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
	lockfile "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo/lock"
	utilmain "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/cmd/ipfs/util"
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
	Thumbs []string `json:"thumbs"`
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
	cfg, err := repo.Config()
	if err != nil {
		return nil, err
	}
	cfg.Swarm.DisableNatPortMap = true
	cfg.Addresses.Swarm = append(cfg.Addresses.Swarm, "/ip4/0.0.0.0/tcp/9005/ws")
	cfg.Addresses.Swarm = append(cfg.Addresses.Swarm, "/ip6/::/tcp/9005/ws")

	// Start assembling node config
	ncfg := &core.BuildCfg{
		Repo:      repo,
		Permanent: true, // It is temporary way to signify that node is permanent
		Online:    true,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
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
	fmt.Println("Starting node...")
	fmt.Println("Repo directory: ", n.config.RepoPath)

	cctx, cancel := context.WithCancel(context.Background())
	n.cancel = cancel

	ctx := oldcmds.Context{}
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
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(tcore.Node.RepoPath, "datastore", "LOCK")
	os.Remove(dsLockFile)
	tcore.Node.IpfsNode.Close()
	return nil
}

func (n *Node) PinPhoto(path string, thumb string) (string, error) {
	// read file from disk
	r, err := os.Open(path)
	if err != nil {
	 	return "", err
	}
	defer r.Close()

	t, err := os.Open(thumb)
	if err != nil {
		return "", err
	}
	defer t.Close()

	fname := filepath.Base(path)

	// pin
	ldn, err := wallet.PinPhoto(r, fname, t, n.node.IpfsNode, n.config.ApiHost)
	if err != nil {
		return "", err
	}
	hash := ldn.Cid().Hash().B58String()

	//byt, err := ioutil.ReadAll(t)
	//if err != nil {
	//	return "", err
	//}
	//bs64 := base64.StdEncoding.EncodeToString(byt)
	bs64 := "HELLO IMAGE"

	// index
	n.node.Datastore.Photos().Put(hash, bs64, time.Now())

	return hash, nil
}

func (n *Node) GetPhotos(offsetId string, limit int) (string, error) {
	// query for available hashes
	list := n.node.Datastore.Photos().GetPhotos(offsetId, limit)

	// return json list of hashes
	res := &PhotoList{
		Hashes: make([]string, len(list)),
		Thumbs: make([]string, len(list)),
	}
	for i := range list {
		res.Hashes[i] = list[i].Cid
		res.Thumbs[i] = list[i].Thumb
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
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	bs64 := base64.StdEncoding.EncodeToString(b)

	return bs64, nil
}

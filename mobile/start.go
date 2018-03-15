package mobile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tcore "github.com/textileio/textile-go/core"
	trepo "github.com/textileio/textile-go/repo"

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

func NewTextile(repoPath string) *Node {

	nodeconfig := MobileConfig{
		RepoPath: repoPath,
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

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		fmt.Errorf("setting file descriptor limit: %s", err)
	}

	// we may be running in an uninitialized state.
	if !fsrepo.IsInitialized(config.RepoPath) {
		err := trepo.DoInit(os.Stdout, config.RepoPath, nil)
		if err != nil {
			return nil, err
		}
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
		RepoPath: config.RepoPath,
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

	defer func() {
		// We wait for the node to close first, as the node has children
		// that it will wait for before closing, such as the API server.
		nd.Close()

		select {
		case <-cctx.Done():
			fmt.Println("Gracefully shut down node")
		default:
		}
	}()

	n.node.Context = ctx
	n.node.IpfsNode = nd

	// construct http gateway - if it is set in the config
	var gwErrc <-chan error
	gwErrc, err = tcore.ServeHTTPGateway(&ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Node is ready\n")
	// collect long-running errors and block for shutdown
	for err := range tcore.Merge(gwErrc) {
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (n *Node) Stop() error {
	repoLockFile := filepath.Join(tcore.Node.RepoPath, lockfile.LockFile)
	os.Remove(repoLockFile)
	tcore.Node.IpfsNode.Close()
	return nil
}

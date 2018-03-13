package mobile

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"path/filepath"

	tcore "github.com/textileio/textile-go/core"
	trepo "github.com/textileio/textile-go/repo"

	oldcmds "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/corehttp"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
	lockfile "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo/lock"
	utilmain "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/cmd/ipfs/util"

	"gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"errors"
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

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		fmt.Errorf("setting file descriptor limit: %s", err)
	}

	// we may be running in an uninitialized state.
	if !fsrepo.IsInitialized(config.RepoPath) {
		err := trepo.InitWithDefaults(os.Stdout, config.RepoPath)
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
		//TODO(Kubuxu): refactor Online vs Offline by adding Permanent vs Ephemeral
	}

	// Textile node setup
	tcore.Node = &tcore.TextileNode{
		RepoPath: config.RepoPath,
	}

	if len(cfg.Addresses.Gateway) <= 0 {
		return nil, errors.New("no gateway addresses configured")
	}

	return &Node{config: config, node: tcore.Node, ipfsConfig: ncfg}, nil
}

func (n *Node) Start() error {
	fmt.Println("Initializing node...")
	fmt.Println("Repo directory: ", n.config.RepoPath)

	cctx, cancel := context.WithCancel(context.Background())
	n.cancel = cancel

	ctx := oldcmds.Context{}
	nd, err := core.NewNode(cctx, n.ipfsConfig)
	if err != nil {
		return err
	}
	nd.SetLocal(false)

	printSwarmAddrs(nd)

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
	//var gwErrc <-chan error
	//gwErrc, err = serveHTTPGateway(&ctx)
	//if err != nil {
	//	fmt.Println(err)
	//}

	fmt.Printf("Node is ready\n")
	// collect long-running errors and block for shutdown
	//// TODO(cryptix): our fuse currently doesnt follow this pattern for graceful shutdown
	//for err := range merge(gwErrc) {
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}

	return nil
}

func (n *Node) Stop() error {
	repoLockFile := filepath.Join(tcore.Node.RepoPath, lockfile.LockFile)
	os.Remove(repoLockFile)
	tcore.Node.IpfsNode.Close()
	return nil
}

// printSwarmAddrs prints the addresses of the host
func printSwarmAddrs(node *core.IpfsNode) {
	if !node.OnlineMode() {
		fmt.Println("Swarm not listening, running in offline mode.")
		return
	}

	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		fmt.Errorf("failed to read listening addresses: %s", err)
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Sort(sort.StringSlice(lisAddrs))
	for _, addr := range lisAddrs {
		fmt.Printf("Swarm listening on %s\n", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Sort(sort.StringSlice(addrs))
	for _, addr := range addrs {
		fmt.Printf("Swarm announcing %s\n", addr)
	}

}

// serveHTTPGateway collects options, creates listener, prints status message and starts serving requests
func serveHTTPGateway(cctx *oldcmds.Context) (<-chan error, error) {
	cfg, err := cctx.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: GetConfig() failed: %s", err)
	}

	gatewayMaddr, err := ma.NewMultiaddr(cfg.Addresses.Gateway)
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: invalid gateway address: %q (err: %s)", cfg.Addresses.Gateway, err)
	}

	gwLis, err := manet.Listen(gatewayMaddr)
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
	}
	// we might have listened to /tcp/0 - lets see what we are listing on
	gatewayMaddr = gwLis.Multiaddr()

	fmt.Printf("Gateway (readonly) server listening on %s\n", gatewayMaddr)

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("gateway"),
		corehttp.CheckVersionOption(),
		corehttp.CommandsROOption(*cctx),
		corehttp.VersionOption(),
		corehttp.IPNSHostnameOption(),
		corehttp.GatewayOption(false, "/ipfs", "/ipns"),
	}

	if len(cfg.Gateway.RootRedirect) > 0 {
		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
	}

	node, err := cctx.ConstructNode()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: ConstructNode() failed: %s", err)
	}

	errc := make(chan error)
	go func() {
		errc <- corehttp.Serve(node, gwLis.NetListener(), opts...)
		close(errc)
	}()
	return errc, nil
}

// merge does fan-in of multiple read-only error channels
// taken from http://blog.golang.org/pipelines
func merge(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	for _, c := range cs {
		if c != nil {
			wg.Add(1)
			go output(c)
		}
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

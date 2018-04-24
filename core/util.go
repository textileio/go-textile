package core

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"

	oldcmds "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/commands"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/corehttp"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/corerepo"

	"gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"strconv"
	"strings"
)

// PrintSwarmAddrs prints the addresses of the host
func printSwarmAddrs(node *core.IpfsNode) error {
	if !node.OnlineMode() {
		log.Info("swarm not listening, running in offline mode")
		return nil
	}

	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		return err
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Sort(sort.StringSlice(lisAddrs))
	for _, addr := range lisAddrs {
		log.Infof("swarm listening on %s\n", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Sort(sort.StringSlice(addrs))
	for _, addr := range addrs {
		log.Infof("swarm announcing %s\n", addr)
	}

	return nil
}

// serveHTTPGateway collects options, creates listener, prints status message and starts serving requests
func serveHTTPGateway(cctx *oldcmds.Context) (<-chan error, error) {
	cfg, err := cctx.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("ServeHTTPGateway: GetConfig() failed: %s", err)
	}

	gatewayMaddr, err := ma.NewMultiaddr(cfg.Addresses.Gateway)
	if err != nil {
		return nil, fmt.Errorf("ServeHTTPGateway: invalid gateway address: %q (err: %s)", cfg.Addresses.Gateway, err)
	}

	gwLis, err := manet.Listen(gatewayMaddr)
	if err != nil {
		return nil, fmt.Errorf("ServeHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
	}
	// we might have listened to /tcp/0 - lets see what we are listing on
	gatewayMaddr = gwLis.Multiaddr()

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
		return nil, fmt.Errorf("ServeHTTPGateway: ConstructNode() failed: %s", err)
	}

	errc := make(chan error)
	go func() {
		errc <- corehttp.Serve(node, gwLis.NetListener(), opts...)
		close(errc)
	}()
	log.Infof("gateway (readonly) server listening on %s\n", gatewayMaddr)

	return errc, nil
}

func ServeHTTPGatewayProxy(node *TextileNode) (<-chan error, error) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := node.GetFile(r.URL.Path)
		if err != nil {
			log.Errorf("error decrypting path %s: %s", r.URL.Path, err)
			w.WriteHeader(400)
			return
		}
		w.Write(b)
	})

	// get config and set proxy address to raw gateway address plus one,
	// so a gateway on 8182 means the proxy will run on 9182
	cfg, err := node.Context.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("ServeHTTPGatewayProxy: GetConfig() failed: %s", err)
	}
	tmp := strings.Split(cfg.Addresses.Gateway, "/")
	gaddrs := tmp[len(tmp)-1]
	gaddr, err := strconv.ParseInt(gaddrs, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ServeHTTPGatewayProxy: get address failed: %s", err)
	}
	addr := gaddr + 1000

	errc := make(chan error)
	go func() {
		errc <- http.ListenAndServe(fmt.Sprintf(":%v", addr), nil)
		close(errc)
	}()
	log.Infof("decrypting gateway (readonly) server listening at http://127.0.0.1:%v\n", addr)

	return errc, nil
}

func runGC(ctx context.Context, node *core.IpfsNode) (<-chan error, error) {
	errc := make(chan error)
	go func() {
		errc <- corerepo.PeriodicGC(ctx, node)
		close(errc)
	}()
	log.Info("auto garbage collection started")

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

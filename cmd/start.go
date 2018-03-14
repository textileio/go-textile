package main

import (
	_ "expvar"
	"fmt"
	_ "net/http/pprof"
	"os"

	trepo "github.com/textileio/textile-go/repo"
	tcore "github.com/textileio/textile-go/core"

	utilmain "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"

	"gx/ipfs/QmabLouZTZwhfALuBcssPvkzhbYGMb4394huT7HY4LQ6d3/go-ipfs-cmds"
	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
)

const (
	repoDirKwd = "dir"
)

var startCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Run a network-connected Textile node.",
		ShortDescription: `
'textile start' runs a persistent textile daemon.
`,
		LongDescription: `
The daemon will start listening on ports on the network, which are
documented in this ipfs config.

Shutdown

To shutdown the daemon, send a SIGINT signal to it (e.g. by pressing 'Ctrl-C')
or send a SIGTERM signal to it (e.g. with 'kill'). It may take a while for the
daemon to shutdown gracefully, but it can be killed forcibly by sending a
second signal.
`,
	},

	Options: []cmdkit.Option{
		cmdkit.StringOption(repoDirKwd, "Repo directory.").WithDefault("~/.ipfs"),
	},
	Subcommands: map[string]*cmds.Command{},
	Run:         daemonFunc,
}

func daemonFunc(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) {
	// let the user know we're going.
	fmt.Printf("Initializing daemon...\n")

	// raise file descriptor limit
	if err := utilmain.ManageFdLimit(); err != nil {
		log.Errorf("setting file descriptor limit: %s", err)
	}

	cctx := env.(*oldcmds.Context)

	go func() {
		<-req.Context.Done()
		fmt.Println("Received interrupt signal, shutting down...")
		fmt.Println("(Hit ctrl-c again to force-shutdown the daemon.)")
	}()

	// we may be running in an uninitialized state.
	repoDir, _ := req.Options[repoDirKwd].(string)
	if !fsrepo.IsInitialized(repoDir) {
		err := trepo.DoInit(os.Stdout, repoDir, nil)
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal)
			return
		}
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(repoDir)
	if err != nil {
		re.SetError(err, cmdkit.ErrNormal)
		return
	}

	cfg, err := cctx.GetConfig()
	if err != nil {
		re.SetError(err, cmdkit.ErrNormal)
		return
	}

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

	node, err := core.NewNode(req.Context, ncfg)
	if err != nil {
		log.Error("error from node construction: ", err)
		re.SetError(err, cmdkit.ErrNormal)
		return
	}
	node.SetLocal(false)

	if err := tcore.PrintSwarmAddrs(node); err != nil {
		log.Errorf("failed to read listening addresses: %s", err)
	}

	defer func() {
		// We wait for the node to close first, as the node has children
		// that it will wait for before closing, such as the API server.
		node.Close()

		select {
		case <-req.Context.Done():
			log.Info("Gracefully shut down daemon")
		default:
		}
	}()

	cctx.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	// construct api endpoint - every time
	apiErrc, err := tcore.ServeHTTPApi(cctx)
	if err != nil {
		re.SetError(err, cmdkit.ErrNormal)
		return
	}

	// construct http gateway - if it is set in the config
	var gwErrc <-chan error
	if len(cfg.Addresses.Gateway) > 0 {
		var err error
		gwErrc, err = tcore.ServeHTTPGateway(cctx)
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal)
			return
		}
	}

	fmt.Printf("Daemon is ready\n")
	// collect long-running errors and block for shutdown
	for err := range tcore.Merge(apiErrc, gwErrc) {
		if err != nil {
			log.Error(err)
			re.SetError(err, cmdkit.ErrNormal)
		}
	}
}

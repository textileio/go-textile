package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	//tcore "github.com/textileio/textile-go/core"
	trepo "github.com/textileio/textile-go/repo"

	utilmain "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/cmd/ipfs/util"
	oldcmds "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/commands"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo"
	lockfile "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo/lock"

	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
	"gx/ipfs/QmfAkMSt9Fwzk48QDJecPcwCUjnf2uG7MLnmCGTp4C6ouL/go-ipfs-cmds"
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

	Options:     []cmdkit.Option{},
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

	// shutdown is not clean here yet, so we have to hackily remove
	// the lockfile that should have been removed on shutdown
	// before we start up again
	repoLockFile := filepath.Join(cctx.ConfigRoot, lockfile.LockFile)
	os.Remove(repoLockFile)
	dsLockFile := filepath.Join(cctx.ConfigRoot, "datastore", "LOCK")
	os.Remove(dsLockFile)

	// we may be running in an uninitialized state.
	if !fsrepo.IsInitialized(cctx.ConfigRoot) {
		err := trepo.DoInit(os.Stdout, cctx.ConfigRoot, false, nil)
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal)
			return
		}
	}

	// acquire the repo lock _before_ constructing a node. we need to make
	// sure we are permitted to access the resources (datastore, etc.)
	repo, err := fsrepo.Open(cctx.ConfigRoot)
	if err != nil {
		re.SetError(err, cmdkit.ErrNormal)
		return
	}

	//cfg, err := cctx.GetConfig()
	//if err != nil {
	//	re.SetError(err, cmdkit.ErrNormal)
	//	return
	//}

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
		Routing: core.DHTOption,
	}

	node, err := core.NewNode(req.Context, ncfg)
	if err != nil {
		log.Error("error from node construction: ", err)
		re.SetError(err, cmdkit.ErrNormal)
		return
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

	node.SetLocal(false)

	//if err := tcore.PrintSwarmAddrs(node); err != nil {
	//	log.Errorf("failed to read listening addresses: %s", err)
	//}

	cctx.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	// construct api endpoint - every time
	//apiErrc, err := tcore.ServeHTTPApi(cctx)
	//if err != nil {
	//	re.SetError(err, cmdkit.ErrNormal)
	//	return
	//}

	// construct http gateway - if it is set in the config
	//var gwErrc <-chan error
	//if len(cfg.Addresses.Gateway) > 0 {
	//	var err error
	//	gwErrc, err = tcore.ServeHTTPGateway(cctx)
	//	if err != nil {
	//		re.SetError(err, cmdkit.ErrNormal)
	//		return
	//	}
	//}

	// tmp setup subscription for testing
	go func() {
		sub, err := node.Floodsub.Subscribe("textile")
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal)
			return
		}
		defer sub.Cancel()

		for {
			select {
			default:
				msg, err := sub.Next(req.Context)
				if err == io.EOF || err == context.Canceled {
					return
				} else if err != nil {
					re.SetError(err, cmdkit.ErrNormal)
					return
				}
				fmt.Printf("Received message: %s\n", msg)
			case <-req.Context.Done():
				return
			}
		}
	}()

	fmt.Printf("Daemon is ready\n")
	// collect long-running errors and block for shutdown
	//for err := range tcore.Merge(apiErrc, gwErrc) {
	//	if err != nil {
	//		log.Error(err)
	//		re.SetError(err, cmdkit.ErrNormal)
	//	}
	//}
}

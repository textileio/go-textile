package core

import (
	"context"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

// createIPFS creates an IPFS node
func (t *Textile) createIPFS(plugins *loader.PluginLoader, online bool) error {
	rep, err := fsrepo.Open(t.repoPath)
	if err != nil {
		return err
	}

	routing := libp2p.DHTClientOption
	if t.Server() {
		routing = libp2p.DHTOption
	}

	cctx, _ := context.WithCancel(context.Background())
	nd, err := core.NewNode(cctx, &core.BuildCfg{
		Repo:      rep,
		Permanent: true, // temporary way to signify that node is permanent
		Online:    online,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
			"mplex":  true,
		},
		Routing: routing,
	})
	if err != nil {
		return err
	}
	nd.IsDaemon = true

	if t.node != nil {
		err = t.node.Close()
		if err != nil {
			return err
		}
	}
	t.node = nd

	return nil
}

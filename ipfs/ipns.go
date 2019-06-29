package ipfs

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	nsopts "github.com/ipfs/interface-go-ipfs-core/options/namesys"
	path "github.com/ipfs/interface-go-ipfs-core/path"
	peer "github.com/libp2p/go-libp2p-core/peer"
	record "github.com/libp2p/go-libp2p-record"
)

// PublishIPNS publishes a content id to ipns
func PublishIPNS(node *core.IpfsNode, id string, key string, timeout time.Duration) (iface.IpnsEntry, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	if key == "" {
		key = "self" // default value in ipns module
	}

	opts := []options.NamePublishOption{
		options.Name.Key(key),
	}

	ctx, cancel := context.WithTimeout(node.Context(), timeout)
	defer cancel()

	return api.Name().Publish(ctx, path.New(id), opts...)
}

// ResolveIPNS resolves an ipns path to an ipfs path
func ResolveIPNS(node *core.IpfsNode, name peer.ID, timeout time.Duration) (path.Path, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("/ipns/%s", name.Pretty())

	opts := []options.NameResolveOption{
		options.Name.ResolveOption(nsopts.Depth(1)),
		options.Name.ResolveOption(nsopts.DhtRecordCount(4)),
		options.Name.ResolveOption(nsopts.DhtTimeout(timeout)),
	}

	ctx, cancel := context.WithTimeout(node.Context(), timeout)
	defer cancel()

	return api.Name().Resolve(ctx, key, opts...)
}

// IpnsSubs shows current name subscriptions
func IpnsSubs(node *core.IpfsNode) ([]string, error) {
	if node.PSRouter == nil {
		return nil, fmt.Errorf("IPNS pubsub subsystem is not enabled")
	}
	var paths []string
	for _, key := range node.PSRouter.GetSubscriptions() {
		ns, k, err := record.SplitKey(key)
		if err != nil || ns != "ipns" {
			// not necessarily an error.
			continue
		}
		pid, err := peer.IDFromString(k)
		if err != nil {
			log.Errorf("ipns key not a valid peer ID: %s", err)
			continue
		}
		paths = append(paths, "/ipns/"+peer.IDB58Encode(pid))
	}
	return paths, nil
}

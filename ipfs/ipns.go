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
	peer "github.com/libp2p/go-libp2p-peer"
	record "github.com/libp2p/go-libp2p-record"
)

const ipnsTimeout = time.Second * 30

// PublishIPNS publishes a content id to ipns
func PublishIPNS(node *core.IpfsNode, id string) (iface.IpnsEntry, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	opts := []options.NamePublishOption{
		options.Name.AllowOffline(true),
		options.Name.ValidTime(time.Hour * 24),
		options.Name.TTL(time.Hour),
	}

	pth, err := iface.ParsePath(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), ipnsTimeout)
	defer cancel()

	return api.Name().Publish(ctx, pth, opts...)
}

// ResolveIPNS resolves an ipns path to an ipfs path
func ResolveIPNS(node *core.IpfsNode, name peer.ID) (iface.Path, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("/ipns/%s", name.Pretty())

	opts := []options.NameResolveOption{
		options.Name.ResolveOption(nsopts.Depth(1)),
		options.Name.ResolveOption(nsopts.DhtRecordCount(4)),
		options.Name.ResolveOption(nsopts.DhtTimeout(ipnsTimeout)),
	}

	ctx, cancel := context.WithTimeout(node.Context(), ipnsTimeout)
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

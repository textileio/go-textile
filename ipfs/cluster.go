package ipfs

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	icid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"github.com/ipfs/go-ipfs/pin"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/ipfs/ipfs-cluster/api"
	corepeer "github.com/libp2p/go-libp2p-core/peer"
	rpc "github.com/libp2p/go-libp2p-gorpc"
	peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
)

type ClusterConnector struct {
	node  *core.IpfsNode
	api   iface.CoreAPI
	peers func(ctx context.Context) []*api.ID
}

func NewClusterConnector(node *core.IpfsNode, peers func(ctx context.Context) []*api.ID) (*ClusterConnector, error) {
	capi, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}
	return &ClusterConnector{
		node:  node,
		api:   capi,
		peers: peers,
	}, nil
}

func (c *ClusterConnector) ID(context.Context) (*api.IPFSID, error) {
	var addrs []api.Multiaddr
	for _, addr := range c.node.PeerHost.Addrs() {
		addrs = append(addrs, api.Multiaddr{Multiaddr: addr})
	}
	return &api.IPFSID{
		ID:        c.node.Identity,
		Addresses: addrs,
	}, nil
}

func (c *ClusterConnector) SetClient(client *rpc.Client) {
	// noop
}

func (c *ClusterConnector) Shutdown(ctx context.Context) error {
	// noop
	return nil
}

// @todo handle maxDepth
func (c *ClusterConnector) Pin(ctx context.Context, cid icid.Cid, maxDepth int) error {
	return c.api.Pin().Add(ctx, path.New(cid.String()))
}

func (c *ClusterConnector) Unpin(ctx context.Context, cid icid.Cid) error {
	return c.api.Pin().Rm(ctx, path.New(cid.String()))
}

func (c *ClusterConnector) PinLsCid(ctx context.Context, cid icid.Cid) (api.IPFSPinStatus, error) {
	pins, err := c.node.Pinning.CheckIfPinned(cid)
	if err != nil {
		return api.IPFSPinStatusError, err
	}
	if len(pins) == 0 {
		return api.IPFSPinStatusError, fmt.Errorf("invalid pin check result")
	}
	return c.pinModeToStatus(pins[0].Mode), nil
}

func (c *ClusterConnector) PinLs(ctx context.Context, typeFilter string) (map[string]api.IPFSPinStatus, error) {
	pins, err := c.api.Pin().Ls(ctx, c.pinFilterToOption(typeFilter))
	if err != nil {
		return nil, err
	}
	statusMap := make(map[string]api.IPFSPinStatus)
	for _, p := range pins {
		mode, ok := pin.StringToMode(p.Type())
		if !ok {
			continue
		}
		statusMap[p.Path().String()] = c.pinModeToStatus(mode)
	}
	return statusMap, nil
}

func (c *ClusterConnector) ConnectSwarms(ctx context.Context) error {
	for _, p := range c.peers(ctx) {
		log.Debugf("cluster dialing %s", p.ID.Pretty())
		var addrs []ma.Multiaddr
		for _, addr := range p.Addresses {
			addrs = append(addrs, addr.Multiaddr)
		}
		err := c.api.Swarm().Connect(ctx, corepeer.AddrInfo{
			ID:    p.ID,
			Addrs: addrs,
		})
		if err != nil {
			return err
		}
		log.Debugf("cluster connected to %s")
	}
	return nil
}

func (c *ClusterConnector) SwarmPeers(ctx context.Context) ([]peer.ID, error) {
	conns, err := c.api.Swarm().Peers(ctx)
	if err != nil {
		return nil, err
	}
	var peers []peer.ID
	for _, c := range conns {
		peers = append(peers, c.ID())
	}
	return peers, nil
}

func (c *ClusterConnector) ConfigKey(keypath string) (interface{}, error) {
	return c.node.Repo.GetConfigKey(keypath)
}

func (c *ClusterConnector) RepoStat(ctx context.Context) (*api.IPFSRepoStat, error) {
	stat, err := corerepo.RepoStat(ctx, c.node)
	if err != nil {
		return nil, err
	}
	return &api.IPFSRepoStat{
		RepoSize:   stat.RepoSize,
		StorageMax: stat.StorageMax,
	}, nil
}

func (c *ClusterConnector) Resolve(ctx context.Context, pth string) (icid.Cid, error) {
	res, err := c.api.ResolvePath(ctx, path.New(pth))
	if err != nil {
		return icid.Undef, err
	}
	return res.Cid(), nil
}

func (c *ClusterConnector) BlockPut(ctx context.Context, b *api.NodeWithMeta) error {
	_, err := c.api.Block().Put(ctx, bytes.NewReader(b.Data), options.Block.Format(b.Format))
	return err
}

func (c *ClusterConnector) BlockGet(ctx context.Context, cid icid.Cid) ([]byte, error) {
	r, err := c.api.Block().Get(ctx, path.New(cid.String()))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func (c *ClusterConnector) pinModeToStatus(mode pin.Mode) api.IPFSPinStatus {
	switch mode {
	case pin.Recursive:
		return api.IPFSPinStatusRecursive
	case pin.Direct:
		return api.IPFSPinStatusDirect
	case pin.Indirect:
		return api.IPFSPinStatusIndirect
	case pin.Internal:
		return api.IPFSPinStatusDirect
	case pin.NotPinned:
		return api.IPFSPinStatusUnpinned
	default:
		return api.IPFSPinStatusError
	}
}

func (c *ClusterConnector) pinFilterToOption(typeFilter string) options.PinLsOption {
	return func(settings *options.PinLsSettings) error {
		settings.Type = typeFilter
		return nil
	}
}

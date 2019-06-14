package ipfs

import (
	"context"

	icid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/ipfs-cluster/api"
	rpc "github.com/libp2p/go-libp2p-gorpc"
	peer "github.com/libp2p/go-libp2p-peer"
)

func NewClusterConnector(node *core.IpfsNode) *ClusterConnector {
	return &ClusterConnector{node: node}
}

type ClusterConnector struct {
	node *core.IpfsNode
}

func (c *ClusterConnector) ID(context.Context) (*api.IPFSID, error) {
	panic("implement me")
}

func (c *ClusterConnector) SetClient(*rpc.Client) {
	panic("implement me")
}

func (c *ClusterConnector) Shutdown(context.Context) error {
	panic("implement me")
}

func (c *ClusterConnector) Pin(context.Context, icid.Cid, int) error {
	panic("implement me")
}

func (c *ClusterConnector) Unpin(context.Context, icid.Cid) error {
	panic("implement me")
}

func (c *ClusterConnector) PinLsCid(context.Context, icid.Cid) (api.IPFSPinStatus, error) {
	panic("implement me")
}

func (c *ClusterConnector) PinLs(ctx context.Context, typeFilter string) (map[string]api.IPFSPinStatus, error) {
	panic("implement me")
}

func (c *ClusterConnector) ConnectSwarms(context.Context) error {
	panic("implement me")
}

func (c *ClusterConnector) SwarmPeers(context.Context) ([]peer.ID, error) {
	panic("implement me")
}

func (c *ClusterConnector) ConfigKey(keypath string) (interface{}, error) {
	panic("implement me")
}

func (c *ClusterConnector) RepoStat(context.Context) (*api.IPFSRepoStat, error) {
	panic("implement me")
}

func (c *ClusterConnector) Resolve(context.Context, string) (icid.Cid, error) {
	panic("implement me")
}

func (c *ClusterConnector) BlockPut(context.Context, *api.NodeWithMeta) error {
	panic("implement me")
}

func (c *ClusterConnector) BlockGet(context.Context, icid.Cid) ([]byte, error) {
	panic("implement me")
}

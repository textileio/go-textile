package ipfs

import (
	"context"
	"fmt"
	"sort"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	inet "github.com/libp2p/go-libp2p-core/network"
)

// SwarmConnect opens a direct connection to a list of peer multi addresses
func SwarmConnect(node *core.IpfsNode, addrs []string) ([]string, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	pis, err := peersWithAddresses(addrs)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), ConnectTimeout)
	defer cancel()

	output := make([]string, len(pis))
	for i, pi := range pis {
		output[i] = "connect " + pi.ID.Pretty()

		err := api.Swarm().Connect(ctx, pi)
		if err != nil {
			return nil, fmt.Errorf("%s failure: %s", output[i], err)
		}
		output[i] += " success"
	}

	return output, nil
}

type streamInfo struct {
	Protocol string `json:"protocol"`
}

type connInfo struct {
	Addr      string         `json:"addr"`
	Peer      string         `json:"peer"`
	Latency   string         `json:"latency,omitempty"`
	Muxer     string         `json:"muxer,omitempty"`
	Direction inet.Direction `json:"direction,omitempty"`
	Streams   []streamInfo   `json:"streams,omitempty"`
}

func (ci *connInfo) Less(i, j int) bool {
	return ci.Streams[i].Protocol < ci.Streams[j].Protocol
}

func (ci *connInfo) Len() int {
	return len(ci.Streams)
}

func (ci *connInfo) Swap(i, j int) {
	ci.Streams[i], ci.Streams[j] = ci.Streams[j], ci.Streams[i]
}

type ConnInfos struct {
	Peers []connInfo `json:"peers"`
}

func (ci ConnInfos) Less(i, j int) bool {
	return ci.Peers[i].Addr < ci.Peers[j].Addr
}

func (ci ConnInfos) Len() int {
	return len(ci.Peers)
}

func (ci ConnInfos) Swap(i, j int) {
	ci.Peers[i], ci.Peers[j] = ci.Peers[j], ci.Peers[i]
}

// SwarmPeers lists the set of peers this node is connected to
func SwarmPeers(node *core.IpfsNode, verbose bool, latency bool, streams bool, direction bool) (*ConnInfos, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), ConnectTimeout)
	defer cancel()

	conns, err := api.Swarm().Peers(ctx)
	if err != nil {
		return nil, err
	}

	var out ConnInfos
	for _, c := range conns {
		ci := connInfo{
			Addr: c.Address().String(),
			Peer: c.ID().Pretty(),
		}

		if verbose || direction {
			// set direction
			ci.Direction = c.Direction()
		}

		if verbose || latency {
			lat, err := c.Latency()
			if err != nil {
				return nil, err
			}

			if lat == 0 {
				ci.Latency = "n/a"
			} else {
				ci.Latency = lat.String()
			}
		}
		if verbose || streams {
			strs, err := c.Streams()
			if err != nil {
				return nil, err
			}

			for _, s := range strs {
				ci.Streams = append(ci.Streams, streamInfo{Protocol: string(s)})
			}
		}
		sort.Sort(&ci)
		out.Peers = append(out.Peers, ci)
	}

	sort.Sort(&out)
	return &out, nil
}

// SwarmConnected returns whether or not the node has the peer in its current swarm
func SwarmConnected(node *core.IpfsNode, peerId string) (bool, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), ConnectTimeout)
	defer cancel()

	conns, err := api.Swarm().Peers(ctx)
	if err != nil {
		return false, err
	}

	for _, c := range conns {
		if c.ID().Pretty() == peerId {
			return true, nil
		}
	}

	return false, nil
}

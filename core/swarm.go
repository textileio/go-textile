package core

import (
	"errors"
	"fmt"

	"gx/ipfs/QmVHhT8NxtApPTndiZPe4JNGNUxGWtJe3ebyxtRz4HnbEp/go-libp2p-swarm"
	inet "gx/ipfs/QmXuRkCR7BNQa9uqfpTiFWsTQLzmTWYg91Ja1w95gnqb6u/go-libp2p-net"

	"github.com/textileio/textile-go/ipfs"
)

// ConnectPeer connect to another ipfs peer (i.e., ipfs swarm connect)
func (t *Textile) ConnectPeer(addrs []string) ([]string, error) {
	swrm, ok := t.node.PeerHost.Network().(*swarm.Swarm)
	if !ok {
		return nil, errors.New("peerhost network was not swarm")
	}

	pis, err := ipfs.PeersWithAddresses(addrs)
	if err != nil {
		return nil, err
	}

	output := make([]string, len(pis))
	for i, pi := range pis {
		swrm.Backoff().Clear(pi.ID)

		output[i] = "connect " + pi.ID.Pretty()

		err := t.node.PeerHost.Connect(t.node.Context(), pi)
		if err != nil {
			return nil, fmt.Errorf("%s failure: %s", output[i], err)
		}
		output[i] += " success"
	}
	return output, nil
}

func (t *Textile) Peers() ([]inet.Conn, error) {
	return t.node.PeerHost.Network().Conns(), nil
}

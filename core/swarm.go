package core

import (
	"errors"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
	libp2pn "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	"gx/ipfs/QmemVjhp1UuWPQqrWSvPcaqH3QJRMjMqNm4T2RULMkDDQe/go-libp2p-swarm"
)

// ConnectPeer connect to another ipfs peer (i.e., ipfs swarm connect)
func (t *Textile) ConnectPeer(addrs []string) ([]string, error) {
	if !t.Online() {
		return nil, ErrOffline
	}
	swrm, ok := t.ipfs.PeerHost.Network().(*swarm.Swarm)
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

		err := t.ipfs.PeerHost.Connect(t.ipfs.Context(), pi)
		if err != nil {
			return nil, fmt.Errorf("%s failure: %s", output[i], err)
		}
		output[i] += " success"
	}
	return output, nil
}

func (t *Textile) Peers() ([]libp2pn.Conn, error) {
	if !t.Online() {
		return nil, ErrOffline
	}
	return t.ipfs.PeerHost.Network().Conns(), nil
}

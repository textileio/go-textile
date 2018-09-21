package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
	"gx/ipfs/QmSwZMWwFZSUpe5muU2xgTUwppH24KfMwdPXiwbEp2c6G5/go-libp2p-swarm"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	libp2pn "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	"time"
)

// ConnectPeer connect to another ipfs peer (i.e., ipfs swarm connect)
func (t *Textile) ConnectPeer(addrs []string) ([]string, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}
	snet, ok := t.ipfs.PeerHost.Network().(*swarm.Network)
	if !ok {
		return nil, errors.New("peerhost network was not swarm")
	}

	swrm := snet.Swarm()
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

// PingPeer pings a peer num times, returning the result to out chan
func (t *Textile) PingPeer(addrs string, num int, out chan string) error {
	if !t.started {
		return ErrStopped
	}
	if !t.IsOnline() {
		return ErrOffline
	}
	addr, pid, err := ipfs.ParsePeerParam(addrs)
	if addr != nil {
		t.ipfs.Peerstore.AddAddr(pid, addr, pstore.TempAddrTTL) // temporary
	}

	if len(t.ipfs.Peerstore.Addrs(pid)) == 0 {
		// Make sure we can find the node in question
		log.Debugf("looking up peer: %s", pid.Pretty())

		ctx, cancel := context.WithTimeout(t.ipfs.Context(), pingTimeout)
		defer cancel()
		p, err := t.ipfs.Routing.FindPeer(ctx, pid)
		if err != nil {
			err = fmt.Errorf("peer lookup error: %s", err)
			log.Errorf(err.Error())
			return err
		}
		t.ipfs.Peerstore.AddAddrs(p.ID, p.Addrs, pstore.TempAddrTTL)
	}

	ctx, cancel := context.WithTimeout(t.ipfs.Context(), pingTimeout*time.Duration(num))
	defer cancel()
	pings, err := t.ipfs.Ping.Ping(ctx, pid)
	if err != nil {
		log.Errorf("error pinging peer %s: %s", pid.Pretty(), err)
		return err
	}

	var done bool
	var total time.Duration
	for i := 0; i < num && !done; i++ {
		select {
		case <-ctx.Done():
			done = true
			close(out)
			break
		case t, ok := <-pings:
			if !ok {
				done = true
				close(out)
				break
			}
			total += t
			msg := fmt.Sprintf("ping %s completed after %f seconds", pid.Pretty(), t.Seconds())
			select {
			case out <- msg:
			default:
			}
			log.Debug(msg)
			time.Sleep(time.Second)
		}
	}
	return nil
}

func (t *Textile) Peers() ([]libp2pn.Conn, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}
	return t.ipfs.PeerHost.Network().Conns(), nil
}

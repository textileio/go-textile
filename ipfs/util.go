package ipfs

import (
	"errors"
	"net"
	"sort"
	"strings"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmQkW9fnCsg9SLHdViiAh6qfBppodsPZVpU92dZLqYtEfs/go-ipfs-cmds"
	ma "gx/ipfs/QmTZBfrPJmjWsCvHEtX5FE6KimVJhsJg5sBbqEFYf4UZtL/go-multiaddr"
	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
	pstore "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore"
	"gx/ipfs/QmbeHtaBy9nZsW4cHRcvgVY4CnDhXudE2Dr6qDxS7yg9rX/go-libp2p-record"
	"gx/ipfs/Qmdf1djucJ1jX5RMF1bDbFg5ybZnupmSAeETQQ3ZV7z6dU/go-ipfs-addr"
)

// PrintSwarmAddrs prints the addresses of the host
func PrintSwarmAddrs(node *core.IpfsNode) error {
	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		return err
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Sort(sort.StringSlice(lisAddrs))
	for _, addr := range lisAddrs {
		log.Infof("swarm listening on %s", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Sort(sort.StringSlice(addrs))
	for _, addr := range addrs {
		log.Infof("swarm announcing %s", addr)
	}
	return nil
}

// PublicIPv4Addr uses the ipfs NAT traveral result to locate a (possibly) public ipv4 address.
// this method is used to inform cafe clients of the http api address
func PublicIPv4Addr(node *core.IpfsNode) (string, error) {
	var pub string
	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		return pub, err
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Sort(sort.StringSlice(lisAddrs))
	for _, addr := range lisAddrs {
		parts := strings.Split(addr, "/")
		if len(parts) < 3 {
			continue
		}
		if publicIPv4(net.ParseIP(parts[2])) {
			pub = parts[2]
			break
		}
	}
	if pub == "" {
		return pub, errors.New("no public ipv4 address found")
	}
	return pub, nil
}

// ShortenID returns the last 7 chars of a string
func ShortenID(id string) string {
	if len(id) < 7 {
		return id
	}
	return id[len(id)-7:]
}

// IpnsSubs shows current name subscriptions
func IpnsSubs(node *core.IpfsNode) ([]string, error) {
	if node.PSRouter == nil {
		return nil, errors.New("IPNS pubsub subsystem is not enabled")
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

// parseAddresses is a function that takes in a slice of string peer addresses
// (multiaddr + peerid) and returns slices of multiaddrs and peerids.
func parseAddresses(addrs []string) (iaddrs []ipfsaddr.IPFSAddr, err error) {
	iaddrs = make([]ipfsaddr.IPFSAddr, len(addrs))
	for i, saddr := range addrs {
		iaddrs[i], err = ipfsaddr.ParseString(saddr)
		if err != nil {
			return nil, cmds.ClientError("invalid peer address: " + err.Error())
		}
	}
	return
}

// peersWithAddresses is a function that takes in a slice of string peer addresses
// (multiaddr + peerid) and returns a slice of properly constructed peers
func peersWithAddresses(addrs []string) ([]pstore.PeerInfo, error) {
	iaddrs, err := parseAddresses(addrs)
	if err != nil {
		return nil, err
	}

	peers := make(map[peer.ID][]ma.Multiaddr, len(iaddrs))
	for _, iaddr := range iaddrs {
		id := iaddr.ID()
		current, ok := peers[id]
		if tpt := iaddr.Transport(); tpt != nil {
			peers[id] = append(current, tpt)
		} else if !ok {
			peers[id] = nil
		}
	}
	pis := make([]pstore.PeerInfo, 0, len(peers))
	for id, maddrs := range peers {
		pis = append(pis, pstore.PeerInfo{
			ID:    id,
			Addrs: maddrs,
		})
	}
	return pis, nil
}

// publicIPv4 returns true if the given ip is not reserved for a private address.
// of course, this only implies that it _might_ be public
// https://stackoverflow.com/a/41670589
func publicIPv4(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

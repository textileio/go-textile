package ipfs

import (
	"fmt"
	"net"
	"sort"
	"strings"

	ipfsaddr "github.com/ipfs/go-ipfs-addr"
	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs/core"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// ShortenID returns the last 7 chars of a string
func ShortenID(id string) string {
	if len(id) < 7 {
		return id
	}
	return id[len(id)-7:]
}

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

// GetPublicIPv4Addr uses the host addresses to return the public ipv4 address of the host machine, if available
func GetPublicIPv4Addr(node *core.IpfsNode) (string, error) {
	var ip string
	for _, addr := range node.PeerHost.Addrs() {
		parts := strings.Split(addr.String(), "/")
		if len(parts) < 3 {
			continue
		}
		parsed := net.ParseIP(parts[2])
		if parsed != nil && publicIPv4(parsed) {
			ip = parts[2]
			break
		}
	}
	if ip == "" {
		return ip, fmt.Errorf("no address was found")
	}
	return ip, nil
}

// GetLANIPv4Addr looks for a LAN IP in the host addresses (192.168.x.x)
func GetLANIPv4Addr(node *core.IpfsNode) (string, error) {
	var ip string
	for _, addr := range node.PeerHost.Addrs() {
		parts := strings.Split(addr.String(), "/")
		if len(parts) < 3 {
			continue
		}
		parsed := net.ParseIP(parts[2])
		if parsed != nil && lanIPv4(parsed) {
			ip = parts[2]
			break
		}
	}
	if ip == "" {
		return ip, fmt.Errorf("no address was found")
	}
	return ip, nil
}

// GetIPv6Addr returns the ipv6 address of the host machine, if available
func GetIPv6Addr(node *core.IpfsNode) (string, error) {
	var ip string
	node.PeerHost.Addrs()
	for _, addr := range node.PeerHost.Addrs() {
		parts := strings.Split(addr.String(), "/")
		if len(parts) < 3 || parts[2] == "::1" {
			continue
		}
		parsed := net.ParseIP(parts[2])
		if parsed != nil && parsed.To4() == nil {
			ip = parts[2]
			break
		}
	}
	if ip == "" {
		return ip, fmt.Errorf("no address was found")
	}
	return ip, nil
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
func peersWithAddresses(addrs []string) ([]peer.AddrInfo, error) {
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
	pis := make([]peer.AddrInfo, 0, len(peers))
	for id, maddrs := range peers {
		pis = append(pis, peer.AddrInfo{
			ID:    id,
			Addrs: maddrs,
		})
	}
	return pis, nil
}

// publicIPv4 returns true if the given ip is not reserved for a private address.
// of course, this only implies that it _might_ be public
// https://stackoverflow.com/a/41670589
func publicIPv4(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
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

// lanIPv4 returns true if the given ip is a LAN IP (192.168.x.x)
func lanIPv4(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 192 && ip4[1] == 168
	}
	return false
}

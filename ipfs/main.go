package ipfs

import (
	"context"
	"fmt"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/archive"
	"gx/ipfs/QmNueRyPRQiV7PUEpnP4GgGLuK1rKQLaRW7sfPvUetYig1/go-ipfs-cmds"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	pstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/Qme4QgoVPyQqxVc4G1c2L2wc9TDa6o294rtspGMnBNRujm/go-ipfs-addr"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface/options"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreunix"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"io"
	"io/ioutil"
	"net"
	"sort"
	"strings"
	"time"
)

var log = logging.MustGetLogger("ipfs")

const pinTimeout = time.Minute * 1
const catTimeout = time.Second * 30

type IpnsEntry struct {
	Name  string
	Value string
}

// GetDataAtPath return bytes under an ipfs path
func GetDataAtPath(ipfs *core.IpfsNode, pth string) ([]byte, error) {
	// convert string to an ipfs path
	ip, err := iface.ParsePath(pth)
	if err != nil {
		return nil, err
	}
	api := coreapi.NewCoreAPI(ipfs)
	ctx, cancel := context.WithTimeout(ipfs.Context(), catTimeout)
	defer cancel()
	reader, err := api.Unixfs().Cat(ctx, ip)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// GetArchiveAtPath builds an archive from directory links under an ipfs path
// NOTE: currently will bork if dir path contains other dirs (depth > 1)
func GetArchiveAtPath(ipfs *core.IpfsNode, path string) (io.Reader, error) {
	// convert string to an ipfs path
	ip, err := iface.ParsePath(path)
	if err != nil {
		return nil, err
	}
	api := coreapi.NewCoreAPI(ipfs)
	ctx, cancel := context.WithTimeout(ipfs.Context(), catTimeout)
	defer cancel()
	links, err := api.Unixfs().Ls(ctx, ip)
	if err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, nil
	}

	// virtual archive
	arch, err := archive.NewArchive(nil)
	for _, link := range links {
		data, err := GetDataAtPath(ipfs, link.Cid.Hash().B58String())
		if err != nil {
			return nil, err
		}
		arch.AddFile(data, link.Name)
	}
	if err := arch.Close(); err != nil {
		return nil, err
	}
	return arch.VirtualReader(), nil
}

// GetLinksAtPath return ipld links under a path
func GetLinksAtPath(ipfs *core.IpfsNode, path string) ([]*ipld.Link, error) {
	// convert string to an ipfs path
	ip, err := iface.ParsePath(path)
	if err != nil {
		return nil, err
	}
	api := coreapi.NewCoreAPI(ipfs)
	ctx, cancel := context.WithTimeout(ipfs.Context(), catTimeout)
	defer cancel()
	links, err := api.Unixfs().Ls(ctx, ip)
	if err != nil {
		return nil, err
	}
	return links, nil
}

// AddFileToDirectory adds bytes as file to a virtual directory (dag) structure
func AddFileToDirectory(ipfs *core.IpfsNode, dir uio.Directory, reader io.Reader, fname string) (*cid.Cid, error) {
	str, err := coreunix.Add(ipfs, reader)
	if err != nil {
		return nil, err
	}
	id, err := cid.Decode(str)
	if err != nil {
		return nil, err
	}
	node, err := ipfs.DAG.Get(ipfs.Context(), id)
	if err != nil {
		return nil, err
	}
	if err := dir.AddChild(ipfs.Context(), fname, node); err != nil {
		return nil, err
	}
	return id, nil
}

// Data pins
func PinData(ipfs *core.IpfsNode, data io.Reader) (*cid.Cid, error) {
	ctx, cancel := context.WithTimeout(ipfs.Context(), pinTimeout)
	defer cancel()
	api := coreapi.NewCoreAPI(ipfs)
	pth, err := api.Unixfs().Add(ctx, data)
	if err != nil {
		return nil, err
	}
	if err := api.Pin().Add(ctx, pth); err != nil {
		return nil, err
	}
	return pth.Cid(), nil
}

// PinPath takes an ipfs path string and pins it
func PinPath(ipfs *core.IpfsNode, path string, recursive bool) error {
	ip, err := iface.ParsePath(path)
	if err != nil {
		log.Errorf("error pinning path: %s, recursive: %t: %s", path, recursive, err)
		return err
	}
	ctx, cancel := context.WithTimeout(ipfs.Context(), pinTimeout)
	defer cancel()
	api := coreapi.NewCoreAPI(ipfs)
	if err := api.Pin().Add(ctx, ip, options.Pin.Recursive(recursive)); err != nil {
		return err
	}
	return nil
}

// UnpinPath takes an ipfs path string and unpins it
func UnpinPath(ipfs *core.IpfsNode, path string) error {
	ip, err := iface.ParsePath(path)
	if err != nil {
		log.Errorf("error unpinning path: %s: %s", path, err)
		return err
	}
	ctx, cancel := context.WithTimeout(ipfs.Context(), pinTimeout)
	defer cancel()
	api := coreapi.NewCoreAPI(ipfs)
	if err := api.Pin().Rm(ctx, ip); err != nil {
		return err
	}
	return nil
}

// PinDirectory pins a directory exluding any provided links
func PinDirectory(ipfs *core.IpfsNode, dir ipld.Node, exclude []string) error {
	ctx, cancel := context.WithTimeout(ipfs.Context(), pinTimeout)
	defer cancel()
	if err := ipfs.Pinning.Pin(ctx, dir, false); err != nil {
		return err
	}
outer:
	for _, item := range dir.Links() {
		for _, ex := range exclude {
			if item.Name == ex {
				continue outer
			}
		}
		node, err := item.GetNode(ctx, ipfs.DAG)
		if err != nil {
			return err
		}
		ipfs.Pinning.Pin(ctx, node, false)
	}
	return ipfs.Pinning.Flush()
}

// Publish publishes a node to ipns
func Publish(ctx context.Context, n *core.IpfsNode, k libp2pc.PrivKey, ref path.Path, dur time.Duration) (*IpnsEntry, error) {
	eol := time.Now().Add(dur)
	err := n.Namesys.PublishWithEOL(ctx, k, ref, eol)
	if err != nil {
		return nil, err
	}
	pid, err := peer.IDFromPrivateKey(k)
	if err != nil {
		return nil, err
	}
	return &IpnsEntry{Name: pid.Pretty(), Value: ref.String()}, nil
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
		log.Infof("swarm listening on %s\n", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Sort(sort.StringSlice(addrs))
	for _, addr := range addrs {
		log.Infof("swarm announcing %s\n", addr)
	}
	return nil
}

// PublicIPv4Addr tries to use the ipfs announced swarm addresses to determine a public ipv4 address
// NOTE: mobile peers must use the http cafe api for file uploads
// this method is used to inform them of that address during p2p registration
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

// publicIPv4 returns true if the given IP is superficially public
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

// PeersWithAddresses is a function that takes in a slice of string peer addresses
// (multiaddr + peerid) and returns a slice of properly constructed peers
func PeersWithAddresses(addrs []string) ([]pstore.PeerInfo, error) {
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

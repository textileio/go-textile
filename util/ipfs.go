package util

import (
	"context"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/core/cafe"
	iaddr "gx/ipfs/QmQViVWBHbU6HmYjXcdNq7tVASCNgdg64ZGcauuDkLCivW/go-ipfs-addr"
	"gx/ipfs/QmTjNRVt2fvaRFu93keEC7z5M1GS1iH6qZ9227htQioTUY/go-ipfs-cmds"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core/coreapi"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core/coreapi/interface/options"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core/coreunix"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/path"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

var log = logging.MustGetLogger("util")

const pinTimeout = time.Minute * 1
const catTimeout = time.Second * 30

type PublishOpts struct {
	VerifyExists bool
	PubValidTime time.Duration
}

type IpnsEntry struct {
	Name  string
	Value string
}

// GetDataAtPath return bytes under an ipfs path
func GetDataAtPath(ipfs *core.IpfsNode, path string) ([]byte, error) {
	// convert string to an ipfs path
	ip, err := coreapi.ParsePath(path)
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
	defer func() {
		if recover() != nil {
			log.Debug("node stopped")
		}
	}()

	return ioutil.ReadAll(reader)
}

// GetArchiveAtPath builds an archive from directory links under an ipfs path
// NOTE: currently will bork if dir path contains other dirs (depth > 1)
func GetArchiveAtPath(ipfs *core.IpfsNode, path string) (io.Reader, error) {
	// convert string to an ipfs path
	ip, err := coreapi.ParsePath(path)
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
	defer func() {
		if recover() != nil {
			log.Debug("node stopped")
		}
	}()
	if len(links) == 0 {
		return nil, nil
	}

	// virtual archive
	archive, err := client.NewArchive(nil)
	for _, link := range links {
		data, err := GetDataAtPath(ipfs, link.Cid.Hash().B58String())
		if err != nil {
			return nil, err
		}
		archive.AddFile(data, link.Name)
	}
	if err := archive.Close(); err != nil {
		return nil, err
	}

	return archive.VirtualReader(), nil
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

// ParsePeerParam takes a peer address string and returns p2p params
func ParsePeerParam(text string) (ma.Multiaddr, peer.ID, error) {
	// to be replaced with just multiaddr parsing, once ptp is a multiaddr protocol
	idx := strings.LastIndex(text, "/")
	if idx == -1 {
		pid, err := peer.IDB58Decode(text)
		if err != nil {
			return nil, "", err
		}

		return nil, pid, nil
	}

	addrS := text[:idx]
	peeridS := text[idx+1:]

	var maddr ma.Multiaddr
	var pid peer.ID

	// make sure addrS parses as a multiaddr.
	if len(addrS) > 0 {
		var err error
		maddr, err = ma.NewMultiaddr(addrS)
		if err != nil {
			return nil, "", err
		}
	}

	// make sure idS parses as a peer.ID
	var err error
	pid, err = peer.IDB58Decode(peeridS)
	if err != nil {
		return nil, "", err
	}

	return maddr, pid, nil
}

// PeersWithAddresses is a function that takes in a slice of string peer addresses
// (multiaddr + peerid) and returns a slice of properly constructed peers
func PeersWithAddresses(addrs []string) (pis []pstore.PeerInfo, err error) {
	iaddrs, err := parseAddresses(addrs)
	if err != nil {
		return nil, err
	}

	for _, a := range iaddrs {
		pis = append(pis, pstore.PeerInfo{
			ID:    a.ID(),
			Addrs: []ma.Multiaddr{a.Transport()},
		})
	}
	return pis, nil
}

// AddFileToDirectory adds bytes as file to a virtual directory (dag) structure
func AddFileToDirectory(ipfs *core.IpfsNode, dirb *uio.Directory, reader io.Reader, fname string) error {
	str, err := coreunix.Add(ipfs, reader)
	if err != nil {
		return err
	}
	id, err := cid.Decode(str)
	if err != nil {
		return err
	}
	node, err := ipfs.DAG.Get(ipfs.Context(), id)
	if err != nil {
		return err
	}
	if err := dirb.AddChild(ipfs.Context(), fname, node); err != nil {
		return err
	}
	return nil
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
	defer func() {
		if recover() != nil {
			log.Debug("node stopped")
		}
	}()
	return pth.Cid(), nil
}

// PinPath takes an ipfs path string and pins it
func PinPath(ipfs *core.IpfsNode, path string, recursive bool) error {
	ip, err := coreapi.ParsePath(path)
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
	defer func() {
		if recover() != nil {
			log.Debug("node stopped")
		}
	}()
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

// MultiaddrFromId creates a multiaddr from an id string
func MultiaddrFromId(id string) (ma.Multiaddr, error) {
	return ma.NewMultiaddr("/ipfs/" + id + "/")
}

func Publish(ctx context.Context, n *core.IpfsNode, k libp2pc.PrivKey, ref path.Path, opts *PublishOpts) (*IpnsEntry, error) {
	eol := time.Now().Add(opts.PubValidTime)
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

// parseAddresses is a function that takes in a slice of string peer addresses
// (multiaddr + peerid) and returns slices of multiaddrs and peerids.
func parseAddresses(addrs []string) (iaddrs []iaddr.IPFSAddr, err error) {
	iaddrs = make([]iaddr.IPFSAddr, len(addrs))
	for i, saddr := range addrs {
		iaddrs[i], err = iaddr.ParseString(saddr)
		if err != nil {
			return nil, cmds.ClientError("invalid peer address: " + err.Error())
		}
	}
	return
}

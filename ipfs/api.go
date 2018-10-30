package ipfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/archive"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface/options"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreunix"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/namesys/opts"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"io"
	"io/ioutil"
	"time"
)

var log = logging.Logger("tex-ipfs")

const pinTimeout = time.Minute
const catTimeout = time.Second * 30
const ipnsTimeout = time.Second * 10

type IpnsEntry struct {
	Name  string
	Value string
}

// DataAtPath return bytes under an ipfs path
func DataAtPath(ipfs *core.IpfsNode, pth string) ([]byte, error) {
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

// ArchiveAtPath builds an archive from directory links under an ipfs path
// NOTE: currently will bork if dir path contains other dirs (depth > 1)
func ArchiveAtPath(ipfs *core.IpfsNode, path string) (io.Reader, error) {
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
		data, err := DataAtPath(ipfs, link.Cid.Hash().B58String())
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

// LinksAtPath return ipld links under a path
func LinksAtPath(ipfs *core.IpfsNode, path string) ([]*ipld.Link, error) {
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

// Publish publishes a content id to ipns
func Publish(node *core.IpfsNode, sk libp2pc.PrivKey, id string, dur time.Duration, cache time.Duration) (*IpnsEntry, error) {
	if node.Mounts.Ipns != nil && node.Mounts.Ipns.IsActive() {
		return nil, errors.New("cannot manually publish while IPNS is mounted")
	}
	pth, err := path.ParsePath(id)
	if err != nil {
		return nil, err
	}
	eol := time.Now().Add(dur)
	ctx, cancel := context.WithTimeout(node.Context(), ipnsTimeout)
	ctx = context.WithValue(ctx, "ipns-publish-ttl", cache)
	defer cancel()
	if err := node.Namesys.PublishWithEOL(ctx, sk, pth, eol); err != nil {
		return nil, err
	}
	pid, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return &IpnsEntry{Name: pid.Pretty(), Value: pth.String()}, nil
}

// Publish publishes a node to ipns
func Resolve(node *core.IpfsNode, name peer.ID) (*path.Path, error) {
	// query options
	key := fmt.Sprintf("/ipns/%s", name.Pretty())
	var ropts []nsopts.ResolveOpt
	ropts = append(ropts, nsopts.Depth(1))
	ropts = append(ropts, nsopts.DhtRecordCount(16))
	ropts = append(ropts, nsopts.DhtTimeout(ipnsTimeout))

	// resolve w/ ipns
	pth, err := node.Namesys.Resolve(node.Context(), key, ropts...)
	if err != nil {
		return nil, err
	}
	return &pth, nil
}

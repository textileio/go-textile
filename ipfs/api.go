package ipfs

import (
	"context"
	"errors"
	"fmt"
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
func DataAtPath(node *core.IpfsNode, pth string) ([]byte, error) {
	ip, err := iface.ParsePath(pth)
	if err != nil {
		return nil, err
	}
	api := coreapi.NewCoreAPI(node)
	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()
	reader, err := api.Unixfs().Cat(ctx, ip)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// LinksAtPath return ipld links under a path
func LinksAtPath(node *core.IpfsNode, path string) ([]*ipld.Link, error) {
	ip, err := iface.ParsePath(path)
	if err != nil {
		return nil, err
	}
	api := coreapi.NewCoreAPI(node)
	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()
	links, err := api.Unixfs().Ls(ctx, ip)
	if err != nil {
		return nil, err
	}
	return links, nil
}

// AddDirectoryFile adds reader bytes to a virtual directory (dag) structure
func AddDirectoryFile(node *core.IpfsNode, dir uio.Directory, reader io.Reader, fname string) (*cid.Cid, error) {
	str, err := coreunix.Add(node, reader)
	if err != nil {
		return nil, err
	}
	id, err := cid.Decode(str)
	if err != nil {
		return nil, err
	}
	n, err := node.DAG.Get(node.Context(), id)
	if err != nil {
		return nil, err
	}
	if err := dir.AddChild(node.Context(), fname, n); err != nil {
		return nil, err
	}
	return id, nil
}

// AddPathToDirectory adds a link to a virtual directory (dag) structure
func AddDirectoryLink(node *core.IpfsNode, dir uio.Directory, fname string, pth string) error {
	id, err := cid.Decode(pth)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()
	n, err := node.DAG.Get(ctx, id)
	if err != nil {
		return err
	}
	ctx2, cancel2 := context.WithTimeout(node.Context(), catTimeout)
	defer cancel2()
	return dir.AddChild(ctx2, fname, n)
}

// AddData takes a reader and adds it, optionally pins it
func AddData(node *core.IpfsNode, data io.Reader, pin bool) (*cid.Cid, error) {
	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()
	api := coreapi.NewCoreAPI(node)
	pth, err := api.Unixfs().Add(ctx, data)
	if err != nil {
		return nil, err
	}
	if !pin {
		return pth.Cid(), nil
	}
	if err := api.Pin().Add(ctx, pth); err != nil {
		return nil, err
	}
	return pth.Cid(), nil
}

// GetNode returns a node behind an ipld link
func GetNode(node *core.IpfsNode, link *ipld.Link) (ipld.Node, error) {
	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()
	return link.GetNode(ctx, node.DAG)
}

// PinPath takes an ipfs path string and pins it
func PinPath(node *core.IpfsNode, path string, recursive bool) error {
	ip, err := iface.ParsePath(path)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()
	api := coreapi.NewCoreAPI(node)
	if err := api.Pin().Add(ctx, ip, options.Pin.Recursive(recursive)); err != nil {
		return err
	}
	return nil
}

// UnpinPath takes an ipfs path string and unpins it
func UnpinPath(node *core.IpfsNode, path string) error {
	ip, err := iface.ParsePath(path)
	if err != nil {
		log.Errorf("error unpinning path: %s: %s", path, err)
		return err
	}
	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()
	api := coreapi.NewCoreAPI(node)
	if err := api.Pin().Rm(ctx, ip); err != nil {
		return err
	}
	return nil
}

// PinDirectory pins a directory structure, not links
func PinNode(node *core.IpfsNode, n ipld.Node) error {
	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()
	if err := node.Pinning.Pin(ctx, n, false); err != nil {
		return err
	}
	return node.Pinning.Flush()
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

// Resolve resolves an ipns path to an ipfs path
func Resolve(node *core.IpfsNode, name peer.ID) (*path.Path, error) {
	key := fmt.Sprintf("/ipns/%s", name.Pretty())
	var ropts []nsopts.ResolveOpt
	ropts = append(ropts, nsopts.Depth(1))
	ropts = append(ropts, nsopts.DhtRecordCount(16))
	ropts = append(ropts, nsopts.DhtTimeout(ipnsTimeout))
	pth, err := node.Namesys.Resolve(node.Context(), key, ropts...)
	if err != nil {
		return nil, err
	}
	return &pth, nil
}

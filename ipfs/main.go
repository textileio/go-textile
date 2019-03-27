package ipfs

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/pin"
	"gx/ipfs/QmQmhotPUzVrMEWNK3x1R5jQ5ZHWyL7tVUrmRPjrBrvyCb/go-ipfs-files"
	"gx/ipfs/QmTbxNB1NwDesLmKTscr4udL2tVP7MaxvXnD1D9yX7g3PN/go-cid"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
	ipld "gx/ipfs/QmZ6nzCLwGLVfRzYLpD7pW6UNuBDKEcA2imJtVpbEx2rxy/go-ipld-format"
	logging "gx/ipfs/QmbkT7eMTyXfpeyB3ZMxxcxg7XH8t6uXp49jqzz4HB7BGF/go-log"
	uio "gx/ipfs/QmcYUTQ7tBZeH1CLsZM2S3xhMEZdvUgXvbjhpMsLDpk3oJ/go-unixfs/io"
)

var log = logging.Logger("tex-ipfs")

const pinTimeout = time.Minute
const catTimeout = time.Minute
const connectTimeout = time.Second * 10

// DataAtPath return bytes under an ipfs path
func DataAtPath(node *core.IpfsNode, pth string) ([]byte, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ip, err := iface.ParsePath(pth)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()

	f, err := api.Unixfs().Get(ctx, ip)
	if err != nil {
		log.Errorf("failed to get data: %s", pth)
		return nil, err
	}
	defer f.Close()

	var file files.File
	switch f := f.(type) {
	case files.File:
		file = f
	case files.Directory:
		return nil, iface.ErrIsDir
	default:
		return nil, iface.ErrNotSupported
	}

	return ioutil.ReadAll(file)
}

// LinksAtPath return ipld links under a path
func LinksAtPath(node *core.IpfsNode, pth string) ([]*ipld.Link, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ip, err := iface.ParsePath(pth)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()

	res, err := api.Unixfs().Ls(ctx, ip)
	if err != nil {
		return nil, err
	}

	links := make([]*ipld.Link, 0)
	for link := range res {
		links = append(links, link.Link)
	}

	return links, nil
}

// AddDataToDirectory adds reader bytes to a virtual dir
func AddDataToDirectory(node *core.IpfsNode, dir uio.Directory, fname string, reader io.Reader) (*cid.Cid, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	id, err := AddData(node, reader, false)
	if err != nil {
		return nil, err
	}

	n, err := api.Dag().Get(node.Context(), *id)
	if err != nil {
		return nil, err
	}

	if err := dir.AddChild(node.Context(), fname, n); err != nil {
		return nil, err
	}

	return id, nil
}

// AddLinkToDirectory adds a link to a virtual dir
func AddLinkToDirectory(node *core.IpfsNode, dir uio.Directory, fname string, pth string) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	id, err := cid.Decode(pth)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()

	nd, err := api.Dag().Get(ctx, id)
	if err != nil {
		return err
	}

	ctx2, cancel2 := context.WithTimeout(node.Context(), catTimeout)
	defer cancel2()

	return dir.AddChild(ctx2, fname, nd)
}

// AddData takes a reader and adds it, optionally pins it
func AddData(node *core.IpfsNode, reader io.Reader, pin bool) (*cid.Cid, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()

	pth, err := api.Unixfs().Add(ctx, files.NewReaderFile(reader))
	if err != nil {
		return nil, err
	}

	if pin {
		if err := api.Pin().Add(ctx, pth, options.Pin.Recursive(false)); err != nil {
			return nil, err
		}
	}
	id := pth.Cid()

	return &id, nil
}

// AddObject takes a reader and adds it as a DAG node, optionally pins it
func AddObject(node *core.IpfsNode, reader io.Reader, pin bool) (*cid.Cid, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()

	pth, err := api.Object().Put(ctx, reader)
	if err != nil {
		return nil, err
	}

	if pin {
		if err := api.Pin().Add(ctx, pth, options.Pin.Recursive(false)); err != nil {
			return nil, err
		}
	}
	id := pth.Cid()

	return &id, nil
}

// NodeAtLink returns the node behind an ipld link
func NodeAtLink(node *core.IpfsNode, link *ipld.Link) (ipld.Node, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()
	return link.GetNode(ctx, api.Dag())
}

// NodeAtCid returns the node behind a cid
func NodeAtCid(node *core.IpfsNode, id cid.Cid) (ipld.Node, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()
	return api.Dag().Get(ctx, id)
}

// NodeAtPath returns the last node under path
func NodeAtPath(node *core.IpfsNode, pth string) (ipld.Node, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	p, err := iface.ParsePath(pth)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()

	return api.ResolveNode(ctx, p)
}

type Node struct {
	Links []Link
	Data  string
}

type Link struct {
	Name, Hash string
	Size       uint64
}

// GetObjectAtPath returns the DAG object at the given path
func GetObjectAtPath(node *core.IpfsNode, pth string) ([]byte, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), catTimeout)
	defer cancel()

	ipth, err := iface.ParsePath(pth)
	if err != nil {
		return nil, err
	}
	nd, err := api.Object().Get(ctx, ipth)
	if err != nil {
		return nil, err
	}

	r, err := api.Object().Data(ctx, ipth)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	out := &Node{
		Links: make([]Link, len(nd.Links())),
		Data:  string(data),
	}

	for i, link := range nd.Links() {
		out.Links[i] = Link{
			Hash: link.Cid.String(),
			Name: link.Name,
			Size: link.Size,
		}
	}

	return json.Marshal(out)
}

// PinNode pins an ipld node
func PinNode(node *core.IpfsNode, nd ipld.Node, recursive bool) error {
	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()

	defer node.Blockstore.PinLock().Unlock()

	if err := node.Pinning.Pin(ctx, nd, recursive); err != nil {
		if strings.Contains(err.Error(), "already pinned recursively") {
			return nil
		}
		return err
	}

	return node.Pinning.Flush()
}

// UnpinNode unpins an ipld node
func UnpinNode(node *core.IpfsNode, nd ipld.Node, recursive bool) error {
	return UnpinCid(node, nd.Cid(), recursive)
}

// UnpinCid unpins a cid
func UnpinCid(node *core.IpfsNode, id cid.Cid, recursive bool) error {
	ctx, cancel := context.WithTimeout(node.Context(), pinTimeout)
	defer cancel()

	err := node.Pinning.Unpin(ctx, id, recursive)
	if err != nil && err != pin.ErrNotPinned {
		return err
	}

	return node.Pinning.Flush()
}

package wallet

import (
	"io"
	"time"
	"bytes"
	"image/jpeg"
	_ "image/png" // register other possible image types for decoding
	_ "image/gif"
	"encoding/json"
	"image"
	"io/ioutil"

	"github.com/disintegration/imaging"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/commands"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreunix"
	uio "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/unixfs/io"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"

	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"net/http"
	"github.com/textileio/textile-go/net"
)

type Photo map[string]string

type WalletData struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Photos []Photo `json:"photos"`
	LastHash string `json:"last_hash"`
}

func (w *WalletData) String() string {
	return "TODO"
}

// TODO: Put this in the sql db, not ipfs
func NewWalletData(node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)

	wallet := &WalletData{
		Created: time.Now(),
		Updated: time.Now(),
		Photos: make([]Photo, 0),
	}

	wb, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	p, err := api.Unixfs().Add(node.Context(), bytes.NewReader(wb))
	if err != nil {
		return err
	}

	// done automatically?
	if err := api.Pin().Add(node.Context(), p); err != nil {
		return err
	}

	err = publish(p, node)
	if err != nil {
		return err
	}

	return nil
}

// PinPhoto takes an io reader pointing to an image file, created a thumbnail, and adds
// both to a new directory, then finally pins that directory.
// TODO: need to "index" this in the sql db wallet
func PinPhoto(reader io.Reader, fname string, nd *core.IpfsNode, apiHost string) (ipld.Node, error) {
	// create thumbnail
	// FIXME: dunno if there's a better way to do this without consuming the fill stream
	// FIXME: into memory... as in, can we split the reader stream or something
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)
	th, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	th = imaging.Thumbnail(th, 100, 100, imaging.CatmullRom)
	thb := new(bytes.Buffer)
	if err = jpeg.Encode(thb, th, nil); err != nil {
		return nil, err
	}

	// rewind source reader for add
	r.Seek(0, 0)

	// top level directory
	dirb := uio.NewDirectory(nd.DAG)

	// add the images
	addFileToDirectory(dirb, r, fname, nd)
	addFileToDirectory(dirb, bytes.NewBuffer(thb.Bytes()), "thumb.jpg", nd)

	// pin the whole thing
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}

	if err := nd.Pinning.Pin(nd.Context(), dir, true); err != nil {
		return nil, err
	}

	if err := nd.Pinning.Flush(); err != nil {
		return nil, err
	}

	// pin it to server
	if apiHost != "" {
		res := &commands.AddPinOutput{}
		client := &http.Client{Timeout: 10 * time.Second}
		args := dir.Cid().Hash().B58String() + "&recursive=true"
		err = net.GetJson(client, apiHost+"/api/v0/pin/add?arg="+args, res)
		if err != nil {
			return dir, err
		}
	}

	return dir, nil
}

func publish(path iface.Path, node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)
	_, err := api.Name().Publish(node.Context(), path)
	return err
}

func addFileToDirectory(dirb *uio.Directory, r io.Reader, fname string, nd *core.IpfsNode) error {
	s, err := coreunix.Add(nd, r)
	if err != nil {
		return err
	}

	c, err := cid.Decode(s)
	if err != nil {
		return err
	}

	node, err := nd.DAG.Get(nd.Context(), c)
	if err != nil {
		return err
	}

	if err := dirb.AddChild(nd.Context(), fname, node); err != nil {
		return err
	}

	return nil
}

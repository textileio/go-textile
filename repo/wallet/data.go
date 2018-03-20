package wallet

import (
	"io"
	"time"
	"bytes"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	_ "github.com/disintegration/imaging"

	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/commands"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreunix"
	uio "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/unixfs/io"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"

	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
)

type PhotoData struct {
	Name string `json:"name"`
	Ext string `json:"extension"`
	Location []float64 `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}

type Data struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Photos []string `json:"photos"`
	LastHash string `json:"last_hash"`
}

func (w *Data) String() string {
	return "TODO"
}

// TODO: Put this in the sql db, not ipfs
func NewWalletData(node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)

	wallet := &Data{
		Created: time.Now(),
		Updated: time.Now(),
		Photos: make([]string, 0),
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

// PinPhoto takes an io reader pointing to an image file, creates a thumbnail, and adds
// both to a new directory, then finally pins that directory.
/* TODO: add an json object to this directory for the photo's metadata
   TODO: add the metadata to the photos table (add some more columns)
   TODO: remove file extensions from file name below
{
	name: sunset
	ext: png, etc.
	location: [lat, lon]
	timestamp: iso8601
}
NOTE: timestamp above would be time taken, whereas timestamp in sql index will be time added,
	maybe we add both to both places?
NOTE: thinking that name and ext should be here so that we can just call the links to the files
	in the directory "photo" and "thumb", thereby removing user private data from link names,
	then on retrieval, we can rename to the original name + ext.
	this also has the benefit of not having to add the filename to the sql db, since we
	will always know its link address: "/photo"
*/
func PinPhoto(reader io.Reader, fname string, thumb io.Reader, nd *core.IpfsNode, apiHost string) (ipld.Node, error) {

	dirb := uio.NewDirectory(nd.DAG)

	// add the image, maintaining the extension type
	ext := filepath.Ext(fname)
	sname := "photo" + ext
	addFileToDirectory(dirb, reader, sname, nd)

	// add the thumbnail
	addFileToDirectory(dirb, thumb, "thumb.jpg", nd)

	// create metadata object
	md := &PhotoData{
		Name: strings.TrimSuffix(fname, ext),
		Ext: ext,
		Location: make([]float64, 0),
		Timestamp: time.Now(),
	}
	wbb, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}
	addFileToDirectory(dirb, bytes.NewReader(wbb), "meta", nd)

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

package wallet

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreapi"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreunix"
	uio "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/unixfs/io"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
)

type PhotoData struct {
	Name      string    `json:"name"`
	Ext       string    `json:"extension"`
	Location  []float64 `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}

type Data struct {
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Photos   []string  `json:"photos"`
	LastHash string    `json:"last_hash"`
}

func (w *Data) String() string {
	return "TODO"
}

//func InitializeWalletData(ctx context.Context, pub Publisher, pins pin.Pinner, key ci.PrivKey) error {
//	wallet := &Data{
//		Created: time.Now(),
//		Updated: time.Now(),
//		Photos:  make([]string, 0),
//	}
//
//	wb, err := json.Marshal(wallet)
//	if err != nil {
//		return err
//	}
//
//	pins.
//
//	// pin recursively because this might already be pinned
//	// and doing a direct pin would throw an error in that case
//	err := pins.Pin(ctx, emptyDir, true)
//	if err != nil {
//		return err
//	}
//
//	err = pins.Flush()
//	if err != nil {
//		return err
//	}
//
//	return pub.Publish(ctx, key, path.FromCid(emptyDir.Cid()))
//}

func NewWalletData(node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)

	wallet := &Data{
		Created: time.Now(),
		Updated: time.Now(),
		Photos:  make([]string, 0),
	}

	wb, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	// add and pin it
	p, err := api.Unixfs().Add(node.Context(), bytes.NewReader(wb))
	if err != nil {
		return err
	}
	if err := api.Pin().Add(node.Context(), p); err != nil {
		return err
	}

	return nil
}

// AddPhoto takes an image file, and optionally a thumbnail file, and adds
// both to a new directory, then finally adds and pins that directory.
func AddPhoto(n *core.IpfsNode, pk libp2p.PubKey, p *os.File, t *os.File) (*net.MultipartRequest, error) {
	// path info
	path := p.Name()
	ext := strings.ToLower(filepath.Ext(path))
	dname := filepath.Dir(t.Name())

	// create a metadata file
	// TODO: get exif data from photo
	md := &PhotoData{
		Name:      strings.TrimSuffix(filepath.Base(path), ext),
		Ext:       ext,
		Location:  make([]float64, 0),
		Timestamp: time.Now(),
	}
	mdb, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}
	cmdb, err := net.Encrypt(pk, mdb)
	if err != nil {
		return nil, err
	}

	// create an empty virtual directory
	dirb := uio.NewDirectory(n.DAG)

	// add the image
	pb, err := getEncryptedReaderBytes(p, pk)
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(n, dirb, pb, "photo")
	if err != nil {
		return nil, err
	}

	// add the thumbnail
	tb, err := getEncryptedReaderBytes(t, pk)
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(n, dirb, tb, "thumb")
	if err != nil {
		return nil, err
	}

	// add the metadata file
	err = addFileToDirectory(n, dirb, cmdb, "meta")
	if err != nil {
		return nil, err
	}

	// pin it
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := pinPhoto(n, dir); err != nil {
		return nil, err
	}

	// create and init a new multipart request
	mr := &net.MultipartRequest{}
	mr.Init(dname, dir.Cid().Hash().B58String())

	// add files
	if err := mr.AddFile(pb, "photo"); err != nil {
		return nil, err
	}
	if err := mr.AddFile(tb, "thumb"); err != nil {
		return nil, err
	}
	if err := mr.AddFile(cmdb, "meta"); err != nil {
		return nil, err
	}

	// finish request
	if err := mr.Finish(); err != nil {
		return nil, err
	}

	return mr, nil
}

func getEncryptedReaderBytes(r io.Reader, pk libp2p.PubKey) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return net.Encrypt(pk, b)
}

func addFileToDirectory(n *core.IpfsNode, dirb *uio.Directory, b []byte, fname string) error {
	r := bytes.NewReader(b)
	s, err := coreunix.Add(n, r)
	if err != nil {
		return err
	}
	c, err := cid.Decode(s)
	if err != nil {
		return err
	}
	node, err := n.DAG.Get(n.Context(), c)
	if err != nil {
		return err
	}
	if err := dirb.AddChild(n.Context(), fname, node); err != nil {
		return err
	}
	return nil
}

func pinPhoto(n *core.IpfsNode, dir ipld.Node) error {
	// pin the top-level structure (recursive: false)
	if err := n.Pinning.Pin(n.Context(), dir, false); err != nil {
		return err
	}
	// pin thumbnail and metadata
	for _, item := range dir.Links() {
		if item.Name == "meta" || item.Name == "thumb" {
			tnode, err := item.GetNode(n.Context(), n.DAG)
			if err != nil {
				return err
			}
			n.Pinning.Pin(n.Context(), tnode, false)
		}
	}
	return n.Pinning.Flush()
}

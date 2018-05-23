package photos

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreunix"
	uio "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/unixfs/io"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
)

var log = logging.MustGetLogger("photos")

type Metadata struct {
	// photo data
	Name      string    `json:"name,omitempty"`
	Ext       string    `json:"extension,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	Added     time.Time `json:"added,omitempty"`
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`

	// user data
	Username string `json:"username,omitempty"`
	PeerID   string `json:"peer_id,omitempty"`
}

// Add adds a photo, it's thumbnail, and it's metadata to ipfs
func Add(n *core.IpfsNode, pk libp2p.PubKey, p *os.File, t *os.File, lc string, un string, cap string) (*net.MultipartRequest, *Metadata, error) {
	// path info
	path := p.Name()
	ext := strings.ToLower(filepath.Ext(path))
	dname := filepath.Dir(t.Name())
	name := strings.TrimSuffix(filepath.Base(path), ext)

	// try to extract exif data
	// TODO: get image size info
	// TODO: break this up into one method with multi sub-methods for testing
	var tm time.Time
	x, ok := exif.Decode(p)
	if ok == nil {
		// time taken
		tmTmp, err := x.DateTime()
		if err == nil {
			tm = tmTmp
		}
	}
	// create a metadata file
	md := &Metadata{
		Name:     name,
		Ext:      ext,
		Username: un,
		PeerID:   n.Identity.Pretty(),
		Created:  tm,
		Added:    time.Now(),
	}

	mdb, err := json.Marshal(md)
	if err != nil {
		return nil, nil, err
	}
	cmdb, err := net.Encrypt(pk, mdb)
	if err != nil {
		return nil, nil, err
	}

	// encrypt the last hash
	clcb, err := net.Encrypt(pk, []byte(lc))
	if err != nil {
		return nil, nil, err
	}

	// encrypt the caption
	ccapb, err := net.Encrypt(pk, []byte(cap))
	if err != nil {
		return nil, nil, err
	}

	// create an empty virtual directory
	dirb := uio.NewDirectory(n.DAG)

	// add the image
	p.Seek(0, 0)
	pr, err := ImagePathWithoutExif(p)
	if err != nil {
		return nil, nil, err
	}
	pb, err := getEncryptedReaderBytes(pr, pk)
	if err != nil {
		return nil, nil, err
	}
	err = addFileToDirectory(n, dirb, pb, "photo")
	if err != nil {
		return nil, nil, err
	}

	// add the thumbnail
	tr, err := ImagePathWithoutExif(t)
	if err != nil {
		return nil, nil, err
	}

	tb, err := getEncryptedReaderBytes(tr, pk)
	if err != nil {
		return nil, nil, err
	}
	err = addFileToDirectory(n, dirb, tb, "thumb")
	if err != nil {
		return nil, nil, err
	}

	// add the metadata file
	err = addFileToDirectory(n, dirb, cmdb, "meta")
	if err != nil {
		return nil, nil, err
	}

	// add caption
	err = addFileToDirectory(n, dirb, ccapb, "caption")
	if err != nil {
		return nil, nil, err
	}

	// add last update's cid
	err = addFileToDirectory(n, dirb, clcb, "last")
	if err != nil {
		return nil, nil, err
	}

	// pin it
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, nil, err
	}
	if err := pinPhoto(n, dir); err != nil {
		return nil, nil, err
	}

	// create and init a new multipart request
	mr := &net.MultipartRequest{}
	mr.Init(dname, dir.Cid().Hash().B58String())

	// add files
	if err := mr.AddFile(pb, "photo"); err != nil {
		return nil, nil, err
	}
	if err := mr.AddFile(tb, "thumb"); err != nil {
		return nil, nil, err
	}
	if err := mr.AddFile(cmdb, "meta"); err != nil {
		return nil, nil, err
	}
	if err := mr.AddFile(ccapb, "caption"); err != nil {
		return nil, nil, err
	}
	if err := mr.AddFile(clcb, "last"); err != nil {
		return nil, nil, err
	}

	// finish request
	if err := mr.Finish(); err != nil {
		return nil, nil, err
	}

	return mr, md, nil
}

// getEncryptedReaderBytes reads reader bytes and returns the encrypted result
func getEncryptedReaderBytes(r io.Reader, pk libp2p.PubKey) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return net.Encrypt(pk, b)
}

// addFileToDirectory adds bytes as file to a virtual directory (dag) structure
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

// pinPhoto pins the entire photo set dag, minues the full res image
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

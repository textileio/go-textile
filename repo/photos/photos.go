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

	"github.com/rwcarlsen/goexif/exif"
	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreunix"
	uio "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/unixfs/io"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
)

type Metadata struct {
	Name      string    `json:"name"`
	Ext       string    `json:"extension"`
	Created   time.Time `json:"created"`
	Added     time.Time `json:"added"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

// Add takes an image file, and optionally a thumbnail file, and adds
// both to a new directory, then finally adds and pins that directory.
func Add(n *core.IpfsNode, pk libp2p.PubKey, p *os.File, t *os.File, lc string) (*net.MultipartRequest, *Metadata, error) {
	// path info
	path := p.Name()
	ext := strings.ToLower(filepath.Ext(path))
	dname := filepath.Dir(t.Name())

	// try to extract exif data
	// TODO: get image size info
	// TODO: break this up into one method with multi sub-methods for testing
	var tm time.Time
	var lat, lon float64 = -1, -1
	x, err := exif.Decode(p)
	if err == nil {
		// time taken
		tmTmp, err := x.DateTime()
		if err == nil {
			tm = tmTmp
		}

		// coords taken
		latTmp, lonTmp, err := x.LatLong()
		if err == nil {
			lat, lon = latTmp, lonTmp
		}
	}

	// create a metadata file
	md := &Metadata{
		Name:      strings.TrimSuffix(filepath.Base(path), ext),
		Ext:       ext,
		Created:   tm,
		Added:     time.Now(),
		Latitude:  lat,
		Longitude: lon,
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

	// create an empty virtual directory
	dirb := uio.NewDirectory(n.DAG)

	// add the image
	p.Seek(0, 0)
	pb, err := getEncryptedReaderBytes(p, pk)
	if err != nil {
		return nil, nil, err
	}
	err = addFileToDirectory(n, dirb, pb, "photo")
	if err != nil {
		return nil, nil, err
	}

	// add the thumbnail
	tb, err := getEncryptedReaderBytes(t, pk)
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
	if err := mr.AddFile(clcb, "last"); err != nil {
		return nil, nil, err
	}

	// finish request
	if err := mr.Finish(); err != nil {
		return nil, nil, err
	}

	return mr, md, nil
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

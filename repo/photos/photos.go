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

	"github.com/tajtiattila/metadata/exif"
	"github.com/tajtiattila/metadata/exif/exiftag"
	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreunix"
	uio "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/unixfs/io"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
)

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
	hasExif := false
	x, err := exif.Decode(p)
	if err == nil {
		hasExif = true
		// time taken
		tmTmp, ok := x.DateTime()
		if ok {
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
	var file *os.File

	if hasExif {
		// strip sensitive GPS tags
		x.Set(exiftag.GPSLatitudeRef, nil)
		x.Set(exiftag.GPSLatitude, nil)
		x.Set(exiftag.GPSLongitudeRef, nil)
		x.Set(exiftag.GPSLongitude, nil)
		x.Set(exiftag.GPSAltitudeRef, nil)
		x.Set(exiftag.GPSAltitude, nil)
		x.Set(exiftag.GPSDateStamp, nil)
		x.Set(exiftag.GPSTimeStamp, nil)
		// copy photo buffer data to file, replacing exif with x
		p.Seek(0, 0) // rewind buffer reader
		ppath := filepath.Join(dname, name+"_tmp"+ext)
		file, err = os.OpenFile(ppath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()

		err = exif.Copy(file, p, x)
		if err != nil {
			return nil, nil, err
		}
	} else {
		file = p
	}

	// add the image
	file.Seek(0, 0)
	pb, err := getEncryptedReaderBytes(file, pk)
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

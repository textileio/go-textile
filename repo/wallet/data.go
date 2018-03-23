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

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreunix"
	uio "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/unixfs/io"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"

	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"mime/multipart"
	"log"
	"fmt"
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

type Hashed struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Bytes int64 `json:"Bytes,omitempty"`
	Size string `json:"Size,omitempty"`
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
	meta := bytes.NewReader(wbb)
	addFileToDirectory(dirb, meta, "meta", nd)

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
		statusCode, hashes := remotePin(reader, thumb, meta, apiHost)
		fmt.Println(statusCode)
		fmt.Println(hashes)
	}

	return dir, nil
}

func remotePin(reader io.Reader, thumb io.Reader, meta io.Reader, apiHost string) (int, []Hashed) {
		// Prepare form to submit to IPFS API.
		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		fw1, err := w.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Fatal(err)
		}
		if _, err = io.Copy(fw1, reader); err != nil {
			log.Fatal(err)
		}

		fw2, err := w.CreateFormFile("file", "thumb.jpg")
		if err != nil {
			log.Fatal(err)
		}
		if _, err = io.Copy(fw2, thumb); err != nil {
			log.Fatal(err)
		}

		fw3, err := w.CreateFormFile("file", "meta")
		if err != nil {
			log.Fatal(err)
		}
		if _, err = io.Copy(fw3, meta); err != nil {
			log.Fatal(err)
		}

		// Don't forget to close multipart writer.
		// If not closed, request will be missing terminating boundary.
		w.Close()

		// Now that we have form, submit it to handler.
		req, err := http.NewRequest("POST", apiHost+"/api/v0/add?wrap-with-directory=true", &b)
		if err != nil {
			log.Fatal(err)
		}
		// Don't forget to set the content type, this will contain the boundary.
		req.Header.Set("Content-Type", w.FormDataContentType())

		// Submit the request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		// Check the response
		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("bad status: %s", res.Status)
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res.StatusCode)
			fmt.Println(body)
		}
		hashes := []Hashed{}
		err = json.NewDecoder(res.Body).Decode(hashes)
		if err != nil {
			log.Fatal(err)
		}
		return res.StatusCode, hashes
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

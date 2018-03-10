package core

import (
	"encoding/base64"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net"
	"os"
	"path"
	"strings"

	"github.com/textileio/mill-go/ipfs"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/unixfs/io"
	"github.com/nfnt/resize"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"time"
)

func (n *TextileNode) addImage(img image.Image, imgPath string) (string, error) {
	out, err := os.Create(imgPath)
	if err != nil {
		return "", err
	}
	jpeg.Encode(out, img, nil)
	out.Close()
	return ipfs.AddFile(n.Context, imgPath)
}

func (n *TextileNode) addResizedImage(img image.Image, imgCfg *image.Config, w, h uint, imgPath string) (string, error) {
	width, height := getImageAttributes(w, h, uint(imgCfg.Width), uint(imgCfg.Height))
	newImg := resize.Resize(width, height, img, resize.Lanczos3)
	return n.addImage(newImg, imgPath)
}

func decodeImageData(base64ImageData string) (image.Image, *image.Config, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64ImageData))
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, nil, err
	}
	reader = base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64ImageData))
	imgCfg, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, nil, err
	}
	return img, &imgCfg, err
}

func getImageAttributes(targetWidth, targetHeight, imgWidth, imgHeight uint) (width, height uint) {
	targetRatio := float32(targetWidth) / float32(targetHeight)
	imageRatio := float32(imgWidth) / float32(imgHeight)
	var h, w float32
	if imageRatio > targetRatio {
		h = float32(targetHeight)
		w = float32(targetHeight) * imageRatio
	} else {
		w = float32(targetWidth)
		h = float32(targetWidth) * (float32(imgHeight) / float32(imgWidth))
	}
	return uint(w), uint(h)
}

func (n *TextileNode) FetchImage(peerId string, imageType string, size string, useCache bool) (io.DagReader, error) {
	fetch := func(rootHash string) (io.DagReader, error) {
		var dr io.DagReader
		var err error

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()
		if rootHash == "" {
			query := "/ipns/" + peerId + "/images/" + size + "/" + imageType
			dr, err = coreunix.Cat(ctx, n.IpfsNode, query)
			if err != nil {
				return dr, err
			}
		} else {
			query := "/ipfs/" + rootHash + "/images/" + size + "/" + imageType
			dr, err = coreunix.Cat(ctx, n.IpfsNode, query)
			if err != nil {
				return dr, err
			}
		}
		return dr, nil
	}

	var dr io.DagReader
	var err error
	//var recordAvailable bool
	//var val interface{}
	if useCache {
		//val, err = n.IpfsNode.Repo.Datastore().Get(ds.NewKey(CachePrefix + peerId))
		//if err != nil { // No record in datastore
		//	dr, err = fetch("")
		//	if err != nil {
		//		return dr, err
		//	}
		//} else { // Record available, let's see how old it is
		//	entry := new(ipnspb.IpnsEntry)
		//	err = proto.Unmarshal(val.([]byte), entry)
		//	if err != nil {
		//		return dr, err
		//	}
		//	p, err := ipnspath.ParsePath(string(entry.GetValue()))
		//	if err != nil {
		//		return dr, err
		//	}
		//	eol, ok := CheckEOL(entry)
		//	if ok && eol.Before(time.Now()) { // Too old, fetch new profile
		//		dr, err = fetch("")
		//	} else { // Relatively new, we can do a standard IPFS query (which should be cached)
		//		dr, err = fetch(strings.TrimPrefix(p.String(), "/ipfs/"))
		//		// Let's now try to get the latest record in a new goroutine so it's available next time
		//		go fetch("")
		//	}
		//	if err != nil {
		//		return dr, err
		//	}
		//	recordAvailable = true
		//}
	} else {
		dr, err = fetch("")
		if err != nil {
			return dr, err
		}
		//recordAvailable = false
	}

	//// Update the record with a new EOL
	//go func() {
	//	if !recordAvailable {
	//		val, err = n.IpfsNode.Repo.Datastore().Get(ds.NewKey(CachePrefix + peerId))
	//		if err != nil {
	//			return
	//		}
	//	}
	//	entry := new(ipnspb.IpnsEntry)
	//	err = proto.Unmarshal(val.([]byte), entry)
	//	if err != nil {
	//		return
	//	}
	//	entry.Validity = []byte(u.FormatRFC3339(time.Now().Add(CachedProfileTime)))
	//	v, err := proto.Marshal(entry)
	//	if err != nil {
	//		return
	//	}
	//	n.IpfsNode.Repo.Datastore().Put(ds.NewKey(CachePrefix+peerId), v)
	//}()
	return dr, nil
}

func (n *TextileNode) GetBase64Image(url string) (base64ImageData, filename string, err error) {
	dial := net.Dial
	tbTransport := &http.Transport{Dial: dial}
	client := &http.Client{Transport: tbTransport, Timeout: time.Second * 30}
	resp, err := client.Get(url)
	if err != nil {
		return "", "", err
	}
	imgBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	img := base64.StdEncoding.EncodeToString(imgBytes)
	u, err := netUrl.Parse(url)
	if err != nil {
		return "", "", err
	}
	_, filename = path.Split(u.Path)
	return img, filename, nil
}

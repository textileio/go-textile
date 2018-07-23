package cafe

import (
	"archive/tar"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/wallet/util"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"io"
	"net/http"
)

func (c *Cafe) pin(g *gin.Context) {
	// handle based on content type
	//var id *cid.Cid
	var hash string
	cType := g.Request.Header.Get("Content-Type")
	switch cType {
	case "application/octet-stream":
		var err error
		id, err := util.PinData(c.Ipfs(), g.Request.Body)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hash = id.Hash().B58String()
	case "application/gzip":
		// create a virtual directory for the photo
		dirb := uio.NewDirectory(c.Ipfs().DAG)
		// unpack archive
		gr, err := gzip.NewReader(g.Request.Body)
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tr := tar.NewReader(gr)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			switch header.Typeflag {
			case tar.TypeDir:
				g.JSON(http.StatusBadRequest, gin.H{"error": "directories are not supported"})
				return
			case tar.TypeReg:
				if err := util.AddFileToDirectory(c.Ipfs(), dirb, tr, header.Name); err != nil {
					g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			default:
				continue
			}
		}

		// pin the directory
		dir, err := dirb.GetNode()
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := util.PinDirectory(c.Ipfs(), dir, []string{}); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hash = dir.Cid().Hash().B58String()
	}

	// ship it
	g.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"id":     hash,
	})
}

package cafe

import (
	"archive/tar"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/util"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"io"
	"net/http"
)

func (c *Cafe) pin(g *gin.Context) {
	var id *cid.Cid

	// handle based on content type
	cType := g.Request.Header.Get("Content-Type")
	switch cType {
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
		id = dir.Cid()

	default:
		var err error
		id, err = util.PinData(c.Ipfs(), g.Request.Body)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// ship it
	g.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"id":     id.Hash().B58String(),
	})
}

package cafe

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/wallet/util"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"io"
	"net/http"
)

func (c *Cafe) pin(g *gin.Context) {
	// handle based on content type
	var id *cid.Cid
	cType := g.Request.Header.Get("Content-Type")
	switch cType {
	case "application/octet-stream":
		var err error
		id, err = util.PinData(c.Ipfs(), g.Request.Body)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "application/gzip":
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
			case tar.TypeReg:
				fmt.Println(header.Name)
			default:
				continue
			}
		}
	}

	// ship it
	g.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"id":     id.Hash().B58String(),
	})
}

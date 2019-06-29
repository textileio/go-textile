package core

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/golang/protobuf/proto"
	cid "github.com/ipfs/go-cid"
	uio "github.com/ipfs/go-unixfs/io"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

// pin take raw data or a tarball and pins it to the local ipfs node.
// request must be authenticated with a token
func (c *cafeApi) pin(g *gin.Context) {
	// handle based on content type
	var id cid.Cid
	cType := g.Request.Header.Get("Content-Type")
	switch cType {
	case "application/gzip":
		dirb := uio.NewDirectory(c.node.Ipfs().DAG)

		gr, err := gzip.NewReader(g.Request.Body)
		if err != nil {
			c.abort(g, http.StatusBadRequest, err)
			return
		}
		tr := tar.NewReader(gr)

		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.abort(g, http.StatusInternalServerError, err)
				return
			}

			switch header.Typeflag {
			case tar.TypeDir:
				c.abort(g, http.StatusBadRequest, fmt.Errorf("directories are not supported"))
				return
			case tar.TypeReg:
				_, err = ipfs.AddDataToDirectory(c.node.Ipfs(), dirb, header.Name, tr)
				if err != nil {
					c.abort(g, http.StatusInternalServerError, err)
					return
				}
			default:
				continue
			}
		}

		// pin the directory
		dir, err := dirb.GetNode()
		if err != nil {
			c.abort(g, http.StatusInternalServerError, err)
			return
		}
		err = ipfs.PinNode(c.node.Ipfs(), dir, true)
		if err != nil {
			c.abort(g, http.StatusInternalServerError, err)
			return
		}
		id = dir.Cid()

	case "application/octet-stream":
		idp, err := ipfs.AddData(c.node.Ipfs(), g.Request.Body, true, false)
		if err != nil {
			c.abort(g, http.StatusInternalServerError, err)
			return
		}
		id = *idp
	default:
		c.abort(g, http.StatusBadRequest, fmt.Errorf("invalid content-type"))
		return
	}
	hash := id.Hash().B58String()

	log.Debugf("pinned request with content type %s: %s", cType, hash)

	g.JSON(http.StatusCreated, gin.H{"id": hash})
}

// service is an HTTP entry point for the cafe service
func (c *cafeApi) service(g *gin.Context) {
	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()

	buf.Grow(bytes.MinRead)
	_, err := buf.ReadFrom(g.Request.Body)
	if err != nil && err != io.EOF {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// parse body as a service envelope
	pmes := new(pb.Envelope)
	err = proto.Unmarshal(buf.Bytes(), pmes)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	peerId := g.Request.Header.Get("X-Textile-Peer")
	if peerId == "" {
		g.String(http.StatusBadRequest, "missing peer ID")
		return
	}
	mPeer, err := peer.IDB58Decode(peerId)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	err = c.node.cafe.service.VerifyEnvelope(pmes, mPeer)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// handle the message as normal
	log.Debugf("received %s from %s", pmes.Message.Type.String(), mPeer.Pretty())
	rpmes, err := c.node.cafe.Handle(pmes, mPeer)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if rpmes != nil {
		res, err := proto.Marshal(rpmes)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}

		g.Render(http.StatusOK, render.Data{Data: res})
		return
	}

	// handle the message as a JSON stream
	rpmesCh, errCh, cancel := c.node.cafe.HandleStream(pmes, mPeer)
	g.Stream(func(w io.Writer) bool {
		select {
		case <-g.Request.Context().Done():
			log.Debug("closing request stream")
			close(cancel)

		case err := <-errCh:
			g.String(http.StatusBadRequest, err.Error())
			return false

		case rpmes, ok := <-rpmesCh:
			if !ok {
				g.Status(http.StatusOK)
				return false
			}
			log.Debugf("responding with %s to %s", rpmes.Message.Type.String(), mPeer.Pretty())

			g.JSON(http.StatusOK, rpmes)
			_, _ = g.Writer.Write([]byte("\n"))
		}
		return true
	})
}

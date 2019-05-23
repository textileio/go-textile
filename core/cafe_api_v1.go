package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/pin"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

func (c *cafeApi) store(g *gin.Context) {
	var err error
	var aid *cid.Cid
	stype := g.Request.Header.Get("X-Textile-Store-Type")

	form, err := g.MultipartForm()
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	files := form.File["file"]

	var f multipart.File
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	for _, file := range files {
		f, err = file.Open()
		if err != nil {
			c.abort(g, http.StatusBadRequest, err)
			return
		}

		switch stype {
		case "data":
			aid, err = ipfs.AddData(c.node.Ipfs(), f, true, false)
		case "object":
			aid, err = ipfs.AddObject(c.node.Ipfs(), f, true)
		default:
			c.abort(g, http.StatusBadRequest, fmt.Errorf("missing store type header"))
			return
		}
		if err != nil {
			c.abort(g, http.StatusBadRequest, err)
			return
		}

		log.Debugf("stored %s", aid.Hash().B58String())

		f.Close()
		f = nil
	}

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) unstore(g *gin.Context) {
	id, err := cid.Decode(g.Param("cid"))
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	pinned, err := c.node.Ipfs().Pinning.CheckIfPinned(id)
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	for _, p := range pinned {
		if p.Mode != pin.NotPinned {
			err = ipfs.UnpinCid(c.node.Ipfs(), p.Key, true)
			if err != nil {
				c.abort(g, http.StatusBadRequest, err)
				return
			}

			log.Debugf("unstored %s", p.Key.Hash().B58String())
		}
	}

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) storeThread(g *gin.Context) {
	pid := g.Request.Header.Get("X-Textile-Peer")
	id := g.Param("id")

	client := c.node.datastore.CafeClients().Get(pid)
	if client == nil {
		c.abort(g, http.StatusForbidden, nil)
		return
	}

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

	err = c.node.datastore.CafeClientThreads().AddOrUpdate(&pb.CafeClientThread{
		Id:         id,
		Client:     client.Id,
		Ciphertext: buf.Bytes(),
	})
	if err != nil {
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	log.Debugf("stored thread %s", id)

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) unstoreThread(g *gin.Context) {
	pid := g.Request.Header.Get("X-Textile-Peer")
	id := g.Param("id")

	client := c.node.datastore.CafeClients().Get(pid)
	if client == nil {
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	err := c.node.datastore.CafeClientThreads().Delete(id, client.Id)
	if err != nil {
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	log.Debugf("unstored thread %s", id)

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) deliverMessage(g *gin.Context) {
	pid := g.Request.Header.Get("X-Textile-Peer")
	clientId := g.Param("pid")

	client := c.node.datastore.CafeClients().Get(clientId)
	if client == nil {
		log.Warningf("received message from %s for unknown client %s", pid, clientId)
		g.Status(http.StatusOK)
		return
	}

	// message id is the request body
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

	mid, err := cid.Decode(string(buf.Bytes()))
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	err = c.node.datastore.CafeClientMessages().AddOrUpdate(&pb.CafeClientMessage{
		Id:     mid.Hash().B58String(),
		Peer:   pid,
		Client: client.Id,
		Date:   ptypes.TimestampNow(),
	})
	if err != nil {
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	go func() {
		cpid, err := peer.IDB58Decode(client.Id)
		if err != nil {
			log.Errorf("error parsing client id %s: %s", client.Id, err)
			return
		}
		err = c.node.cafe.notifyClient(cpid)
		if err != nil {
			log.Debugf("unable to notify client: %s", client.Id)
		}
	}()

	log.Debugf("delivered message %s", mid.Hash().B58String())

	g.Status(http.StatusOK)
}

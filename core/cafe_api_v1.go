package core

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
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

		aid, err = ipfs.AddObject(c.node.Ipfs(), f, true)
		if err != nil {
			_, _ = f.Seek(0, 0)
			aid, err = ipfs.AddData(c.node.Ipfs(), f, true, false)
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
	from := g.GetString("from")
	id := g.Param("id")

	client := c.node.datastore.CafeClients().Get(from)
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
	from := g.GetString("from")
	id := g.Param("id")

	client := c.node.datastore.CafeClients().Get(from)
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
	from := g.Param("from")
	clientId := g.Param("to")

	client := c.node.datastore.CafeClients().Get(clientId)
	if client == nil {
		log.Warningf("received message for unknown client %s", clientId)
		g.Status(http.StatusOK)
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
	body := buf.Bytes()

	// pin inner node
	nenv := new(pb.Envelope)
	err = proto.Unmarshal(body, nenv)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	tenv := new(pb.ThreadEnvelope)
	err = ptypes.UnmarshalAny(nenv.Message.Payload, tenv)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	oid, err := ipfs.AddObject(c.node.Ipfs(), bytes.NewReader(tenv.Node), true)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	node, err := ipfs.NodeAtCid(c.node.Ipfs(), *oid)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	_, err = extractNode(c.node.Ipfs(), node)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// pin envelope
	id, err := ipfs.AddData(c.node.Ipfs(), bytes.NewReader(body), true, false)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	msgId := id.Hash().B58String()
	err = c.node.datastore.CafeClientMessages().AddOrUpdate(&pb.CafeClientMessage{
		Id:     msgId,
		Peer:   from,
		Client: client.Id,
		Date:   ptypes.TimestampNow(),
	})
	if err != nil {
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	go func() {
		err = c.node.cafe.notifyClient(client.Id)
		if err != nil {
			log.Debugf("unable to notify client: %s", client.Id)
		}
	}()

	log.Debugf("delivered message %s", msgId)

	g.Status(http.StatusOK)
}

func (c *cafeApi) search(g *gin.Context) {
	from := g.GetString("from")

	pid, err := peer.IDB58Decode(from)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()

	buf.Grow(bytes.MinRead)
	_, err = buf.ReadFrom(g.Request.Body)
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

	// handle the message as a JSON stream
	rpmesCh, errCh, cancel := c.node.cafe.HandleStream(pmes, pid)
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
			log.Debugf("responding with %s", rpmes.Message.Type.String())

			g.JSON(http.StatusOK, rpmes)
			_, _ = g.Writer.Write([]byte("\n"))
		}
		return true
	})
}

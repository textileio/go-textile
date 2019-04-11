package core

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/pin"
	"gx/ipfs/QmTbxNB1NwDesLmKTscr4udL2tVP7MaxvXnD1D9yX7g3PN/go-cid"
	peer "gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

func (c *cafeApi) store(g *gin.Context) {
	id, err := cid.Decode(g.Param("cid"))
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	hash := id.Hash().B58String()

	var aid *cid.Cid
	switch g.Request.Header.Get("X-Textile-Store-Type") {
	case "data":
		aid, err = ipfs.AddData(c.node.Ipfs(), g.Request.Body, true)
	case "object":
		aid, err = ipfs.AddObject(c.node.Ipfs(), g.Request.Body, true)
	default:
		c.abort(g, http.StatusBadRequest, fmt.Errorf("missing store type header"))
		return
	}
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	rhash := aid.Hash().B58String()

	log.Debugf("stored %s", rhash)

	if rhash != hash {
		c.abort(g, http.StatusBadRequest, fmt.Errorf("cids do not match (received %s, resolved %s)", hash, rhash))
		return
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
			if err := ipfs.UnpinCid(c.node.Ipfs(), p.Key, true); err != nil {
				c.abort(g, http.StatusBadRequest, err)
				return
			}
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

	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	thrd := &pb.CafeClientThread{
		Id:         id,
		Client:     client.Id,
		Ciphertext: body,
	}
	if err := c.node.datastore.CafeClientThreads().AddOrUpdate(thrd); err != nil {
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

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

	if err := c.node.datastore.CafeClientThreads().Delete(id, client.Id); err != nil {
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

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
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	mid, err := cid.Decode(string(body))
	if err != nil {
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	message := &pb.CafeClientMessage{
		Id:     mid.Hash().B58String(),
		Peer:   pid,
		Client: client.Id,
		Date:   ptypes.TimestampNow(),
	}
	if err := c.node.datastore.CafeClientMessages().AddOrUpdate(message); err != nil {
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	go func() {
		cpid, err := peer.IDB58Decode(client.Id)
		if err != nil {
			log.Errorf("error parsing client id %s: %s", client.Id, err)
			return
		}
		if err := c.node.cafe.notifyClient(cpid); err != nil {
			log.Debugf("unable to notify offline client: %s", client.Id)
		}
	}()

	g.Status(http.StatusOK)
}

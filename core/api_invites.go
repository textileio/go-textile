package core

import (
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"net/http"

	"github.com/mr-tron/base58/base58"

	"github.com/gin-gonic/gin"
)

func (a *api) createInvite(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	var pid peer.ID
	if opts["peer"] != "" {
		pid, err = peer.IDB58Decode(opts["peer"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	threadId := opts["thread"]
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	thrd := a.node.Thread(threadId)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}

	result := make(map[string]string)
	if pid != "" {
		hash, err := thrd.AddInvite(pid)
		if err != nil {
			a.abort500(g, err)
			return
		}
		result["invite"] = hash.B58String()
	} else {
		hash, key, err := thrd.AddExternalInvite()
		if err != nil {
			a.abort500(g, err)
			return
		}
		result["invite"] = hash.B58String()
		result["key"] = base58.FastBase58Encoding(key)
	}

	g.JSON(http.StatusCreated, result)
}

func (a *api) acceptInvite(g *gin.Context) {
	id := g.Param("id")
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	var hash mh.Multihash
	if opts["key"] != "" {
		key, err := base58.Decode(opts["key"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		hash, err = a.node.AcceptExternalThreadInvite(id, key)
	} else {
		hash, err = a.node.AcceptThreadInvite(id)
	}
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if hash == nil {
		g.String(http.StatusConflict, "thread already exists")
		return
	}

	info, err := a.node.BlockInfo(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, info)
}

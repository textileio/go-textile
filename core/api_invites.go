package core

import (
	"net/http"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/pb"
)

// createInvites godoc
// @Summary Create an invite to a thread
// @Description Creates a direct peer-to-peer or external invite to a thread
// @Tags invites
// @Produce application/json
// @Param X-Textile-Opts header string false "thread: Thread ID (can also use 'default'), peer: Peer ID (omit to create an external invite)" default(thread=,peer=)
// @Success 201 {object} pb.NewInvite "invite"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /invites [post]
func (a *api) createInvites(g *gin.Context) {
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

	result := &pb.NewInvite{}
	if pid != "" {
		hash, err := thrd.AddInvite(pid)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		result.Id = hash.B58String()
	} else {
		hash, key, err := thrd.AddExternalInvite()
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		result.Id = hash.B58String()
		result.Key = base58.FastBase58Encoding(key)
	}

	pbJSON(g, http.StatusCreated, result)
}

// lsInvites godoc
// @Summary List invites
// @Description Lists all pending thread invites
// @Tags invites
// @Produce application/json
// @Success 200 {object} pb.InviteViewList "invites"
// @Router /invites [get]
func (a *api) lsInvites(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.node.Invites())
}

// acceptInvites godoc
// @Summary Accept a thread invite
// @Description Accepts a direct peer-to-peer or external invite to a thread. Use the key option
// @Description with an external invite
// @Tags invites
// @Produce application/json
// @Param id path string true "invite id"
// @Param X-Textile-Opts header string false "key: key for an external invite" default(key=)
// @Success 201 {object} pb.Block "join block"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /invites/{id}/accept [post]
func (a *api) acceptInvites(g *gin.Context) {
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
		hash, err = a.node.AcceptExternalInvite(id, key)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		hash, err = a.node.AcceptInvite(id)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}
	if hash == nil {
		g.String(http.StatusConflict, "thread already exists")
		return
	}

	block, err := a.node.BlockView(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, block)
}

// ignoreInvites godoc
// @Summary Ignore a thread invite
// @Description Ignores a direct peer-to-peer invite to a thread
// @Tags invites
// @Produce application/json
// @Param id path string true "invite id"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Router /invites/{id}/ignore [post]
func (a *api) ignoreInvites(g *gin.Context) {
	id := g.Param("id")

	if err := a.node.IgnoreInvite(id); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.String(http.StatusOK, "ok")
}

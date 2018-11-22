package core

import (
	"net/http"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/gin-gonic/gin"
)

func (a *api) getProfile(g *gin.Context) {
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
	} else {
		pid = a.node.node.Identity
	}

	profile, err := a.node.Profile(pid)
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusOK, profile)
}

func (a *api) setUsername(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing username")
		return
	}
	if err := a.node.SetUsername(args[0]); err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusCreated, "ok")
}

func (a *api) setAvatar(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing avatar")
		return
	}
	if err := a.node.SetAvatar(args[0]); err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusCreated, "ok")
}

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) getProfile(g *gin.Context) {
	profile := a.node.Profile()
	if profile == nil {
		g.String(http.StatusBadRequest, "profile is not set")
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

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/pb"
)

func (a *api) addCafes(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing cafe host")
		return
	}
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	token := opts["token"]
	if token == "" {
		g.String(http.StatusBadRequest, "missing access token")
		return
	}
	session, err := a.node.RegisterCafe(args[0], token)
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusCreated, session)
}

func (a *api) lsCafes(g *gin.Context) {
	sessions, err := a.node.CafeSessions()
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(sessions) == 0 {
		sessions = make([]*pb.CafeSession, 0)
	}
	g.JSON(http.StatusOK, sessions)
}

func (a *api) getCafes(g *gin.Context) {
	id := g.Param("id")
	session, err := a.node.CafeSession(id)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if session == nil {
		g.String(http.StatusNotFound, "cafe not found")
		return
	}
	g.JSON(http.StatusOK, session)
}

func (a *api) rmCafes(g *gin.Context) {
	id := g.Param("id")
	if err := a.node.DeregisterCafe(id); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

func (a *api) checkCafeMessages(g *gin.Context) {
	if err := a.node.CheckCafeMessages(); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

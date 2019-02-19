package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) addBlockComments(g *gin.Context) {
	id := g.Param("id")

	thrd := a.getBlockThread(g, id)
	if thrd == nil {
		return
	}

	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing comment body")
		return
	}

	hash, err := thrd.AddComment(id, args[0])
	if err != nil {
		a.abort500(g, err)
		return
	}

	info, err := a.node.Comment(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, info)
}

func (a *api) lsBlockComments(g *gin.Context) {
	id := g.Param("id")

	comments, err := a.node.Comments(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	pbJSON(g, http.StatusOK, comments)
}

func (a *api) getBlockComment(g *gin.Context) {
	info, err := a.node.Comment(g.Param("id"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}

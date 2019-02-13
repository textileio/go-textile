package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) addBlockLikes(g *gin.Context) {
	id := g.Param("id")

	thrd := a.getBlockThread(g, id)
	if thrd == nil {
		return
	}

	hash, err := thrd.AddLike(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	block, err := a.node.Block(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	info, err := a.node.FeedLike(block, true)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, info)
}

func (a *api) lsBlockLikes(g *gin.Context) {
	id := g.Param("id")

	likes, err := a.node.Likes(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	pbJSON(g, http.StatusOK, likes)
}

func (a *api) getBlockLike(g *gin.Context) {
	id := g.Param("id")

	block, err := a.node.Block(id)
	if err != nil {
		g.String(http.StatusNotFound, "block not found")
		return
	}

	info, err := a.node.FeedLike(block, true)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}

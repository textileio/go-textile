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

	info, err := a.node.Like(hash.B58String())
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
	info, err := a.node.Like(g.Param("id"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}

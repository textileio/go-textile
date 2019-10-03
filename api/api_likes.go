package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// addBlockLikes godoc
// @Summary Add a like
// @Description Adds a like to a thread block
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 201 {object} pb.Like "like"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks/{id}/likes [post]
func (a *Api) addBlockLikes(g *gin.Context) {
	id := g.Param("id")

	thread, err, code := getBlockThread(a.Node, id)
	if err != nil {
		sendError(g, err, code)
		return
	}

	hash, err := thread.AddLike(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	like, err := a.Node.Like(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	pbJSON(g, http.StatusCreated, like)
}

// lsBlockLikes godoc
// @Summary List likes
// @Description Lists likes on a thread block
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.LikeList "likes"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks/{id}/likes [get]
func (a *Api) lsBlockLikes(g *gin.Context) {
	id := g.Param("id")

	likes, err := a.Node.Likes(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	pbJSON(g, http.StatusOK, likes)
}

// getBlockLike godoc
// @Summary Get thread like
// @Description Gets a thread like by block ID
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.Like "like"
// @Failure 400 {string} string "Bad Request"
// @Router /blocks/{id}/like [get]
func (a *Api) getBlockLike(g *gin.Context) {
	info, err := a.Node.Like(g.Param("id"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}

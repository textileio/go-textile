package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// addBlockComments godoc
// @Summary Add a comment
// @Description Adds a comment to a thread block
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Param X-Textile-Args header string true "urlescaped comment body"
// @Success 201 {object} pb.Comment "comment"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks/{id}/comments [post]
func (a *Api) addBlockComments(g *gin.Context) {
	id := g.Param("id")

	thread, err, code := getBlockThread(a.Node, id)
	if err != nil {
		sendError(g, err, code)
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

	hash, err := thread.AddComment(id, args[0])
	if err != nil {
		a.abort500(g, err)
		return
	}

	comment, err := a.Node.Comment(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	pbJSON(g, http.StatusCreated, comment)
}

// lsBlockComments godoc
// @Summary List comments
// @Description Lists comments on a thread block
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.CommentList "comments"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks/{id}/comments [get]
func (a *Api) lsBlockComments(g *gin.Context) {
	id := g.Param("id")

	comments, err := a.Node.Comments(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	pbJSON(g, http.StatusOK, comments)
}

// getBlocks godoc
// @Summary Get thread comment
// @Description Gets a thread comment by block ID
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.Comment "comment"
// @Failure 400 {string} string "Bad Request"
// @Router /blocks/{id}/comment [get]
func (a *Api) getBlockComment(g *gin.Context) {
	info, err := a.Node.Comment(g.Param("id"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}

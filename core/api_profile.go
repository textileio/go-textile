package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getProfile godoc
// @Summary Get public profile
// @Description Gets the local node's public profile
// @Tags profile
// @Produce application/json
// @Success 200 {object} pb.Peer "peer"
// @Failure 400 {string} string "Bad Request"
// @Router /profile [get]
func (a *api) getProfile(g *gin.Context) {
	profile := a.node.Profile()
	if profile == nil {
		g.String(http.StatusBadRequest, "profile is not set")
		return
	}
	pbJSON(g, http.StatusOK, profile)
}

// setName godoc
// @Summary Set display name
// @Description Sets public profile display name to given string
// @Tags profile
// @Produce text/plain
// @Param X-Textile-Args header string true "name"
// @Success 201 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /profile/name [post]
func (a *api) setName(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing name")
		return
	}
	if err := a.node.SetName(args[0]); err != nil {
		a.abort500(g, err)
		return
	}

	a.node.FlushCafes()

	g.JSON(http.StatusCreated, "ok")
}

// setAvatar godoc
// @Summary Set avatar
// @Description Forces local node to update avatar image to latest image added to 'account' thread
// @Tags profile
// @Produce text/plain
// @Success 201 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /profile/avatar [post]
func (a *api) setAvatar(g *gin.Context) {
	if err := a.node.SetAvatar(); err != nil {
		a.abort500(g, err)
		return
	}

	a.node.FlushCafes()

	g.JSON(http.StatusCreated, "ok")
}

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
// @Success 200 {object} pb.Contact "contact"
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

// setUsername godoc
// @Summary Set username
// @Description Sets public profile username to given string
// @Tags profile
// @Produce text/plain
// @Param X-Textile-Args header string true "username"
// @Success 201 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /profile/username [post]
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

// setAvatar godoc
// @Summary Set avatar
// @Description Sets public profile avatar by specifying an existing image file hash
// @Tags profile
// @Produce text/plain
// @Param X-Textile-Args header string true "hash"
// @Success 201 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /profile/avatar [post]
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

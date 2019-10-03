package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// accountGet godoc
// @Summary Show account contact
// @Description Shows the local peer's account info as a contact
// @Tags account
// @Produce application/json
// @Success 200 {object} pb.Contact "contact"
// @Failure 400 {string} string "Bad Request"
// @Router /account [get]
func (a *Api) accountGet(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.Node.AccountContact())
}

// accountSeed godoc
// @Summary Show account seed
// @Description Shows the local peer's account seed
// @Tags account
// @Produce text/plain
// @Success 200 {string} string "seed"
// @Router /account/seed [get]
func (a *Api) accountSeed(g *gin.Context) {
	g.String(http.StatusOK, a.Node.Account().Seed())
}

// accountAddress godoc
// @Summary Show account address
// @Description Shows the local peer's account address
// @Tags account
// @Produce text/plain
// @Success 200 {string} string "address"
// @Router /account/address [get]
func (a *Api) accountAddress(g *gin.Context) {
	g.String(http.StatusOK, a.Node.Account().Address())
}

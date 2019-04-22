package core

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
func (a *api) accountGet(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.node.AccountContact())
}

// accountSeed godoc
// @Summary Show account seed
// @Description Shows the local peer's account seed
// @Tags account
// @Produce text/plain
// @Success 200 {string} string "seed"
// @Router /account/seed [get]
func (a *api) accountSeed(g *gin.Context) {
	g.String(http.StatusOK, a.node.account.Seed())
}

// accountAddress godoc
// @Summary Show account address
// @Description Shows the local peer's account address
// @Tags account
// @Produce text/plain
// @Success 200 {string} string "address"
// @Router /account/address [get]
func (a *api) accountAddress(g *gin.Context) {
	g.String(http.StatusOK, a.node.account.Address())
}

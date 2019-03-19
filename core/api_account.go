package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// accountAddress godoc
// @Summary Show account address
// @Description Shows the local node's account address
// @Tags account
// @Produce text/plain
// @Success 200 {string} string "address"
// @Router /account/address [get]
func (a *api) accountAddress(g *gin.Context) {
	g.String(http.StatusOK, a.node.account.Address())
}

// accountContact godoc
// @Summary Show own contact
// @Description Shows own contact
// @Tags account
// @Produce application/json
// @Success 200 {object} pb.Contact "contact"
// @Failure 400 {string} string "Bad Request"
// @Router /account/peers [get]
func (a *api) accountContact(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.node.AccountContact())
}

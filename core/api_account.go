package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) accountAddress(g *gin.Context) {
	g.String(http.StatusOK, a.node.account.Address())
}

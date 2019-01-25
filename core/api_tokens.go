package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58/base58"
)

func (a *api) addTokens(g *gin.Context) {
	token, err := a.node.CreateCafeToken()
	if err != nil {
		a.abort500(g, err)
		return
	}

	result := make(map[string]string)
	result["id"] = token.Id
	result["created"] = token.Created
	result["token"] = base58.FastBase58Encoding(token.Token)

	g.JSON(http.StatusCreated, result)
}

func (a *api) lsTokens(g *gin.Context) {
	tokens, err := a.node.CafeTokens()
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(tokens) == 0 {
		tokens = make([]*TokenInfo, 0)
	}
	g.JSON(http.StatusOK, tokens)
}

// func (a *api) getCafes(g *gin.Context) {
// 	id := g.Param("id")
// 	session, err := a.node.CafeSession(id)
// 	if err != nil {
// 		a.abort500(g, err)
// 		return
// 	}
// 	if session == nil {
// 		g.String(http.StatusNotFound, "cafe not found")
// 		return
// 	}
// 	g.JSON(http.StatusOK, session)
// }

// func (a *api) rmCafes(g *gin.Context) {
// 	id := g.Param("id")
// 	if err := a.node.DeregisterCafe(id); err != nil {
// 		a.abort500(g, err)
// 		return
// 	}
// 	g.String(http.StatusOK, "ok")
// }

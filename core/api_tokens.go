package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) createTokens(g *gin.Context) {
	token, err := a.node.CreateCafeToken()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusCreated, token)
}

func (a *api) lsTokens(g *gin.Context) {
	tokens, err := a.node.CafeDevTokens()
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(tokens) == 0 {
		tokens = make([]string, 0)
	}
	g.JSON(http.StatusOK, tokens)
}

func (a *api) compareTokens(g *gin.Context) {
	token := g.Param("id")
	ok, err := a.node.CompareCafeDevToken(token)
	if err != nil || !ok {
		g.String(http.StatusUnauthorized, "invlaid credentials")
		return
	}
	g.String(http.StatusOK, "ok")
}

func (a *api) rmTokens(g *gin.Context) {
	token := g.Param("id")
	if err := a.node.RemoveCafeDevToken(token); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

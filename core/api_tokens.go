package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) createTokens(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	token, err := a.node.CreateCafeToken(opts["token"], opts["store"] == "true")
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusCreated, token)
}

func (a *api) lsTokens(g *gin.Context) {
	tokens, err := a.node.CafeTokens()
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(tokens) == 0 {
		tokens = make([]string, 0)
	}
	g.JSON(http.StatusOK, tokens)
}

func (a *api) validateTokens(g *gin.Context) {
	token := g.Param("id")
	ok, err := a.node.ValidateCafeToken(token)
	if err != nil || !ok {
		g.String(http.StatusUnauthorized, "invalid credentials")
		return
	}
	g.String(http.StatusOK, "ok")
}

func (a *api) rmTokens(g *gin.Context) {
	token := g.Param("id")
	if err := a.node.RemoveCafeToken(token); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

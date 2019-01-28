package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/repo"
)

func (a *api) createTokens(g *gin.Context) {
	token, err := a.node.CreateCafeToken()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusCreated, token)
}

func (a *api) lsTokens(g *gin.Context) {
	tokens, err := a.node.CafeDevTokens()
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(tokens) == 0 {
		tokens = make([]repo.CafeDevToken, 0)
	} else {
		for _, token := range tokens {
			token.Token = base58.FastBase58Encoding([]byte(token.Token))
		}
	}
	g.JSON(http.StatusOK, tokens)
}

func (a *api) compareTokens(g *gin.Context) {
	id := g.Param("id")

	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing dev token")
		return
	}
	ok, err := a.node.CompareCafeDevToken(id, args[0])
	if err != nil {
		a.abort500(g, err)
		return
	}
	if !ok {
		g.String(http.StatusUnauthorized, "invlaid credentials")
		return
	}
	g.String(http.StatusOK, "ok")
}

func (a *api) rmTokens(g *gin.Context) {
	id := g.Param("id")
	if err := a.node.RemoveCafeDevToken(id); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

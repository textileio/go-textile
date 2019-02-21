package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// createTokens godoc
// @Summary Create an access token
// @Description Generates an access token (44 random bytes) and saves a bcrypt hashed version for
// @Description future lookup. The response contains a base58 encoded version of the random bytes
// @Description token. If the 'store' option is set to false, the token is generated, but not
// @Description stored in the local Cafe db. Alternatively, an existing token can be added using
// @Description by specifying the 'token' option.
// @Description Tokens allow other peers to register with a Cafe peer.
// @Tags tokens
// @Produce application/json
// @Param X-Textile-Opts header string false "token: Use existing token, rather than creating a new one, store: Whether to store the added/generated token to the local db" default(token=,store="true")
// @Success 201 {string} string "token"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tokens [post]
func (a *api) createTokens(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	token, err := a.node.CreateCafeToken(opts["token"], opts["store"] == "true")
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	g.String(http.StatusCreated, token)
}

// lsTokens godoc
// @Summary List local tokens
// @Description List info about all stored cafe tokens
// @Tags tokens
// @Produce application/json
// @Success 200 {array} string "tokens"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tokens [get]
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

// validateTokens godoc
// @Summary Check token validity
// @Description Check validity of existing cafe access token
// @Tags tokens
// @Produce application/json
// @Param token path string true "invite id"
// @Success 200 {string} string "ok"
// @Failure 401 {string} string "Unauthorized"
// @Router /tokens/{id} [get]
func (a *api) validateTokens(g *gin.Context) {
	token := g.Param("token")
	ok, err := a.node.ValidateCafeToken(token)
	if err != nil || !ok {
		g.String(http.StatusUnauthorized, "invalid credentials")
		return
	}
	g.String(http.StatusOK, "ok")
}

// rmTokens godoc
// @Summary Removes a cafe token
// @Description Removes an existing cafe token
// @Tags tokens
// @Produce application/json
// @Param token path string true "token"
// @Success 200 {string} string "ok"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tokens/{id} [delete]
func (a *api) rmTokens(g *gin.Context) {
	token := g.Param("token")
	if err := a.node.RemoveCafeToken(token); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

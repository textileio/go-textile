package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// addCafes godoc
// @Summary Register with a Cafe
// @Description Registers with a cafe and saves an expiring service session token. An access
// @Description token is required to register, and should be obtained separately from the target
// @Description Cafe
// @Tags cafes
// @Produce application/json
// @Param X-Textile-Args header string true "cafe id"
// @Param X-Textile-Opts header string false "token: An access token supplied by the Cafe" default(token=)
// @Success 201 {object} pb.CafeSession "cafe session"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /cafes [post]
func (a *Api) addCafes(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing cafe id")
		return
	}

	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	token := opts["token"]
	if token == "" {
		g.String(http.StatusBadRequest, "missing access token")
		return
	}

	session, err := a.Node.RegisterCafe(args[0], token)
	if err != nil {
		a.abort500(g, err)
		return
	}

	a.Node.FlushCafes()

	pbJSON(g, http.StatusCreated, session)
}

// lsCafes godoc
// @Summary List info about all active cafe sessions
// @Description List info about all active cafe sessions. Cafes are other peers on the network
// @Description who offer pinning, backup, and inbox services
// @Tags cafes
// @Produce application/json
// @Success 200 {object} pb.CafeSessionList "cafe sessions"
// @Failure 500 {string} string "Internal Server Error"
// @Router /cafes [get]
func (a *Api) lsCafes(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.Node.CafeSessions())
}

// getCafes godoc
// @Summary Gets and displays info about a cafe session
// @Description Gets and displays info about a cafe session. Cafes are other peers on the network
// @Description who offer pinning, backup, and inbox services
// @Tags cafes
// @Produce application/json
// @Param id path string true "cafe id"
// @Success 200 {object} pb.CafeSession "cafe session"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /cafes/{id} [get]
func (a *Api) getCafes(g *gin.Context) {
	id := g.Param("id")

	session, err := a.Node.CafeSession(id)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if session == nil {
		g.String(http.StatusNotFound, "cafe not found")
		return
	}

	pbJSON(g, http.StatusOK, session)
}

// rmCafes godoc
// @Summary Deregisters a cafe
// @Description Deregisters with a cafe (content will expire based on the cafe's service rules)
// @Tags cafes
// @Param id path string true "cafe id"
// @Success 204 {string} string "ok"
// @Failure 500 {string} string "Internal Server Error"
// @Router /cafes/{id} [delete]
func (a *Api) rmCafes(g *gin.Context) {
	id := g.Param("id")

	err := a.Node.DeregisterCafe(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	a.Node.FlushCafes()

	g.Status(http.StatusNoContent)
}

// checkCafeMessages godoc
// @Summary Check for messages at all cafes
// @Description Check for messages at all cafes. New messages are downloaded and processed
// @Description opportunistically.
// @Tags cafes
// @Produce text/plain
// @Success 200 {string} string "ok"
// @Failure 500 {string} string "Internal Server Error"
// @Router /cafes/messages [post]
func (a *Api) checkCafeMessages(g *gin.Context) {
	err := a.Node.CheckCafeMessages()
	if err != nil {
		a.abort500(g, err)
		return
	}

	a.Node.FlushCafes()

	g.String(http.StatusOK, "ok")
}

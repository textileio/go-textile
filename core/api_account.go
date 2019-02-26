package core

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/pb"
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

// accountPeers godoc
// @Summary Show account peers
// @Description Shows all known account peers
// @Tags account
// @Produce application/json
// @Success 200 {object} pb.ContactList "peers"
// @Failure 400 {string} string "Bad Request"
// @Router /account/peers [get]
func (a *api) accountPeers(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.node.AccountPeers())
}

// accountPeers godoc
// @Summary Show account peers
// @Description Shows all known account peers
// @Tags account
// @Produce application/json
// @Param X-Textile-Opts header string false "wait: Stops searching after 'wait' seconds have elapsed (max 10s), events: Whether to emit Server-Sent Events (SSEvent) or plain JSON" default(wait=5,events="false")
// @Success 200 {object} pb.QueryResult "results stream"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /account/backups [post]
func (a *api) accountBackups(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	wait, err := strconv.Atoi(opts["wait"])
	if err != nil {
		wait = 5
	}

	query := &pb.ThreadBackupQuery{
		Address: a.node.account.Address(),
	}
	options := &pb.QueryOptions{
		Local: false,
		Limit: -1,
		Wait:  int32(wait),
	}

	resCh, errCh, cancel, err := a.node.FindThreadBackups(query, options)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	handleSearchStream(g, resCh, errCh, cancel, opts["events"] == "true")
}

package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/pb"
)

// createThreadSnapshots godoc
// @Summary Create thread snapshots
// @Description Snapshots all threads and pushes to registered cafes
// @Tags threads
// @Produce application/json
// @Success 201 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /snapshots [post]
func (a *Api) createThreadSnapshots(g *gin.Context) {
	if err := a.Node.SnapshotThreads(); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	g.String(http.StatusCreated, "ok")
}

// searchThreadSnapshots godoc
// @Summary Search for thread snapshots
// @Description Searches the network for thread snapshots
// @Tags threads
// @Produce application/json
// @Param X-Textile-Opts header string false "wait: Stops searching after 'wait' seconds have elapsed (max 30s), events: Whether to emit Server-Sent Events (SSEvent) or plain JSON" default(wait=5,events="false")
// @Success 200 {object} pb.QueryResult "results stream"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /snapshots/search [post]
func (a *Api) searchThreadSnapshots(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	wait, err := strconv.Atoi(opts["wait"])
	if err != nil {
		wait = 5
	}

	query := &pb.ThreadSnapshotQuery{
		Address: a.Node.Account().Address(),
	}
	options := &pb.QueryOptions{
		Limit: -1,
		Wait:  int32(wait),
	}

	resCh, errCh, cancel, err := a.Node.SearchThreadSnapshots(query, options)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	handleSearchStream(g, resCh, errCh, cancel, opts["events"] == "true")
}

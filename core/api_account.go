package core

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/pb"
)

func (a *api) accountAddress(g *gin.Context) {
	g.String(http.StatusOK, a.node.account.Address())
}

func (a *api) accountPeers(g *gin.Context) {
	peers, err := a.node.AccountPeers()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, peers)
}

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

func (a *api) accountSync(g *gin.Context) {
	backup := new(pb.Thread)
	if err := g.BindJSON(&backup); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if backup == nil {
		g.String(http.StatusBadRequest, "missing backup")
		return
	}
	if backup.Id == "" || len(backup.Sk) == 0 {
		g.String(http.StatusBadRequest, "invalid backup")
		return
	}

	fmt.Println(backup)
	//
	//if err := a.node.ApplyThreadBackup(backup); err != nil {
	//	g.String(http.StatusBadRequest, err.Error())
	//	return
	//}

	g.String(http.StatusOK, "ok")
}

package core

import (
	"net/http"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/gin-gonic/gin"
)

func (a *api) ping(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing peer id")
		return
	}

	pid, err := peer.IDB58Decode(args[0])
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	status, err := a.node.Ping(pid)
	if err != nil {
		a.abort500(g, err)
		return
	}

	g.String(http.StatusOK, string(status))
}

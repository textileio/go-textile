package core

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/repo"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"net/http"
)

func (a *api) addThreads(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing thread name")
		return
	}
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	ttype := repo.OpenThread
	if opts["type"] != "" {
		var err error
		ttype, err = repo.ThreadTypeFromString(opts["type"])
		if err != nil {
			g.String(http.StatusBadRequest, "invalid thread type")
			return
		}
	}
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		a.abort500(g, err)
		return
	}
	thrd, err := a.node.AddThread(args[0], sk, ttype, true)
	if err != nil {
		a.abort500(g, err)
		return
	}
	info, err := thrd.Info()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusCreated, info)
}

func (a *api) lsThreads(g *gin.Context) {
	infos := make([]*ThreadInfo, 0)
	for _, thrd := range a.node.Threads() {
		info, err := thrd.Info()
		if err != nil {
			a.abort500(g, err)
			return
		}
		infos = append(infos, info)
	}
	g.JSON(http.StatusOK, infos)
}

func (a *api) getThreads(g *gin.Context) {
	id := g.Param("id")
	thrd := a.node.Thread(id)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}
	info, err := thrd.Info()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusOK, info)
}

func (a *api) rmThreads(g *gin.Context) {
	id := g.Param("id")
	thrd := a.node.Thread(id)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}
	if _, err := a.node.RemoveThread(id); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

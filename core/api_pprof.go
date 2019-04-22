package core

import (
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *api) pprofIndex(g *gin.Context) {
	pprof.Index(g.Writer, g.Request)
}

func (a *api) pprofCmdline(g *gin.Context) {
	pprof.Cmdline(g.Writer, g.Request)
}

func (a *api) pprofProfile(g *gin.Context) {
	pprof.Profile(g.Writer, g.Request)
}

func (a *api) pprofSymbol(g *gin.Context) {
	pprof.Symbol(g.Writer, g.Request)
}

func (a *api) pprofTrace(g *gin.Context) {
	pprof.Trace(g.Writer, g.Request)
}

func (a *api) pprofGoroutine(g *gin.Context) {
	pprof.Handler("goroutine")
}

func (a *api) pprofHeap(g *gin.Context) {
	pprof.Handler("heap")
}

func (a *api) pprofThreadCreate(g *gin.Context) {
	pprof.Handler("threadcreate")
}

func (a *api) pprofBlock(g *gin.Context) {
	pprof.Handler("block")
}

func (a *api) mutexFractionOption(g *gin.Context) {
	str := g.DefaultQuery("fraction", "0")
	fraction, err := strconv.Atoi(str)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Infof("Setting MutexProfileFraction to %d", fraction)
	runtime.SetMutexProfileFraction(fraction)

	g.Status(http.StatusOK)
}

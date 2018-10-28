package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// apiVersion is the api version
const apiVersion = "v0"

// apiHost is the instance used by the daemon
var apiHost *api

// api is a limited HTTP REST API for the cmd tool
type api struct {
	addr   string
	server *http.Server
	node   *Textile
}

// StartApi starts the host instance
func (t *Textile) StartApi(addr string) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = t.writer
	apiHost = &api{addr: addr, node: t}
	apiHost.Start()
}

// StopApi starts the host instance
func (t *Textile) StopApi() error {
	return apiHost.Stop()
}

// ApiAddr returns the api address
func (t *Textile) ApiAddr() string {
	if apiHost == nil {
		return ""
	}
	return apiHost.addr
}

// Start starts the http api
func (a *api) Start() {
	// setup router
	router := gin.Default()
	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"cafe_version": apiVersion,
			"node_version": Version,
		})
	})
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})

	// v0 routes
	v0 := router.Group("/api/v0")
	{
		v0.GET("/peer", a.peer)
		v0.GET("/address", a.address)
	}
	a.server = &http.Server{
		Addr:    a.addr,
		Handler: router,
	}

	// start listening
	errc := make(chan error)
	go func() {
		errc <- a.server.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err != http.ErrServerClosed {
					log.Errorf("api error: %s", err)
				}
				if !ok {
					log.Info("api was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("api listening at %s\n", a.server.Addr)
}

// Stop stops the http api
func (a *api) Stop() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down api: %s", err)
		return err
	}
	cancel()
	return nil
}

func (a *api) peer(g *gin.Context) {
	pid, err := a.node.PeerId()
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.String(http.StatusOK, pid.Pretty())
}

func (a *api) address(g *gin.Context) {
	addr, err := a.node.Address()
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.String(http.StatusOK, addr)
}

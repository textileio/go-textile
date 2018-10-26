package cafe

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"net/http"
)

var log = logging.MustGetLogger("cafe")

// Version is the api version
const Version = "v0"

// Host is the instance used by the daemon
var Host *Cafe

// Cafe is a limited HTTP API to the cafe service
type Cafe struct {
	Ipfs        func() *core.IpfsNode
	NodeVersion string
	Protocol    protocol.ID
	server      *http.Server
}

// Start starts the cafe api
func (c *Cafe) Start(addr string) {
	// setup router
	router := gin.Default()
	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"api_version":  Version,
			"node_version": c.NodeVersion,
		})
	})
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})

	// v0 routes
	v0 := router.Group("/api/v0")
	{
		v0.POST("/pin", c.pin)
	}
	c.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// start listening
	errc := make(chan error)
	go func() {
		errc <- c.server.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err != http.ErrServerClosed {
					log.Errorf("cafe error: %s", err)
				}
				if !ok {
					log.Info("cafe was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("cafe listening at %s\n", c.server.Addr)
}

// Stop stops the cafe api
func (c *Cafe) Stop() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := c.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down cafe: %s", err)
		return err
	}
	cancel()
	return nil
}

// Cafe returns the cafe address
func (c *Cafe) Addr() string {
	return c.server.Addr
}

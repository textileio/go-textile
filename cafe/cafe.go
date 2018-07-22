package cafe

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	cdao "github.com/textileio/textile-go/cafe/dao"
	"github.com/textileio/textile-go/cafe/middleware"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"net/http"
)

var log = logging.MustGetLogger("core")

const Version = "v0"

var Host *Cafe

type Cafe struct {
	Ipfs        func() *core.IpfsNode
	Dao         *cdao.DAO
	TokenSecret string
	ReferralKey string
	NodeVersion string
	server      *http.Server
}

// Start starts the cafe api
func (c *Cafe) Start(addr string) {
	// init db connection
	cdao.Dao = c.Dao
	cdao.Dao.Connect()
	cdao.Dao.Index()

	// setup router
	router := gin.Default()
	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"cafe_version": Version,
			"node_version": c.NodeVersion,
		})
	})
	router.GET("/health", c.health)

	// api routes
	v0 := router.Group("/api/v0")
	v0.Use(middleware.Auth(c.TokenSecret))
	{
		v0.PUT("/users", c.signUp)
		v0.POST("/users", c.signIn)
		//v0.POST("/pin", c.pin)

		v0.POST("/referrals", c.createReferral)
		v0.GET("/referrals", c.listReferrals)
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

// GetCafeAddress returns the cafe address
func (c *Cafe) GetCafeAddress() string {
	return c.server.Addr
}

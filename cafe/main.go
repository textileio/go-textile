package cafe

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/cafe/dao"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"net/http"
	"time"
)

var log = logging.MustGetLogger("cafe")

const Version = "v1"
const oneMonth = time.Hour * 24 * 7 * 4

var Host *Cafe

type Cafe struct {
	Ipfs        func() *core.IpfsNode
	Dao         *dao.DAO
	TokenSecret string
	ReferralKey string
	NodeVersion string
	server      *http.Server
}

// Start starts the cafe api
func (c *Cafe) Start(addr string) {
	// init db connection
	dao.Dao = c.Dao
	dao.Dao.Connect()
	dao.Dao.Index()

	// setup router
	router := gin.Default()
	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"cafe_version": Version,
			"node_version": c.NodeVersion,
		})
	})
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})

	// v0 routes
	v0 := router.Group("/api/v0")
	{
		v0.PUT("/users", c.signUpUser)
		v0.POST("/users", c.signInUser)
		v0.POST("/referrals", c.createReferral)
		v0.GET("/referrals", c.listReferrals)
		v0.POST("/tokens", c.authSession, c.refreshSession)
		v0.POST("/pin", c.authSession, c.pin)
	}

	// v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/profiles/challenge", c.profileChallenge)
		v1.PUT("/profiles", c.registerProfile)
		v1.POST("/profiles", c.loginProfile)
		v1.POST("/referrals", c.createReferral)
		v1.GET("/referrals", c.listReferrals)
		v1.POST("/tokens", c.authSession, c.refreshSession)
		v1.POST("/pin", c.authSession, c.pin)
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

package core

import (
	"context"
	"net/http"
	"strings"
	"time"

	njwt "github.com/dgrijalva/jwt-go"
	limit "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/jwt"
	"github.com/textileio/go-textile/pb"
)

// cafeApiVersion is the cafe api version
const cafeApiVersion = "v0"

// cafeApiHost is the instance used by the core instance
var cafeApiHost *cafeApi

// cafeApi is a limited HTTP API to the cafe service
type cafeApi struct {
	addr   string
	server *http.Server
	node   *Textile
}

// CafeApiAddr returns the cafe api address
func (t *Textile) CafeApiAddr() string {
	if cafeApiHost == nil {
		return ""
	}
	return cafeApiHost.addr
}

// startCafeApi starts the host instance
func (t *Textile) startCafeApi(addr string) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = t.writer
	cafeApiHost = &cafeApi{addr: addr, node: t}
	cafeApiHost.start()
}

// StopCafeApi stops the host instance
func (t *Textile) stopCafeApi() error {
	if cafeApiHost == nil {
		return nil
	}
	return cafeApiHost.stop()
}

// CafeInfo returns info about this cafe
func (t *Textile) CafeInfo() *pb.Cafe {
	return t.cafe.info
}

// start starts the cafe api
func (c *cafeApi) start() {
	router := gin.Default()
	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, c.node.CafeInfo())
	})
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})

	conf := c.node.Config()
	if conf.Cafe.Host.SizeLimit > 0 {
		router.Use(limit.RequestSizeLimiter(conf.Cafe.Host.SizeLimit))
	}

	// v0 routes
	v0 := router.Group("/cafe/v0")
	{
		v0.POST("/pin", c.pin)
		v0.POST("/service", c.service)
	}

	// v1 routes
	v1 := router.Group("/api/v1")

	store := v1.Group("/store", c.validate)
	{
		store.PUT("/:cid", c.store)
		store.DELETE("/:cid", c.unstore)
	}

	threads := v1.Group("/threads", c.validate)
	{
		threads.PUT("/:id", c.storeThread)
		threads.DELETE("/:id", c.unstoreThread)
	}

	inbox := v1.Group("/inbox")
	{
		inbox.POST("/:pid", c.deliverMessage)
	}

	c.server = &http.Server{
		Addr:    c.addr,
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
					log.Errorf("cafe api error: %s", err)
				}
				if !ok {
					log.Info("cafe api was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("cafe api listening at %s\n", c.server.Addr)
}

// stop stops the cafe api
func (c *cafeApi) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := c.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down cafe api: %s", err)
		return err
	}
	return nil
}

// validate checks if the node is running and validates auth token
func (c *cafeApi) validate(g *gin.Context) {
	if !c.node.Started() {
		g.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"error": "node is stopped",
		})
		return
	}

	// validate request token
	if !c.tokenValid(g) {
		return
	}
}

// validToken aborts the request if the token is invalid
func (c *cafeApi) tokenValid(g *gin.Context) bool {
	auth := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(auth) < 2 {
		g.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedResponse)
		return false
	}
	token := auth[1]

	protocol := string(c.node.cafe.Protocol())
	if err := jwt.Validate(token, c.verifyKeyFunc, false, protocol, nil); err != nil {
		switch err {
		case jwt.ErrNoToken, jwt.ErrExpired:
			g.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedResponse)
		case jwt.ErrInvalid:
			g.AbortWithStatusJSON(http.StatusForbidden, forbiddenResponse)
		}
		return false
	}
	return true
}

// verifyKeyFunc returns the correct key for token verification
func (c *cafeApi) verifyKeyFunc(token *njwt.Token) (interface{}, error) {
	return c.node.Ipfs().PrivateKey.GetPublic(), nil
}

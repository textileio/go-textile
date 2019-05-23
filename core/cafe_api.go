package core

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"sync"
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

// bodyPool handles service payloads
var bodyPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

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
		v0.POST("/pin", c.validateToken, c.pin)
		v0.POST("/service", c.service)
	}

	// v1 routes
	v1 := router.Group("/api/v1")

	store := v1.Group("/store", c.validateToken)
	{
		store.PUT("", c.store)
		store.DELETE("/:cid", c.unstore)
	}

	threads := v1.Group("/threads", c.validateToken)
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

// validateToken aborts the request if the token is invalid
func (c *cafeApi) validateToken(g *gin.Context) {
	auth := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(auth) < 2 {
		c.abort(g, http.StatusUnauthorized, nil)
		return
	}
	token := auth[1]

	protocol := string(c.node.cafe.Protocol())
	if err := jwt.Validate(token, c.verifyKeyFunc, false, protocol, nil); err != nil {
		switch err {
		case jwt.ErrNoToken, jwt.ErrExpired:
			c.abort(g, http.StatusUnauthorized, nil)
		case jwt.ErrInvalid:
			c.abort(g, http.StatusForbidden, nil)
		}
	}
}

// verifyKeyFunc returns the correct key for token verification
func (c *cafeApi) verifyKeyFunc(token *njwt.Token) (interface{}, error) {
	return c.node.Ipfs().PrivateKey.GetPublic(), nil
}

// abort aborts the request with the given status code and error
func (c *cafeApi) abort(g *gin.Context, status int, err error) {
	if err != nil {
		g.AbortWithStatusJSON(status, gin.H{
			"error": err.Error(),
		})
	} else {
		g.AbortWithStatus(status)
	}
}

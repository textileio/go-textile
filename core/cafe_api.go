package core

import (
	"bytes"
	"context"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	njwt "github.com/dgrijalva/jwt-go"
	limit "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/go-textile/jwt"
	"github.com/textileio/go-textile/pb"
	"golang.org/x/crypto/bcrypt"
)

// CafeApiVersion is the cafe api version
const CafeApiVersion = "v1"

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

// CafeInfo returns info about this cafe
func (t *Textile) CafeInfo() *pb.Cafe {
	return t.cafe.info
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
	router.POST("/api/v0/search", func(g *gin.Context) {
		g.Redirect(http.StatusPermanentRedirect, "/api/v1/search")
	})

	// v1 routes
	v1 := router.Group("/api/v1")

	sessions := v1.Group("/sessions")
	{
		sessions.GET("/challenge", c.validateChallengeToken, c.getSessionChallenge)
		sessions.POST("/:pid", c.validateChallengeToken, c.createSession)
		sessions.POST("/:pid/refresh", c.validateToken, c.refreshSession)
		sessions.DELETE("/:pid", c.validateToken, c.deleteSession)
	}

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
		inbox.GET("/:pid", c.validateToken, c.checkMessages)
		inbox.DELETE("/:pid", c.validateToken, c.deleteMessages)
		inbox.POST("/:from/:to", c.deliverMessage)
	}

	search := v1.Group("/search", c.validateToken)
	{
		search.POST("", c.search)
	}

	// Enables bots on cafes
	bots := v1.Group("/bots")
	{
		bots.PUT("/id/:id", c.reverseProxyBotAPI("PUT"))
		bots.POST("/id/:id", c.reverseProxyBotAPI("POST"))
		bots.GET("/id/:id", c.reverseProxyBotAPI("GET"))
		bots.DELETE("/id/:id", c.reverseProxyBotAPI("DELETE"))
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
		log.Warning("missing token")
		c.abort(g, http.StatusUnauthorized, nil)
		return
	}
	token := auth[1]

	var subject *string
	pid, err := peer.IDB58Decode(g.Param("pid"))
	if err == nil {
		tmp := pid.Pretty()
		subject = &tmp
	}

	protocol := string(c.node.cafe.Protocol())
	refreshing := strings.Contains(g.Request.URL.Path, "refresh")
	claims, err := jwt.Validate(token, c.verifyKeyFunc, refreshing, protocol, subject)
	if err != nil {
		switch err {
		case jwt.ErrNoToken, jwt.ErrExpired:
			log.Warning("bad or expired token")
			c.abort(g, http.StatusUnauthorized, nil)
		case jwt.ErrInvalid:
			log.Warning("invalid token")
			c.abort(g, http.StatusForbidden, nil)
		}
		return
	}

	g.Set("from", claims.Subject)
	g.Set("token", token)
}

// validateChallengeToken aborts the request if the token is invalid
func (c *cafeApi) validateChallengeToken(g *gin.Context) {
	auth := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(auth) < 2 {
		log.Warning("missing token pass")
		c.abort(g, http.StatusUnauthorized, nil)
		return
	}
	token := auth[1]

	// does the provided token match?
	// dev tokens are actually base58(id+token)
	plainBytes, err := base58.FastBase58Decoding(token)
	if err != nil || len(plainBytes) < 44 {
		log.Warning("error decoding token pass")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	encodedToken := c.node.datastore.CafeTokens().Get(hex.EncodeToString(plainBytes[:12]))
	if encodedToken == nil {
		log.Warning("token not found")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	err = bcrypt.CompareHashAndPassword(encodedToken.Value, plainBytes[12:])
	if err != nil {
		log.Warning("bad token pass")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	g.Set("token", encodedToken.Id)
}

// verifyKeyFunc returns the correct key for token verification
func (c *cafeApi) verifyKeyFunc(token *njwt.Token) (interface{}, error) {
	return c.node.Ipfs().PrivateKey.GetPublic(), nil
}

// CafeError represents a cafe request error
type CafeError struct {
	Error string `json:"error"`
}

// abort aborts the request with the given status code and error
func (c *cafeApi) abort(g *gin.Context, status int, err error) {
	if err != nil {
		g.AbortWithStatusJSON(status, CafeError{
			Error: err.Error(),
		})
	} else {
		g.AbortWithStatus(status)
	}
}

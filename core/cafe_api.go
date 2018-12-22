package core

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	uio "gx/ipfs/QmfB3oNXGGq9S4B2a9YeCajoATms3Zw2VvDm8fK7VeLSV8/go-unixfs/io"

	njwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/jwt"
	"github.com/textileio/textile-go/pb"
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

// start starts the cafe api
func (c *cafeApi) start() {
	// setup router
	router := gin.Default()
	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"cafe_version": cafeApiVersion,
			"node_version": Version,
		})
	})
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})

	// v0 routes
	v0 := router.Group("/cafe/v0")
	{
		v0.POST("/pin", c.pin)
		v0.POST("/service", c.service)
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

// PinResponse is the json response from a pin request
type PinResponse struct {
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// forbiddenResponse is used for bad tokens
var forbiddenResponse = PinResponse{
	Error: errForbidden,
}

// unauthorizedResponse is used when a token is expired or not present
var unauthorizedResponse = PinResponse{
	Error: errUnauthorized,
}

// pin take raw data or a tarball and pins it to the local ipfs node.
// request must be authenticated with a token
func (c *cafeApi) pin(g *gin.Context) {
	if !c.node.Started() {
		g.AbortWithStatusJSON(http.StatusInternalServerError, PinResponse{
			Error: "node is stopped",
		})
		return
	}

	// validate request token
	if !c.tokenValid(g) {
		return
	}

	// handle based on content type
	var id cid.Cid
	cType := g.Request.Header.Get("Content-Type")
	switch cType {
	case "application/gzip":
		dirb := uio.NewDirectory(c.node.Ipfs().DAG)

		gr, err := gzip.NewReader(g.Request.Body)
		if err != nil {
			log.Errorf("error creating gzip reader %s", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tr := tar.NewReader(gr)

		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Errorf("error getting tar next %s", err)
				g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			switch header.Typeflag {
			case tar.TypeDir:
				log.Error("got nested directory, aborting")
				g.JSON(http.StatusBadRequest, gin.H{"error": "directories are not supported"})
				return
			case tar.TypeReg:
				if _, err := ipfs.AddDataToDirectory(c.node.Ipfs(), dirb, header.Name, tr); err != nil {
					log.Errorf("error adding file to dir %s", err)
					g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			default:
				continue
			}
		}

		// pin the directory
		dir, err := dirb.GetNode()
		if err != nil {
			log.Errorf("error creating dir node %s", err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := ipfs.PinNode(c.node.Ipfs(), dir, true); err != nil {
			log.Errorf("error pinning dir node %s", err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id = dir.Cid()

	case "application/octet-stream":
		idp, err := ipfs.AddData(c.node.Ipfs(), g.Request.Body, true)
		if err != nil {
			log.Errorf("error pinning raw body %s", err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id = *idp
	default:
		log.Errorf("got bad content type %s", cType)
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid content-type"})
		return
	}
	hash := id.Hash().B58String()

	log.Debugf("pinned request with content type %s: %s", cType, hash)

	// ship it
	g.JSON(http.StatusCreated, PinResponse{
		Id: hash,
	})
}

func (c *cafeApi) service(g *gin.Context) {
	if !c.node.Online() {
		g.AbortWithStatusJSON(http.StatusInternalServerError, PinResponse{
			Error: "node is offline",
		})
		return
	}

	// validate request token
	if !c.tokenValid(g) {
		return
	}

	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// parse body as a service envelope
	pmes := new(pb.Envelope)
	if err := proto.Unmarshal(body, pmes); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	peerId := g.Request.Header.Get("X-Textile-Peer")
	if peerId == "" {
		g.AbortWithStatusJSON(http.StatusBadRequest, "missing peer ID")
		return
	}
	mPeer, err := peer.IDB58Decode(peerId)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := c.node.cafeService.service.VerifyEnvelope(pmes, mPeer); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	}

	// handle the message as normal
	log.Debugf("received %s from %s", pmes.Message.Type.String(), mPeer.Pretty())
	rpmes, err := c.node.cafeService.Handle(mPeer, pmes)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if rpmes == nil {
		g.Status(200)
		return
	}

	// send out response msg
	log.Debugf("responding with %s to %s", rpmes.Message.Type.String(), mPeer.Pretty())

	res, err := proto.Marshal(rpmes)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// ship it
	g.Render(200, render.Data{Data: res})
}

// validToken aborts the request if the token is invalid
func (c *cafeApi) tokenValid(g *gin.Context) bool {
	auth := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(auth) < 2 {
		g.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedResponse)
		return false
	}
	token := auth[1]

	protocol := string(c.node.cafeService.Protocol())
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

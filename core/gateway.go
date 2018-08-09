package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/textileio/textile-go/crypto"
	"net/http"
	"strings"
)

// StartGateway starts the gateway
func (t *TextileNode) StartGateway(addr string) {
	// setup router
	router := gin.Default()
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})
	router.GET("/ipfs/:root", gatewayHandler)
	router.GET("/ipfs/:root/*path", gatewayHandler)
	router.GET("/ipns/:root", profileHandler)
	router.GET("/ipns/:root/*path", profileHandler)
	t.gateway = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// start listening
	errc := make(chan error)
	go func() {
		errc <- t.gateway.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err != http.ErrServerClosed {
					log.Errorf("gateway error: %s", err)
				}
				if !ok {
					log.Info("gateway was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("gateway listening at %s\n", t.gateway.Addr)
}

// StopGateway stops the gateway
func (t *TextileNode) StopGateway() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := t.gateway.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down gateway: %s", err)
		return err
	}
	cancel()
	return nil
}

// GetGatewayAddr returns the gateway's address
func (t *TextileNode) GetGatewayAddr() string {
	return t.gateway.Addr
}

// gatewayHandler handles gateway http requests
func gatewayHandler(c *gin.Context) {
	contentPath := c.Param("root") + c.Param("path")

	// look for block id
	blockId, exists := c.GetQuery("block")
	if exists {
		block, err := Node.Wallet.GetBlock(blockId)
		if err != nil {
			log.Errorf("error finding block %s: %s", blockId, err)
			c.Status(404)
			return
		}
		_, thrd := Node.Wallet.GetThread(block.ThreadId)
		if thrd == nil {
			log.Errorf("could not find thread for block: %s", block.Id)
			c.Status(404)
			return
		}
		data, err := thrd.GetBlockData(contentPath, block)
		if err != nil {
			log.Errorf("error decrypting path %s: %s", contentPath, err)
			c.Status(404)
			return
		}
		c.Render(200, render.Data{Data: data})
		return
	}

	// get raw data
	data, err := Node.Wallet.GetDataAtPath(contentPath)
	if err != nil {
		log.Errorf("error getting raw path %s: %s", contentPath, err)
		c.Status(404)
		return
	}

	// if key is provided, try to decrypt the data with it
	key, exists := c.GetQuery("key")
	if exists {
		plain, err := crypto.DecryptAES(data, []byte(key))
		if err != nil {
			log.Errorf("error decrypting %s: %s", contentPath, err)
			c.Status(404)
			return
		}
		c.Render(200, render.Data{Data: plain})
		return
	}

	// lastly, just return the raw bytes (standard gateway)
	c.Render(200, render.Data{Data: data})
}

// profileHandler handles requests for profile info hosted on ipns
// NOTE: avatar is a magic path, will return data behind link at avatar_id
func profileHandler(c *gin.Context) {
	pathp := c.Param("path")
	var isAvatar bool
	if pathp == "/avatar" {
		pathp += "_id"
		isAvatar = true
	}

	pth, err := Node.Wallet.ResolveProfile(c.Param("root"))
	if err != nil {
		log.Errorf("error resolving profile %s: %s", c.Param("root"), err)
		c.Status(404)
		return
	}

	// get data
	contentPath := pth.String() + pathp
	data, err := Node.Wallet.GetDataAtPath(contentPath)
	if err != nil {
		log.Errorf("error getting data at profile path %s: %s", contentPath, err)
		c.Status(404)
		return
	}

	// if this is an avatar request, fetch and return the linked image
	if isAvatar {
		// cheap migration from older formats...
		location := string(data)
		if !strings.HasPrefix(location, "/ipfs") || !strings.Contains(location, "?key=") {
			location = fmt.Sprintf("https://avatars.dicebear.com/v2/identicon/%s.svg", c.Param("root"))
		}
		c.Redirect(302, location)
		return
	}

	c.Render(200, render.Data{Data: data})
}

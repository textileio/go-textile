package gateway

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/crypto"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface"
	"net/http"
	"strings"
)

var log = logging.Logger("gateway")

// Host is the instance used by the daemon
var Host *Gateway

// Gateway is a HTTP API for getting files and links from IPFS
type Gateway struct {
	server *http.Server
}

// Start creates a gateway server
func (g *Gateway) Start(addr string) {
	gin.SetMode(gin.ReleaseMode)
	if core.Node != nil {
		gin.DefaultWriter = core.Node.Writer()
	}

	router := gin.Default()
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})
	router.GET("/ipfs/:root", gatewayHandler)
	router.GET("/ipfs/:root/*path", gatewayHandler)
	router.GET("/ipns/:root", profileHandler)
	router.GET("/ipns/:root/*path", profileHandler)

	g.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	errc := make(chan error)
	go func() {
		errc <- g.server.ListenAndServe()
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
	log.Infof("gateway listening at %s\n", g.server.Addr)
}

// Stop stops the gateway
func (g *Gateway) Stop() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := g.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down gateway: %s", err)
		return err
	}
	cancel()
	return nil
}

// Addr returns the gateway's address
func (g *Gateway) Addr() string {
	return g.server.Addr
}

// gatewayHandler handles gateway http requests
func gatewayHandler(c *gin.Context) {
	contentPath := c.Param("root") + c.Param("path")

	// look for block id
	// NOTE: this only works for the local node, but very useful for desktop
	//blockId, exists := c.GetQuery("block")
	//if exists {
	//	block, err := core.Node.Block(blockId)
	//	if err != nil {
	//		log.Errorf("error finding block %s: %s", blockId, err)
	//		c.Status(404)
	//		return
	//	}
	//	data, err := core.Node.BlockData(contentPath, block)
	//	if err != nil {
	//		log.Errorf("error decrypting path %s: %s", contentPath, err)
	//		c.Status(404)
	//		return
	//	}
	//	c.Render(200, render.Data{Data: data})
	//	return
	//}

	// get data behind path
	data := getDataAtPath(c, contentPath)

	// if key is provided, try to decrypt the data with it
	key, exists := c.GetQuery("key")
	if exists {
		keyb, err := base58.Decode(key)
		if err != nil {
			log.Errorf("error decoding key %s: %s", key, err)
			c.Status(404)
			return
		}
		plain, err := crypto.DecryptAES(data, keyb)
		if err != nil {
			log.Errorf("error decrypting %s: %s", contentPath, err)
			c.Status(404)
			return
		}
		c.Render(200, render.Data{Data: plain})
		return
	}

	// lastly, just return the raw bytes
	c.Render(200, render.Data{Data: data})
}

// profileHandler handles requests for profile info hosted on ipns
// NOTE: avatar is a magic path, will return data behind link at avatar_uri
func profileHandler(c *gin.Context) {
	pathp := c.Param("path")
	var isAvatar bool
	if pathp == "/avatar" {
		pathp += "_uri"
		isAvatar = true
	}

	// resolve the actual content
	rootId, err := peer.IDB58Decode(c.Param("root"))
	if err != nil {
		log.Errorf("error resolving profile %s: %s", c.Param("root"), err)
		c.Status(404)
		return
	}
	pth, err := core.Node.ResolveProfile(rootId)
	if err != nil {
		log.Errorf("error resolving profile %s: %s", c.Param("root"), err)
		c.Status(404)
		return
	}

	// get data behind path
	contentPath := pth.String() + pathp
	data := getDataAtPath(c, contentPath)

	// if this is an avatar request, fetch and return the linked image
	if isAvatar {
		location := string(data)
		if location == "" {
			fallback, _ := c.GetQuery("fallback")
			if fallback == "true" {
				location = fmt.Sprintf("https://avatars.dicebear.com/v2/identicon/%s.svg", c.Param("root"))
				c.Redirect(307, location)
				return
			} else {
				c.Status(404)
				return
			}
		}

		// parse ipfs link, must have key present
		parsed := strings.Split(location, "?key=")
		if len(parsed) != 2 {
			log.Errorf("invalid raw avatar path: %s", location)
			c.Status(404)
			return
		}
		cipher, err := core.Node.DataAtPath(parsed[0])
		if err != nil {
			log.Errorf("error getting raw avatar path %s: %s", parsed[0], err)
			c.Status(404)
			return
		}
		keyb, err := base58.Decode(parsed[1])
		if err != nil {
			log.Errorf("error decoding key %s: %s", parsed[1], err)
			c.Status(404)
			return
		}
		data, err = crypto.DecryptAES(cipher, keyb)
		if err != nil {
			log.Errorf("error decrypting %s: %s", parsed[0], err)
			c.Status(404)
			return
		}

		c.Header("Content-Type", "image/jpeg")
		c.Header("Cache-Control", "public, max-age=172800") // 2 days
	}

	c.Render(200, render.Data{Data: data})
}

// getDataAtPath get raw data or directory links at path
func getDataAtPath(c *gin.Context, pth string) []byte {
	data, err := core.Node.DataAtPath(pth)
	if err != nil {
		if err == iface.ErrIsDir {
			links, err := core.Node.LinksAtPath(pth)
			if err != nil {
				log.Errorf("error getting raw path %s: %s", pth, err)
				c.Status(404)
				return nil
			}
			var list []string
			for _, link := range links {
				list = append(list, "/"+link.Name)
			}
			c.String(200, "%s", strings.Join(list, "\n"))
			return nil
		}
		log.Errorf("error getting raw path %s: %s", pth, err)
		c.Status(404)
		return nil
	}
	return data
}

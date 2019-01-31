package gateway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	iface "gx/ipfs/QmX9YciaxRii8TARoEbmavzaeTUAe7BozeAgydsThNcTpy/go-ipfs/core/coreapi/interface"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"
	"gx/ipfs/QmZMWMvWMVKCbHetJ4RgndbuEF1io2UpUxwQwtNjtYPzSC/go-ipfs-files"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/gateway/static/css"
	"github.com/textileio/textile-go/gateway/templates"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
)

var log = logging.Logger("tex-gateway")

// Host is the instance used by the daemon
var Host *Gateway

// Gateway is a HTTP API for getting files and links from IPFS
type Gateway struct {
	Node   *core.Textile
	server *http.Server
}

// Start creates a gateway server
func (g *Gateway) Start(addr string) {
	gin.SetMode(gin.ReleaseMode)
	if g.Node != nil {
		gin.DefaultWriter = g.Node.Writer()
	}

	router := gin.Default()
	router.SetHTMLTemplate(parseTemplates())

	router.GET("/health", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusNoContent)
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		img, err := base64.StdEncoding.DecodeString(favicon)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		c.Header("Cache-Control", "public, max-age=172800")
		c.Render(http.StatusOK, render.Data{Data: img})
	})
	router.GET("/static/css/style.css", func(c *gin.Context) {
		c.Header("Content-Type", "text/css; charset=utf-8")
		c.Header("Cache-Control", "public, max-age=172800")
		c.String(http.StatusOK, css.Style)
	})

	router.GET("/ipfs/:root", g.gatewayHandler)
	router.GET("/ipfs/:root/*path", g.gatewayHandler)
	router.GET("/ipns/:root", g.profileHandler)
	router.GET("/ipns/:root/*path", g.profileHandler)

	router.GET("/cafe", g.cafeHandler)
	router.GET("/cafes", g.cafesHandler)

	router.NoRoute(func(c *gin.Context) {
		render404(c)
	})

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
	log.Infof("gateway listening at %s", g.server.Addr)
}

// Stop stops the gateway
func (g *Gateway) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := g.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down gateway: %s", err)
		return err
	}
	return nil
}

// Addr returns the gateway's address
func (g *Gateway) Addr() string {
	return g.server.Addr
}

// gatewayHandler handles gateway http requests
func (g *Gateway) gatewayHandler(c *gin.Context) {
	contentPath := c.Param("root") + c.Param("path")

	data := g.getDataAtPath(c, contentPath)

	// attempt decrypt if key present
	key, exists := c.GetQuery("key")
	if exists {
		keyb, err := base58.Decode(key)
		if err != nil {
			log.Errorf("error decoding key %s: %s", key, err)
			render404(c)
			return
		}
		plain, err := crypto.DecryptAES(data, keyb)
		if err != nil {
			log.Errorf("error decrypting %s: %s", contentPath, err)
			render404(c)
			return
		}
		c.Render(200, render.Data{Data: plain})
		return
	}

	c.Render(200, render.Data{Data: data})
}

var avatarRx = regexp.MustCompile(`/avatar($|/small$|/large$)`)

// profileHandler handles requests for profile info hosted on ipns
// NOTE: avatar is a magic path, will return data behind link at avatar_uri
// NOTICE: This method has been deprecated and is only here temporarily for backward compatibility
func (g *Gateway) profileHandler(c *gin.Context) {
	pathp := c.Param("path")
	if len(pathp) > 0 && pathp[len(pathp)-1] == '/' {
		pathp = pathp[:len(pathp)-1]
	}
	var isAvatar bool
	var avatarSize string

	matches := avatarRx.FindStringSubmatch(pathp)
	if len(matches) == 2 {
		pathp = "/avatar_uri"
		isAvatar = true

		switch matches[1] {
		case "/large":
			avatarSize = "large"
		default:
			avatarSize = "small"
		}
	}

	rootId, err := peer.IDB58Decode(c.Param("root"))
	if err != nil {
		log.Errorf("error decoding root %s: %s", c.Param("root"), err)
		render404(c)
		return
	}

	pth, err := ipfs.ResolveIPNS(g.Node.Ipfs(), rootId)
	if err != nil {
		log.Errorf("error resolving profile %s: %s", c.Param("root"), err)
		render404(c)
		return
	}

	contentPath := pth.String() + pathp
	data := g.getDataAtPath(c, contentPath)

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
				render404(c)
				return
			}
		}

		// old style w/ key
		parsed := strings.Split(location, "?key=")
		if len(parsed) == 2 {
			keyb, err := base58.Decode(parsed[1])
			if err != nil {
				log.Errorf("error decoding key %s: %s", parsed[1], err)
				render404(c)
				return
			}

			ciphertext, err := g.Node.DataAtPath(parsed[0])
			if err != nil {
				render404(c)
				return
			}

			data, err = crypto.DecryptAES(ciphertext, keyb)
			if err != nil {
				log.Errorf("error decrypting %s: %s", parsed[0], err)
				render404(c)
				return
			}

			c.Header("Content-Type", "image/jpeg")

		} else {
			pth := fmt.Sprintf("%s/0/%s/d", location, avatarSize)
			data, err = g.Node.DataAtPath(pth)
			if err != nil {
				render404(c)
				return
			}

			var stop int
			if len(data) < 512 {
				stop = len(data)
			} else {
				stop = 512
			}
			media := http.DetectContentType(data[:stop])
			if media != "" {
				c.Header("Content-Type", media)
			}
		}

		c.Header("Cache-Control", "public, max-age=172800") // 2 days
	}

	c.Render(200, render.Data{Data: data})
}

// cafeHandler returns this peer's cafe info
func (g *Gateway) cafeHandler(c *gin.Context) {
	conf := g.Node.Config()
	if !conf.Cafe.Host.Open {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.JSON(http.StatusOK, g.Node.CafeInfo())
}

// cafesHandler returns this peer's and it's zone neighbor's cafe info
func (g *Gateway) cafesHandler(c *gin.Context) {
	conf := g.Node.Config()
	if !conf.Cafe.Host.Open {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if conf.Cafe.Host.NeighborURL == "" {
		c.JSON(http.StatusOK, gin.H{
			"primary": g.Node.CafeInfo(),
		})
		return
	}

	secondary, err := getCafeInfo(conf.Cafe.Host.NeighborURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"primary":   g.Node.CafeInfo(),
		"secondary": secondary,
	})
}

type link struct {
	Path iface.Path
	Link *ipld.Link
	Size string
}

// getDataAtPath get raw data or directory links at path
func (g *Gateway) getDataAtPath(c *gin.Context, pth string) []byte {
	data, err := g.Node.DataAtPath(pth)
	if err != nil {
		if err == files.ErrNotReader {
			root, err := iface.ParsePath(pth)
			if err != nil {
				log.Errorf("error parsing path %s: %s", pth, err)
				render404(c)
				return nil
			}

			var back string
			parts := strings.Split(root.String(), "/")
			if len(parts) > 0 {
				last := parts[:len(parts)-1]
				back = strings.Join(last, "/")
			}
			if back == "/ipfs" || back == "/ipns" {
				back = root.String()
			}

			ilinks, err := g.Node.LinksAtPath(pth)
			if err != nil {
				log.Errorf("error getting links %s: %s", pth, err)
				render404(c)
				return nil
			}

			var links []link
			for _, l := range ilinks {
				ipath, err := iface.ParsePath(pth + "/" + l.Name)
				if err != nil {
					log.Errorf("error parsing path %s: %s", pth, err)
					render404(c)
					return nil
				}
				links = append(links, link{
					Path: ipath,
					Link: l,
					Size: byteCountDecimal(int64(l.Size)),
				})
			}

			c.HTML(http.StatusOK, "index", gin.H{
				"root":  root,
				"back":  back,
				"links": links,
			})
			return nil
		}

		log.Errorf("error getting path %s: %s", pth, err)
		render404(c)
		return nil
	}
	return data
}

func render404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", nil)
}

func parseTemplates() *template.Template {
	temp, err := template.New("index").Parse(templates.Index)
	if err != nil {
		panic(err)
	}
	temp, err = temp.New("404").Parse(templates.NotFound)
	if err != nil {
		panic(err)
	}
	return temp
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// getCafeInfo returns info about a fellow cafe
func getCafeInfo(addr string) (*repo.Cafe, error) {
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s returned bad status: %d", addr, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, fmt.Errorf("no response from %s", addr)
	}

	var info repo.Cafe
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

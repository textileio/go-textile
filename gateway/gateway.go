package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/golang/protobuf/jsonpb"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
	ipfspath "github.com/ipfs/go-path"
	iface "github.com/ipfs/interface-go-ipfs-core"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/mr-tron/base58/base58"
	gincors "github.com/rs/cors/wrapper/gin"
	"github.com/textileio/go-textile/bots"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/gateway/static/css"
	"github.com/textileio/go-textile/gateway/templates"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

var log = logging.Logger("tex-gateway")

// Host is the instance used by the daemon
var Host *Gateway

// Gateway is a HTTP API for getting files and links from IPFS
type Gateway struct {
	Node   *core.Textile
	Bots   *bots.Service
	server *http.Server
}

// Start creates a gateway server
func (g *Gateway) Start(addr string) {
	gin.SetMode(gin.ReleaseMode)
	if g.Node != nil {
		gin.DefaultWriter = g.Node.Writer()
	}
	conf := g.Node.Config()

	router := gin.Default()
	router.Use(location.Default())

	// Add the CORS middleware
	// Merges the API HTTPHeaders (from config/init) into blank/default CORS configuration
	router.Use(gincors.New(core.ConvertHeadersToCorsOptions(conf.API.HTTPHeaders)))

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

	router.GET("/ipfs/:root", g.ipfsHandler)
	router.GET("/ipfs/:root/*path", g.ipfsHandler)
	router.GET("/ipns/:root", g.ipnsHandler)
	router.GET("/ipns/:root/*path", g.ipnsHandler)

	router.GET("/", g.cafeHandler)
	router.GET("/cafe", g.cafeHandler)
	router.GET("/cafes", g.cafesHandler)

	router.GET("/bots/:root", g.botsHandler)

	router.NoRoute(func(c *gin.Context) {
		g.render404(c)
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

// ipfsHandler renders and optionally decrypts data behind an IPFS address
func (g *Gateway) ipfsHandler(c *gin.Context) {
	contentPath := c.Param("root") + c.Param("path")

	data := g.getDataAtPath(c, contentPath)
	if data == nil {
		return
	}

	// attempt decrypt if key present
	key, exists := c.GetQuery("key")
	if exists {
		keyb, err := base58.Decode(key)
		if err != nil {
			log.Debugf("error decoding key %s: %s", key, err)
			g.render404(c)
			return
		}
		plain, err := crypto.DecryptAES(data, keyb)
		if err != nil {
			log.Debugf("error decrypting %s: %s", contentPath, err)
			g.render404(c)
			return
		}
		c.Render(200, render.Data{Data: plain})
		return
	}

	c.Render(200, render.Data{Data: data})
}

// ipnsHandler renders data behind an IPNS address
func (g *Gateway) ipnsHandler(c *gin.Context) {
	pathp := c.Param("path")
	if len(pathp) > 0 && pathp[len(pathp)-1] == '/' {
		pathp = pathp[:len(pathp)-1]
	}

	rootId, err := peer.IDB58Decode(c.Param("root"))
	if err != nil {
		log.Debugf("error decoding root %s: %s", c.Param("root"), err)
		g.render404(c)
		return
	}

	pth, err := ipfs.ResolveIPNS(g.Node.Ipfs(), rootId, time.Second*30)
	if err != nil {
		log.Debugf("error resolving profile %s: %s", c.Param("root"), err)
		g.render404(c)
		return
	}

	data := g.getDataAtPath(c, pth.String()+pathp)
	if data == nil {
		return
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

// ipnsHandler renders data behind an IPNS address
func (g *Gateway) botsHandler(c *gin.Context) {
	botID := c.Param("root")
	if g.Bots == (&bots.Service{}) { // bot doesn't exist yet
		log.Errorf("error no bots: %s", botID)
		g.render404(c)
		return
	}

	if !g.Bots.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
		g.render404(c)
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := g.Bots.Get(botID, qbytes)
	if err != nil {
		log.Errorf("bad bot response: %s", botID)
		g.render404(c)
		return
	}
	statusInt := int(botResponse.Status)

	c.Render(statusInt, render.Data{Data: botResponse.Body})
}

// link represents a node link for HTML rendering
type link struct {
	Path ipfspath.Path
	Link *ipld.Link
	Size string
}

// getDataAtPath get raw data or directory links at path
func (g *Gateway) getDataAtPath(c *gin.Context, pth string) []byte {
	data, err := g.Node.DataAtPath(pth)
	if err != nil {
		if err == iface.ErrIsDir {
			root, err := ipfspath.ParsePath(pth)
			if err != nil {
				log.Debugf("error parsing path %s: %s", pth, err)
				g.render404(c)
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
				log.Debugf("error getting links %s: %s", pth, err)
				g.render404(c)
				return nil
			}

			var links []link
			for _, l := range ilinks {
				ipath, err := ipfspath.ParsePath(pth + "/" + l.Name)
				if err != nil {
					log.Debugf("error parsing path %s: %s", pth, err)
					g.render404(c)
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

		log.Debugf("error getting path %s: %s", pth, err)
		g.render404(c)
		return nil
	}
	return data
}

// render404 renders the 404 template
func (g *Gateway) render404(c *gin.Context) {
	if strings.Contains(c.Request.URL.String(), "small/content") ||
		strings.Contains(c.Request.URL.String(), "large/content") {
		var url string
		if g.Node.Config().Cafe.Host.URL != "" {
			url = strings.TrimRight(g.Node.Config().Cafe.Host.URL, "/")
		} else {
			loc := location.Get(c)
			url = fmt.Sprintf("%s://%s", loc.Scheme, loc.Host)
		}
		pth := strings.Replace(c.Request.URL.String(), "/content", "/d", 1)
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s", url, pth))
		return
	}

	c.HTML(http.StatusNotFound, "404", nil)
}

// parseTemplates loads HTML templates
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

// byteCountDecimal formats bytes
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
func getCafeInfo(addr string) (*pb.Cafe, error) {
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

	var info pb.Cafe
	if err := jsonpb.Unmarshal(res.Body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

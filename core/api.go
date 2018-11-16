package core

import (
	"context"
	"errors"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	m "github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"

	"github.com/gin-gonic/gin"
)

// apiVersion is the api version
const apiVersion = "v0"

// apiHost is the instance used by the daemon
var apiHost *api

// api is a limited HTTP REST API for the cmd tool
type api struct {
	addr   string
	server *http.Server
	node   *Textile
}

// StartApi starts the host instance
func (t *Textile) StartApi(addr string) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = t.writer
	apiHost = &api{addr: addr, node: t}
	apiHost.Start()
}

// StopApi starts the host instance
func (t *Textile) StopApi() error {
	return apiHost.Stop()
}

// ApiAddr returns the api address
func (t *Textile) ApiAddr() string {
	if apiHost == nil {
		return ""
	}
	return apiHost.addr
}

// Start starts the http api
func (a *api) Start() {
	// setup router
	router := gin.Default()
	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"cafe_version": apiVersion,
			"node_version": Version,
		})
	})
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})

	// v0 routes
	v0 := router.Group("/api/v0")
	{
		v0.GET("/peer", a.peer)
		v0.GET("/address", a.address)
		v0.GET("/ping", a.ping)

		mills := v0.Group("/mills")
		mills.POST("/schema", a.schemaMill)
		mills.POST("/blob", a.blobMill)
		mills.POST("/image/resize", a.imageResizeMill)
		mills.POST("/image/exif", a.imageExifMill)

		threads := v0.Group("/threads")
		threads.POST("", a.addThreads)
		threads.GET("", a.lsThreads)
		threads.GET("/:id", a.getThreads)
		threads.DELETE("/:id", a.rmThreads)

		threads.POST("/:id/files", a.addThreadFiles)
		threads.GET("/:id/files", a.lsThreadFiles)
		threads.GET("/:id/files/:block", a.getThreadFiles)

		invite := v0.Group("/invite")
		invite.POST("", a.createInvite)
		invite.POST("/:id/accept", a.acceptInvite)
		invite.POST("/:id/ignore", a.ignoreInvite)

		cafes := v0.Group("/cafes")
		cafes.POST("", a.addCafes)
		cafes.GET("", a.lsCafes)
		cafes.GET("/:id", a.getCafes)
		cafes.DELETE("/:id", a.rmCafes)
		cafes.POST("/check_mail", a.checkMailCafes)
	}
	a.server = &http.Server{
		Addr:    a.addr,
		Handler: router,
	}

	// start listening
	errc := make(chan error)
	go func() {
		errc <- a.server.ListenAndServe()
		close(errc)
	}()
	go func() {
		for {
			select {
			case err, ok := <-errc:
				if err != nil && err != http.ErrServerClosed {
					log.Errorf("api error: %s", err)
				}
				if !ok {
					log.Info("api was shutdown")
					return
				}
			}
		}
	}()
	log.Infof("api listening at %s\n", a.server.Addr)
}

// Stop stops the http api
func (a *api) Stop() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down api: %s", err)
		return err
	}
	cancel()
	return nil
}

// -- UTILITY ENDPOINTS -- //

func (a *api) peer(g *gin.Context) {
	pid, err := a.node.PeerId()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, pid.Pretty())
}

func (a *api) address(g *gin.Context) {
	g.String(http.StatusOK, a.node.account.Address())
}

func (a *api) ping(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing peer id")
		return
	}
	pid, err := peer.IDB58Decode(args[0])
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	status, err := a.node.Ping(pid)
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, string(status))
}

func (a *api) readArgs(g *gin.Context) ([]string, error) {
	header := g.Request.Header.Get("X-Textile-Args")
	var args []string
	for _, a := range strings.Split(header, ",") {
		arg := strings.TrimSpace(a)
		if arg != "" {
			args = append(args, arg)
		}
	}
	return args, nil
}

func (a *api) readOpts(g *gin.Context) (map[string]string, error) {
	header := g.Request.Header.Get("X-Textile-Opts")
	opts := make(map[string]string)
	for _, o := range strings.Split(header, ",") {
		opt := strings.TrimSpace(o)
		if opt != "" {
			parts := strings.Split(opt, "=")
			if len(parts) == 2 {
				opts[parts[0]] = parts[1]
			}
		}
	}
	return opts, nil
}

func (a *api) openFile(g *gin.Context) (multipart.File, string, error) {
	form, err := g.MultipartForm()
	if err != nil {
		return nil, "", err
	}
	if len(form.File["file"]) == 0 {
		return nil, "", errors.New("no file attached")
	}
	header := form.File["file"][0]
	file, err := header.Open()
	if err != nil {
		return nil, "", err
	}
	return file, header.Filename, nil
}

func (a *api) getFileConfig(g *gin.Context, mill m.Mill, use string) (*AddFileConfig, error) {
	var reader io.ReadSeeker
	conf := &AddFileConfig{}

	if use == "" {
		file, fn, err := a.openFile(g)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
		conf.Name = fn
	} else {
		var file *repo.File
		var err error
		reader, file, err = a.node.filePlaintext(use)
		if err != nil {
			return nil, err
		}
		conf.Name = file.Name
		conf.Use = file.Checksum
	}

	media, err := a.node.MediaType(reader, mill)
	if err != nil {
		return nil, err
	}
	conf.Media = media
	reader.Seek(0, 0)

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf.Input = data

	return conf, nil
}

func (a *api) abort500(g *gin.Context, err error) {
	g.String(http.StatusInternalServerError, err.Error())
}

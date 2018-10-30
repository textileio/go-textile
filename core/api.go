package core

import (
	"context"
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"net/http"
	"strings"
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

		v0.POST("/threads/add", a.addThread)
		v0.GET("/threads", a.getThreads)
		v0.GET("/threads/:id", a.getThread)
		//v0.DELETE("/threads/:id", a.delThread)

		v0.POST("/images/add", a.addImage)
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

// -- INFO -- //

func (a *api) peer(g *gin.Context) {
	pid, err := a.node.PeerId()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, pid.Pretty())
}

func (a *api) address(g *gin.Context) {
	addr, err := a.node.Address()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, addr)
}

// -- NETWORK -- //

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

// -- IMAGES -- //

func (a *api) addImage(g *gin.Context) {
	form, err := g.MultipartForm()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	fileHeaders := form.File["file"]
	var adds []*AddDataResult
	for _, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		added, err := a.node.AddImage(file, header.Filename)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		file.Close()
		adds = append(adds, added)
	}
	g.JSON(http.StatusCreated, gin.H{"items": adds})
}

// -- THREADS -- //

func (a *api) addThread(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing thread name")
		return
	}
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		a.abort500(g, err)
		return
	}
	thrd, err := a.node.AddThread(args[0], sk, true)
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusCreated, thrd.Id)
}

func (a *api) getThreads(g *gin.Context) {
	var infos []*ThreadInfo
	for _, thrd := range a.node.Threads() {
		info, err := thrd.Info()
		if err != nil {
			a.abort500(g, err)
			return
		}
		infos = append(infos, info)
	}
	g.JSON(http.StatusOK, gin.H{"items": infos})
}

func (a *api) getThread(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing thread id")
		return
	}
	_, thrd := a.node.Thread(args[0])
	if thrd == nil {
		g.String(http.StatusNotFound, "thread not found")
		return
	}
	info, err := thrd.Info()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusOK, info)
}

// -- HELPERS -- //

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

func (a *api) abort500(g *gin.Context, err error) {
	g.String(http.StatusInternalServerError, err.Error())
}

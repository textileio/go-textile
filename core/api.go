package core

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/gin-contrib/cors"
	limit "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	m "github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/config"
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

	// middleware
	conf := a.node.Config()
	// CORS
	router.Use(cors.New(getCORSSettings(conf)))
	// size limits
	if conf.API.SizeLimit > 0 {
		router.Use(limit.RequestSizeLimiter(conf.API.SizeLimit))
	}

	// v0 routes
	v0 := router.Group("/api/v0")
	{
		v0.GET("/peer", a.peer)
		v0.GET("/address", a.address)
		v0.GET("/ping", a.ping)

		profile := v0.Group("/profile")
		{
			profile.GET("", a.getProfile)
			profile.POST("/username", a.setUsername)
			profile.POST("/avatar", a.setAvatar)
		}

		mills := v0.Group("/mills")
		{
			mills.POST("/schema", a.schemaMill)
			mills.POST("/blob", a.blobMill)
			mills.POST("/image/resize", a.imageResizeMill)
			mills.POST("/image/exif", a.imageExifMill)
			mills.POST("/json", a.jsonMill)
		}

		threads := v0.Group("/threads")
		{
			threads.POST("", a.addThreads)
			threads.GET("", a.lsThreads)
			threads.GET("/:id", a.getThreads)
			threads.GET("/:id/peers", a.peersThreads)
			threads.DELETE("/:id", a.rmThreads)
			threads.POST("/:id/messages", a.addThreadMessages)
			threads.POST("/:id/files", a.addThreadFiles)
		}

		blocks := v0.Group("/blocks")
		{
			blocks.GET("", a.lsBlocks)

			block := blocks.Group("/:id")
			{
				block.GET("", a.getBlocks)
				block.DELETE("", a.rmBlocks)

				block.GET("/comment", a.getBlockComment)
				comments := block.Group("/comments")
				{
					comments.POST("", a.addBlockComments)
					comments.GET("", a.lsBlockComments)
				}

				block.GET("/like", a.getBlockLike)
				likes := block.Group("/likes")
				{
					likes.POST("", a.addBlockLikes)
					likes.GET("", a.lsBlockLikes)
				}
			}
		}

		feed := v0.Group("/feed")
		{
			feed.GET("", a.lsThreadFeed)
		}

		messages := v0.Group("/messages")
		{
			messages.GET("", a.lsThreadMessages)
			messages.GET("/:block", a.getThreadMessages)
		}

		files := v0.Group("/files")
		{
			files.GET("", a.lsThreadFiles)
			files.GET("/:block", a.getThreadFiles)
		}

		keys := v0.Group("/keys")
		{
			keys.GET("/:target", a.lsThreadFileTargetKeys)
		}

		sub := v0.Group("/sub")
		{
			sub.GET("", a.getThreadsSub)
			sub.GET("/:id", a.getThreadsSub)
		}

		invites := v0.Group("/invites")
		{
			invites.POST("", a.createInvites)
			invites.GET("", a.lsInvites)
			invites.POST("/:id/accept", a.acceptInvites)
			invites.POST("/:id/ignore", a.ignoreInvites)
		}

		notifs := v0.Group("/notifications")
		{
			notifs.GET("", a.lsNotifications)
			notifs.POST("/:id/read", a.readNotifications)
		}

		cafes := v0.Group("/cafes")
		{
			cafes.POST("", a.addCafes)
			cafes.GET("", a.lsCafes)
			cafes.GET("/:id", a.getCafes)
			cafes.DELETE("/:id", a.rmCafes)
			cafes.POST("/messages", a.checkCafeMessages)
		}

		swarm := v0.Group("/swarm")
		{
			swarm.POST("/connect", a.swarmConnect)
			swarm.GET("/peers", a.swarmPeers)
		}

		contacts := v0.Group("/contacts")
		{
			contacts.GET("", a.lsContacts)
			contacts.GET("/:id", a.getContacts)
			contacts.DELETE("/:id", a.rmContacts)
			contacts.POST("/search", a.searchContacts)
		}

		ipfs := v0.Group("/ipfs")
		{
			ipfs.GET("/:cid", a.ipfsCat)
		}

		logs := v0.Group("/logs")
		{
			logs.POST("", a.logsCall)
			logs.GET("", a.logsCall)
			logs.POST("/:subsystem", a.logsCall)
			logs.GET("/:subsystem", a.logsCall)
		}

		conf := v0.Group("/config")
		{
			conf.GET("", a.getConfig)
			conf.PUT("", a.setConfig)
			conf.GET("/*path", a.getConfig)
			conf.PATCH("", a.patchConfig)
		}

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
	log.Infof("api listening at %s", a.server.Addr)
}

// Stop stops the http api
func (a *api) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down api: %s", err)
		return err
	}
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
		arg, err := url.PathUnescape(strings.TrimSpace(a))
		if err != nil {
			return nil, err
		}
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
				v, err := url.PathUnescape(parts[1])
				if err != nil {
					return nil, err
				}
				opts[parts[0]] = v
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

func (a *api) getFileConfig(g *gin.Context, mill m.Mill, use string, plaintext bool) (*AddFileConfig, error) {
	var reader io.ReadSeeker
	conf := &AddFileConfig{}

	if use == "" {
		f, fn, err := a.openFile(g)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
		conf.Name = fn

	} else {
		var file *repo.File
		var err error
		reader, file, err = a.node.FileData(use)
		if err != nil {
			return nil, err
		}
		conf.Name = file.Name
		conf.Use = file.Checksum
	}

	media, err := a.node.GetMedia(reader, mill)
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
	conf.Plaintext = plaintext

	return conf, nil
}

func (a *api) abort500(g *gin.Context, err error) {
	g.String(http.StatusInternalServerError, err.Error())
}

// getCORSSettings returns custom CORS settings given HTTPHeaders config options
func getCORSSettings(config *config.Config) cors.Config {
	headers := config.API.HTTPHeaders
	cconfig := cors.DefaultConfig()

	control, ok := headers["Access-Control-Allow-Origin"]
	if ok && len(control) > 0 {
		cconfig.AllowOrigins = control
		for _, origin := range control {
			if origin == "*" {
				cconfig.AllowAllOrigins = true
				cconfig.AllowOrigins = nil
				break
			}
		}
	} else {
		defaultHost := config.Addresses.API
		match, _ := regexp.MatchString("^https?://", defaultHost)
		if !match {
			defaultHost = "http://" + defaultHost
		}
		cconfig.AllowOrigins = []string{defaultHost}
	}

	control, ok = headers["Access-Control-Allow-Methods"]
	if ok && len(control) > 0 {
		cconfig.AllowMethods = control
	}

	control, ok = headers["Access-Control-Allow-Headers"]
	if ok && len(control) > 0 {
		cconfig.AllowHeaders = control
	}

	return cconfig
}

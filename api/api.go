package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	limit "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	logging "github.com/ipfs/go-log"
	ipfspath "github.com/ipfs/go-path"
	gincors "github.com/rs/cors/wrapper/gin"
	swagger "github.com/swaggo/gin-swagger"
	sfiles "github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/textileio/go-textile/api/docs"
	"github.com/textileio/go-textile/bots"
	"github.com/textileio/go-textile/common"
	"github.com/textileio/go-textile/core"
	ipfsutil "github.com/textileio/go-textile/ipfs"
	m "github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
)

// apiVersion is the api version
const apiVersion = "v0"

// TODO: create api logger
var log = logging.Logger("tex-gateway")

// Host is the instance used by the daemon
var Host *Api

// Gateway is a HTTP API for getting files and links from IPFS
type Api struct {
	Node     *core.Textile
	Bots     *bots.Service
	PinCode  string
	RepoPath string
	server   *http.Server
	addr     string
	docs     bool
}

// pbUnmarshaler is used to unmarshal JSON protobufs
var pbUnmarshaler = jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

// Start starts the host instance
func (a *Api) Start(addr string, serveDocs bool) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = a.Node.Writer()
	a.addr = addr
	a.docs = serveDocs
	a.Run()
}

// Addr returns the api address
func (a *Api) Addr() string {
	if Host == nil {
		return ""
	}
	return Host.addr
}

// @title Textile REST API
// @version 0
// @description Textile's HTTP REST API Documentation
// @termsOfService https://github.com/textileio/go-textile/blob/master/TERMS

// @contact.name Textile
// @contact.url https://textile.io/
// @contact.email contact@textile.io

// @license.name MIT License
// @license.url https://github.com/textileio/go-textile/blob/master/LICENSE

// @securityDefinitions.basic BasicAuth
// @Security BasicAuth
// @BasePath /api/v0
func (a *Api) Run() {
	// Dynamically set the swagger 'host' value
	docs.SwaggerInfo.Host = a.addr

	router := gin.Default()

	conf := a.Node.Config()

	// middleware setup

	// Add the CORS middleware
	// Merges the API HTTPHeaders (from config/init) into blank/default CORS configuration
	router.Use(gincors.New(core.ConvertHeadersToCorsOptions(conf.API.HTTPHeaders)))

	// Add size limits
	if conf.API.SizeLimit > 0 {
		router.Use(limit.RequestSizeLimiter(conf.API.SizeLimit))
	}

	router.GET("/", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"cafe_version": apiVersion,
			"node_version": common.Version,
		})
	})
	router.GET("/health", func(g *gin.Context) {
		g.Writer.WriteHeader(http.StatusNoContent)
	})

	// API docs
	if a.docs {
		router.GET("/docs/*any", swagger.WrapHandler(sfiles.Handler))
	}

	// If given a passcode use it, else leave API wide open
	var auth gin.HandlerFunc
	if a.PinCode != "" {
		auth = gin.BasicAuth(gin.Accounts{a.Node.Account().Address(): a.PinCode})
	} else {
		auth = func(c *gin.Context) {
			// noop handler function
			c.Next()
		}
	}

	// v0 routes
	v0 := router.Group("/api/v0", auth)
	{
		v0.GET("/summary", a.nodeSummary)
		v0.GET("/ping", a.ping)
		v0.POST("/publish", a.publish)

		account := v0.Group("/account")
		{
			account.GET("", a.accountGet)
			account.GET("/seed", a.accountSeed)
			account.GET("/address", a.accountAddress)
		}

		profile := v0.Group("/profile")
		{
			profile.GET("", a.getProfile)
			profile.POST("/name", a.setName)
			profile.POST("/avatar", a.setAvatar)
		}

		contacts := v0.Group("/contacts")
		{
			contacts.PUT(":address", a.addContacts)
			contacts.GET("", a.lsContacts)
			contacts.GET("/:address", a.getContacts)
			contacts.DELETE("/:address", a.rmContacts)
			contacts.POST("/search", a.searchContacts)
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
			threads.PUT(":id", a.addOrUpdateThreads)
			threads.PUT(":id/name", a.renameThreads)
			threads.GET("", a.lsThreads)
			threads.GET("/:id", a.getThreads)
			threads.GET("/:id/peers", a.peersThreads)
			threads.DELETE("/:id", a.rmThreads)
			threads.POST("/:id/messages", a.addThreadMessages)
			threads.POST("/:id/files", a.addThreadFiles)
		}

		snapshots := v0.Group("/snapshots")
		{
			snapshots.POST("", a.createThreadSnapshots)
			snapshots.POST("/search", a.searchThreadSnapshots)
		}

		blocks := v0.Group("/blocks")
		{
			blocks.GET("", a.lsBlocks)

			block := blocks.Group("/:id")
			{
				block.GET("/meta", a.getBlockMeta)
				files := block.Group("/files")
				{
					files.GET("", a.getBlockFiles)
					file := files.Group("/:index/:path")
					{
						file.GET("/meta", a.getBlockFileMeta)
						file.GET("/content", a.getBlockFileContent)
					}
				}

				block.GET("", func(g *gin.Context) {
					id := g.Param("id")
					g.Redirect(http.StatusPermanentRedirect, "/api/v0/blocks/"+id+"/meta")
				})
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

		messages := v0.Group("/messages")
		{
			messages.GET("", a.lsThreadMessages)
			messages.GET("/:block", a.getThreadMessages)
		}

		files := v0.Group("/files")
		{
			files.GET("", a.lsThreadFiles)
			files.GET("/:block", func(g *gin.Context) {
				g.Redirect(http.StatusPermanentRedirect, "/api/v0/blocks/"+g.Param("block")+"/files")
			})
		}

		file := v0.Group("/file")
		{
			hash := file.Group("/:hash")
			{
				hash.GET("", func(g *gin.Context) {
					g.Redirect(http.StatusPermanentRedirect, "/api/v0/file/"+g.Param("hash")+"/meta")
				})
				hash.GET("/data", func(g *gin.Context) {
					g.Redirect(http.StatusPermanentRedirect, "/api/v0/file/"+g.Param("hash")+"/content")
				})
				hash.GET("/meta", a.getFileMeta)
				hash.GET("/content", a.getFileContent)
			}
		}

		feed := v0.Group("/feed")
		{
			feed.GET("", a.lsThreadFeed)
		}

		keys := v0.Group("/keys")
		{
			keys.GET("/:target", a.lsThreadFileTargetKeys)
		}

		observe := v0.Group("/observe")
		{
			observe.GET("", a.getThreadsObserve)
			observe.GET("/:thread", a.getThreadsObserve)
		}
		// alias
		subscribe := v0.Group("/subscribe")
		{
			subscribe.GET("", func(g *gin.Context) {
				g.Redirect(http.StatusPermanentRedirect, "/api/v0/observe")
			})
			subscribe.GET("/:thread", func(g *gin.Context) {
				g.Redirect(http.StatusPermanentRedirect, "/api/v0/observe/"+g.Param("thread"))
			})
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

		tokens := v0.Group("/tokens")
		{
			tokens.POST("", a.createTokens)
			tokens.GET("", a.lsTokens)
			tokens.GET("/:token", a.validateTokens)
			tokens.DELETE("/:token", a.rmTokens)
		}

		ipfs := v0.Group("/ipfs")
		{
			ipfs.GET("/id", a.ipfsId)
			ipfs.GET("/cat/*path", a.ipfsCat)

			swarm := ipfs.Group("/swarm")
			{
				swarm.POST("/connect", a.ipfsSwarmConnect)
				swarm.GET("/peers", a.ipfsSwarmPeers)
			}
		}

		bots := v0.Group("/bots")
		{
			bots.GET("/list", a.botsList)
			bots.POST("/disable", a.botsDisable)
			bots.POST("/enable", a.botsEnable)
			bots.GET("/id/:id", a.botsGet)
			bots.POST("/id/:id", a.botsPost)
			bots.DELETE("/id/:id", a.botsDelete)
			bots.PUT("/id/:id", a.botsPut)
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
func (a *Api) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("error shutting down api: %s", err)
		return err
	}
	return nil
}

// summary godoc
// @Summary Get a summary of node data
// @Tags utils
// @Produce application/json
// @Success 200 {object} pb.Summary "summary"
// @Router /summary [get]
func (a *Api) nodeSummary(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.Node.Summary())
}

func (a *Api) abort500(g *gin.Context, err error) {
	sendError(g, err, http.StatusInternalServerError)
}

func (a *Api) readArgs(g *gin.Context) ([]string, error) {
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

func (a *Api) readOpts(g *gin.Context) (map[string]string, error) {
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

func (a *Api) openFile(g *gin.Context) (multipart.File, string, error) {
	form, err := g.MultipartForm()
	if err != nil {
		return nil, "", err
	}
	if len(form.File["file"]) == 0 {
		return nil, "", fmt.Errorf("no file attached")
	}
	header := form.File["file"][0]
	file, err := header.Open()
	if err != nil {
		return nil, "", err
	}
	return file, header.Filename, nil
}

func (a *Api) getFileConfig(g *gin.Context, mill m.Mill, use string, plaintext bool) (*core.AddFileConfig, error) {
	var reader io.ReadSeeker
	conf := &core.AddFileConfig{}

	if use == "" {
		f, fn, err := a.openFile(g)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
		conf.Name = fn
	} else {
		ref, err := ipfspath.ParsePath(use)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(ref.String(), "/")
		hash := parts[len(parts)-1]
		var file *pb.FileIndex
		reader, file, err = a.Node.FileContent(hash)
		if err != nil {
			if err == core.ErrFileNotFound {
				// just cat the data from ipfs
				b, err := ipfsutil.DataAtPath(a.Node.Ipfs(), ref.String())
				if err != nil {
					return nil, err
				}
				reader = bytes.NewReader(b)
				conf.Use = ref.String()
			} else {
				return nil, err
			}
		} else {
			conf.Use = file.Checksum
		}
	}

	media, err := a.Node.GetMillMedia(reader, mill)
	if err != nil {
		return nil, err
	}
	conf.Media = media
	_, _ = reader.Seek(0, 0)

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf.Input = data
	conf.Plaintext = plaintext

	return conf, nil
}

// pbMarshaler is used to marshal protobufs to JSON
var pbMarshaler = jsonpb.Marshaler{
	OrigName: true,
}

// pbJSON responds with a JSON rendered protobuf message
func pbJSON(g *gin.Context, status int, msg proto.Message) {
	str, err := pbMarshaler.MarshalToString(msg)
	if err != nil {
		sendError(g, err, http.StatusBadRequest)
		return
	}
	g.Data(status, "application/json", []byte(str))
}

// pbValForEnumString returns the int value of a case-insensitive string representation of a pb enum
func pbValForEnumString(vals map[string]int32, str string) int32 {
	for v, i := range vals {
		if strings.ToLower(v) == strings.ToLower(str) {
			return i
		}
	}
	return 0
}

// sendError sends the error to the gin context
func sendError(g *gin.Context, err error, statusCode int) {
	g.String(statusCode, err.Error())
}

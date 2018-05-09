package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"

	"github.com/textileio/textile-go/central/controllers"
	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/middleware"
	tcore "github.com/textileio/textile-go/core"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
)

var log = logging.MustGetLogger("main")

var updateCache []string

const (
	cacheSize     = 32
	relayInterval = time.Second * 30
)

var tid = "QmPnQL1qT4cxWhUDwYHFKXVU9zXiwnU4anT22JTEj8cXnC"

func init() {
	// establish a connection to DB
	dao.Dao = &dao.DAO{
		Hosts:    os.Getenv("DB_HOSTS"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		TLS:      os.Getenv("DB_TLS") == "yes",
	}
	dao.Dao.Connect()

	// ensure we're indexed
	dao.Dao.Index()
}

func main() {
	go func() {
		// create a pubsub relay node
		config := tcore.NodeConfig{
			RepoPath:      ".ipfs",
			CentralApiURL: "https://api.textile.io",
			IsMobile:      false,
			LogLevel:      logging.DEBUG,
			LogFiles:      false,
		}
		node, err := tcore.NewNode(config)
		if err != nil {
			log.Fatal(err)
		}

		// bring it online
		err = node.Start()
		if err != nil {
			log.Fatal(err)
		}
		self := node.IpfsNode.Identity.Pretty()

		// create ticker for relaying updates
		ticker := time.NewTicker(relayInterval)
		go func() {
			for range ticker.C {
				relayLatest(node.IpfsNode)
			}
		}()

		// create the subscription
		sub, err := node.IpfsNode.Floodsub.Subscribe(tid)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("joined room %s as relay buddy\n", tid)

		ctx, _ := context.WithCancel(context.Background())
		for {
			// unload new message
			msg, err := sub.Next(ctx)
			if err == io.EOF || err == context.Canceled {
				return
			} else if err != nil {
				return
			}

			// unpack message
			from := msg.GetFrom().Pretty()
			if from == self {
				continue
			}
			hash := string(msg.GetData())
			log.Infof("adding update %s from %s to relay", hash, from)

			// ignore if exists
			var exists bool
		inner:
			for _, u := range updateCache {
				if hash == u {
					exists = true
					break inner
				}
			}
			if exists {
				continue
			}

			// add update to cache
			if len(updateCache) == cacheSize {
				updateCache = updateCache[1:]
			}
			updateCache = append(updateCache, hash)

			// relay now
			relayLatest(node.IpfsNode)
		}
	}()

	// build http router
	router := gin.Default()
	router.GET("/", controllers.Info)
	router.GET("/health", controllers.Health)

	// api routes
	v1 := router.Group("/api/v1")
	v1.Use(middleware.Auth(os.Getenv("TOKEN_SECRET")))
	{
		v1.PUT("/users", controllers.SignUp)
		v1.POST("/users", controllers.SignIn)
		v1.POST("/referrals", controllers.CreateReferral)
		v1.GET("/referrals", controllers.ListReferrals)
	}
	router.Run(fmt.Sprintf("%s", os.Getenv("BIND")))
}

func relayLatest(ipfs *core.IpfsNode) {
	for _, update := range updateCache {
		log.Debugf("relaying update %s to %s", update, tid)
		if err := ipfs.Floodsub.Publish(tid, []byte(update)); err != nil {
			log.Errorf("error relaying update: %s", err)
		}
	}
}

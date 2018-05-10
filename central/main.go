package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"

	"github.com/textileio/textile-go/central/controllers"
	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/middleware"
	tcore "github.com/textileio/textile-go/core"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"path/filepath"
)

var log = logging.MustGetLogger("main")

var updateCache = make(map[string]string)

const (
	relayInterval = time.Second * 30
)

var relayThread = os.Getenv("RELAY")

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
		// get home dir
		hd, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// create a pubsub relay node
		config := tcore.NodeConfig{
			RepoPath:      filepath.Join(hd, ".textile_central"),
			CentralApiURL: os.Getenv("BIND"),
			IsServer:      true,
			LogLevel:      logging.DEBUG,
			LogFiles:      false,
			SwarmPort:     "4001",
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
		sub, err := node.IpfsNode.Floodsub.Subscribe(relayThread)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("joined room %s as relay buddy\n", relayThread)

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

			// ignore if the latest from this peer has not changed
			if updateCache[from] == hash {
				continue
			}

			// add update to cache
			updateCache[from] = hash
			log.Infof("added update %s from %s to relay", hash, from)

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
	router.Run(os.Getenv("BIND"))
}

func relayLatest(ipfs *core.IpfsNode) {
	for from, update := range updateCache {
		log.Debugf("relaying update %s from %s", update, from)
		if err := ipfs.Floodsub.Publish(relayThread, []byte(update)); err != nil {
			log.Errorf("error relaying update: %s", err)
		}
	}
}

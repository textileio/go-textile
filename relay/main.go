package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"

	tcore "github.com/textileio/textile-go/core"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
)

var (
	log         = logging.MustGetLogger("main")
	updateCache = make(map[string]string)
	relayThread = os.Getenv("RELAY")
)

const (
	relayInterval = time.Second * 30
)

func main() {
	// get home dir
	hd, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	// create a pubsub relay node
	config := tcore.NodeConfig{
		RepoPath:  filepath.Join(hd, ".textile_central"),
		IsServer:  true,
		LogLevel:  logging.DEBUG,
		LogFiles:  false,
		SwarmPort: "4001",
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
			log.Errorf("subscription ended with known error: %s", err)
			return
		} else if err != nil {
			log.Errorf("subscription ended with unknown error: %s", err)
			return
		}

		// unpack message
		from := msg.GetFrom().Pretty()
		hash := string(msg.GetData())

		// ignore if from us
		if from == self {
			continue
		}

		// add new updates to cache
		if updateCache[from] != hash {
			updateCache[from] = hash
			log.Infof("added new update %s from %s to relay", hash, from)
		}

		// relay now
		relayLatest(node.IpfsNode)
	}
}

func relayLatest(ipfs *core.IpfsNode) {
	for from, update := range updateCache {
		log.Debugf("relaying update %s from %s", update, from)
		if err := ipfs.Floodsub.Publish(relayThread, []byte(update)); err != nil {
			log.Errorf("error relaying update: %s", err)
		}
	}
}

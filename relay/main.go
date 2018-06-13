package main

import (
	"context"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	tcore "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		LogLevel: logging.DEBUG,
		LogFiles: false,
		WalletConfig: wallet.Config{
			RepoPath:  filepath.Join(hd, fmt.Sprintf(".relay_%s", relayThread)),
			IsServer:  true,
			SwarmPort: os.Getenv("SWARM_PORT"),
		},
	}
	node, err := tcore.NewNode(config)
	if err != nil {
		log.Fatal(err)
	}

	// bring it online
	online, err := node.StartWallet()
	if err != nil {
		log.Fatal(err)
	}
	<-online
	self, err := node.Wallet.GetIPFSPeerID()
	if err != nil {
		log.Fatal(err)
	}

	var relay = func() {
		for from, update := range updateCache {
			go func(from string, update string) {
				log.Debug("starting relay...")
				msg := fmt.Sprintf("relay:%s", update)
				if err := node.Wallet.Publish(relayThread, []byte(msg)); err != nil {
					log.Errorf("error relaying update: %s", err)
				}
				log.Debugf("relayed update %s from %s", update, from)
			}(from, update)
		}
	}

	// create ticker for relaying updates
	ticker := time.NewTicker(relayInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			relay()
		}
	}()

	// create the subscription
	sub, err := node.Wallet.Subscribe(relayThread)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("joined room %s as relay buddy\n", relayThread)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

		// ignore if from another relay
		tmp := strings.Split(hash, ":")
		if len(tmp) > 1 && tmp[0] == "relay" {
			log.Debugf("got update from fellow relay: %s, aborting", from)
			continue
		}

		// add new updates to cache
		if hash == "ping" {
			log.Infof("got ping from %s", from)
		} else if updateCache[from] != hash {
			updateCache[from] = hash
			log.Infof("added new update %s from %s to relay", hash, from)
		}

		// relay now
		relay()
	}
}

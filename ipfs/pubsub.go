package ipfs

import (
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

const PublishTimeout = time.Second * 5

// Publish publishes data to a topic
func Publish(node *core.IpfsNode, topic string, data []byte) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	ctx, pcancel := context.WithTimeout(node.Context(), PublishTimeout)
	defer pcancel()

	log.Debugf("publishing to topic %s", topic)
	return api.PubSub().Publish(ctx, topic, data)
}

// Subscribe subscribes to a topic
func Subscribe(node *core.IpfsNode, ctx context.Context, topic string, discover bool, msgs chan iface.PubSubMessage) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	sub, err := api.PubSub().Subscribe(ctx, topic, options.PubSub.Discover(discover))
	if err != nil {
		return err
	}
	defer sub.Close()

	for {
		msg, err := sub.Next(node.Context())
		if err == io.EOF || err == context.Canceled {
			return nil
		} else if err != nil {
			return err
		}
		msgs <- msg
	}
}

// connectToTopicReceiver attempts to connect with a pubsub topic's receiver
func connectToTopicReceiver(node *core.IpfsNode, ctx context.Context, topic string) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	blk, err := api.Block().Put(ctx, strings.NewReader("floodsub:"+topic))
	if err != nil {
		return err
	}

	fctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Debugf("looking for peers for topic %s", topic)
	provs := node.Routing.FindProvidersAsync(fctx, blk.Path().Cid(), 10)
	var wg sync.WaitGroup
	for p := range provs {
		log.Debugf("found topic provider %s", p.ID.Pretty())
		if !strings.Contains(topic, p.ID.Pretty()) {
			continue
		}
		wg.Add(1)
		go func(pi peer.AddrInfo) {
			defer wg.Done()
			err := node.PeerHost.Connect(ctx, pi)
			if err != nil {
				log.Info("pubsub discover: ", err)
				return
			}
			log.Info("connected to pubsub peer:", pi.ID)
			cancel()
		}(p)
	}

	wg.Wait()
	return nil
}

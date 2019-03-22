package ipfs

import (
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
	pstore "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore"
)

const publishTimeout = time.Second * 5

// Publish publishes data to a topic
func Publish(node *core.IpfsNode, topic string, data []byte, connect bool) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(node.Context(), publishTimeout)
	defer cancel()

	if connect {
		if err := connectToPubSubPeers(node, ctx, topic); err != nil {
			return err
		}
	}

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

// connectToPubSubPeers attempts to connect with a pubsub topic's peers
func connectToPubSubPeers(node *core.IpfsNode, ctx context.Context, topic string) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	blk, err := api.Block().Put(ctx, strings.NewReader("floodsub:"+topic))
	if err != nil {
		return err
	}

	provs := node.Routing.FindProvidersAsync(ctx, blk.Path().Cid(), 10)
	var wg sync.WaitGroup
	for p := range provs {
		if !strings.Contains(topic, p.ID.Pretty()) {
			continue
		}
		wg.Add(1)
		go func(pi pstore.PeerInfo) {
			defer wg.Done()
			//ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			//defer cancel()
			err := node.PeerHost.Connect(ctx, pi)
			if err != nil {
				log.Info("pubsub discover: ", err)
				return
			}
			log.Info("connected to pubsub peer:", pi.ID)
		}(p)
	}

	wg.Wait()
	return nil
}

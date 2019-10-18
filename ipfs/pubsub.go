package ipfs

import (
	"context"
	"io"
	"time"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
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

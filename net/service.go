package net

import (
	"context"
	"errors"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/pb"
	inet "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

var (
	log               = logging.MustGetLogger("net")
	OutOfOrderMessage = errors.New("message arrived out of order")
	DuplicateMessage  = errors.New("duplicate message")
)

type NetworkService interface {
	// Handle incoming streams
	HandleNewStream(s inet.Stream)

	// Get handler for mesage type
	HandlerForMsgType(t pb.Message_MessageType) func(peer.ID, *pb.Message, interface{}) (*pb.Message, error)

	// Send request to a peer and wait for the response
	SendRequest(ctx context.Context, p peer.ID, pmes *pb.Message) (*pb.Message, error)

	// Send a message to a peer without requiring a response
	SendMessage(ctx context.Context, p peer.ID, pmes *pb.Message) error

	// Disconnect from the given peer
	DisconnectFromPeer(p peer.ID) error
}

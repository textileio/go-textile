package net

import (
	"context"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/pb"
	inet "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

var log = logging.MustGetLogger("net")

type NetworkService interface {
	HandleNewStream(s inet.Stream)
	HandlerForMsgType(t pb.Message_Type) func(peer.ID, *pb.Message, interface{}) (*pb.Message, error)
	SendRequest(ctx context.Context, p peer.ID, pmes *pb.Message) (*pb.Message, error)
	SendMessage(ctx context.Context, p peer.ID, pmes *pb.Message) error
	DisconnectFromPeer(p peer.ID) error
}

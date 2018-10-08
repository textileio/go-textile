package net

import (
	"context"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/pb"
	inet "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

var log = logging.MustGetLogger("net")

type NetworkService interface {
	HandleNewStream(s inet.Stream)
	HandlerForMsgType(t pb.Message_Type) func(peer.ID, *pb.Envelope, interface{}) (*pb.Envelope, error)
	SendRequest(ctx context.Context, p peer.ID, pmes *pb.Envelope) (*pb.Envelope, error)
	SendMessage(ctx context.Context, p peer.ID, pmes *pb.Envelope) error
	DisconnectFromPeer(p peer.ID) error
}

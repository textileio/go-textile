package net

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/net/service"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"time"
)

// CafeService is a libp2p service for orchestrating a collection of files with annotations
// amongst a group of peers
type CafeService struct {
	service *service.Service
}

// NewCafeService returns a new threads service
func NewCafeService(node *core.IpfsNode, datastore repo.Datastore) *CafeService {
	handler := &CafeService{}
	handler.service = service.NewService(handler, node, datastore)
	return handler
}

// Protocol returns the handler protocol
func (h *CafeService) Protocol() protocol.ID {
	return protocol.ID("/textile/cafe/1.0.0")
}

// Handle is called by the underlying service handler method
func (h *CafeService) Handle(mtype pb.Message_Type) func(*service.Service, peer.ID, *pb.Envelope) (*pb.Envelope, error) {
	switch mtype {
	case pb.Message_CAFE_CHALLENGE_REQUEST:
		return h.handleChallenge
	default:
		return nil
	}
}

func (h *CafeService) RequestChallenge(kp *keypair.Full, pid peer.ID) error {
	env, err := h.service.NewEnvelope(&pb.CafeChallengeRequest{
		Address: kp.Address(),
	}, pb.Message_CAFE_CHALLENGE_RESPONSE)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err = h.service.SendRequest(ctx, pid, env)
	if err != nil {
		return err
	}
	return nil
}

// handleChallenge receives a challenge request
func (h *CafeService) handleChallenge(s *service.Service, pid peer.ID, pmes *pb.Envelope) (*pb.Envelope, error) {
	log.Debug("received CAFE_CHALLENGE message")
	signed, err := unpackThreadMessage(pmes)
	if err != nil {
		return nil, err
	}
	req := new(pb.CafeChallengeRequest)
	if err := proto.Unmarshal(signed.Block, req); err != nil {
		return nil, err
	}

	// validate address
	accnt, err := keypair.Parse(req.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we dont want to handle account seeds, just addresses
		errMsg, err := s.NewErrorMessage(400, "invalid address")
		if err != nil {
			return nil, err
		}
		return errMsg, nil
	}

	// return a wrapped response
	return s.NewEnvelope(&pb.CafeChallengeResponse{Value: ksuid.New().String()}, pb.Message_CAFE_CHALLENGE_RESPONSE)
}

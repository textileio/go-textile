package service

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	inet "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	"gx/ipfs/QmTKsRYeY4simJyf37K93juSq75Lo8MVCDJ7owjmf46u8W/go-context/io"
	ggio "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/io"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"io"
	"sync"
	"time"
)

var log = logging.MustGetLogger("net")

// Service represents a libp2p service
type Service struct {
	node      *core.IpfsNode
	datastore repo.Datastore
	handler   Handler
	sender    map[peer.ID]*sender
	senderMux sync.Mutex
}

// Handler is used to handle messages for a specific protocol
type Handler interface {
	Protocol() protocol.ID
	Handle(mtype pb.Message_Type) func(*Service, peer.ID, *pb.Envelope) (*pb.Envelope, error)
}

// NewService returns a service for the given config
func NewService(
	handler Handler,
	node *core.IpfsNode,
	datastore repo.Datastore,
) *Service {
	service := &Service{
		node:      node,
		datastore: datastore,
		handler:   handler,
		sender:    make(map[peer.ID]*sender),
	}
	node.PeerHost.SetStreamHandler(handler.Protocol(), service.handleNewStream)
	log.Infof("registered service: %s", handler.Protocol())
	return service
}

// Node returns the underlying ipfs node
func (s *Service) Node() *core.IpfsNode {
	return s.node
}

// Datastore returns the underlying datastore
func (s *Service) Datastore() repo.Datastore {
	return s.datastore
}

// SendMessage sends a message to a peer
func (s *Service) SendMessage(ctx context.Context, pid peer.ID, pmes *pb.Envelope) error {
	log.Debugf("sending %s message to %s", pmes.Message.Type.String(), pid.Pretty())
	ms, err := s.messageSenderForPeer(pid, s.handler.Protocol())
	if err != nil {
		return err
	}
	if err := ms.SendMessage(ctx, pmes); err != nil {
		return err
	}
	return nil
}

// SendRequest sends a request to a peer
func (s *Service) SendRequest(ctx context.Context, pid peer.ID, pmes *pb.Envelope) (*pb.Envelope, error) {
	log.Debugf("sending %s request to %s", pmes.Message.Type.String(), pid.Pretty())
	ms, err := s.messageSenderForPeer(pid, s.handler.Protocol())
	if err != nil {
		return nil, err
	}
	rpmes, err := ms.SendRequest(ctx, pmes)
	if err != nil {
		log.Debugf("no response from %s", pid.Pretty())
		return nil, err
	}
	if rpmes == nil {
		log.Debugf("no response from %s", pid.Pretty())
		return nil, errors.New("no response from peer")
	}
	log.Debugf("received response from %s", pid.Pretty())
	return rpmes, nil
}

// DisconnectFromPeer attempts to disconnect from the given peer
func (s *Service) DisconnectFromPeer(pid peer.ID) error {
	log.Debugf("disconnecting from %s", pid.Pretty())
	s.senderMux.Lock()
	defer s.senderMux.Unlock()
	ms, ok := s.sender[pid]
	if !ok {
		return nil
	}
	if ms != nil && ms.stream != nil {
		ms.stream.Close()
	}
	delete(s.sender, pid)
	return nil
}

// NewEnvelope returns a signed pb message
func (s *Service) NewEnvelope(message proto.Message, mtype pb.Message_Type) (*pb.Envelope, error) {
	payload, err := ptypes.MarshalAny(message)
	if err != nil {
		return nil, err
	}
	msg := &pb.Message{Type: mtype, Payload: payload}
	serialized, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	authorSig, err := s.node.PrivateKey.Sign(serialized)
	if err != nil {
		return nil, err
	}
	authorPk, err := s.node.PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	return &pb.Envelope{Message: msg, Pk: authorPk, Sig: authorSig}, nil
}

// NewErrorMessage returns a signed pb error message
func (s *Service) NewErrorMessage(code int, msg string) (*pb.Envelope, error) {
	return s.NewEnvelope(&pb.Error{Code: uint32(code), Message: msg}, pb.Message_ERROR)
}

// handleNewStream handles a p2p net stream in the background
func (s *Service) handleNewStream(stream inet.Stream) {
	go s.handleNewMessage(stream)
}

// handleNewMessage handles a p2p net stream
func (s *Service) handleNewMessage(stream inet.Stream) {
	defer stream.Close()

	// setup reader
	ctxReader := ctxio.NewReader(s.node.Context(), stream)
	reader := ggio.NewDelimitedReader(ctxReader, inet.MessageSizeMax)

	// get sender
	mPeer := stream.Conn().RemotePeer()
	ms, err := s.messageSenderForPeer(mPeer, s.handler.Protocol())
	if err != nil {
		log.Error("error getting message sender")
		return
	}

	// start listening for messages from this sender
	for {
		select {
		// end loop on context close
		case <-s.node.Context().Done():
			return
		default:
		}

		// receive msg
		pmes := new(pb.Envelope)
		if err := reader.ReadMsg(pmes); err != nil {
			stream.Reset()
			if err == io.EOF {
				log.Debugf("disconnected from peer %s", mPeer.Pretty())
			}
			return
		}

		// validate the signature
		ser, err := proto.Marshal(pmes.Message)
		if err != nil {
			log.Errorf("marshal error %s", err)
			return
		}
		pk, err := libp2pc.UnmarshalPublicKey(pmes.Pk)
		if err != nil {
			log.Errorf("unmarshal pubkey error %s", err)
			return
		}
		if err := crypto.Verify(pk, ser, pmes.Sig); err != nil {
			log.Warningf("invalid signature %s", err)
			return
		}

		// check if the message is a response
		if pmes.Message.IsResponse {
			ms.requestMux.Lock()
			ch, ok := ms.requests[pmes.Message.RequestId]
			if ok {
				// this is a request response
				select {
				case ch <- pmes:
					// message returned to requester
				case <-time.After(time.Second):
					// in case ch is closed on the other end - the lock should prevent this happening
					log.Debug("request id was not removed from map on timeout")
				}
				close(ch)
				delete(ms.requests, pmes.Message.RequestId)
			} else {
				log.Debug("received response message with unknown request id: requesting function may have timed out")
			}
			ms.requestMux.Unlock()
			stream.Reset()
			return
		}

		// try a generic handler for this msg type
		handler := s.handleGeneric(pmes.Message.Type)
		if handler == nil {
			// get service specific handler for this msg type
			handler := s.handler.Handle(pmes.Message.Type)
			if handler == nil {
				stream.Reset()
				log.Debug("got back nil handler")
				return
			}
		}

		// dispatch handler
		rpmes, err := handler(s, mPeer, pmes)
		if err != nil {
			log.Errorf("%s handle message error: %s", pmes.Message.Type.String(), err)
		}

		// if nil response, return it before serializing
		if rpmes == nil {
			continue
		}

		// give back request id
		rpmes.Message.RequestId = pmes.Message.RequestId
		rpmes.Message.IsResponse = true

		// send out response msg
		if err := ms.SendMessage(s.node.Context(), rpmes); err != nil {
			stream.Reset()
			log.Errorf("send response error: %s", err)
			return
		}
	}
}

// handleGeneric provides service level handlers for common message types
func (s *Service) handleGeneric(mtype pb.Message_Type) func(*Service, peer.ID, *pb.Envelope) (*pb.Envelope, error) {
	switch mtype {
	case pb.Message_PING:
		return handlePing
	case pb.Message_ERROR:
		return handleError
	default:
		return nil
	}
}

// handlePing receives a PING message
func handlePing(_ *Service, pid peer.ID, pmes *pb.Envelope) (*pb.Envelope, error) {
	log.Debugf("received PING message from %h", pid.Pretty())
	return pmes, nil
}

// handleError receives an ERROR message
func handleError(_ *Service, pid peer.ID, pmes *pb.Envelope) (*pb.Envelope, error) {
	log.Debugf("received ERROR message from %h", pid.Pretty())
	if pmes.Message.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	errorMessage := new(pb.Error)
	err := ptypes.UnmarshalAny(pmes.Message.Payload, errorMessage)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

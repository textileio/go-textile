package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"gx/ipfs/QmTKsRYeY4simJyf37K93juSq75Lo8MVCDJ7owjmf46u8W/go-context/io"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core"
	inet "gx/ipfs/QmXuRkCR7BNQa9uqfpTiFWsTQLzmTWYg91Ja1w95gnqb6u/go-libp2p-net"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	ggio "gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/io"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
)

var log = logging.Logger("tex-service")

// service represents a libp2p service
type Service struct {
	Account   *keypair.Full
	Node      *core.IpfsNode
	handler   Handler
	sender    map[peer.ID]*sender
	senderMux sync.Mutex
}

// defaultTimeout is the context timeout for sending / requesting messages
const defaultTimeout = time.Second * 5

// PeerStatus is the possible results from pinging another peer
type PeerStatus string

const (
	PeerOnline  PeerStatus = "online"
	PeerOffline PeerStatus = "offline"
)

// Handler is used to handle messages for a specific protocol
type Handler interface {
	Protocol() protocol.ID
	Ping(pid peer.ID) (PeerStatus, error)
	Handle(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error)
}

// NewService returns a service for the given config
func NewService(account *keypair.Full, handler Handler, node *core.IpfsNode) *Service {
	service := &Service{
		Account: account,
		Node:    node,
		handler: handler,
		sender:  make(map[peer.ID]*sender),
	}
	node.PeerHost.SetStreamHandler(handler.Protocol(), service.handleNewStream)
	log.Debugf("registered service: %s", handler.Protocol())
	return service
}

// SendMessage sends a message to a peer
func (s *Service) SendMessage(pid peer.ID, env *pb.Envelope) error {
	log.Debugf("sending %s to %s", env.Message.Type.String(), pid.Pretty())
	ms, err := s.messageSenderForPeer(pid, s.handler.Protocol())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if err := ms.SendMessage(ctx, env); err != nil {
		return err
	}
	return nil
}

// SendRequest sends a request to a peer
func (s *Service) SendRequest(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	log.Debugf("sending %s to %s", env.Message.Type.String(), pid.Pretty())
	ms, err := s.messageSenderForPeer(pid, s.handler.Protocol())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	renv, err := ms.SendRequest(ctx, env)
	if err != nil {
		return nil, err
	}
	if renv == nil {
		log.Debugf("no response from %s", pid.Pretty())
		return nil, errors.New("no response from peer")
	}
	log.Debugf("received %s response from %s", renv.Message.Type.String(), pid.Pretty())
	if err := s.handleError(renv); err != nil {
		return nil, err
	}
	return renv, nil
}

// Ping pings another peer and returns status
func (s *Service) Ping(pid peer.ID) (PeerStatus, error) {
	id := rand.Int31()
	env, err := s.NewEnvelope(pb.Message_PING, nil, &id, false)
	if err != nil {
		return "", err
	}
	if _, err := s.SendRequest(pid, env); err != nil {
		return PeerOffline, nil
	}
	return PeerOnline, nil
}

// NewEnvelope returns a signed pb message for transport
func (s *Service) NewEnvelope(mtype pb.Message_Type, msg proto.Message, id *int32, response bool) (*pb.Envelope, error) {
	var payload *any.Any
	if msg != nil {
		var err error
		payload, err = ptypes.MarshalAny(msg)
		if err != nil {
			return nil, err
		}
	}
	message := &pb.Message{Type: mtype, Payload: payload}
	if id != nil {
		message.RequestId = *id
	}
	if response {
		message.IsResponse = true
	}
	ser, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	if s.Node.PrivateKey == nil {
		if err := s.Node.LoadPrivateKey(); err != nil {
			return nil, err
		}
	}
	sig, err := s.Node.PrivateKey.Sign(ser)
	if err != nil {
		return nil, err
	}
	return &pb.Envelope{Message: message, Sig: sig}, nil
}

// NewError returns a signed pb error message
func (s *Service) NewError(code int, msg string, id int32) (*pb.Envelope, error) {
	return s.NewEnvelope(pb.Message_ERROR, &pb.Error{
		Code:    uint32(code),
		Message: msg,
	}, &id, true)
}

// VerifyEnvelope verifies the authenticity of an envelope
func (s *Service) VerifyEnvelope(env *pb.Envelope, pid peer.ID) error {
	ser, err := proto.Marshal(env.Message)
	if err != nil {
		return err
	}
	pk, err := pid.ExtractPublicKey()
	if err != nil {
		return err
	}
	return crypto.Verify(pk, ser, env.Sig)
}

// handleError receives an error response
func (s *Service) handleError(env *pb.Envelope) error {
	if env.Message.Payload == nil && env.Message.Type != pb.Message_PONG {
		err := fmt.Sprintf("message payload with type %s is nil", env.Message.Type.String())
		log.Error(err)
		return errors.New(err)
	}
	if env.Message.Type != pb.Message_ERROR {
		return nil
	} else {
		errMsg := new(pb.Error)
		if err := ptypes.UnmarshalAny(env.Message.Payload, errMsg); err != nil {
			return err
		}
		return errors.New(errMsg.Message)
	}
}

// handleNewStream handles a p2p net stream in the background
func (s *Service) handleNewStream(stream inet.Stream) {
	go s.handleNewMessage(stream)
}

// handleNewMessage handles a p2p net stream
func (s *Service) handleNewMessage(stream inet.Stream) {
	defer stream.Close()

	// setup reader
	ctxReader := ctxio.NewReader(s.Node.Context(), stream)
	reader := ggio.NewDelimitedReader(ctxReader, inet.MessageSizeMax)

	// get sender
	rpid := stream.Conn().RemotePeer()
	ms, err := s.messageSenderForPeer(rpid, s.handler.Protocol())
	if err != nil {
		log.Errorf("error getting message sender: %s", err)
		return
	}

	// start listening for messages from this sender
	for {
		select {
		// end loop on context close
		case <-s.Node.Context().Done():
			return
		default:
		}

		env := new(pb.Envelope)
		if err := reader.ReadMsg(env); err != nil {
			stream.Reset()
			if err == io.EOF {
				log.Debugf("disconnected from peer %s", rpid.Pretty())
			}
			return
		}

		if err := s.VerifyEnvelope(env, rpid); err != nil {
			log.Warningf("error verifying message: %s", err)
			continue
		}

		// check if the message is a response
		if env.Message.IsResponse {
			ms.requestMux.Lock()

			ch, ok := ms.requests[env.Message.RequestId]
			if ok {
				// this is a request response
				select {
				case ch <- env:
					// message returned to requester
				case <-time.After(defaultTimeout):
					// in case ch is closed on the other end - the lock should prevent this happening
					log.Debug("request id was not removed from map on timeout")
				}

				close(ch)
				delete(ms.requests, env.Message.RequestId)
			} else {
				log.Debug("unknown request id: requesting function may have timed out")
			}

			ms.requestMux.Unlock()
			stream.Reset()
			return
		}

		// try a core handler for this msg type
		handler := s.handleCore(env.Message.Type)
		if handler == nil {
			// get service specific handler
			handler = s.handler.Handle
		}

		log.Debugf("received %s from %s", env.Message.Type.String(), rpid.Pretty())
		renv, err := handler(rpid, env)
		if err != nil {
			log.Errorf("%s handle message error: %s", env.Message.Type.String(), err)
		}
		if renv == nil {
			continue
		}

		log.Debugf("responding with %s to %s", renv.Message.Type.String(), rpid.Pretty())
		if err := ms.SendMessage(s.Node.Context(), renv); err != nil {
			stream.Reset()
			log.Errorf("send response error: %s", err)
			return
		}
	}
}

// handleCore provides service level handlers for common message types
func (s *Service) handleCore(mtype pb.Message_Type) func(peer.ID, *pb.Envelope) (*pb.Envelope, error) {
	switch mtype {
	case pb.Message_PING:
		return s.handlePing
	default:
		return nil
	}
}

// handlePing receives a PING message
func (s *Service) handlePing(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	log.Debugf("received PING message from %s", pid.Pretty())
	return s.NewEnvelope(pb.Message_PONG, nil, &env.Message.RequestId, true)
}

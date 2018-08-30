package service

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/thread"
	"gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
	"gx/ipfs/QmTKsRYeY4simJyf37K93juSq75Lo8MVCDJ7owjmf46u8W/go-context/io"
	ps "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	inet "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	ggio "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/io"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"io"
	"sync"
	"time"
)

var log = logging.MustGetLogger("service")

var TextileProtocol protocol.ID = "/textile/app/1.0.0"

type TextileService struct {
	host      host.Host
	self      peer.ID
	peerstore ps.Peerstore
	ctx       context.Context
	datastore repo.Datastore
	node      *core.IpfsNode
	getThread func(string) (*int, *thread.Thread)
	notify    func(notificaiton *repo.Notification) error
	sender    map[peer.ID]*sender
	senderlk  sync.Mutex
}

func NewService(
	node *core.IpfsNode,
	datastore repo.Datastore,
	getThread func(string) (*int, *thread.Thread),
	notify func(notificaiton *repo.Notification) error,
) *TextileService {
	service := &TextileService{
		host:      node.PeerHost.(host.Host),
		self:      node.Identity,
		peerstore: node.PeerHost.Peerstore(),
		ctx:       node.Context(),
		datastore: datastore,
		node:      node,
		getThread: getThread,
		notify:    notify,
		sender:    make(map[peer.ID]*sender),
	}
	node.PeerHost.SetStreamHandler(TextileProtocol, service.HandleNewStream)
	log.Infof("textile service running at %s", TextileProtocol)
	return service
}

func (s *TextileService) DisconnectFromPeer(p peer.ID) error {
	log.Debugf("disconnecting from %s", p.Pretty())
	s.senderlk.Lock()
	defer s.senderlk.Unlock()
	ms, ok := s.sender[p]
	if !ok {
		return nil
	}
	if ms != nil && ms.s != nil {
		ms.s.Close()
	}
	delete(s.sender, p)
	return nil
}

func (s *TextileService) HandleNewStream(stream inet.Stream) {
	go s.handleNewMessage(stream, true)
}

func (s *TextileService) handleNewMessage(stream inet.Stream, incoming bool) {
	defer func() {
		if recover() != nil {
			log.Error("recovered from handleNewMessage")
		}
	}()
	defer stream.Close()
	cr := ctxio.NewReader(s.ctx, stream) // ok to use. we defer close stream in this func
	r := ggio.NewDelimitedReader(cr, inet.MessageSizeMax)
	mPeer := stream.Conn().RemotePeer()

	ms, err := s.messageSenderForPeer(mPeer)
	if err != nil {
		log.Error("error getting message sender")
		return
	}

	for {
		select {
		// end loop on context close
		case <-s.ctx.Done():
			return
		default:
		}
		// Receive msg
		pmes := new(pb.Envelope)
		if err := r.ReadMsg(pmes); err != nil {
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

		if pmes.Message.IsResponse {
			ms.requestlk.Lock()
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
			ms.requestlk.Unlock()
			stream.Reset()
			return
		}

		// Get handler for this msg type
		handler := s.HandlerForMsgType(pmes.Message.Type)
		if handler == nil {
			stream.Reset()
			log.Debug("got back nil handler from handlerForMsgType")
			return
		}

		// Dispatch handler
		rpmes, err := handler(mPeer, pmes, nil)
		if err != nil {
			log.Errorf("%s handle message error: %s", pmes.Message.Type.String(), err)
		}

		// If nil response, return it before serializing
		if rpmes == nil {
			continue
		}

		// give back request id
		rpmes.Message.RequestId = pmes.Message.RequestId
		rpmes.Message.IsResponse = true

		// send out response msg
		if err := ms.SendMessage(s.ctx, rpmes); err != nil {
			stream.Reset()
			log.Errorf("send response error: %s", err)
			return
		}
	}
}

func (s *TextileService) SendRequest(ctx context.Context, p peer.ID, pmes *pb.Envelope) (*pb.Envelope, error) {
	defer func() {
		if recover() != nil {
			log.Error("recovered from service.SendRequest")
		}
	}()
	log.Debugf("sending %s request to %s", pmes.Message.Type.String(), p.Pretty())
	ms, err := s.messageSenderForPeer(p)
	if err != nil {
		return nil, err
	}

	rpmes, err := ms.SendRequest(ctx, pmes)
	if err != nil {
		log.Debugf("no response from %s", p.Pretty())
		return nil, err
	}

	if rpmes == nil {
		log.Debugf("no response from %s", p.Pretty())
		return nil, errors.New("no response from peer")
	}

	log.Debugf("received response from %s", p.Pretty())
	return rpmes, nil
}

func (s *TextileService) SendMessage(ctx context.Context, p peer.ID, pmes *pb.Envelope) error {
	defer func() {
		if recover() != nil {
			log.Error("recovered from service.SendMessage")
		}
	}()
	log.Debugf("sending %s message to %s", pmes.Message.Type.String(), p.Pretty())
	ms, err := s.messageSenderForPeer(p)
	if err != nil {
		return err
	}

	if err := ms.SendMessage(ctx, pmes); err != nil {
		return err
	}
	return nil
}

func (s *TextileService) newEnvelope(message *pb.Message) (*pb.Envelope, error) {
	// sign it
	serialized, err := proto.Marshal(message)
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
	return &pb.Envelope{Message: message, Pk: authorPk, Sig: authorSig}, nil
}

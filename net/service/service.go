package service

import (
	"context"
	"errors"
	"github.com/op/go-logging"
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
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	"io"
	"sync"
	"time"
)

var log = logging.MustGetLogger("service")

var ProtocolTextile protocol.ID = "/textile/app/1.0.0"

type TextileService struct {
	host      host.Host
	self      peer.ID
	peerstore ps.Peerstore
	ctx       context.Context
	datastore repo.Datastore
	node      *core.IpfsNode
	getThread func(string) *thread.Thread
	addThread func(string, libp2pc.PrivKey) (*thread.Thread, error)
	sender    map[peer.ID]*sender
	senderlk  sync.Mutex
}

func NewService(
	node *core.IpfsNode,
	datastore repo.Datastore,
	getThread func(string) *thread.Thread,
	addThread func(string, libp2pc.PrivKey) (*thread.Thread, error),
) *TextileService {
	service := &TextileService{
		host:      node.PeerHost.(host.Host),
		self:      node.Identity,
		peerstore: node.PeerHost.Peerstore(),
		ctx:       node.Context(),
		datastore: datastore,
		node:      node,
		getThread: getThread,
		addThread: addThread,
		sender:    make(map[peer.ID]*sender),
	}
	node.PeerHost.SetStreamHandler(ProtocolTextile, service.HandleNewStream)
	log.Infof("textile service running at %s", ProtocolTextile)
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
		pmes := new(pb.Message)
		if err := r.ReadMsg(pmes); err != nil {
			stream.Reset()
			if err == io.EOF {
				log.Debugf("disconnected from peer %s", mPeer.Pretty())
			}
			return
		}

		if pmes.IsResponse {
			ms.requestlk.Lock()
			ch, ok := ms.requests[pmes.RequestId]
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
				delete(ms.requests, pmes.RequestId)
			} else {
				log.Debug("received response message with unknown request id: requesting function may have timed out")
			}
			ms.requestlk.Unlock()
			stream.Reset()
			return
		}

		// Get handler for this msg type
		handler := s.HandlerForMsgType(pmes.MessageType)
		if handler == nil {
			stream.Reset()
			log.Debug("got back nil handler from handlerForMsgType")
			return
		}

		// Dispatch handler
		rpmes, err := handler(mPeer, pmes, nil)
		if err != nil {
			log.Debugf("%s handle message error: %s", pmes.MessageType.String(), err)
		}

		// If nil response, return it before serializing
		if rpmes == nil {
			continue
		}

		// give back request id
		rpmes.RequestId = pmes.RequestId
		rpmes.IsResponse = true

		// send out response msg
		if err := ms.SendMessage(s.ctx, rpmes); err != nil {
			stream.Reset()
			log.Debugf("send response error: %s", err)
			return
		}
	}
}

func (s *TextileService) SendRequest(ctx context.Context, p peer.ID, pmes *pb.Message) (*pb.Message, error) {
	log.Debugf("sending %s request to %s", pmes.MessageType.String(), p.Pretty())
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

func (s *TextileService) SendMessage(ctx context.Context, p peer.ID, pmes *pb.Message) error {
	log.Debugf("sending %s message to %s", pmes.MessageType.String(), p.Pretty())
	ms, err := s.messageSenderForPeer(p)
	if err != nil {
		return err
	}

	if err := ms.SendMessage(ctx, pmes); err != nil {
		return err
	}
	return nil
}

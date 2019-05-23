package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	ggio "github.com/gogo/protobuf/io"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/ipfs/go-ipfs/core"
	logging "github.com/ipfs/go-log"
	iface "github.com/ipfs/interface-go-ipfs-core"
	ctxio "github.com/jbenet/go-context/io"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

var log = logging.Logger("tex-service")

// service represents a libp2p service
type Service struct {
	Account *keypair.Full
	Node    func() *core.IpfsNode

	handler Handler

	strmap map[peer.ID]*messageSender
	smlk   sync.Mutex
}

// defaultTimeout is the context timeout for sending / requesting messages
const defaultTimeout = time.Second * 30

// PeerStatus is the possible results from pinging another peer
type PeerStatus string

const (
	PeerOnline  PeerStatus = "online"
	PeerOffline PeerStatus = "offline"
)

// Handler is used to handle messages for a specific protocol
type Handler interface {
	Protocol() protocol.ID
	Start()
	Ping(pid peer.ID) (PeerStatus, error)
	Handle(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error)
	HandleStream(pid peer.ID, env *pb.Envelope) (chan *pb.Envelope, chan error, chan interface{})
}

// NewService returns a service for the given config
func NewService(account *keypair.Full, handler Handler, node func() *core.IpfsNode) *Service {
	return &Service{
		Account: account,
		Node:    node,
		handler: handler,
		strmap:  make(map[peer.ID]*messageSender),
	}
}

// Start sets the peer host stream handler
func (srv *Service) Start() {
	srv.Node().PeerHost.SetStreamHandler(srv.handler.Protocol(), srv.handleNewStream)
	go srv.listen("")
	go srv.listen(srv.Node().Identity.Pretty())
}

// Ping pings another peer and returns status
func (srv *Service) Ping(p peer.ID) (PeerStatus, error) {
	id := rand.Int31()
	env, err := srv.NewEnvelope(pb.Message_PING, nil, &id, false)
	if err != nil {
		return "", err
	}

	if _, err := srv.SendRequest(p, env); err != nil {
		return PeerOffline, nil
	}

	return PeerOnline, nil
}

// SendRequest sends out a request
func (srv *Service) SendRequest(p peer.ID, pmes *pb.Envelope) (*pb.Envelope, error) {
	log.Debugf("sending %s to %s", pmes.Message.Type.String(), p.Pretty())

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	ms, err := srv.messageSenderForPeer(ctx, p)
	if err != nil {
		return nil, err
	}

	rpmes, err := ms.SendRequest(ctx, pmes)
	if err != nil {
		return nil, err
	}

	if rpmes == nil {
		err := fmt.Errorf("no response from %s", p.Pretty())
		log.Debug(err.Error())
		return nil, err
	}

	log.Debugf("received %s response from %s", rpmes.Message.Type.String(), p.Pretty())
	if err := srv.handleError(rpmes); err != nil {
		return nil, err
	}

	return rpmes, nil
}

// SendHTTPRequest sends a request over HTTP
func (srv *Service) SendHTTPRequest(addr string, pmes *pb.Envelope) (*pb.Envelope, error) {
	log.Debugf("sending %s to %s", pmes.Message.Type.String(), addr)

	payload, err := proto.Marshal(pmes)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", addr, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Textile-Peer", srv.Node().Identity.Pretty())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		res, err := util.UnmarshalString(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(res)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	rpmes := new(pb.Envelope)
	if err := proto.Unmarshal(body, rpmes); err != nil {
		return nil, err
	}

	if rpmes.Message == nil {
		err := fmt.Errorf("no response from %s", addr)
		log.Debug(err.Error())
		return nil, err
	}

	log.Debugf("received %s response from %s", rpmes.Message.Type.String(), addr)
	if err := srv.handleError(rpmes); err != nil {
		return nil, err
	}

	return rpmes, nil
}

// SendHTTPStreamRequest sends a request over HTTP
func (srv *Service) SendHTTPStreamRequest(addr string, pmes *pb.Envelope) (chan *pb.Envelope, chan error, *func()) {
	rpmesCh := make(chan *pb.Envelope)
	errCh := make(chan error)

	var cancel func()
	go func() {
		defer close(rpmesCh)
		log.Debugf("sending %s to %s", pmes.Message.Type.String(), addr)

		payload, err := proto.Marshal(pmes)
		if err != nil {
			errCh <- err
			return
		}

		req, err := http.NewRequest("POST", addr, bytes.NewReader(payload))
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("X-Textile-Peer", srv.Node().Identity.Pretty())

		tr := &http.Transport{}
		client := &http.Client{Transport: tr}

		cancel = func() {
			tr.CancelRequest(req)
		}

		res, err := client.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			res, err := util.UnmarshalString(res.Body)
			if err != nil {
				errCh <- err
			} else {
				errCh <- fmt.Errorf(res)
			}
			return
		}

		decoder := json.NewDecoder(res.Body)
		for decoder.More() {
			var rpmes *pb.Envelope
			if err := decoder.Decode(&rpmes); err == io.EOF {
				return
			} else if err != nil {
				errCh <- err
				return
			}

			if rpmes == nil || rpmes.Message == nil {
				err := fmt.Errorf("no response from %s", addr)
				log.Debug(err.Error())
				errCh <- err
				return
			}

			log.Debugf("received %s response from %s", rpmes.Message.Type.String(), addr)
			if err := srv.handleError(rpmes); err != nil {
				errCh <- err
				return
			}
			rpmesCh <- rpmes
		}
	}()

	return rpmesCh, errCh, &cancel
}

// SendMessage sends out a message
func (srv *Service) SendMessage(ctx context.Context, p peer.ID, pmes *pb.Envelope) error {
	log.Debugf("sending %s to %s", pmes.Message.Type.String(), p.Pretty())

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
	}

	ms, err := srv.messageSenderForPeer(ctx, p)
	if err != nil {
		return err
	}

	if err := ms.SendMessage(ctx, pmes); err != nil {
		return err
	}

	return nil
}

// SendHTTPMessage sends a message over HTTP
func (srv *Service) SendHTTPMessage(addr string, pmes *pb.Envelope) error {
	log.Debugf("sending %s to %s", pmes.Message.Type.String(), addr)

	payload, err := proto.Marshal(pmes)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", addr, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("X-Textile-Peer", srv.Node().Identity.Pretty())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		res, err := util.UnmarshalString(req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(res)
	}

	return nil
}

// NewEnvelope returns a signed pb message for transport
func (srv *Service) NewEnvelope(mtype pb.Message_Type, msg proto.Message, id *int32, response bool) (*pb.Envelope, error) {
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
		message.Request = *id
	}
	if response {
		message.Response = true
	}

	ser, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	sig, err := srv.Node().PrivateKey.Sign(ser)
	if err != nil {
		return nil, err
	}

	return &pb.Envelope{Message: message, Sig: sig}, nil
}

// NewError returns a signed pb error message
func (srv *Service) NewError(code int, msg string, id int32) (*pb.Envelope, error) {
	return srv.NewEnvelope(pb.Message_ERROR, &pb.Error{
		Code:    uint32(code),
		Message: msg,
	}, &id, true)
}

// VerifyEnvelope verifies the authenticity of an envelope
func (srv *Service) VerifyEnvelope(env *pb.Envelope, pid peer.ID) error {
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
func (srv *Service) handleError(env *pb.Envelope) error {
	if env.Message.Payload == nil && env.Message.Type != pb.Message_PONG {
		err := fmt.Errorf("message payload with type %s is nil", env.Message.Type.String())
		log.Error(err.Error())
		return err
	}

	if env.Message.Type != pb.Message_ERROR {
		return nil
	} else {
		errMsg := new(pb.Error)
		if err := ptypes.UnmarshalAny(env.Message.Payload, errMsg); err != nil {
			return err
		}
		return fmt.Errorf(errMsg.Message)
	}
}

// handleCore provides service level handlers for common message types
func (srv *Service) handleCore(mtype pb.Message_Type) func(peer.ID, *pb.Envelope) (*pb.Envelope, error) {
	switch mtype {
	case pb.Message_PING:
		return srv.handlePing
	default:
		return nil
	}
}

// handlePing receives a PING message
func (srv *Service) handlePing(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	return srv.NewEnvelope(pb.Message_PONG, nil, &env.Message.Request, true)
}

var dhtReadMessageTimeout = time.Minute
var ErrReadTimeout = fmt.Errorf("timed out reading response")

type bufferedWriteCloser interface {
	ggio.WriteCloser
	Flush() error
}

// The Protobuf writer performs multiple small writes when writing a message.
// We need to buffer those writes, to make sure that we're not sending a new
// packet for every single write.
type bufferedDelimitedWriter struct {
	*bufio.Writer
	ggio.WriteCloser
}

func newBufferedDelimitedWriter(str io.Writer) bufferedWriteCloser {
	w := bufio.NewWriter(str)
	return &bufferedDelimitedWriter{
		Writer:      w,
		WriteCloser: ggio.NewDelimitedWriter(w),
	}
}

func (w *bufferedDelimitedWriter) Flush() error {
	return w.Writer.Flush()
}

// handleNewStream implements the inet.StreamHandler
func (srv *Service) handleNewStream(s inet.Stream) {
	go srv.handleNewMessage(s)
}

func (srv *Service) handleNewMessage(s inet.Stream) {
	ctx := srv.Node().Context()
	cr := ctxio.NewReader(ctx, s) // ok to use. we defer close stream in this func
	cw := ctxio.NewWriter(ctx, s) // ok to use. we defer close stream in this func
	r := ggio.NewDelimitedReader(cr, inet.MessageSizeMax)
	w := newBufferedDelimitedWriter(cw)
	mPeer := s.Conn().RemotePeer()

	for {
		select {
		// end loop on context close
		case <-srv.Node().Context().Done():
			return
		default:
		}

		// receive msg
		pmes := new(pb.Envelope)
		switch err := r.ReadMsg(pmes); err {
		case io.EOF:
			_ = s.Close()
			return
		case nil:
		default:
			_ = s.Reset()
			log.Debugf("error unmarshaling data: %s", err)
			return
		}

		if err := srv.VerifyEnvelope(pmes, mPeer); err != nil {
			log.Warningf("error verifying message: %s", err)
			continue
		}

		// try a core handler for this msg type
		handler := srv.handleCore(pmes.Message.Type)
		if handler == nil {
			// get service specific handler
			handler = srv.handler.Handle
			// TODO: handle stream requests over p2p?
		}

		log.Debugf("received %s from %s", pmes.Message.Type.String(), mPeer.Pretty())
		rpmes, err := handler(mPeer, pmes)
		if err != nil {
			_ = s.Reset()
			log.Errorf("%s handle message error: %s", pmes.Message.Type.String(), err)
			return
		}

		// if nil response, return it before serializing
		if rpmes == nil {
			continue
		}

		// send out response msg
		log.Debugf("responding with %s to %s", rpmes.Message.Type.String(), mPeer.Pretty())

		// send out response msg
		err = w.WriteMsg(rpmes)
		if err == nil {
			err = w.Flush()
		}
		if err != nil {
			_ = s.Reset()
			log.Errorf("send response error: %s", err)
			return
		}
	}
}

// listen subscribes to the protocol tag for network-wide requests
func (srv *Service) listen(tag string) {
	topic := string(srv.handler.Protocol())
	if tag != "" {
		topic += "/" + tag
	}

	msgs := make(chan iface.PubSubMessage, 10)
	ctx := srv.Node().Context()
	go func() {
		if err := ipfs.Subscribe(srv.Node(), ctx, topic, true, msgs); err != nil {
			close(msgs)
			log.Errorf("pubsub service listener stopped with error: %s")
			return
		}
	}()
	log.Infof("pubsub service listener started for %s", topic)

	for {
		select {
		// end loop on context close
		case <-srv.Node().Context().Done():
			log.Debugf("pubsub listener shutdown for %s", topic)
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Debugf("pubsub listener shutdown for %s", topic)
				return
			}

			mPeer := msg.From()
			if mPeer.Pretty() == srv.Node().Identity.Pretty() {
				continue
			}

			pmes := new(pb.Envelope)
			if err := proto.Unmarshal(msg.Data(), pmes); err != nil {
				log.Errorf("error unmarshaling pubsub message data from %s: %s", mPeer.Pretty(), err)
				continue
			}

			if err := srv.VerifyEnvelope(pmes, mPeer); err != nil {
				log.Warningf("error verifying message: %s", err)
				continue
			}

			// try a core handler for this msg type
			handler := srv.handleCore(pmes.Message.Type)
			if handler == nil {
				// get service specific handler
				handler = srv.handler.Handle
			}

			log.Debugf("received pubsub %s from %s", pmes.Message.Type.String(), mPeer.Pretty())
			rpmes, err := handler(mPeer, pmes)
			if err != nil {
				log.Errorf("%s handle message error: %s", pmes.Message.Type.String(), err)
				continue
			}

			// if nil response, return it before serializing
			if rpmes == nil {
				continue
			}

			// send out response msg
			log.Debugf("responding with %s to %s", rpmes.Message.Type.String(), mPeer.Pretty())

			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			if err := srv.SendMessage(ctx, mPeer, rpmes); err != nil {
				log.Errorf("error sending message response to %s: %s", mPeer, err)
			}
			cancel()
		}
	}
}

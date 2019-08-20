package core

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	njwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	icid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	iface "github.com/ipfs/interface-go-ipfs-core"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	"github.com/mr-tron/base58/base58"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/common"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/jwt"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/repo/config"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/service"
	"golang.org/x/crypto/bcrypt"
)

// defaultSessionDuration after which session token expires
const defaultSessionDuration = time.Hour * 24 * 7 * 4

// maxRequestAttempts is the number of times a request can fail before being deleted
const maxRequestAttempts = 5

// inboxMessagePageSize is the page size used when checking messages
const inboxMessagePageSize = 10

// maxQueryWaitSeconds is used to limit a query request's max wait time
const maxQueryWaitSeconds = 30

// defaultQueryWaitSeconds is a query request's default wait time
const defaultQueryWaitSeconds = 5

// defaultQueryResultsLimit is a query request's default results limit
const defaultQueryResultsLimit = 5

// cafeOutFlushGroupSize is the size of concurrently processed requests
// note: reqs from this group are batched to each cafe
const cafeOutFlushGroupSize = 32

// validation errors
const (
	errInvalidAddress = "invalid address"
	errUnauthorized   = "unauthorized"
	errForbidden      = "forbidden"
	errBadRequest     = "bad request"
)

// cafeServiceProtocol is the current protocol tag
const cafeServiceProtocol = protocol.ID("/textile/cafe/1.0.0")

// CafeService is a libp2p pinning and offline message service
type CafeService struct {
	service         *service.Service
	datastore       repo.Datastore
	inbox           *CafeInbox
	info            *pb.Cafe
	online          bool
	open            bool
	queryResults    *broadcast.Broadcaster
	inFlightQueries map[string]struct{}
}

// NewCafeService returns a new threads service
func NewCafeService(
	account *keypair.Full,
	node func() *core.IpfsNode,
	datastore repo.Datastore,
	inbox *CafeInbox,
) *CafeService {
	handler := &CafeService{
		datastore:       datastore,
		inbox:           inbox,
		queryResults:    broadcast.NewBroadcaster(10),
		inFlightQueries: make(map[string]struct{}),
	}
	handler.service = service.NewService(account, handler, node)
	return handler
}

// Protocol returns the handler protocol
func (h *CafeService) Protocol() protocol.ID {
	return cafeServiceProtocol
}

// Start begins online services
func (h *CafeService) Start() {
	h.service.Start()
}

// Ping pings another peer
func (h *CafeService) Ping(pid peer.ID) (service.PeerStatus, error) {
	return h.service.Ping(pid.Pretty())
}

// Handle is called by the underlying service handler method
func (h *CafeService) Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	switch env.Message.Type {
	case pb.Message_CAFE_CHALLENGE:
		return h.handleChallenge(env, pid)
	case pb.Message_CAFE_REGISTRATION:
		return h.handleRegistration(env, pid)
	case pb.Message_CAFE_DEREGISTRATION:
		return h.handleDeregistration(env, pid)
	case pb.Message_CAFE_REFRESH_SESSION:
		return h.handleRefreshSession(env, pid)
	case pb.Message_CAFE_STORE:
		return h.handleStore(env, pid)
	case pb.Message_CAFE_UNSTORE:
		return h.handleUnstore(env, pid)
	case pb.Message_CAFE_OBJECT:
		return h.handleObject(env, pid)
	case pb.Message_CAFE_STORE_THREAD:
		return h.handleStoreThread(env, pid)
	case pb.Message_CAFE_UNSTORE_THREAD:
		return h.handleUnstoreThread(env, pid)
	case pb.Message_CAFE_DELIVER_MESSAGE:
		return h.handleDeliverMessage(env, pid)
	case pb.Message_CAFE_CHECK_MESSAGES:
		return h.handleCheckMessages(env, pid)
	case pb.Message_CAFE_DELETE_MESSAGES:
		return h.handleDeleteMessages(env, pid)
	case pb.Message_CAFE_YOU_HAVE_MAIL:
		return h.handleNotifyClient(env, pid)
	case pb.Message_CAFE_PUBLISH_PEER:
		return h.handlePublishPeer(env, pid)
	case pb.Message_CAFE_PUBSUB_QUERY:
		return h.handlePubSubQuery(env, pid)
	case pb.Message_CAFE_PUBSUB_QUERY_RES:
		return h.handlePubSubQueryResults(env, pid)
	default:
		return nil, nil
	}
}

// HandleStream is called by the underlying service handler method
func (h *CafeService) HandleStream(env *pb.Envelope, pid peer.ID) (chan *pb.Envelope, chan error, chan interface{}) {
	renvCh := make(chan *pb.Envelope)
	errCh := make(chan error)
	cancelCh := make(chan interface{})

	go func() {
		defer close(renvCh)

		var err error
		switch env.Message.Type {
		case pb.Message_CAFE_QUERY:
			err = h.handleQuery(env, pid, renvCh, cancelCh)
		}
		if err != nil {
			errCh <- err
		}
	}()

	return renvCh, errCh, cancelCh
}

// Register creates a session with a cafe
func (h *CafeService) Register(cafeId string, token string) (*pb.CafeSession, error) {
	accnt, err := h.datastore.Config().GetAccount()
	if err != nil {
		return nil, err
	}
	challenge, err := h.challenge(cafeId, accnt)
	if err != nil {
		return nil, err
	}

	// complete the challenge
	cnonce := ksuid.New().String()
	sig, err := accnt.Sign([]byte(challenge.Value + cnonce))
	if err != nil {
		return nil, err
	}
	reg := &pb.CafeRegistration{
		Address: accnt.Address(),
		Value:   challenge.Value,
		Nonce:   cnonce,
		Sig:     sig,
		Token:   token,
	}

	env, err := h.service.NewEnvelope(pb.Message_CAFE_REGISTRATION, reg, nil, false)
	if err != nil {
		return nil, err
	}
	renv, err := h.service.SendRequest(cafeId, env)
	if err != nil {
		return nil, err
	}

	session := new(pb.CafeSession)
	err = ptypes.UnmarshalAny(renv.Message.Payload, session)
	if err != nil {
		return nil, err
	}

	err = h.datastore.CafeSessions().AddOrUpdate(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Deregister removes this peer from a cafe
func (h *CafeService) Deregister(cafeId string) error {
	// @todo: We need to send a retryable de-register request that can be queued,
	// @todo: which will unblock callers. This will require a new request type, target can be the token.
	// @todo: Perhaps registration should also move this way.

	// cleanup
	err := h.datastore.CafeRequests().DeleteByCafe(cafeId)
	if err != nil {
		return err
	}
	err = h.datastore.CafeSessions().Delete(cafeId)
	if err != nil {
		return err
	}

	return nil
}

// Flush begins handling requests recursively
func (h *CafeService) Flush() {
	h.batchRequests(h.datastore.CafeRequests().List("", cafeOutFlushGroupSize))
}

// CheckMessages asks each session's inbox for new messages
func (h *CafeService) CheckMessages(cafeId string) error {
	renv, err := h.sendCafeRequest(cafeId, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_CHECK_MESSAGES, &pb.CafeCheckMessages{
			Token: session.Access,
		}, nil, false)
	})
	if err != nil {
		return err
	}

	res := new(pb.CafeMessages)
	err = ptypes.UnmarshalAny(renv.Message.Payload, res)
	if err != nil {
		return err
	}

	// save messages to inbox
	for _, msg := range res.Messages {
		err = h.inbox.Add(msg)
		if err != nil {
			if !db.ConflictError(err) {
				return err
			}
		}
	}

	h.inbox.Flush()

	// delete them from the remote so that more can be fetched
	if len(res.Messages) > 0 {
		return h.DeleteMessages(cafeId)
	}
	return nil
}

// DeleteMessages deletes a page of messages from a cafe
func (h *CafeService) DeleteMessages(cafeId string) error {
	renv, err := h.sendCafeRequest(cafeId, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_DELETE_MESSAGES, &pb.CafeDeleteMessages{
			Token: session.Access,
		}, nil, false)
	})
	if err != nil {
		return err
	}

	res := new(pb.CafeDeleteMessagesAck)
	err = ptypes.UnmarshalAny(renv.Message.Payload, res)
	if err != nil {
		return err
	}
	if !res.More {
		return nil
	}

	// apparently there are more new messages waiting...
	return h.CheckMessages(cafeId)
}

// PublishPeer publishes the local peer's info
func (h *CafeService) PublishPeer(peer *pb.Peer, cafeId string) error {
	_, err := h.sendCafeRequest(cafeId, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_PUBLISH_PEER, &pb.CafePublishPeer{
			Token: session.Access,
			Peer:  peer,
		}, nil, false)
	})
	if err != nil {
		return err
	}
	return nil
}

// Search performs a query via a cafe
func (h *CafeService) Search(query *pb.Query, cafeId string, reply func(*pb.QueryResult), cancelCh <-chan interface{}) error {
	session := h.datastore.CafeSessions().Get(cafeId)
	if session == nil {
		return fmt.Errorf("could not find session for cafe %s", cafeId)
	}

	env, err := h.service.NewEnvelope(pb.Message_CAFE_QUERY, query, nil, false)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s/api/%s/search", session.Cafe.Url, session.Cafe.Api)
	renvCh, errCh, cancel := h.service.SendHTTPStreamRequest(addr, env, session.Access)
	cancelFn := func() {
		if cancel != nil {
			fn := *cancel
			if fn != nil {
				fn()
			}
		}
	}

	for {
		select {
		case <-cancelCh:
			cancelFn()
			return nil
		case err := <-errCh:
			return err
		case renv, ok := <-renvCh:
			if !ok {
				return nil
			}

			res := new(pb.QueryResults)
			err := ptypes.UnmarshalAny(renv.Message.Payload, res)
			if err != nil {
				return err
			}
			for _, item := range res.Items {
				reply(item)
			}
		}
	}
}

// notifyClient attempts to ping a client that has messages waiting to download
func (h *CafeService) notifyClient(peerId string) error {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_YOU_HAVE_MAIL, nil, nil, false)
	if err != nil {
		return err
	}
	client := string(cafeServiceProtocol) + "/" + peerId

	log.Debugf("sending pubsub %s to %s", env.Message.Type.String(), client)

	payload, err := proto.Marshal(env)
	if err != nil {
		return err
	}
	return ipfs.Publish(h.service.Node(), client, payload)
}

// sendCafeRequest sends an authenticated request, retrying once after a session refresh
func (h *CafeService) sendCafeRequest(cafeId string, envFactory func(*pb.CafeSession) (*pb.Envelope, error)) (*pb.Envelope, error) {
	session := h.datastore.CafeSessions().Get(cafeId)
	if session == nil {
		return nil, fmt.Errorf("could not find session for cafe %s", cafeId)
	}

	env, err := envFactory(session)
	if err != nil {
		return nil, err
	}

	renv, err := h.service.SendRequest(session.Cafe.Peer, env)
	if err != nil {
		if err.Error() == errUnauthorized {
			refreshed, err := h.refresh(session)
			if err != nil {
				return nil, err
			}

			env, err := envFactory(refreshed)
			if err != nil {
				return nil, err
			}

			renv, err = h.service.SendRequest(refreshed.Cafe.Peer, env)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return renv, nil
}

// challenge asks a fellow peer for a cafe challenge
func (h *CafeService) challenge(cafeId string, kp *keypair.Full) (*pb.CafeNonce, error) {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_CHALLENGE, &pb.CafeChallenge{
		Address: kp.Address(),
	}, nil, false)
	if err != nil {
		return nil, err
	}
	renv, err := h.service.SendRequest(cafeId, env)
	if err != nil {
		return nil, err
	}
	res := new(pb.CafeNonce)
	err = ptypes.UnmarshalAny(renv.Message.Payload, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// refresh refreshes a session with a cafe
func (h *CafeService) refresh(session *pb.CafeSession) (*pb.CafeSession, error) {
	refresh := &pb.CafeRefreshSession{
		Access:  session.Access,
		Refresh: session.Refresh,
	}
	env, err := h.service.NewEnvelope(pb.Message_CAFE_REFRESH_SESSION, refresh, nil, false)
	if err != nil {
		return nil, err
	}

	renv, err := h.service.SendRequest(session.Cafe.Peer, env)
	if err != nil {
		return nil, err
	}

	refreshed := new(pb.CafeSession)
	err = ptypes.UnmarshalAny(renv.Message.Payload, refreshed)
	if err != nil {
		return nil, err
	}

	err = h.datastore.CafeSessions().AddOrUpdate(refreshed)
	if err != nil {
		return nil, err
	}
	return refreshed, nil
}

// sendObject sends data or an object by cid to a cafe peer
func (h *CafeService) sendObject(id icid.Cid, cafeId string, token string) error {
	hash := id.Hash().B58String()
	obj := &pb.CafeObject{
		Token: token,
		Cid:   hash,
	}

	data, err := ipfs.DataAtPath(h.service.Node(), hash)
	if err != nil {
		if err == iface.ErrIsDir {
			data, err := ipfs.ObjectAtPath(h.service.Node(), hash)
			if err != nil {
				return err
			}
			obj.Node = data
		} else {
			return err
		}
	} else {
		obj.Data = data
	}

	// send over the raw object data
	env, err := h.service.NewEnvelope(pb.Message_CAFE_OBJECT, obj, nil, false)
	if err != nil {
		return err
	}
	_, err = h.service.SendRequest(cafeId, env)
	if err != nil {
		return err
	}
	return nil
}

// searchLocal searches the local index based on the given query
func (h *CafeService) searchLocal(qtype pb.Query_Type, options *pb.QueryOptions, payload *any.Any, local bool) (*queryResultSet, error) {
	results := newQueryResultSet(options)

	switch qtype {
	case pb.Query_THREAD_SNAPSHOTS:
		q := new(pb.ThreadSnapshotQuery)
		err := ptypes.UnmarshalAny(payload, q)
		if err != nil {
			return nil, err
		}

		clients := h.datastore.CafeClients().ListByAddress(q.Address)
		for _, client := range clients {
			snapshots := h.datastore.CafeClientThreads().ListByClient(client.Id)
			for _, s := range snapshots {
				value, err := proto.Marshal(&pb.CafeClientThread{
					Id:         s.Id,
					Client:     s.Client,
					Ciphertext: s.Ciphertext,
				})
				if err != nil {
					return nil, err
				}
				results.Add(&pb.QueryResult{
					Id:    s.Id,
					Local: local,
					Value: &any.Any{
						TypeUrl: "/CafeClientThread",
						Value:   value,
					},
				})
			}
		}

		// return own threads (encrypted) if query is from an account peer
		if q.Address == h.service.Account.Address() {
			self := h.service.Node().Identity.Pretty()
			for _, t := range h.datastore.Threads().List().Items {
				plaintext, err := proto.Marshal(t)
				if err != nil {
					return nil, err
				}
				ciphertext, err := h.service.Account.Encrypt(plaintext)
				if err != nil {
					return nil, err
				}

				value, err := proto.Marshal(&pb.CafeClientThread{
					Id:         t.Id,
					Client:     self,
					Ciphertext: ciphertext,
				})
				if err != nil {
					return nil, err
				}
				results.Add(&pb.QueryResult{
					Id:    t.Id,
					Local: local,
					Value: &any.Any{
						TypeUrl: "/CafeClientThread",
						Value:   value,
					},
				})
			}

		}

	case pb.Query_CONTACTS:
		q := new(pb.ContactQuery)
		err := ptypes.UnmarshalAny(payload, q)
		if err != nil {
			return nil, err
		}

		peers := h.datastore.Peers().Find(q.Address, q.Name, options.Exclude)
		for _, p := range peers {
			value, err := proto.Marshal(p)
			if err != nil {
				return nil, err
			}
			results.Add(&pb.QueryResult{
				Id:    p.Id,
				Date:  p.Updated,
				Local: local,
				Value: &any.Any{
					TypeUrl: "/Peer",
					Value:   value,
				},
			})
		}
	}

	return results, nil
}

// searchPubSub performs a network-wide search for the given query
func (h *CafeService) searchPubSub(query *pb.Query, reply func(*pb.QueryResults) bool, cancelCh <-chan interface{}, fromCafe bool) error {
	h.inFlightQueries[query.Id] = struct{}{}
	defer func() {
		delete(h.inFlightQueries, query.Id)
	}()

	// respond pubsub if this is a cafe and the request is not from a cafe
	var rtype pb.PubSubQuery_ResponseType
	if h.open && !fromCafe {
		rtype = pb.PubSubQuery_PUBSUB
	} else {
		rtype = pb.PubSubQuery_P2P
	}

	err := h.publishQuery(&pb.PubSubQuery{
		Id:           query.Id,
		Type:         query.Type,
		Payload:      query.Payload,
		ResponseType: rtype,
		Exclude:      query.Options.Exclude,
		Topic:        string(cafeServiceProtocol) + "/" + h.service.Node().Identity.Pretty(),
		Timeout:      query.Options.Wait,
	})
	if err != nil {
		return err
	}

	timer := time.NewTimer(time.Second * time.Duration(query.Options.Wait))
	listener := h.queryResults.Listen()
	doneCh := make(chan struct{})

	done := func() {
		listener.Close()
		close(doneCh)
	}

	go func() {
		<-timer.C
		done()
	}()

	for {
		select {
		case <-cancelCh:
			if timer.Stop() {
				done()
			}
		case <-doneCh:
			return nil
		case value, ok := <-listener.Ch:
			if !ok {
				return nil
			}
			if r, ok := value.(*pb.PubSubQueryResults); ok && r.Id == query.Id && r.Results.Type == query.Type {
				if reply(r.Results) {
					if timer.Stop() {
						done()
					}
				}
			}
		}
	}
}

// publishQuery publishes a search request to the network
func (h *CafeService) publishQuery(req *pb.PubSubQuery) error {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_PUBSUB_QUERY, req, nil, false)
	if err != nil {
		return err
	}
	topic := string(cafeServiceProtocol)

	log.Debugf("sending pubsub %s to %s", env.Message.Type.String(), topic)

	payload, err := proto.Marshal(env)
	if err != nil {
		return err
	}
	return ipfs.Publish(h.service.Node(), topic, payload)
}

// handleChallenge receives a challenge request
func (h *CafeService) handleChallenge(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	req := new(pb.CafeChallenge)
	err := ptypes.UnmarshalAny(env.Message.Payload, req)
	if err != nil {
		return nil, err
	}

	accnt, err := keypair.Parse(req.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		return h.service.NewError(400, errInvalidAddress, env.Message.Request)
	}

	// generate a new random nonce
	nonce := &pb.CafeClientNonce{
		Value:   ksuid.New().String(),
		Address: req.Address,
		Date:    ptypes.TimestampNow(),
	}
	err = h.datastore.CafeClientNonces().Add(nonce)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	return h.service.NewEnvelope(pb.Message_CAFE_NONCE, &pb.CafeNonce{
		Value: nonce.Value,
	}, &env.Message.Request, true)
}

// handleRegistration receives a registration request
func (h *CafeService) handleRegistration(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	reg := new(pb.CafeRegistration)
	err := ptypes.UnmarshalAny(env.Message.Payload, reg)
	if err != nil {
		return nil, err
	}

	// are we open?
	if !h.open {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	// does the provided token match?
	// dev tokens are actually base58(id+token)
	plainBytes, err := base58.FastBase58Decoding(reg.Token)
	if err != nil || len(plainBytes) < 44 {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	encodedToken := h.datastore.CafeTokens().Get(hex.EncodeToString(plainBytes[:12]))
	if encodedToken == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	err = bcrypt.CompareHashAndPassword(encodedToken.Value, plainBytes[12:])
	if err != nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	// check nonce
	snonce := h.datastore.CafeClientNonces().Get(reg.Value)
	if snonce == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}
	if snonce.Address != reg.Address {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	accnt, err := keypair.Parse(reg.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		return h.service.NewError(400, errInvalidAddress, env.Message.Request)
	}

	payload := []byte(reg.Value + reg.Nonce)
	err = accnt.Verify(payload, reg.Sig)
	if err != nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	now := ptypes.TimestampNow()
	client := &pb.CafeClient{
		Id:      pid.Pretty(),
		Address: reg.Address,
		Created: now,
		Seen:    now,
		Token:   encodedToken.Id,
	}
	err = h.datastore.CafeClients().Add(client)
	if err != nil {
		// check if already exists
		client = h.datastore.CafeClients().Get(pid.Pretty())
		if client == nil {
			return h.service.NewError(500, "get or create client failed", env.Message.Request)
		}
	}

	session, err := jwt.NewSession(
		h.service.Node().PrivateKey,
		pid,
		h.Protocol(),
		defaultSessionDuration,
		h.info,
	)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	err = h.datastore.CafeClientNonces().Delete(snonce.Value)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.Request, true)
}

// handleDeregistration receives a deregistration request
func (h *CafeService) handleDeregistration(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	dreg := new(pb.CafeDeregistration)
	err := ptypes.UnmarshalAny(env.Message.Payload, dreg)
	if err != nil {
		return nil, err
	}

	// cleanup
	peerId := pid.Pretty()
	err = h.datastore.CafeClientThreads().DeleteByClient(peerId)
	if err != nil {
		return h.service.NewError(500, "delete client threads failed", env.Message.Request)
	}
	err = h.datastore.CafeClientMessages().DeleteByClient(peerId, -1)
	if err != nil {
		return h.service.NewError(500, "delete client messages failed", env.Message.Request)
	}
	err = h.datastore.CafeClients().Delete(peerId)
	if err != nil {
		return h.service.NewError(500, "delete client failed", env.Message.Request)
	}

	res := &pb.CafeDeregistrationAck{
		Id: peerId,
	}
	return h.service.NewEnvelope(pb.Message_CAFE_DEREGISTRATION_ACK, res, &env.Message.Request, true)
}

// handleRefreshSession receives a refresh session request
func (h *CafeService) handleRefreshSession(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	ref := new(pb.CafeRefreshSession)
	err := ptypes.UnmarshalAny(env.Message.Payload, ref)
	if err != nil {
		return nil, err
	}

	// are we _still_ open?
	if !h.open {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	rerr, err := h.authToken(pid, ref.Refresh, true, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// ensure access and refresh are a valid pair
	access, _ := njwt.Parse(ref.Access, h.verifyKeyFunc)
	if access == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}
	refresh, _ := njwt.Parse(ref.Refresh, h.verifyKeyFunc)
	if refresh == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}
	accessClaims, err := jwt.ParseClaims(access.Claims)
	if err != nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}
	refreshClaims, err := jwt.ParseClaims(refresh.Claims)
	if err != nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}
	if refreshClaims.Id[1:] != accessClaims.Id {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}
	if refreshClaims.Subject != accessClaims.Subject {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	// get a new session
	spid, err := peer.IDB58Decode(accessClaims.Subject)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}
	session, err := jwt.NewSession(
		h.service.Node().PrivateKey,
		spid,
		h.Protocol(),
		defaultSessionDuration,
		h.info,
	)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.Request, true)
}

// handleStore receives a store request
func (h *CafeService) handleStore(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	store := new(pb.CafeStore)
	err := ptypes.UnmarshalAny(env.Message.Payload, store)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, store.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// ignore cids for data already pinned
	list, err := ipfs.NotPinned(h.service.Node(), store.Cids)
	if err != nil {
		return nil, err
	}
	var need []string
	for _, p := range list {
		need = append(need, p.Hash().B58String())
	}

	res := &pb.CafeObjectList{Cids: need}
	return h.service.NewEnvelope(pb.Message_CAFE_OBJECT_LIST, res, &env.Message.Request, true)
}

// handleUnstore receives an unstore request
func (h *CafeService) handleUnstore(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	unstore := new(pb.CafeUnstore)
	err := ptypes.UnmarshalAny(env.Message.Payload, unstore)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, unstore.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// ignore cids for data not pinned
	list, err := ipfs.Pinned(h.service.Node(), unstore.Cids)
	if err != nil {
		return nil, err
	}
	var unstored []string
	for _, p := range list {
		err := ipfs.UnpinCid(h.service.Node(), p, true)
		if err != nil {
			return nil, err
		}
		unstored = append(unstored, p.Hash().B58String())
	}

	res := &pb.CafeUnstoreAck{Cids: unstored}
	return h.service.NewEnvelope(pb.Message_CAFE_UNSTORE_ACK, res, &env.Message.Request, true)
}

// handleObject receives an object request
func (h *CafeService) handleObject(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	obj := new(pb.CafeObject)
	err := ptypes.UnmarshalAny(env.Message.Payload, obj)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, obj.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	var aid *icid.Cid
	if obj.Data != nil {
		aid, err = ipfs.AddData(h.service.Node(), bytes.NewReader(obj.Data), true, false)
	} else if obj.Node != nil {
		aid, err = ipfs.AddObject(h.service.Node(), bytes.NewReader(obj.Node), true)
	} else {
		return h.service.NewError(400, errBadRequest, env.Message.Request)
	}
	if err != nil {
		return nil, err
	}
	rhash := aid.Hash().B58String()

	log.Debugf("stored %s", rhash)

	if rhash != obj.Cid {
		log.Warningf("cids do not match (received %s, resolved %s)", obj.Cid, rhash)
	}

	res := &pb.CafeStoreAck{Id: obj.Cid}
	return h.service.NewEnvelope(pb.Message_CAFE_STORE_ACK, res, &env.Message.Request, true)
}

// handleStoreThread receives a thread store request
func (h *CafeService) handleStoreThread(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	store := new(pb.CafeStoreThread)
	err := ptypes.UnmarshalAny(env.Message.Payload, store)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, store.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	thrd := &pb.CafeClientThread{
		Id:         store.Id,
		Client:     client.Id,
		Ciphertext: store.Ciphertext,
	}
	err = h.datastore.CafeClientThreads().AddOrUpdate(thrd)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	res := &pb.CafeStoreThreadAck{Id: store.Id}
	return h.service.NewEnvelope(pb.Message_CAFE_STORE_THREAD_ACK, res, &env.Message.Request, true)
}

// handleUnstoreThread receives a thread unstore request
func (h *CafeService) handleUnstoreThread(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	unstore := new(pb.CafeUnstoreThread)
	err := ptypes.UnmarshalAny(env.Message.Payload, unstore)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, unstore.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	err = h.datastore.CafeClientThreads().Delete(unstore.Id, client.Id)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	res := &pb.CafeUnstoreThreadAck{Id: unstore.Id}
	return h.service.NewEnvelope(pb.Message_CAFE_UNSTORE_THREAD_ACK, res, &env.Message.Request, true)
}

// handleDeliverMessage receives an inbox message for a client
func (h *CafeService) handleDeliverMessage(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	msg := new(pb.CafeDeliverMessage)
	err := ptypes.UnmarshalAny(env.Message.Payload, msg)
	if err != nil {
		return nil, err
	}

	client := h.datastore.CafeClients().Get(msg.Client)
	if client == nil {
		log.Warningf("received message from %s for unknown client %s", pid.Pretty(), msg.Client)
		return nil, nil
	}

	if msg.Env != nil {
		// pin inner node
		nenv := new(pb.Envelope)
		err = proto.Unmarshal(msg.Env, nenv)
		if err != nil {
			log.Warningf("error unmarshaling envelope: %s", err)
			return nil, err
		}
		tenv := new(pb.ThreadEnvelope)
		err = ptypes.UnmarshalAny(nenv.Message.Payload, tenv)
		if err != nil {
			log.Warningf("error unmarshaling payload: %s", err)
			return nil, err
		}
		oid, err := ipfs.AddObject(h.service.Node(), bytes.NewReader(tenv.Node), true)
		if err != nil {
			log.Warningf("error adding object: %s", err)
			return nil, err
		}
		node, err := ipfs.NodeAtCid(h.service.Node(), *oid)
		if err != nil {
			log.Warningf("error getting node: %s", err)
			return nil, err
		}
		if tenv.Block != nil {
			_, err = ipfs.AddData(h.service.Node(), bytes.NewReader(tenv.Block), true, false)
			if err != nil {
				log.Warningf("error adding block: %s", err)
				return nil, err
			}
		}
		_, err = extractNode(h.service.Node(), node, tenv.Block == nil)
		if err != nil {
			log.Warningf("error extracting node: %s", err)
			return nil, err
		}

		// pin envelope
		id, err := ipfs.AddData(h.service.Node(), bytes.NewReader(msg.Env), true, false)
		if err != nil {
			log.Warningf("error pinning envelope: %s", err)
			return nil, err
		}
		msg.Id = id.Hash().B58String()
	}

	err = h.datastore.CafeClientMessages().AddOrUpdate(&pb.CafeClientMessage{
		Id:     msg.Id,
		Peer:   pid.Pretty(),
		Client: client.Id,
		Date:   ptypes.TimestampNow(),
	})
	if err != nil {
		log.Errorf("error adding message: %s", err)
		return nil, nil
	}
	log.Debugf("added message for %s: %s", client.Id, msg.Id)

	go func() {
		err = h.notifyClient(client.Id)
		if err != nil {
			log.Debugf("unable to notify offline client: %s", client.Id)
		}
	}()
	return nil, nil
}

// handleCheckMessages receives a check inbox messages request
func (h *CafeService) handleCheckMessages(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	check := new(pb.CafeCheckMessages)
	err := ptypes.UnmarshalAny(env.Message.Payload, check)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, check.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	err = h.datastore.CafeClients().UpdateLastSeen(client.Id, time.Now())
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	res := &pb.CafeMessages{
		Messages: make([]*pb.CafeMessage, 0),
	}
	msgs := h.datastore.CafeClientMessages().ListByClient(client.Id, inboxMessagePageSize)
	for _, msg := range msgs {
		res.Messages = append(res.Messages, &pb.CafeMessage{
			Id:   msg.Id,
			Peer: msg.Peer,
			Date: msg.Date,
		})
	}

	return h.service.NewEnvelope(pb.Message_CAFE_MESSAGES, res, &env.Message.Request, true)
}

// handleDeleteMessages receives a message delete request
func (h *CafeService) handleDeleteMessages(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	del := new(pb.CafeDeleteMessages)
	err := ptypes.UnmarshalAny(env.Message.Payload, del)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, del.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	// delete the most recent page
	err = h.datastore.CafeClientMessages().DeleteByClient(client.Id, inboxMessagePageSize)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.Request)
	}

	// check for more
	remaining := h.datastore.CafeClientMessages().CountByClient(client.Id)

	res := &pb.CafeDeleteMessagesAck{More: remaining > 0}
	return h.service.NewEnvelope(pb.Message_CAFE_DELETE_MESSAGES_ACK, res, &env.Message.Request, true)
}

// handleNotifyClient receives a message informing this peer that it has new messages waiting
func (h *CafeService) handleNotifyClient(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	session := h.datastore.CafeSessions().Get(pid.Pretty())
	if session == nil {
		log.Warningf("received message from unknown cafe %s", pid.Pretty())
		return nil, nil
	}

	err := h.CheckMessages(pid.Pretty())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handlePublishPeer indexes a client's peer info for others to search
func (h *CafeService) handlePublishPeer(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	pub := new(pb.CafePublishPeer)
	err := ptypes.UnmarshalAny(env.Message.Payload, pub)
	if err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, pub.Token, false, env.Message.Request)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.Request)
	}

	err = h.datastore.Peers().AddOrUpdate(pub.Peer)
	if err != nil {
		return nil, err
	}

	res := &pb.CafePublishPeerAck{
		Id: pub.Peer.Id,
	}
	return h.service.NewEnvelope(pb.Message_CAFE_PUBLISH_PEER_ACK, res, &env.Message.Request, true)
}

// handleQuery receives a query request
func (h *CafeService) handleQuery(env *pb.Envelope, pid peer.ID, renvs chan *pb.Envelope, cancelCh <-chan interface{}) error {
	query := new(pb.Query)
	err := ptypes.UnmarshalAny(env.Message.Payload, query)
	if err != nil {
		return err
	}
	query = queryDefaults(query)

	results := newQueryResultSet(query.Options)
	reply := func(res *pb.QueryResults) bool {
		added := results.Add(res.Items...)
		if len(added) == 0 {
			return false
		}

		renv, err := h.service.NewEnvelope(pb.Message_CAFE_QUERY_RES, res, &env.Message.Request, true)
		if err != nil {
			log.Errorf("error replying with query results: %s", err)
			return false
		}
		renvs <- renv

		return results.Full()
	}

	// search local
	localResults, err := h.searchLocal(query.Type, query.Options, query.Payload, false)
	if err != nil {
		return err
	}
	if reply(&pb.QueryResults{
		Type:  query.Type,
		Items: localResults.List(),
	}) {
		return nil
	}

	// search network
	return h.searchPubSub(query, reply, cancelCh, true)
}

// handlePubSubQuery receives a query request over pubsub and responds with a direct message
func (h *CafeService) handlePubSubQuery(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	query := new(pb.PubSubQuery)
	err := ptypes.UnmarshalAny(env.Message.Payload, query)
	if err != nil {
		return nil, err
	}

	if _, ok := h.inFlightQueries[query.Id]; ok {
		return nil, nil
	}

	// return results, if any
	options := &pb.QueryOptions{
		Filter:  pb.QueryOptions_NO_FILTER,
		Exclude: query.Exclude,
	}
	results, err := h.searchLocal(query.Type, options, query.Payload, false)
	if err != nil {
		return nil, err
	}
	if len(results.items) == 0 {
		return nil, nil
	}

	res := &pb.PubSubQueryResults{
		Id: query.Id,
		Results: &pb.QueryResults{
			Type:  query.Type,
			Items: results.List(),
		},
	}
	renv, err := h.service.NewEnvelope(pb.Message_CAFE_PUBSUB_QUERY_RES, res, nil, false)
	if err != nil {
		return nil, err
	}

	switch query.ResponseType {
	case pb.PubSubQuery_P2P:
		log.Debugf("responding with %s to %s", renv.Message.Type.String(), pid.Pretty())
		ctx, cancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
		defer cancel()

		err = h.service.SendMessage(ctx, pid.Pretty(), renv)
		if err != nil {
			log.Warningf("error sending message response to %s: %s", pid, err)
		}
	case pb.PubSubQuery_PUBSUB:
		log.Debugf("responding with %s to %s", renv.Message.Type.String(), pid.Pretty())
		if query.Topic == "" {
			query.Topic = query.Id
		}

		payload, err := proto.Marshal(renv)
		if err != nil {
			return nil, err
		}
		err = ipfs.Publish(h.service.Node(), query.Topic, payload)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// handlePubSubQueryResults handles search results received from a pubsub query
func (h *CafeService) handlePubSubQueryResults(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	res := new(pb.PubSubQueryResults)
	err := ptypes.UnmarshalAny(env.Message.Payload, res)
	if err != nil {
		return nil, err
	}

	h.queryResults.Send(res)
	return nil, nil
}

// authToken verifies a request token from a peer
func (h *CafeService) authToken(pid peer.ID, token string, refreshing bool, requestId int32) (*pb.Envelope, error) {
	subject := pid.Pretty()
	_, err := jwt.Validate(token, h.verifyKeyFunc, refreshing, string(h.Protocol()), &subject)
	if err != nil {
		switch err {
		case jwt.ErrNoToken, jwt.ErrExpired:
			return h.service.NewError(401, errUnauthorized, requestId)
		case jwt.ErrInvalid:
			return h.service.NewError(403, errForbidden, requestId)
		}
	}
	return nil, nil
}

// verifyKeyFunc returns the correct key for token verification
func (h *CafeService) verifyKeyFunc(token *njwt.Token) (interface{}, error) {
	return h.service.Node().PrivateKey.GetPublic(), nil
}

// setAddrs sets addresses used in sessions generated by this host
func (h *CafeService) setAddrs(conf *config.Config) {
	url := strings.TrimRight(conf.Cafe.Host.URL, "/")
	if url == "" {
		ip4, err := ipfs.GetLANIPv4Addr(h.service.Node())
		if err != nil {
			ip4, err = h.getPublicIPv4Addr(time.Now().Add(10 * time.Second))
			if err != nil {
				ip4 = "127.0.0.1"
			}
		}
		url = "http://" + ip4
		parts := strings.Split(conf.Addresses.CafeAPI, ":")
		if len(parts) == 2 {
			url += ":" + parts[1]
		}
	}
	log.Infof("cafe url: %s", url)

	h.info = &pb.Cafe{
		Peer:     h.service.Node().Identity.Pretty(),
		Address:  conf.Account.Address,
		Api:      CafeApiVersion,
		Protocol: string(cafeServiceProtocol),
		Node:     common.Version,
		Url:      url,
	}
}

// getPublicIPv4Addr polls for a public ipv4 address until deadline is reached
func (h *CafeService) getPublicIPv4Addr(deadline time.Time) (string, error) {
	ip, err := ipfs.GetPublicIPv4Addr(h.service.Node())
	if err != nil {
		if time.Now().After(deadline) {
			return "", fmt.Errorf("no public ipv4 address was found")
		}
		timer := time.NewTimer(time.Second)
		<-timer.C
		return h.getPublicIPv4Addr(deadline)
	}
	return ip, nil
}

// batchRequests flushes a batch of requests
func (h *CafeService) batchRequests(reqs *pb.CafeRequestList) {
	log.Debugf("handling %d cafe requests", len(reqs.Items))
	if len(reqs.Items) == 0 {
		return
	}

	// group reqs by cafe
	groups := make(map[string][]*pb.CafeRequest)
	for _, req := range reqs.Items {
		groups[req.Cafe.Peer] = append(groups[req.Cafe.Peer], req)
	}

	// process each cafe group concurrently
	var toComplete, toFail, toUnpin []string
	wg := sync.WaitGroup{}
	for cafeId, group := range groups {
		wg.Add(1)
		go func(cafeId string, reqs []*pb.CafeRequest) {
			// group by type
			types := make(map[pb.CafeRequest_Type][]*pb.CafeRequest)
			for _, req := range reqs {
				types[req.Type] = append(types[req.Type], req)
			}
			for t, group := range types {
				handled, failed, err := h.handleRequests(group, t, cafeId)
				if err != nil {
					log.Warningf("error handling requests of type %s: %s", t.String(), err)
				}
				for _, id := range handled {
					toComplete = append(toComplete, id)
					if t == pb.CafeRequest_INBOX {
						toUnpin = append(toUnpin, id)
					}
				}
				for _, id := range failed {
					toFail = append(toFail, id)
				}
			}
			wg.Done()
		}(cafeId, group)
	}
	wg.Wait()

	// next batch
	offset := reqs.Items[len(reqs.Items)-1].Id
	next := h.datastore.CafeRequests().List(offset, cafeOutFlushGroupSize)

	for _, id := range toUnpin {
		req := h.datastore.CafeRequests().Get(id)
		if req == nil {
			continue
		}
		cid, err := icid.Decode(req.Target)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = ipfs.UnpinCid(h.service.Node(), cid, false)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	var err error
	for _, id := range toFail {
		req := h.datastore.CafeRequests().Get(id)
		if req == nil {
			continue
		}
		if req.Attempts+1 >= maxRequestAttempts {
			err = h.datastore.CafeRequests().Delete(id)
			if err != nil {
				log.Error(err.Error())
				return
			}

			// delete queued block
			// @todo: Uncomment this when sync can only be handled by a single cafe session
			//err = h.datastore.Blocks().Delete(req.SyncGroup)
		} else {
			err = h.datastore.CafeRequests().AddAttempt(id)
		}
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	var completed []string
	for _, id := range toComplete {
		err = h.datastore.CafeRequests().UpdateStatus(id, pb.CafeRequest_COMPLETE)
		if err != nil {
			log.Error(err.Error())
			return
		}
		completed = append(completed, id)
	}
	log.Debugf("handled %d cafe requests, %d next", len(completed), len(next.Items))

	h.batchRequests(next)
}

// handleRequest handles a group of requests for a single cafe
func (h *CafeService) handleRequests(reqs []*pb.CafeRequest, rtype pb.CafeRequest_Type, cafeId string) ([]string, []string, error) {
	var handled, failed []string
	var herr error
	switch rtype {

	// store requests are handled in bulk
	case pb.CafeRequest_STORE:
		var cids []string
		for _, req := range reqs {
			cids = append(cids, req.Target)
		}

		stored, err := h.store(cids, cafeId)
		for _, s := range stored {
			for _, r := range reqs {
				if r.Target == s {
					handled = append(handled, r.Id)
				}
			}
		}
		if err != nil {
			log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafeId, err)
			herr = err
			for _, r := range reqs {
				failed = append(failed, r.Id)
			}
		}

	case pb.CafeRequest_UNSTORE:
		var cids []string
		for _, req := range reqs {
			cids = append(cids, req.Target)
		}

		unstored, err := h.unstore(cids, cafeId)
		for _, u := range unstored {
			for _, r := range reqs {
				if r.Target == u {
					handled = append(handled, r.Id)
				}
			}
		}
		if err != nil {
			log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafeId, err)
			herr = err
			for _, r := range reqs {
				failed = append(failed, r.Id)
			}
		}

	case pb.CafeRequest_STORE_THREAD:
		for _, req := range reqs {
			thrd := h.datastore.Threads().Get(req.Target)
			if thrd == nil {
				log.Warningf("could not find thread: %s", req.Target)
				handled = append(handled, req.Id)
				continue
			}

			err := h.storeThread(thrd, cafeId)
			if err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafeId, err)
				herr = err
				failed = append(failed, req.Id)
				continue
			}
			handled = append(handled, req.Id)
		}

	case pb.CafeRequest_UNSTORE_THREAD:
		var err error
		for _, req := range reqs {
			err = h.unstoreThread(req.Target, cafeId)
			if err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafeId, err)
				herr = err
				failed = append(failed, req.Id)
				continue
			}
			handled = append(handled, req.Id)
		}

	case pb.CafeRequest_INBOX:
		var err error
		for _, req := range reqs {
			err = h.deliverMessage(req.Target, req.Peer, req.Cafe)
			if err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafeId, err)
				herr = err
				failed = append(failed, req.Id)
				continue
			}
			handled = append(handled, req.Id)
		}

	}

	return handled, failed, herr
}

// store stores (pins) content on a cafe and returns a list of successful cids
func (h *CafeService) store(cids []string, cafeId string) ([]string, error) {
	var stored []string

	var accessToken string
	renv, err := h.sendCafeRequest(cafeId, func(session *pb.CafeSession) (*pb.Envelope, error) {
		store := &pb.CafeStore{
			Token: session.Access,
			Cids:  cids,
		}
		accessToken = session.Access
		return h.service.NewEnvelope(pb.Message_CAFE_STORE, store, nil, false)
	})
	if err != nil {
		return stored, err
	}

	// unpack response as a request list of cids the cafe is able/willing to store
	req := new(pb.CafeObjectList)
	err = ptypes.UnmarshalAny(renv.Message.Payload, req)
	if err != nil {
		return stored, err
	}
	if len(req.Cids) == 0 {
		log.Debugf("peer %s requested zero objects", cafeId)
		return cids, nil
	}

	// include not-requested (already stored) cids in result
loop:
	for _, i := range cids {
		for _, j := range req.Cids {
			if j == i {
				continue loop
			}
		}
		stored = append(stored, i)
	}

	log.Debugf("sending %d objects to %s", len(req.Cids), cafeId)

	// send each object
	for _, id := range req.Cids {
		decoded, err := icid.Decode(id)
		if err != nil {
			return stored, err
		}
		err = h.sendObject(decoded, cafeId, accessToken)
		if err != nil {
			return stored, err
		}
		stored = append(stored, id)
	}
	return stored, nil
}

// unstore unstores (unpins) content on a cafe and returns a list of successful cids
func (h *CafeService) unstore(cids []string, cafeId string) ([]string, error) {
	renv, err := h.sendCafeRequest(cafeId, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_UNSTORE, &pb.CafeUnstore{
			Token: session.Access,
			Cids:  cids,
		}, nil, false)
	})
	if err != nil {
		return nil, err
	}

	req := new(pb.CafeUnstoreAck)
	err = ptypes.UnmarshalAny(renv.Message.Payload, req)
	if err != nil {
		return nil, err
	}
	return req.Cids, nil
}

// storeThread pushes a thread to a cafe snapshot
func (h *CafeService) storeThread(thrd *pb.Thread, cafeId string) error {
	plaintext, err := proto.Marshal(thrd)
	if err != nil {
		return err
	}
	ciphertext, err := h.service.Account.Encrypt(plaintext)
	if err != nil {
		return err
	}

	_, err = h.sendCafeRequest(cafeId, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_STORE_THREAD, &pb.CafeStoreThread{
			Token:      session.Access,
			Id:         thrd.Id,
			Ciphertext: ciphertext,
		}, nil, false)
	})
	if err != nil {
		return err
	}
	return nil
}

// unstoreThread removes a cafe's thread snapshot
func (h *CafeService) unstoreThread(id string, cafeId string) error {
	renv, err := h.sendCafeRequest(cafeId, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_UNSTORE_THREAD, &pb.CafeUnstoreThread{
			Token: session.Access,
			Id:    id,
		}, nil, false)
	})
	if err != nil {
		return err
	}

	req := new(pb.CafeUnstoreThreadAck)
	return ptypes.UnmarshalAny(renv.Message.Payload, req)
}

// deliverMessage delivers a message content id to a peer's cafe inbox
func (h *CafeService) deliverMessage(mid string, peerId string, cafe *pb.Cafe) error {
	body, err := ipfs.DataAtPath(h.service.Node(), mid)
	if err != nil {
		return err
	}

	env, err := h.service.NewEnvelope(pb.Message_CAFE_DELIVER_MESSAGE, &pb.CafeDeliverMessage{
		Id:     mid,
		Client: peerId,
		Env:    body,
	}, nil, false)
	if err != nil {
		return err
	}

	return h.service.SendMessage(nil, cafe.Peer, env)
}

// queryDefaults ensures the query is within the expected bounds
func queryDefaults(query *pb.Query) *pb.Query {
	if query.Options == nil {
		query.Options = &pb.QueryOptions{
			LocalOnly:  false,
			RemoteOnly: false,
			Limit:      defaultQueryResultsLimit,
			Wait:       defaultQueryWaitSeconds,
			Filter:     pb.QueryOptions_NO_FILTER,
		}
	}

	if query.Options.Limit <= 0 {
		query.Options.Limit = math.MaxInt32
	}

	if query.Options.Wait <= 0 {
		query.Options.Wait = defaultQueryWaitSeconds
	} else if query.Options.Wait > maxQueryWaitSeconds {
		query.Options.Wait = maxQueryWaitSeconds
	}

	return query
}

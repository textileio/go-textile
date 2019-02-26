package core

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/pin"
	"gx/ipfs/QmZMWMvWMVKCbHetJ4RgndbuEF1io2UpUxwQwtNjtYPzSC/go-ipfs-files"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"

	njwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/mr-tron/base58/base58"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/common"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/jwt"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/service"
	"golang.org/x/crypto/bcrypt"
)

// defaultSessionDuration after which session token expires
const defaultSessionDuration = time.Hour * 24 * 7 * 4

// inboxMessagePageSize is the page size used when checking messages
const inboxMessagePageSize = 10

// maxQueryWaitSeconds is used to limit a query request's max wait time
const maxQueryWaitSeconds = 10

// defaultQueryWaitSeconds is a query request's default wait time
const defaultQueryWaitSeconds = 5

// defaultQueryResultsLimit is a query request's default results limit
const defaultQueryResultsLimit = 5

// validation errors
const (
	errInvalidAddress = "invalid address"
	errUnauthorized   = "unauthorized"
	errForbidden      = "forbidden"
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
	clientTopic := string(cafeServiceProtocol) + "/" + h.service.Node().Identity.Pretty()
	go h.service.ListenFor(clientTopic, true, h.handleNotifyClient)
}

// Ping pings another peer
func (h *CafeService) Ping(pid peer.ID) (service.PeerStatus, error) {
	return h.service.Ping(pid)
}

// Handle is called by the underlying service handler method
func (h *CafeService) Handle(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	switch env.Message.Type {
	case pb.Message_CAFE_CHALLENGE:
		return h.handleChallenge(pid, env)
	case pb.Message_CAFE_REGISTRATION:
		return h.handleRegistration(pid, env)
	case pb.Message_CAFE_REFRESH_SESSION:
		return h.handleRefreshSession(pid, env)
	case pb.Message_CAFE_STORE:
		return h.handleStore(pid, env)
	case pb.Message_CAFE_OBJECT:
		return h.handleObject(pid, env)
	case pb.Message_CAFE_STORE_THREAD:
		return h.handleStoreThread(pid, env)
	case pb.Message_CAFE_DELIVER_MESSAGE:
		return h.handleDeliverMessage(pid, env)
	case pb.Message_CAFE_CHECK_MESSAGES:
		return h.handleCheckMessages(pid, env)
	case pb.Message_CAFE_DELETE_MESSAGES:
		return h.handleDeleteMessages(pid, env)
	case pb.Message_CAFE_YOU_HAVE_MAIL:
		return h.handleNotifyClient(pid, env)
	case pb.Message_CAFE_PUBLISH_CONTACT:
		return h.handlePublishContact(pid, env)
	case pb.Message_CAFE_PUBSUB_QUERY:
		return h.handlePubSubQuery(pid, env)
	case pb.Message_CAFE_PUBSUB_QUERY_RES:
		return h.handlePubSubQueryResults(pid, env)
	default:
		return nil, nil
	}
}

// HandleStream is called by the underlying service handler method
func (h *CafeService) HandleStream(pid peer.ID, env *pb.Envelope) (chan *pb.Envelope, chan error, chan interface{}) {
	renvCh := make(chan *pb.Envelope)
	errCh := make(chan error)
	cancelCh := make(chan interface{})

	go func() {
		defer close(renvCh)

		var err error
		switch env.Message.Type {
		case pb.Message_CAFE_QUERY:
			err = h.handleQuery(pid, env, renvCh, cancelCh)
		}
		if err != nil {
			errCh <- err
		}
	}()

	return renvCh, errCh, cancelCh
}

// Register creates a session with a cafe
func (h *CafeService) Register(host string, token string) (*pb.CafeSession, error) {
	host = strings.TrimRight(host, "/")
	addr := fmt.Sprintf("%s/cafe/%s/service", host, cafeApiVersion)

	accnt, err := h.datastore.Config().GetAccount()
	if err != nil {
		return nil, err
	}
	challenge, err := h.challenge(addr, accnt)
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
	renv, err := h.service.SendHTTPRequest(addr, env)
	if err != nil {
		return nil, err
	}

	session := new(pb.CafeSession)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, session); err != nil {
		return nil, err
	}

	if err := h.datastore.CafeSessions().AddOrUpdate(session); err != nil {
		return nil, err
	}

	return session, nil
}

// Store stores (pins) content on a cafe and returns a list of successful cids
func (h *CafeService) Store(cids []string, cafe peer.ID) ([]string, error) {
	var stored []string

	var accessToken string
	var addr string
	renv, err := h.sendCafeRequest(cafe, func(session *pb.CafeSession) (*pb.Envelope, error) {
		store := &pb.CafeStore{
			Token: session.Access,
			Cids:  cids,
		}
		accessToken = session.Access
		addr = getCafeHTTPAddr(session)
		return h.service.NewEnvelope(pb.Message_CAFE_STORE, store, nil, false)
	})
	if err != nil {
		return stored, err
	}

	// unpack response as a request list of cids the cafe is able/willing to store
	req := new(pb.CafeObjectList)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, req); err != nil {
		return stored, err
	}
	if len(req.Cids) == 0 {
		log.Debugf("peer %s requested zero objects", cafe.Pretty())
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

	log.Debugf("sending %d objects to %s", len(req.Cids), cafe.Pretty())

	// send each object
	for _, id := range req.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			return stored, err
		}
		if err := h.sendObject(decoded, addr, accessToken); err != nil {
			return stored, err
		}
		stored = append(stored, id)
	}
	return stored, nil
}

// StoreThread pushes a thread to a cafe backup
func (h *CafeService) StoreThread(thrd *pb.Thread, cafe peer.ID) error {
	plaintext, err := proto.Marshal(thrd)
	if err != nil {
		return err
	}
	ciphertext, err := h.service.Account.Encrypt(plaintext)
	if err != nil {
		return err
	}

	if _, err := h.sendCafeRequest(cafe, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_STORE_THREAD, &pb.CafeStoreThread{
			Token:      session.Access,
			Id:         thrd.Id,
			Ciphertext: ciphertext,
		}, nil, false)
	}); err != nil {
		return err
	}
	return nil
}

// DeliverMessage delivers a message content id to a peer's cafe inbox
// TODO: unpin message locally after it's delivered
func (h *CafeService) DeliverMessage(mid string, pid peer.ID, cafe *pb.Cafe) error {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_DELIVER_MESSAGE, &pb.CafeDeliverMessage{
		Id:     mid,
		Client: pid.Pretty(),
	}, nil, false)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s/cafe/%s/service", cafe.Url, cafe.Api)
	return h.service.SendHTTPMessage(addr, env)
}

// CheckMessages asks each session's inbox for new messages
func (h *CafeService) CheckMessages(cafe peer.ID) error {
	renv, err := h.sendCafeRequest(cafe, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_CHECK_MESSAGES, &pb.CafeCheckMessages{
			Token: session.Access,
		}, nil, false)
	})
	if err != nil {
		return err
	}

	res := new(pb.CafeMessages)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, res); err != nil {
		return err
	}

	// save messages to inbox
	for _, msg := range res.Messages {
		if err := h.inbox.Add(msg); err != nil {
			if !repo.ConflictError(err) {
				return err
			}
		}
	}

	go h.inbox.Flush()

	// delete them from the remote so that more can be fetched
	if len(res.Messages) > 0 {
		return h.DeleteMessages(cafe)
	}
	return nil
}

// DeleteMessages deletes a page of messages from a cafe
func (h *CafeService) DeleteMessages(cafe peer.ID) error {
	renv, err := h.sendCafeRequest(cafe, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_DELETE_MESSAGES, &pb.CafeDeleteMessages{
			Token: session.Access,
		}, nil, false)
	})
	if err != nil {
		return err
	}

	res := new(pb.CafeDeleteMessagesAck)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, res); err != nil {
		return err
	}
	if !res.More {
		return nil
	}

	// apparently there's more new messages waiting...
	return h.CheckMessages(cafe)
}

// PublishContact publishes the local peer's contact info
func (h *CafeService) PublishContact(contact *pb.Contact, cafe peer.ID) error {
	if _, err := h.sendCafeRequest(cafe, func(session *pb.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_PUBLISH_CONTACT, &pb.CafePublishContact{
			Token:   session.Access,
			Contact: contact,
		}, nil, false)
	}); err != nil {
		return err
	}
	return nil
}

// Search performs a query via a cafe
func (h *CafeService) Search(query *pb.Query, cafe peer.ID, reply func(*pb.QueryResult), cancelCh <-chan interface{}) error {
	h.inFlightQueries[query.Id] = struct{}{}
	defer func() {
		delete(h.inFlightQueries, query.Id)
	}()

	envFactory := func(session *pb.CafeSession) (*pb.Envelope, error) {
		query.Token = session.Access
		return h.service.NewEnvelope(pb.Message_CAFE_QUERY, query, nil, false)
	}

	session := h.datastore.CafeSessions().Get(cafe.Pretty())
	if session == nil {
		return errors.New(fmt.Sprintf("could not find session for cafe %s", cafe.Pretty()))
	}

	env, err := envFactory(session)
	if err != nil {
		return err
	}

	renvCh, errCh, cancel := h.service.SendHTTPStreamRequest(getCafeHTTPAddr(session), env)
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
			if err := ptypes.UnmarshalAny(renv.Message.Payload, res); err != nil {
				return err
			}
			for _, item := range res.Items {
				reply(item)
			}
		}
	}
}

// notifyClient attempts to ping a client that has messages waiting to download
func (h *CafeService) notifyClient(pid peer.ID) error {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_YOU_HAVE_MAIL, nil, nil, false)
	if err != nil {
		return err
	}
	client := string(cafeServiceProtocol) + "/" + pid.Pretty()

	log.Debugf("sending pubsub %s to %s", env.Message.Type.String(), client)

	payload, err := proto.Marshal(env)
	if err != nil {
		return err
	}

	return ipfs.Publish(h.service.Node(), client, payload)
}

// sendCafeRequest sends an authenticated request, retrying once after a session refresh
func (h *CafeService) sendCafeRequest(
	cafe peer.ID, envFactory func(*pb.CafeSession) (*pb.Envelope, error)) (*pb.Envelope, error) {

	session := h.datastore.CafeSessions().Get(cafe.Pretty())
	if session == nil {
		return nil, errors.New(fmt.Sprintf("could not find session for cafe %s", cafe.Pretty()))
	}

	env, err := envFactory(session)
	if err != nil {
		return nil, err
	}

	renv, err := h.service.SendHTTPRequest(getCafeHTTPAddr(session), env)
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

			renv, err = h.service.SendHTTPRequest(getCafeHTTPAddr(refreshed), env)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return renv, nil
}

// getCafeHTTPAddr returns the http address of a cafe from a session
func getCafeHTTPAddr(session *pb.CafeSession) string {
	return fmt.Sprintf("%s/cafe/%s/service", session.Cafe.Url, session.Cafe.Api)
}

// challenge asks a fellow peer for a cafe challenge
func (h *CafeService) challenge(cafeAddr string, kp *keypair.Full) (*pb.CafeNonce, error) {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_CHALLENGE, &pb.CafeChallenge{
		Address: kp.Address(),
	}, nil, false)
	if err != nil {
		return nil, err
	}
	renv, err := h.service.SendHTTPRequest(cafeAddr, env)
	if err != nil {
		return nil, err
	}
	res := new(pb.CafeNonce)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, res); err != nil {
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

	renv, err := h.service.SendHTTPRequest(getCafeHTTPAddr(session), env)
	if err != nil {
		return nil, err
	}

	refreshed := new(pb.CafeSession)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, refreshed); err != nil {
		return nil, err
	}

	if err := h.datastore.CafeSessions().AddOrUpdate(refreshed); err != nil {
		return nil, err
	}
	return refreshed, nil
}

// sendObject sends data or an object by cid to a peer
func (h *CafeService) sendObject(id cid.Cid, addr string, token string) error {
	obj := &pb.CafeObject{
		Token: token,
		Cid:   id.Hash().B58String(),
	}

	data, err := ipfs.DataAtPath(h.service.Node(), id.Hash().B58String())
	if err != nil {
		if err == files.ErrNotReader {
			data, err := ipfs.GetObjectAtPath(h.service.Node(), id.Hash().B58String())
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
	if _, err := h.service.SendHTTPRequest(addr, env); err != nil {
		return err
	}
	return nil
}

// searchLocal searches the local index based on the given query
func (h *CafeService) searchLocal(qtype pb.Query_Type, options *pb.QueryOptions, payload *any.Any, local bool) (*queryResultSet, error) {
	results := newQueryResultSet(options)

	switch qtype {
	case pb.Query_THREAD_BACKUPS:
		q := new(pb.ThreadBackupQuery)
		if err := ptypes.UnmarshalAny(payload, q); err != nil {
			return nil, err
		}

		clients := h.datastore.CafeClients().ListByAddress(q.Address)
		for _, client := range clients {
			backups := h.datastore.CafeClientThreads().ListByClient(client.Id)
			for _, b := range backups {
				value, err := proto.Marshal(&pb.CafeClientThread{
					Id:         b.Id,
					Client:     b.Client,
					Ciphertext: b.Ciphertext,
				})
				if err != nil {
					return nil, err
				}
				results.Add(&pb.QueryResult{
					Id:    b.Id,
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
		if err := ptypes.UnmarshalAny(payload, q); err != nil {
			return nil, err
		}

		contacts := h.datastore.Contacts().Find(q.Id, q.Address, q.Username, options.Exclude).Items
		for _, c := range contacts {
			c.Username = toName(c)

			value, err := proto.Marshal(c)
			if err != nil {
				return nil, err
			}
			results.Add(&pb.QueryResult{
				Id:    c.Id,
				Date:  c.Updated,
				Local: local,
				Value: &any.Any{
					TypeUrl: "/Contact",
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

	// if caller needs results over pubsub, start a tmp subscription
	var rtype pb.PubSubQuery_ResponseType
	var psCancel context.CancelFunc
	if fromCafe {
		rtype = pb.PubSubQuery_P2P
	} else {
		rtype = pb.PubSubQuery_PUBSUB
		go func() {
			psCancel = h.service.ListenFor(query.Id, false, h.handlePubSubQueryResults)
		}()
	}

	if err := h.publishQuery(&pb.PubSubQuery{
		Id:           query.Id,
		Type:         query.Type,
		Payload:      query.Payload,
		ResponseType: rtype,
	}); err != nil {
		return err
	}

	timer := time.NewTimer(time.Second * time.Duration(query.Options.Wait))
	listener := h.queryResults.Listen()
	doneCh := make(chan struct{})

	done := func() {
		listener.Close()
		if psCancel != nil {
			psCancel()
		}
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
func (h *CafeService) handleChallenge(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	req := new(pb.CafeChallenge)
	if err := ptypes.UnmarshalAny(env.Message.Payload, req); err != nil {
		return nil, err
	}

	accnt, err := keypair.Parse(req.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		return h.service.NewError(400, errInvalidAddress, env.Message.RequestId)
	}

	// generate a new random nonce
	nonce := &pb.CafeClientNonce{
		Value:   ksuid.New().String(),
		Address: req.Address,
		Date:    ptypes.TimestampNow(),
	}
	if err := h.datastore.CafeClientNonces().Add(nonce); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	return h.service.NewEnvelope(pb.Message_CAFE_NONCE, &pb.CafeNonce{
		Value: nonce.Value,
	}, &env.Message.RequestId, true)
}

// handleRegistration receives a registration request
func (h *CafeService) handleRegistration(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	reg := new(pb.CafeRegistration)
	if err := ptypes.UnmarshalAny(env.Message.Payload, reg); err != nil {
		return nil, err
	}

	// are we open?
	if !h.open {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// does the provided token match?
	// dev tokens are actually base58(id+token)
	plainBytes, err := base58.FastBase58Decoding(reg.Token)
	if err != nil || len(plainBytes) < 44 {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	encodedToken := h.datastore.CafeTokens().Get(hex.EncodeToString(plainBytes[:12]))
	if encodedToken == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	err = bcrypt.CompareHashAndPassword(encodedToken.Value, plainBytes[12:])
	if err != nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// check nonce
	snonce := h.datastore.CafeClientNonces().Get(reg.Value)
	if snonce == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}
	if snonce.Address != reg.Address {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	accnt, err := keypair.Parse(reg.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		return h.service.NewError(400, errInvalidAddress, env.Message.RequestId)
	}

	payload := []byte(reg.Value + reg.Nonce)
	if err := accnt.Verify(payload, reg.Sig); err != nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	now := ptypes.TimestampNow()
	client := &pb.CafeClient{
		Id:      pid.Pretty(),
		Address: reg.Address,
		Created: now,
		Seen:    now,
		Token:   encodedToken.Id,
	}
	if err := h.datastore.CafeClients().Add(client); err != nil {
		// check if already exists
		client = h.datastore.CafeClients().Get(pid.Pretty())
		if client == nil {
			return h.service.NewError(500, "get or create client failed", env.Message.RequestId)
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
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	if err := h.datastore.CafeClientNonces().Delete(snonce.Value); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.RequestId, true)
}

// handleRefreshSession receives a refresh session request
func (h *CafeService) handleRefreshSession(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	ref := new(pb.CafeRefreshSession)
	if err := ptypes.UnmarshalAny(env.Message.Payload, ref); err != nil {
		return nil, err
	}

	// are we _still_ open?
	if !h.open {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	rerr, err := h.authToken(pid, ref.Refresh, true, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// ensure access and refresh are a valid pair
	access, _ := njwt.Parse(ref.Access, h.verifyKeyFunc)
	if access == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}
	refresh, _ := njwt.Parse(ref.Refresh, h.verifyKeyFunc)
	if refresh == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}
	accessClaims, err := jwt.ParseClaims(access.Claims)
	if err != nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}
	refreshClaims, err := jwt.ParseClaims(refresh.Claims)
	if err != nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}
	if refreshClaims.Id[1:] != accessClaims.Id {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}
	if refreshClaims.Subject != accessClaims.Subject {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// get a new session
	spid, err := peer.IDB58Decode(accessClaims.Subject)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}
	session, err := jwt.NewSession(
		h.service.Node().PrivateKey,
		spid,
		h.Protocol(),
		defaultSessionDuration,
		h.info,
	)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.RequestId, true)
}

// handleStore receives a store request
func (h *CafeService) handleStore(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	store := new(pb.CafeStore)
	if err := ptypes.UnmarshalAny(env.Message.Payload, store); err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, store.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// ignore cids for data already pinned
	var decoded []cid.Cid
	for _, id := range store.Cids {
		dec, err := cid.Decode(id)
		if err != nil {
			return nil, err
		}
		decoded = append(decoded, dec)
	}

	pinned, err := h.service.Node().Pinning.CheckIfPinned(decoded...)
	if err != nil {
		return nil, err
	}

	var need []string
	for _, p := range pinned {
		if p.Mode == pin.NotPinned {
			need = append(need, p.Key.Hash().B58String())
		}
	}

	res := &pb.CafeObjectList{Cids: need}
	return h.service.NewEnvelope(pb.Message_CAFE_OBJECT_LIST, res, &env.Message.RequestId, true)
}

// handleObject receives an object request
func (h *CafeService) handleObject(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	obj := new(pb.CafeObject)
	if err := ptypes.UnmarshalAny(env.Message.Payload, obj); err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, obj.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	var id string
	if obj.Data != nil {
		aid, err := ipfs.AddData(h.service.Node(), bytes.NewReader(obj.Data), true)
		if err != nil {
			return nil, err
		}
		id = aid.Hash().B58String()

		log.Debugf("pinned object %s", id)

	} else if obj.Node != nil {
		aid, err := ipfs.AddObject(h.service.Node(), bytes.NewReader(obj.Node), true)
		if err != nil {
			return nil, err
		}
		id = aid.Hash().B58String()

		log.Debugf("pinned node %s", id)
	}

	if id != obj.Cid {
		log.Warningf("cids do not match (received %s, resolved %s)", obj.Cid, id)
	}

	res := &pb.CafeStored{Id: obj.Cid}
	return h.service.NewEnvelope(pb.Message_CAFE_STORED, res, &env.Message.RequestId, true)
}

// handleStoreThread receives a thread store request
func (h *CafeService) handleStoreThread(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	store := new(pb.CafeStoreThread)
	if err := ptypes.UnmarshalAny(env.Message.Payload, store); err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, store.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	thrd := &pb.CafeClientThread{
		Id:         store.Id,
		Client:     client.Id,
		Ciphertext: store.Ciphertext,
	}
	if err := h.datastore.CafeClientThreads().AddOrUpdate(thrd); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	res := &pb.CafeStored{Id: store.Id}
	return h.service.NewEnvelope(pb.Message_CAFE_STORED, res, &env.Message.RequestId, true)
}

// handleDeliverMessage receives an inbox message for a client
func (h *CafeService) handleDeliverMessage(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	msg := new(pb.CafeDeliverMessage)
	if err := ptypes.UnmarshalAny(env.Message.Payload, msg); err != nil {
		return nil, err
	}

	client := h.datastore.CafeClients().Get(msg.Client)
	if client == nil {
		log.Warningf("received message from %s for unknown client %s", pid.Pretty(), msg.Client)
		return nil, nil
	}

	message := &pb.CafeClientMessage{
		Id:     msg.Id,
		Peer:   pid.Pretty(),
		Client: client.Id,
		Date:   ptypes.TimestampNow(),
	}
	if err := h.datastore.CafeClientMessages().AddOrUpdate(message); err != nil {
		log.Errorf("error adding message: %s", err)
		return nil, nil
	}

	go func() {
		pid, err := peer.IDB58Decode(client.Id)
		if err != nil {
			log.Errorf("error parsing client id %s: %s", client.Id, err)
			return
		}
		if err := h.notifyClient(pid); err != nil {
			log.Debugf("unable to notify offline client: %s", client.Id)
		}
	}()
	return nil, nil
}

// handleCheckMessages receives a check inbox messages request
func (h *CafeService) handleCheckMessages(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	check := new(pb.CafeCheckMessages)
	if err := ptypes.UnmarshalAny(env.Message.Payload, check); err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, check.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	if err := h.datastore.CafeClients().UpdateLastSeen(client.Id, time.Now()); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
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

	return h.service.NewEnvelope(pb.Message_CAFE_MESSAGES, res, &env.Message.RequestId, true)
}

// handleDeleteMessages receives a message delete request
func (h *CafeService) handleDeleteMessages(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	del := new(pb.CafeDeleteMessages)
	if err := ptypes.UnmarshalAny(env.Message.Payload, del); err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, del.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// delete the most recent page
	if err := h.datastore.CafeClientMessages().DeleteByClient(client.Id, inboxMessagePageSize); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// check for more
	remaining := h.datastore.CafeClientMessages().CountByClient(client.Id)

	res := &pb.CafeDeleteMessagesAck{More: remaining > 0}
	return h.service.NewEnvelope(pb.Message_CAFE_DELETE_MESSAGES_ACK, res, &env.Message.RequestId, true)
}

// handleNotifyClient receives a message informing this peer that it has new messages waiting
func (h *CafeService) handleNotifyClient(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	session := h.datastore.CafeSessions().Get(pid.Pretty())
	if session == nil {
		log.Warningf("received message from unknown cafe %s", pid.Pretty())
		return nil, nil
	}

	if err := h.CheckMessages(pid); err != nil {
		return nil, err
	}

	return nil, nil
}

// handlePublishContact indexes a client's contact info for others to search
func (h *CafeService) handlePublishContact(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	pub := new(pb.CafePublishContact)
	if err := ptypes.UnmarshalAny(env.Message.Payload, pub); err != nil {
		return nil, err
	}

	rerr, err := h.authToken(pid, pub.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	if err := h.datastore.Contacts().AddOrUpdate(pub.Contact); err != nil {
		return nil, err
	}

	res := &pb.CafePublishContactAck{
		Id: pub.Contact.Id,
	}
	return h.service.NewEnvelope(pb.Message_CAFE_PUBLISH_CONTACT_ACK, res, &env.Message.RequestId, true)
}

// handleQuery receives a query request
func (h *CafeService) handleQuery(pid peer.ID, env *pb.Envelope, renvs chan *pb.Envelope, cancelCh <-chan interface{}) error {
	query := new(pb.Query)
	if err := ptypes.UnmarshalAny(env.Message.Payload, query); err != nil {
		return err
	}
	query = queryDefaults(query)

	rerr, err := h.authToken(pid, query.Token, false, env.Message.RequestId)
	if err != nil {
		return err
	}
	if rerr != nil {
		renvs <- rerr
		return nil
	}

	results := newQueryResultSet(query.Options)
	reply := func(res *pb.QueryResults) bool {
		added := results.Add(res.Items...)
		if len(added) == 0 {
			return false
		}

		renv, err := h.service.NewEnvelope(pb.Message_CAFE_QUERY_RES, res, &env.Message.RequestId, true)
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
func (h *CafeService) handlePubSubQuery(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	query := new(pb.PubSubQuery)
	if err := ptypes.UnmarshalAny(env.Message.Payload, query); err != nil {
		return nil, err
	}

	if _, ok := h.inFlightQueries[query.Id]; ok {
		return nil, nil
	}

	// return results, if any
	options := &pb.QueryOptions{Filter: pb.QueryOptions_NO_FILTER}
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
		return renv, nil
	case pb.PubSubQuery_PUBSUB:
		log.Debugf("responding with %s to %s", renv.Message.Type.String(), pid.Pretty())

		payload, err := proto.Marshal(renv)
		if err != nil {
			return nil, err
		}
		if err := ipfs.Publish(h.service.Node(), query.Id, payload); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// handlePubSubQueryResults handles search results received from a pubsub query
func (h *CafeService) handlePubSubQueryResults(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	res := new(pb.PubSubQueryResults)
	if err := ptypes.UnmarshalAny(env.Message.Payload, res); err != nil {
		return nil, err
	}

	h.queryResults.Send(res)
	return nil, nil
}

// authToken verifies a request token from a peer
func (h *CafeService) authToken(pid peer.ID, token string, refreshing bool, requestId int32) (*pb.Envelope, error) {
	subject := pid.Pretty()
	if err := jwt.Validate(token, h.verifyKeyFunc, refreshing, string(h.Protocol()), &subject); err != nil {
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
func (h *CafeService) setAddrs(conf *config.Config, swarmPorts config.SwarmPorts) {
	publicIP := conf.Cafe.Host.PublicIP
	if publicIP == "" {
		publicIP = "127.0.0.1"
	}

	url := strings.TrimRight(conf.Cafe.Host.URL, "/")
	if url == "" {
		url = "http://" + publicIP
		parts := strings.Split(conf.Addresses.CafeAPI, ":")
		if len(parts) == 2 {
			url += ":" + parts[1]
		}
	}
	log.Infof("cafe url: %s", url)

	swarm := []string{
		fmt.Sprintf("/ip4/%s/tcp/%s", publicIP, swarmPorts.TCP),
	}
	if swarmPorts.WS != "" {
		swarm = append(swarm, fmt.Sprintf("/ip4/%s/tcp/%s/ws", publicIP, swarmPorts.WS))
	}
	log.Infof("cafe multiaddresses: %s", swarm)

	h.info = &pb.Cafe{
		Peer:     h.service.Node().Identity.Pretty(),
		Address:  conf.Account.Address,
		Api:      cafeApiVersion,
		Protocol: string(cafeServiceProtocol),
		Node:     common.Version,
		Url:      url,
		Swarm:    swarm,
	}
}

// queryDefaults ensures the query is within the expected bounds
func queryDefaults(query *pb.Query) *pb.Query {
	if query.Options == nil {
		query.Options = &pb.QueryOptions{
			Local:  false,
			Limit:  defaultQueryResultsLimit,
			Wait:   defaultQueryWaitSeconds,
			Filter: pb.QueryOptions_NO_FILTER,
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

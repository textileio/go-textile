package core

import (
	"context"
	"errors"
	"fmt"
	njwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/jwt"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/service"
	"gx/ipfs/QmVzK524a2VWLqyvtBeiHKsUAWYgeAk4DBeZoY7vpNPNRx/go-block-format"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"time"
)

// defaultSessionDuration after which session token expires
const defaultSessionDuration = time.Hour * 24 * 7 * 4

// blockTimeout is the context timeout for getting a local block
const blockTimeout = time.Second * 5

// inboxMessagePageSize is the page size used when checking messages
const inboxMessagePageSize = 10

// validation errors
const (
	errInvalidAddress = "invalid address"
	errUnauthorized   = "unauthorized"
	errForbidden      = "forbidden"
)

// CafeService is a libp2p pinning and offline message service
type CafeService struct {
	service   *service.Service
	datastore repo.Datastore
	inbox     *CafeInbox
}

// NewCafeService returns a new threads service
func NewCafeService(
	account *keypair.Full,
	node *core.IpfsNode,
	datastore repo.Datastore,
	inbox *CafeInbox,
) *CafeService {
	handler := &CafeService{
		datastore: datastore,
		inbox:     inbox,
	}
	handler.service = service.NewService(account, handler, node)
	return handler
}

// Protocol returns the handler protocol
func (h *CafeService) Protocol() protocol.ID {
	return protocol.ID("/textile/cafe/1.0.0")
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
	case pb.Message_CAFE_BLOCK:
		return h.handleBlock(pid, env)
	case pb.Message_CAFE_STORE_THREAD:
		return h.handleStoreThread(pid, env)
	case pb.Message_CAFE_DELIVER_MESSAGE:
		return h.handleDeliverMessage(pid, env)
	case pb.Message_CAFE_CHECK_MESSAGES:
		return h.handleCheckMessages(pid, env)
	case pb.Message_CAFE_DELETE_MESSAGES:
		return h.handleDeleteMessages(pid, env)
	default:
		return nil, nil
	}
}

// Register creates a session with a cafe
func (h *CafeService) Register(cafe peer.ID) error {
	// get a challenge
	accnt, err := h.datastore.Config().GetAccount()
	if err != nil {
		return err
	}
	challenge, err := h.challenge(accnt, cafe)
	if err != nil {
		return err
	}

	// complete the challenge
	cnonce := ksuid.New().String()
	sig, err := accnt.Sign([]byte(challenge.Value + cnonce))
	if err != nil {
		return err
	}
	reg := &pb.CafeRegistration{
		Address: accnt.Address(),
		Value:   challenge.Value,
		Nonce:   cnonce,
		Sig:     sig,
	}

	// register
	env, err := h.service.NewEnvelope(pb.Message_CAFE_REGISTRATION, reg, nil, false)
	if err != nil {
		return err
	}
	renv, err := h.service.SendRequest(cafe, env)
	if err != nil {
		return err
	}

	// unpack session
	res := new(pb.CafeSession)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, res); err != nil {
		return err
	}

	// local login
	exp, err := ptypes.Timestamp(res.Exp)
	if err != nil {
		return err
	}
	session := &repo.CafeSession{
		CafeId:  cafe.Pretty(),
		Access:  res.Access,
		Refresh: res.Refresh,
		Expiry:  exp,
	}
	return h.datastore.CafeSessions().AddOrUpdate(session)
}

// Store stores (pins) content on a cafe and returns a list of successful cids
func (h *CafeService) Store(cids []string, cafe peer.ID) ([]string, error) {
	var stored []string

	// ask cafe if it can store these cids
	var accessToken string
	renv, err := h.sendCafeRequest(cafe, func(session *repo.CafeSession) (*pb.Envelope, error) {
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
	req := new(pb.CafeBlockList)
	err = ptypes.UnmarshalAny(renv.Message.Payload, req)
	if err != nil {
		return stored, err
	}
	if len(req.Cids) == 0 {
		log.Debugf("peer %s requested zero blocks", cafe.Pretty())
		return cids, nil
	}
	log.Debugf("sending %d blocks to %s", len(req.Cids), cafe.Pretty())

	// send each block
	for _, id := range req.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			continue
		}
		if err := h.sendBlock(*decoded, cafe, accessToken); err != nil {
			log.Errorf("error sending block: %s", err)
			continue
		}
		stored = append(stored, id)
	}
	return stored, nil
}

// StoreThread pushes a thread to a cafe backup
func (h *CafeService) StoreThread(thrd *repo.Thread, cafe peer.ID) error {
	// encrypt thread components
	skCipher, err := h.service.Account.Encrypt(thrd.PrivKey)
	if err != nil {
		return err
	}
	var headCipher []byte
	if thrd.Head != "" {
		headCipher, err = h.service.Account.Encrypt([]byte(thrd.Head))
		if err != nil {
			return err
		}
	}
	var nameCipher []byte
	if thrd.Name != "" {
		nameCipher, err = h.service.Account.Encrypt([]byte(thrd.Name))
		if err != nil {
			return err
		}
	}

	// build request
	if _, err := h.sendCafeRequest(cafe, func(session *repo.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_STORE_THREAD, &pb.CafeStoreThread{
			Token:      session.Access,
			Id:         thrd.Id,
			SkCipher:   skCipher,
			HeadCipher: headCipher,
			NameCipher: nameCipher,
		}, nil, false)
	}); err != nil {
		return err
	}
	return nil
}

// DeliverMessage delivers a message content id to a peer's cafe inbox
// TODO: unpin message locally after it's delivered
func (h *CafeService) DeliverMessage(mid string, pid peer.ID, cafe peer.ID) error {
	// build request
	env, err := h.service.NewEnvelope(pb.Message_CAFE_DELIVER_MESSAGE, &pb.CafeDeliverMessage{
		Id:       mid,
		ClientId: pid.Pretty(),
	}, nil, false)
	if err != nil {
		return err
	}
	return h.service.SendMessage(cafe, env)
}

// CheckMessages asks each session's inbox for new messages
func (h *CafeService) CheckMessages(cafe peer.ID) error {
	// build an authenticated request
	renv, err := h.sendCafeRequest(cafe, func(session *repo.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_CHECK_MESSAGES, &pb.CafeCheckMessages{
			Token: session.Access,
		}, nil, false)
	})
	if err != nil {
		return err
	}

	// unpack response as a message list of cids
	res := new(pb.CafeMessages)
	err = ptypes.UnmarshalAny(renv.Message.Payload, res)
	if err != nil {
		return err
	}

	// save messages to inbox
	for _, msg := range res.Messages {
		if err := h.inbox.Add(msg); err != nil {
			return err
		}
	}

	// flush inbox
	go h.inbox.Flush()

	// delete them from the remote so that more can be fetched
	if len(res.Messages) > 0 {
		return h.DeleteMessages(cafe)
	}
	return nil
}

// DeleteMessages deletes a page of messages from a cafe
func (h *CafeService) DeleteMessages(cafe peer.ID) error {
	// build an authenticated request
	renv, err := h.sendCafeRequest(cafe, func(session *repo.CafeSession) (*pb.Envelope, error) {
		return h.service.NewEnvelope(pb.Message_CAFE_DELETE_MESSAGES, &pb.CafeDeleteMessages{
			Token: session.Access,
		}, nil, false)
	})
	if err != nil {
		return err
	}

	// unpack response as an ack
	res := new(pb.CafeDeleteMessagesAck)
	err = ptypes.UnmarshalAny(renv.Message.Payload, res)
	if err != nil {
		return err
	}
	if !res.More {
		return nil
	}

	// apparently there's more new messages waiting...
	return h.CheckMessages(cafe)
}

// sendCafeRequest sends an authenticated request, retrying once after a session refresh
func (h *CafeService) sendCafeRequest(cafe peer.ID, envFactory func(*repo.CafeSession) (*pb.Envelope, error)) (*pb.Envelope, error) {
	// find access token for this cafe
	session := h.datastore.CafeSessions().Get(cafe.Pretty())
	if session == nil {
		return nil, errors.New(fmt.Sprintf("could not find session for cafe %s", cafe.Pretty()))
	}
	env, err := envFactory(session)
	if err != nil {
		return nil, err
	}
	renv, err := h.service.SendRequest(cafe, env)
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
			renv, err = h.service.SendRequest(cafe, env)
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
func (h *CafeService) challenge(kp *keypair.Full, cafe peer.ID) (*pb.CafeNonce, error) {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_CHALLENGE, &pb.CafeChallenge{
		Address: kp.Address(),
	}, nil, false)
	if err != nil {
		return nil, err
	}
	renv, err := h.service.SendRequest(cafe, env)
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
func (h *CafeService) refresh(session *repo.CafeSession) (*repo.CafeSession, error) {
	refresh := &pb.CafeRefreshSession{
		Access:  session.Access,
		Refresh: session.Refresh,
	}
	env, err := h.service.NewEnvelope(pb.Message_CAFE_REFRESH_SESSION, refresh, nil, false)
	if err != nil {
		return nil, err
	}
	pid, err := peer.IDB58Decode(session.CafeId)
	if err != nil {
		return nil, err
	}
	renv, err := h.service.SendRequest(pid, env)
	if err != nil {
		return nil, err
	}

	// unpack session
	res := new(pb.CafeSession)
	if err := ptypes.UnmarshalAny(renv.Message.Payload, res); err != nil {
		return nil, err
	}

	// local login
	exp, err := ptypes.Timestamp(res.Exp)
	if err != nil {
		return nil, err
	}
	refreshed := &repo.CafeSession{
		CafeId:  session.CafeId,
		Access:  res.Access,
		Refresh: res.Refresh,
		Expiry:  exp,
	}
	if err := h.datastore.CafeSessions().AddOrUpdate(refreshed); err != nil {
		return nil, err
	}
	return refreshed, nil
}

// sendBlock sends a block by cid to a peer
func (h *CafeService) sendBlock(id cid.Cid, pid peer.ID, token string) error {
	// get block locally
	ctx, cancel := context.WithTimeout(context.Background(), blockTimeout)
	defer cancel()
	block, err := h.service.Node.Blocks.GetBlock(ctx, &id)
	if err != nil {
		return err
	}

	// send over the raw block data
	pblock := &pb.CafeBlock{
		Token:   token,
		Cid:     block.Cid().String(),
		RawData: block.RawData(),
	}
	env, err := h.service.NewEnvelope(pb.Message_CAFE_BLOCK, pblock, nil, false)
	if err != nil {
		return err
	}
	if _, err := h.service.SendRequest(pid, env); err != nil {
		return err
	}
	return nil
}

// handleChallenge receives a challenge request
func (h *CafeService) handleChallenge(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	req := new(pb.CafeChallenge)
	if err := ptypes.UnmarshalAny(env.Message.Payload, req); err != nil {
		return nil, err
	}

	// validate address
	accnt, err := keypair.Parse(req.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		return h.service.NewError(400, errInvalidAddress, env.Message.RequestId)
	}

	// generate a new random nonce
	nonce := &repo.CafeClientNonce{
		Value:   ksuid.New().String(),
		Address: req.Address,
		Date:    time.Now(),
	}
	if err := h.datastore.CafeClientNonces().Add(nonce); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// return a wrapped response
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

	// lookup the nonce
	snonce := h.datastore.CafeClientNonces().Get(reg.Value)
	if snonce == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}
	if snonce.Address != reg.Address {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// validate address
	accnt, err := keypair.Parse(reg.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		return h.service.NewError(400, errInvalidAddress, env.Message.RequestId)
	}

	// verify
	payload := []byte(reg.Value + reg.Nonce)
	if err := accnt.Verify(payload, reg.Sig); err != nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// create new
	now := time.Now()
	client := &repo.CafeClient{
		Id:       pid.Pretty(),
		Address:  reg.Address,
		Created:  now,
		LastSeen: now,
	}
	if err := h.datastore.CafeClients().Add(client); err != nil {
		// check if already exists
		client = h.datastore.CafeClients().Get(pid.Pretty())
		if client == nil {
			return h.service.NewError(500, "get or create client failed", env.Message.RequestId)
		}
	}

	// get a session
	session, err := jwt.NewSession(h.service.Node.PrivateKey, pid, h.Protocol(), defaultSessionDuration)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// delete the nonce
	if err := h.datastore.CafeClientNonces().Delete(snonce.Value); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// return a wrapped response
	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.RequestId, true)
}

// handleRefreshSession receives a refresh session request
func (h *CafeService) handleRefreshSession(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	ref := new(pb.CafeRefreshSession)
	if err := ptypes.UnmarshalAny(env.Message.Payload, ref); err != nil {
		return nil, err
	}

	// validate refresh token
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
	session, err := jwt.NewSession(h.service.Node.PrivateKey, spid, h.Protocol(), defaultSessionDuration)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// return a wrapped response
	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.RequestId, true)
}

// handleSession receives a store request
func (h *CafeService) handleStore(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	store := new(pb.CafeStore)
	err := ptypes.UnmarshalAny(env.Message.Payload, store)
	if err != nil {
		return nil, err
	}

	// validate access token
	rerr, err := h.authToken(pid, store.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// ignore cids for blocks already present in local datastore
	var need []string
	for _, id := range store.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			continue
		}
		has, err := h.service.Node.Blockstore.Has(decoded)
		if err != nil || !has {
			need = append(need, decoded.String())
		}
	}

	// return a wrapped response
	res := &pb.CafeBlockList{Cids: need}
	return h.service.NewEnvelope(pb.Message_CAFE_BLOCKLIST, res, &env.Message.RequestId, true)
}

// handleBlock receives a block request
func (h *CafeService) handleBlock(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	block := new(pb.CafeBlock)
	err := ptypes.UnmarshalAny(env.Message.Payload, block)
	if err != nil {
		return nil, err
	}

	// validate access token
	rerr, err := h.authToken(pid, block.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// add a new block to the local datastore
	id, err := cid.Decode(block.Cid)
	if err != nil {
		return nil, err
	}
	bblock, err := blocks.NewBlockWithCid(block.RawData, id)
	if err != nil {
		return nil, err
	}
	if err := h.service.Node.Blocks.AddBlock(bblock); err != nil {
		return nil, err
	}

	// return a wrapped response
	res := &pb.CafeStored{Id: block.Cid}
	return h.service.NewEnvelope(pb.Message_CAFE_STORED, res, &env.Message.RequestId, true)
}

// handleStoreThread receives a thread request
func (h *CafeService) handleStoreThread(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	store := new(pb.CafeStoreThread)
	err := ptypes.UnmarshalAny(env.Message.Payload, store)
	if err != nil {
		return nil, err
	}

	// validate access token
	rerr, err := h.authToken(pid, store.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// lookup client
	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// add or update
	thrd := &repo.CafeClientThread{
		Id:         store.Id,
		ClientId:   client.Id,
		SkCipher:   store.SkCipher,
		HeadCipher: store.HeadCipher,
		NameCipher: store.NameCipher,
	}
	if err := h.datastore.CafeClientThreads().AddOrUpdate(thrd); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// return a wrapped response
	res := &pb.CafeStored{Id: store.Id}
	return h.service.NewEnvelope(pb.Message_CAFE_STORED, res, &env.Message.RequestId, true)
}

// handleDeliverMessage receives an inbox message for a client
func (h *CafeService) handleDeliverMessage(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	msg := new(pb.CafeDeliverMessage)
	err := ptypes.UnmarshalAny(env.Message.Payload, msg)
	if err != nil {
		return nil, err
	}

	// lookup client
	client := h.datastore.CafeClients().Get(msg.ClientId)
	if client == nil {
		log.Warningf("received message from %s for unknown client %s", pid.Pretty(), msg.ClientId)
		return nil, nil
	}

	// add or update
	message := &repo.CafeClientMessage{
		Id:       msg.Id,
		PeerId:   pid.Pretty(),
		ClientId: client.Id,
		Date:     time.Now(),
	}
	if err := h.datastore.CafeClientMessages().AddOrUpdate(message); err != nil {
		log.Errorf("error adding message: %s", err)
		return nil, nil
	}
	return nil, nil
}

// handleCheckMessages receives a check inbox messages request
func (h *CafeService) handleCheckMessages(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	check := new(pb.CafeCheckMessages)
	err := ptypes.UnmarshalAny(env.Message.Payload, check)
	if err != nil {
		return nil, err
	}

	// validate access token
	rerr, err := h.authToken(pid, check.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// lookup client
	client := h.datastore.CafeClients().Get(pid.Pretty())
	if client == nil {
		return h.service.NewError(403, errForbidden, env.Message.RequestId)
	}

	// update last seen
	if err := h.datastore.CafeClients().UpdateLastSeen(client.Id, time.Now()); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// list new messages
	res := &pb.CafeMessages{
		Messages: make([]*pb.CafeMessage, 0),
	}
	msgs := h.datastore.CafeClientMessages().ListByClient(client.Id, inboxMessagePageSize)
	for _, msg := range msgs {
		date, err := ptypes.TimestampProto(msg.Date)
		if err != nil {
			return h.service.NewError(500, err.Error(), env.Message.RequestId)
		}
		res.Messages = append(res.Messages, &pb.CafeMessage{
			Id:     msg.Id,
			PeerId: msg.PeerId,
			Date:   date,
		})
	}

	// return a wrapped response
	return h.service.NewEnvelope(pb.Message_CAFE_MESSAGES, res, &env.Message.RequestId, true)
}

// handleDeleteMessages receives a message delete request
func (h *CafeService) handleDeleteMessages(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	del := new(pb.CafeDeleteMessages)
	err := ptypes.UnmarshalAny(env.Message.Payload, del)
	if err != nil {
		return nil, err
	}

	// validate access token
	rerr, err := h.authToken(pid, del.Token, false, env.Message.RequestId)
	if err != nil {
		return nil, err
	}
	if rerr != nil {
		return rerr, nil
	}

	// lookup client
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

	// return a wrapped response
	res := &pb.CafeDeleteMessagesAck{More: remaining > 0}
	return h.service.NewEnvelope(pb.Message_CAFE_DELETE_MESSAGES_ACK, res, &env.Message.RequestId, true)
}

// authToken verifies a request token from a peer
func (h *CafeService) authToken(pid peer.ID, tokenString string, refreshing bool, requestId int32) (*pb.Envelope, error) {
	// parse it
	token, pErr := njwt.Parse(tokenString, h.verifyKeyFunc)
	if token == nil {
		return h.service.NewError(401, errUnauthorized, requestId)
	}

	// pull out claims
	claims, err := jwt.ParseClaims(token.Claims)
	if err != nil {
		return h.service.NewError(403, errForbidden, requestId)
	}

	// check valid
	if pErr != nil {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			// 401 indicates a retry is expected after a token refresh
			return h.service.NewError(401, errUnauthorized, requestId)
		}
		return h.service.NewError(403, errForbidden, requestId)
	}

	// check scope
	switch claims.Scope {
	case jwt.Access:
		if refreshing {
			return h.service.NewError(403, errForbidden, requestId)
		}
	case jwt.Refresh:
		if !refreshing {
			return h.service.NewError(403, errForbidden, requestId)
		}
	default:
		return h.service.NewError(403, errForbidden, requestId)
	}

	// verify owner
	if claims.Subject != pid.Pretty() {
		return h.service.NewError(403, errForbidden, requestId)
	}

	// verify protocol
	if !claims.VerifyAudience(string(h.Protocol()), true) {
		return h.service.NewError(403, errForbidden, requestId)
	}
	return nil, nil
}

// verifyKeyFunc returns the correct key for token verification
func (h *CafeService) verifyKeyFunc(token *njwt.Token) (interface{}, error) {
	return h.service.Node.PrivateKey.GetPublic(), nil
}

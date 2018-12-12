package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	"gx/ipfs/QmRcHuYzAyswytBuMF78rj3LTChYszomRFXNg4685ZN1WM/go-block-format"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"

	njwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/jwt"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/service"
)

// defaultSessionDuration after which session token expires
const defaultSessionDuration = time.Hour * 24 * 7 * 4

// blockTimeout is the context timeout for getting a local block
const blockTimeout = time.Second * 10

// inboxMessagePageSize is the page size used when checking messages
const inboxMessagePageSize = 10

// validation errors
const (
	errInvalidAddress = "invalid address"
	errUnauthorized   = "unauthorized"
	errForbidden      = "forbidden"
)

const CafeServiceProtocol = protocol.ID("/textile/cafe/1.0.0")

// CafeService is a libp2p pinning and offline message service
type CafeService struct {
	service    *service.Service
	datastore  repo.Datastore
	inbox      *CafeInbox
	httpAddr   string
	swarmAddrs []string
	open       bool
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
	return CafeServiceProtocol
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
	case pb.Message_CAFE_YOU_HAVE_MAIL:
		return h.handleNotifyClient(pid, env)
	default:
		return nil, nil
	}
}

// Register creates a session with a cafe
func (h *CafeService) Register(cafe peer.ID) error {
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

	env, err := h.service.NewEnvelope(pb.Message_CAFE_REGISTRATION, reg, nil, false)
	if err != nil {
		return err
	}
	renv, err := h.service.SendRequest(cafe, env)
	if err != nil {
		return err
	}

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
		CafeId:     cafe.Pretty(),
		Access:     res.Access,
		Refresh:    res.Refresh,
		Expiry:     exp,
		HttpAddr:   res.HttpAddr,
		SwarmAddrs: res.SwarmAddrs,
	}
	return h.datastore.CafeSessions().AddOrUpdate(session)
}

// Store stores (pins) content on a cafe and returns a list of successful cids
func (h *CafeService) Store(cids []string, cafe peer.ID) ([]string, error) {
	var stored []string

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
		if err := h.sendBlock(decoded, cafe, accessToken); err != nil {
			log.Errorf("error sending block: %s", err)
			continue
		}
		stored = append(stored, id)
	}
	return stored, nil
}

// StoreThread pushes a thread to a cafe backup
func (h *CafeService) StoreThread(thrd *repo.Thread, cafe peer.ID) error {
	plaintext, err := proto.Marshal(&pb.CafeThread{
		Key:       thrd.Key,
		Sk:        thrd.PrivKey,
		Name:      thrd.Name,
		Schema:    thrd.Schema,
		Initiator: thrd.Initiator,
		Type:      int32(thrd.Type),
		State:     int32(thrd.State),
		Head:      thrd.Head,
	})
	if err != nil {
		return err
	}
	ciphertext, err := h.service.Account.Encrypt(plaintext)
	if err != nil {
		return err
	}

	if _, err := h.sendCafeRequest(cafe, func(session *repo.CafeSession) (*pb.Envelope, error) {
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
func (h *CafeService) DeliverMessage(mid string, pid peer.ID, cafe peer.ID) error {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_DELIVER_MESSAGE, &pb.CafeDeliverMessage{
		Id:       mid,
		ClientId: pid.Pretty(),
	}, nil, false)
	if err != nil {
		return err
	}
	return h.service.SendMessage(nil, cafe, env)
}

// CheckMessages asks each session's inbox for new messages
func (h *CafeService) CheckMessages(cafe peer.ID) error {
	renv, err := h.sendCafeRequest(cafe, func(session *repo.CafeSession) (*pb.Envelope, error) {
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
		if err := h.inbox.Add(msg); err != nil {
			return err
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
	renv, err := h.sendCafeRequest(cafe, func(session *repo.CafeSession) (*pb.Envelope, error) {
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

	// apparently there's more new messages waiting...
	return h.CheckMessages(cafe)
}

// notifyClient attempts to ping a client that has messages waiting to download
func (h *CafeService) notifyClient(pid peer.ID) error {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_YOU_HAVE_MAIL, nil, nil, false)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), service.DirectTimeout) // fail fast
	defer cancel()
	return h.service.SendMessage(ctx, pid, env)
}

// sendCafeRequest sends an authenticated request, retrying once after a session refresh
func (h *CafeService) sendCafeRequest(
	cafe peer.ID, envFactory func(*repo.CafeSession) (*pb.Envelope, error)) (*pb.Envelope, error) {
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
		CafeId:     session.CafeId,
		Access:     res.Access,
		Refresh:    res.Refresh,
		Expiry:     exp,
		HttpAddr:   res.HttpAddr,
		SwarmAddrs: res.SwarmAddrs,
	}
	if err := h.datastore.CafeSessions().AddOrUpdate(refreshed); err != nil {
		return nil, err
	}
	return refreshed, nil
}

// sendBlock sends a block by cid to a peer
func (h *CafeService) sendBlock(id cid.Cid, pid peer.ID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), blockTimeout)
	defer cancel()
	block, err := h.service.Node.Blocks.GetBlock(ctx, id)
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

	session, err := jwt.NewSession(
		h.service.Node.PrivateKey,
		pid,
		h.Protocol(),
		defaultSessionDuration,
		h.httpAddr,
		h.swarmAddrs,
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
		h.service.Node.PrivateKey,
		spid,
		h.Protocol(),
		defaultSessionDuration,
		h.httpAddr,
		h.swarmAddrs,
	)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.RequestId, true)
}

// handleSession receives a store request
func (h *CafeService) handleStore(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	store := new(pb.CafeStore)
	err := ptypes.UnmarshalAny(env.Message.Payload, store)
	if err != nil {
		return nil, err
	}

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

	thrd := &repo.CafeClientThread{
		Id:         store.Id,
		ClientId:   client.Id,
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
	err := ptypes.UnmarshalAny(env.Message.Payload, msg)
	if err != nil {
		return nil, err
	}

	client := h.datastore.CafeClients().Get(msg.ClientId)
	if client == nil {
		log.Warningf("received message from %s for unknown client %s", pid.Pretty(), msg.ClientId)
		return nil, nil
	}

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

// handleDeliverMessage receives an inbox message for a client
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

// handleCheckMessages receives a check inbox messages request
func (h *CafeService) handleCheckMessages(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	check := new(pb.CafeCheckMessages)
	err := ptypes.UnmarshalAny(env.Message.Payload, check)
	if err != nil {
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

	return h.service.NewEnvelope(pb.Message_CAFE_MESSAGES, res, &env.Message.RequestId, true)
}

// handleDeleteMessages receives a message delete request
func (h *CafeService) handleDeleteMessages(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	del := new(pb.CafeDeleteMessages)
	err := ptypes.UnmarshalAny(env.Message.Payload, del)
	if err != nil {
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
	return h.service.Node.PrivateKey.GetPublic(), nil
}

// setAddrs sets addresses used in sessions generated by this host
func (h *CafeService) setAddrs(bindedHttpAddr string, publicIP string, swarmPorts config.SwarmPorts) {
	if publicIP == "" {
		publicIP = "127.0.0.1"
	}

	// set the http address where peers can reach this cafe
	parts := strings.Split(bindedHttpAddr, ":")
	if len(parts) == 2 {
		h.httpAddr = fmt.Sprintf("%s:%s", publicIP, parts[1])
		log.Infof("cafe http api address: %s", h.httpAddr)
	}

	// set the swarm multiaddress(es) where other peers can reach this cafe
	h.swarmAddrs = []string{
		fmt.Sprintf("/ip4/%s/tcp/%s", publicIP, swarmPorts.TCP),
		fmt.Sprintf("/ip4/%s/tcp/%s/ws", publicIP, swarmPorts.WS),
	}
	log.Infof("cafe multiaddresses: %s, %s", h.swarmAddrs[0], h.swarmAddrs[1])
}

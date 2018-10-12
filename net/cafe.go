package net

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/jwt"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/net/service"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmVzK524a2VWLqyvtBeiHKsUAWYgeAk4DBeZoY7vpNPNRx/go-block-format"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"time"
)

const defaultSessionDuration = time.Hour * 24 * 7 * 4

// CafeService is a libp2p service for proxing
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

// Node returns the underlying ipfs Node
func (h *CafeService) Node() *core.IpfsNode {
	return h.service.Node
}

// Datastore returns the underlying datastore
func (h *CafeService) Datastore() repo.Datastore {
	return h.service.Datastore
}

// Ping pings another peer
func (h *CafeService) Ping(pid peer.ID) (service.PeerStatus, error) {
	return h.service.Ping(pid)
}

// VerifyEnvelope calls service verify
func (h *CafeService) VerifyEnvelope(env *pb.Envelope) error {
	return h.service.VerifyEnvelope(env)
}

// Handle is called by the underlying service handler method
func (h *CafeService) Handle(mtype pb.Message_Type) func(peer.ID, *pb.Envelope) (*pb.Envelope, error) {
	switch mtype {
	case pb.Message_CAFE_CHALLENGE:
		return h.handleChallenge
	case pb.Message_CAFE_REGISTRATION:
		return h.handleRegistration
	case pb.Message_CAFE_STORE:
		return h.handleStore
	case pb.Message_CAFE_BLOCK:
		return h.handleBlock
	default:
		return nil
	}
}

// Challenge asks a fellow peer for a cafe challenge
func (h *CafeService) Challenge(kp *keypair.Full, cafe peer.ID) (*pb.CafeNonce, error) {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_CHALLENGE, &pb.CafeChallenge{
		Address: kp.Address(),
	}, nil, false)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
	defer cancel()
	renv, err := h.service.SendRequest(ctx, cafe, env)
	if err != nil {
		return nil, err
	}
	return h.handleNonce(cafe, renv)
}

// Register registers a peer with a cafe
func (h *CafeService) Register(reg *pb.CafeRegistration, cafe peer.ID) (*pb.CafeSession, error) {
	env, err := h.service.NewEnvelope(pb.Message_CAFE_REGISTRATION, reg, nil, false)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
	defer cancel()
	renv, err := h.service.SendRequest(ctx, cafe, env)
	if err != nil {
		return nil, err
	}
	return h.handleSession(cafe, renv)
}

// Store stores (pins) content on a cafe
func (h *CafeService) Store(cids []cid.Cid, cafe peer.ID) error {
	var ids []string
	for _, c := range cids {
		ids = append(ids, c.String())
	}
	store := &pb.CafeCidList{Cids: ids}
	env, err := h.service.NewEnvelope(pb.Message_CAFE_STORE, store, nil, false)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
	defer cancel()
	renv, err := h.service.SendRequest(ctx, cafe, env)
	if err != nil {
		return err
	}
	if err := h.service.HandleError(cafe, env); err != nil {
		return err
	}

	// unpack response as a request list
	req := new(pb.CafeCidList)
	err = ptypes.UnmarshalAny(renv.Message.Payload, req)
	if err != nil {
		return err
	}
	if len(req.Cids) == 0 {
		log.Debugf("peer %s requested no blocks", cafe.Pretty())
		return nil
	}
	log.Debugf("sending %d blocks to %s", len(req.Cids), cafe.Pretty())

	// send each block
	for _, id := range req.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			continue
		}
		h.SendBlock(*decoded, cafe)
	}
	return nil
}

// SendBlock sends a block by cid to a peer
func (h *CafeService) SendBlock(id cid.Cid, pid peer.ID) error {
	// get block locally
	ctx, cancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
	defer cancel()
	block, err := h.Node().Blocks.GetBlock(ctx, &id)
	if err != nil {
		return err
	}

	// send it
	pblock := &pb.CafeBlock{
		Cid:     block.Cid().String(),
		RawData: block.RawData(),
	}
	env, err := h.service.NewEnvelope(pb.Message_CAFE_BLOCK, pblock, nil, false)
	if err != nil {
		return err
	}
	sctx, scancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
	defer scancel()
	return h.service.SendMessage(sctx, pid, env)
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
		return h.service.NewError(400, "invalid address", env.Message.RequestId)
	}

	// generate a new random nonce
	nonce := &repo.Nonce{
		Value:   ksuid.New().String(),
		Address: req.Address,
		Date:    time.Now(),
	}
	if err := h.Datastore().Nonces().Add(nonce); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// return a wrapped response
	return h.service.NewEnvelope(pb.Message_CAFE_NONCE, &pb.CafeNonce{
		Value: nonce.Value,
	}, &env.Message.RequestId, true)
}

// handleNonce receives a challenge response
func (h *CafeService) handleNonce(pid peer.ID, env *pb.Envelope) (*pb.CafeNonce, error) {
	if err := h.service.HandleError(pid, env); err != nil {
		return nil, err
	}
	res := new(pb.CafeNonce)
	if err := ptypes.UnmarshalAny(env.Message.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// handleRegistration receives a registration request
func (h *CafeService) handleRegistration(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	reg := new(pb.CafeRegistration)
	if err := ptypes.UnmarshalAny(env.Message.Payload, reg); err != nil {
		return nil, err
	}

	// lookup the nonce
	snonce := h.Datastore().Nonces().Get(reg.Value)
	if snonce == nil {
		return h.service.NewError(403, "forbidden", env.Message.RequestId)
	}
	if snonce.Address != reg.Address {
		return h.service.NewError(403, "forbidden", env.Message.RequestId)
	}

	// validate address
	accnt, err := keypair.Parse(reg.Address)
	if err != nil {
		return nil, err
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		return h.service.NewError(400, "invalid address", env.Message.RequestId)
	}

	// verify
	payload := []byte(reg.Value + reg.Nonce)
	if err := accnt.Verify(payload, reg.Sig); err != nil {
		return h.service.NewError(403, "forbidden", env.Message.RequestId)
	}

	// create new
	now := time.Now()
	account := &repo.Account{
		Id:       pid.Pretty(),
		Address:  reg.Address,
		Created:  now,
		LastSeen: now,
	}
	if err := h.Datastore().Accounts().Add(account); err != nil {
		return h.service.NewError(409, "conflict", env.Message.RequestId)
	}

	// get a session
	session, err := jwt.NewSession(h.Node().PrivateKey, pid, h.Protocol(), defaultSessionDuration)
	if err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// delete the nonce
	if err := h.Datastore().Nonces().Delete(snonce.Value); err != nil {
		return h.service.NewError(500, err.Error(), env.Message.RequestId)
	}

	// return a wrapped response
	return h.service.NewEnvelope(pb.Message_CAFE_SESSION, session, &env.Message.RequestId, true)
}

// handleSession receives a session response
func (h *CafeService) handleSession(pid peer.ID, env *pb.Envelope) (*pb.CafeSession, error) {
	if err := h.service.HandleError(pid, env); err != nil {
		return nil, err
	}
	res := new(pb.CafeSession)
	if err := ptypes.UnmarshalAny(env.Message.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// handleSession receives a store request
func (h *CafeService) handleStore(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	store := new(pb.CafeCidList)
	err := ptypes.UnmarshalAny(env.Message.Payload, store)
	if err != nil {
		return nil, err
	}

	// ignore cids for blocks already present in local datastore
	var need []string
	for _, id := range store.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			continue
		}
		has, err := h.Node().Blockstore.Has(decoded)
		if err != nil || !has {
			need = append(need, decoded.String())
		}
	}
	res := &pb.CafeCidList{Cids: need}
	return h.service.NewEnvelope(pb.Message_CAFE_STORE, res, &env.Message.RequestId, true)
}

// handleBlock receives a block message
func (h *CafeService) handleBlock(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	pblock := new(pb.CafeBlock)
	err := ptypes.UnmarshalAny(env.Message.Payload, pblock)
	if err != nil {
		return nil, err
	}
	id, err := cid.Decode(pblock.Cid)
	if err != nil {
		return nil, err
	}

	// add a new block to the local datastore
	block, err := blocks.NewBlockWithCid(pblock.RawData, id)
	if err != nil {
		return nil, err
	}
	if err := h.Node().Blocks.AddBlock(block); err != nil {
		return nil, err
	}
	log.Debugf("pinned %s from %s", pblock.Cid, pid.Pretty())
	return nil, nil
}

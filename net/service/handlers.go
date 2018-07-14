package service

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net/common"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"gx/ipfs/Qmej7nf81hi2x2tvjRBF3mcp74sQyuDH4VMYDGd1YtXjb2/go-block-format"
)

func (s *TextileService) HandlerForMsgType(t pb.Message_Type) func(peer.ID, *pb.Message, interface{}) (*pb.Message, error) {
	switch t {
	case pb.Message_PING:
		return s.handlePing
	case pb.Message_THREAD_INVITE:
		return s.handleThreadInvite
	case pb.Message_THREAD_JOIN:
		return s.handleThreadJoin
	case pb.Message_THREAD_LEAVE:
		return s.handleThreadLeave
	case pb.Message_THREAD_DATA:
		return s.handleThreadData
	case pb.Message_THREAD_ANNOTATION:
		return s.handleThreadAnnotation
	case pb.Message_OFFLINE_ACK:
		return s.handleOfflineAck
	case pb.Message_OFFLINE_RELAY:
		return s.handleOfflineRelay
	case pb.Message_BLOCK:
		return s.handleBlock
	case pb.Message_STORE:
		return s.handleStore
	case pb.Message_ERROR:
		return s.handleError
	default:
		return nil
	}
}

func (s *TextileService) handlePing(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	log.Debugf("received PING message from %s", pid.Pretty())
	return pmes, nil
}

func (s *TextileService) handleThreadInvite(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	log.Debug("received THREAD_INVITE message")
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	signed := new(pb.SignedThreadBlock)
	err := ptypes.UnmarshalAny(pmes.Payload, signed)
	if err != nil {
		return nil, err
	}
	invite := new(pb.ThreadInvite)
	err = proto.Unmarshal(signed.Data, invite)
	if err != nil {
		return nil, err
	}
	issuerpk, err := libp2pc.UnmarshalPublicKey(signed.IssuerPubKey)
	if err != nil {
		return nil, err
	}

	// load thread
	thrd := s.getThread(signed.ThreadId)
	if thrd != nil {
		// known thread and invite meant for us
		if invite.Block.Target == s.self.Pretty() {
			return nil, errors.New("thread already exists")
		}

		// verify thread sig
		if err := thrd.Verify(signed.Data, signed.ThreadSignature); err != nil {
			return nil, err
		}

		// verify issuer sig
		if err := verify(issuerpk, signed.Data, signed.IssuerSignature); err != nil {
			return nil, err
		}

		return nil, thrd.HandleBlock(signed.Id)
	} else {
		// unknown thread and external invite... shouldn't happen
		if invite.Type == pb.ThreadInvite_EXTERNAL {
			return nil, errors.New("invalid invite block")
		}
		// unknown thread and invite not meant for us... shouldn't happen
		if invite.Block.Target != s.self.Pretty() {
			return nil, errors.New("invalid invite block")
		}
	}

	// lastly, unknown thread and invite meant for us
	// unpack new thread secret that should be encrypted with our key
	skb, err := crypto.Decrypt(s.node.PrivateKey, invite.Block.TargetKey)
	if err != nil {
		return nil, err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}

	// verify thread sig
	if err := verify(sk.GetPublic(), signed.Data, signed.ThreadSignature); err != nil {
		return nil, err
	}

	// verify issuer sig
	if err := verify(issuerpk, signed.Data, signed.IssuerSignature); err != nil {
		return nil, err
	}

	// add the new thread (name will bump if already exists, e.g., cats -> cats_1)
	thrd, err = s.addThread(signed.ThreadName, sk)
	if err != nil {
		return nil, err
	}

	// accept it, yolo
	// TODO: Don't auto accept. Need to show some UI with pending invites.
	thrd.Join(issuerpk, invite.Block.Id)

	log.Debugf("accepted invite to thread %s with name %s", thrd.Id, thrd.Name)

	return nil, thrd.HandleBlock(signed.Id)
}

func (s *TextileService) handleThreadJoin(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	log.Debug("received THREAD_JOIN message")
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	signed := new(pb.SignedThreadBlock)
	err := ptypes.UnmarshalAny(pmes.Payload, signed)
	if err != nil {
		return nil, err
	}
	join := new(pb.ThreadJoin)
	err = proto.Unmarshal(signed.Data, join)
	if err != nil {
		return nil, err
	}
	issuerpk, err := libp2pc.UnmarshalPublicKey(signed.IssuerPubKey)
	if err != nil {
		return nil, err
	}

	// load thread
	thrd := s.getThread(signed.ThreadId)
	if thrd == nil {
		return nil, errors.New("invalid join block")
	}

	// verify thread sig
	if err := thrd.Verify(signed.Data, signed.ThreadSignature); err != nil {
		return nil, err
	}

	// verify issuer sig
	if err := verify(issuerpk, signed.Data, signed.IssuerSignature); err != nil {
		return nil, err
	}

	// get the issuer id
	issuerId, err := peer.IDFromPublicKey(issuerpk)
	if err != nil {
		return nil, err
	}

	// add issuer as a new local peer
	newPeer := &repo.Peer{
		Row:      ksuid.New().String(),
		Id:       issuerId.Pretty(),
		ThreadId: join.Block.ThreadPubKey,
		PubKey:   signed.IssuerPubKey,
	}
	if err := s.datastore.Peers().Add(newPeer); err != nil {
		return nil, err
	}

	return nil, thrd.HandleBlock(signed.Id)
}

func (s *TextileService) handleThreadLeave(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	log.Debug("received THREAD_LEAVE message")
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	signed := new(pb.SignedThreadBlock)
	err := ptypes.UnmarshalAny(pmes.Payload, signed)
	if err != nil {
		return nil, err
	}
	leave := new(pb.ThreadLeave)
	err = proto.Unmarshal(signed.Data, leave)
	if err != nil {
		return nil, err
	}
	issuerpk, err := libp2pc.UnmarshalPublicKey(signed.IssuerPubKey)
	if err != nil {
		return nil, err
	}

	// load thread
	thrd := s.getThread(signed.ThreadId)
	if thrd == nil {
		return nil, errors.New("invalid leave block")
	}

	// verify thread sig
	if err := thrd.Verify(signed.Data, signed.ThreadSignature); err != nil {
		return nil, err
	}

	// verify issuer sig
	if err := verify(issuerpk, signed.Data, signed.IssuerSignature); err != nil {
		return nil, err
	}

	// remove peer
	if err := s.datastore.Peers().Delete(leave.Block.Target, leave.Block.ThreadPubKey); err != nil {
		return nil, err
	}

	return nil, thrd.HandleBlock(signed.Id)
}

func (s *TextileService) handleThreadData(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	log.Debug("received THREAD_DATA message")
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	signed := new(pb.SignedThreadBlock)
	err := ptypes.UnmarshalAny(pmes.Payload, signed)
	if err != nil {
		return nil, err
	}
	data := new(pb.ThreadData)
	err = proto.Unmarshal(signed.Data, data)
	if err != nil {
		return nil, err
	}

	// load thread
	thrd := s.getThread(signed.ThreadId)
	if thrd == nil {
		return nil, common.OutOfOrderMessage
	}

	// verify
	if err := thrd.Verify(signed.Data, signed.ThreadSignature); err != nil {
		return nil, err
	}

	return nil, thrd.HandleBlock(signed.Id)
}

func (s *TextileService) handleThreadAnnotation(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	log.Debug("received THREAD_ANNOTATION message")
	return nil, errors.New("TODO")
}

func (s *TextileService) handleOfflineAck(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	id, err := peer.IDB58Decode(string(pmes.Payload.Value))
	if err != nil {
		return nil, err
	}
	pointer, err := s.datastore.Pointers().Get(id)
	if err != nil {
		return nil, err
	}
	if pointer.CancelId == nil || pointer.CancelId.Pretty() != pid.Pretty() {
		return nil, errors.New("peer is not authorized to delete pointer")
	}
	err = s.datastore.Pointers().Delete(id)
	if err != nil {
		return nil, err
	}
	log.Debugf("received OFFLINE_ACK message from %s", pid.Pretty())
	return nil, nil
}

func (s *TextileService) handleOfflineRelay(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	plaintext, err := crypto.Decrypt(s.node.PrivateKey, pmes.Payload.Value)
	if err != nil {
		return nil, err
	}

	// Unmarshal plaintext
	env := pb.Envelope{}
	err = proto.Unmarshal(plaintext, &env)
	if err != nil {
		return nil, err
	}

	// Validate the signature
	ser, err := proto.Marshal(env.Message)
	if err != nil {
		return nil, err
	}
	pubkey, err := libp2pc.UnmarshalPublicKey(env.Pubkey)
	if err != nil {
		return nil, err
	}
	valid, err := pubkey.Verify(ser, env.Signature)
	if err != nil || !valid {
		return nil, err
	}

	id, err := peer.IDFromPublicKey(pubkey)
	if err != nil {
		return nil, err
	}

	// Get handler for this message type
	handler := s.HandlerForMsgType(env.Message.Type)
	if handler == nil {
		log.Debug("got back nil handler from HandlerForMsgType")
		return nil, nil
	}

	// Dispatch handler
	_, err = handler(id, env.Message, true)
	if err != nil {
		log.Errorf("handle message error: %s", err)
		return nil, err
	}
	log.Debugf("received OFFLINE_RELAY message from %s", pid.Pretty())
	return nil, nil
}

func (s *TextileService) handleBlock(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	pbblock := new(pb.Block)
	err := ptypes.UnmarshalAny(pmes.Payload, pbblock)
	if err != nil {
		return nil, err
	}
	id, err := cid.Decode(pbblock.Cid)
	if err != nil {
		return nil, err
	}
	block, err := blocks.NewBlockWithCid(pbblock.RawData, id)
	if err != nil {
		return nil, err
	}
	if err := s.node.Blocks.AddBlock(block); err != nil {
		return nil, err
	}
	log.Debugf("received IPFS_BLOCK message from %s", pid.Pretty())
	return nil, nil
}

func (s *TextileService) handleStore(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	errorResponse := func(error string) *pb.Message {
		payload := &any.Any{Value: []byte(error)}
		message := &pb.Message{
			Type:    pb.Message_ERROR,
			Payload: payload,
		}
		return message
	}

	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	cList := new(pb.CidList)
	err := ptypes.UnmarshalAny(pmes.Payload, cList)
	if err != nil {
		return errorResponse("could not unmarshall message"), err
	}
	var need []string
	for _, id := range cList.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			continue
		}
		has, err := s.node.Blockstore.Has(decoded)
		if err != nil || !has {
			need = append(need, decoded.String())
		}
	}
	log.Debugf("received STORE message from %s", pid.Pretty())
	log.Debugf("requesting %d blocks from %s", len(need), pid.Pretty())

	resp := new(pb.CidList)
	resp.Cids = need
	payload, err := ptypes.MarshalAny(resp)
	if err != nil {
		return errorResponse("error marshalling response"), err
	}
	message := &pb.Message{
		Type:    pb.Message_STORE,
		Payload: payload,
	}
	return message, nil
}

func (s *TextileService) handleError(peer peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	errorMessage := new(pb.Error)
	err := ptypes.UnmarshalAny(pmes.Payload, errorMessage)
	if err != nil {
		return nil, err
	}

	// TODO

	return nil, nil
}

// verify verifies a signature
func verify(pk libp2pc.PubKey, data []byte, sig []byte) error {
	good, err := pk.Verify(data, sig)
	if err != nil || !good {
		return errors.New("bad signature")
	}
	return nil
}

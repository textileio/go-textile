package service

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net/common"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/util"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"gx/ipfs/Qmej7nf81hi2x2tvjRBF3mcp74sQyuDH4VMYDGd1YtXjb2/go-block-format"
)

func (s *TextileService) HandlerForMsgType(t pb.Message_Type) func(peer.ID, *pb.Envelope, interface{}) (*pb.Envelope, error) {
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
	case pb.Message_THREAD_IGNORE:
		return s.handleThreadIgnore
	case pb.Message_THREAD_MERGE:
		return s.handleThreadMerge
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

func (s *TextileService) handlePing(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	log.Debugf("received PING message from %s", pid.Pretty())
	return pmes, nil
}

func (s *TextileService) handleThreadInvite(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	log.Debug("received THREAD_INVITE message")
	signed, err := unpackMessage(pmes)
	if err != nil {
		return nil, err
	}
	invite := new(pb.ThreadInvite)
	if err := proto.Unmarshal(signed.Block, invite); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(invite.Header.ThreadPk)
	_, thrd := s.getThread(threadId)
	if thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	// check if it's meant for us (should be, but safety first)
	if invite.InviteeId != s.self.Pretty() {
		return nil, errors.New("invalid invite block")
	}

	// unknown thread and invite meant for us
	// unpack new thread secret that should be encrypted with our key
	skb, err := crypto.Decrypt(s.node.PrivateKey, invite.SkCipher)
	if err != nil {
		return nil, err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}

	// verify thread sig
	if err := crypto.Verify(sk.GetPublic(), signed.Block, signed.ThreadSig); err != nil {
		return nil, err
	}

	// add to local ipfs for later use when joining
	pmesb, err := proto.Marshal(pmes)
	if err != nil {
		return nil, err
	}
	if _, err := util.PinData(s.node, bytes.NewReader(pmesb)); err != nil {
		return nil, err
	}

	// TODO: add invite notification that can be used to accept

	return nil, nil
}

func (s *TextileService) handleThreadJoin(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	log.Debug("received THREAD_JOIN message")
	signed, err := unpackMessage(pmes)
	if err != nil {
		return nil, err
	}
	join := new(pb.ThreadJoin)
	if err := proto.Unmarshal(signed.Block, join); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(join.Header.ThreadPk)
	_, thrd := s.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid join block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	if _, _, err := thrd.HandleJoinBlock(&pid, pmes, signed, join, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *TextileService) handleThreadMerge(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	log.Debug("received THREAD_MERGE message")
	signed, err := unpackMessage(pmes)
	if err != nil {
		return nil, err
	}
	merge := new(pb.ThreadMerge)
	if err := proto.Unmarshal(signed.Block, merge); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(merge.Header.ThreadPk)
	_, thrd := s.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid merge block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	if _, err := thrd.HandleMergeBlock(&pid, pmes, signed, merge, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *TextileService) handleThreadLeave(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	log.Debug("received THREAD_LEAVE message")
	signed, err := unpackMessage(pmes)
	if err != nil {
		return nil, err
	}
	leave := new(pb.ThreadLeave)
	if err := proto.Unmarshal(signed.Block, leave); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(leave.Header.ThreadPk)
	_, thrd := s.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid leave block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	if _, err := thrd.HandleLeaveBlock(&pid, pmes, signed, leave, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *TextileService) handleThreadData(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	log.Debug("received THREAD_DATA message")
	signed, err := unpackMessage(pmes)
	if err != nil {
		return nil, err
	}
	data := new(pb.ThreadData)
	if err := proto.Unmarshal(signed.Block, data); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(data.Header.ThreadPk)
	_, thrd := s.getThread(threadId)
	if thrd == nil {
		return nil, common.OutOfOrderMessage
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	if _, err := thrd.HandleDataBlock(&pid, pmes, signed, data, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *TextileService) handleThreadIgnore(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	log.Debug("received THREAD_IGNORE message")
	signed, err := unpackMessage(pmes)
	if err != nil {
		return nil, err
	}
	ignore := new(pb.ThreadIgnore)
	if err := proto.Unmarshal(signed.Block, ignore); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(ignore.Header.ThreadPk)
	_, thrd := s.getThread(threadId)
	if thrd == nil {
		return nil, common.OutOfOrderMessage
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	if _, err := thrd.HandleIgnoreBlock(&pid, pmes, signed, ignore, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *TextileService) handleOfflineAck(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	if pmes.Message.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	id, err := peer.IDB58Decode(string(pmes.Message.Payload.Value))
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
	if err := s.datastore.Pointers().Delete(id); err != nil {
		return nil, err
	}
	log.Debugf("received OFFLINE_ACK message from %s", pid.Pretty())
	return nil, nil
}

func (s *TextileService) handleOfflineRelay(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	if pmes.Message.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	plaintext, err := crypto.Decrypt(s.node.PrivateKey, pmes.Message.Payload.Value)
	if err != nil {
		return nil, err
	}

	// unmarshal plaintext
	env := &pb.Envelope{}
	err = proto.Unmarshal(plaintext, env)
	if err != nil {
		return nil, err
	}

	// validate envelope signature
	ser, err := proto.Marshal(env.Message)
	if err != nil {
		return nil, err
	}
	pk, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return nil, err
	}
	if err := crypto.Verify(pk, ser, env.Sig); err != nil {
		return nil, err
	}

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return nil, err
	}

	// get handler for this message type
	handler := s.HandlerForMsgType(env.Message.Type)
	if handler == nil {
		log.Debug("got back nil handler from HandlerForMsgType")
		return nil, nil
	}

	// dispatch handler
	_, err = handler(id, env, true)
	if err != nil {
		log.Errorf("handle message error: %s", err)
		return nil, err
	}
	log.Debugf("received OFFLINE_RELAY message from %s", pid.Pretty())
	return nil, nil
}

func (s *TextileService) handleBlock(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	if pmes.Message.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	pbblock := new(pb.Block)
	err := ptypes.UnmarshalAny(pmes.Message.Payload, pbblock)
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

func (s *TextileService) handleStore(pid peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	errorResponse := func(error string) *pb.Message {
		payload := &any.Any{Value: []byte(error)}
		message := &pb.Message{
			Type:    pb.Message_ERROR,
			Payload: payload,
		}
		return message
	}

	if pmes.Message.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	cList := new(pb.CidList)
	err := ptypes.UnmarshalAny(pmes.Message.Payload, cList)
	if err != nil {
		return s.newEnvelope(errorResponse("could not unmarshall message"))
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
		return s.newEnvelope(errorResponse("error marshalling response"))
	}
	message := &pb.Message{
		Type:    pb.Message_STORE,
		Payload: payload,
	}
	return s.newEnvelope(message)
}

func (s *TextileService) handleError(peer peer.ID, pmes *pb.Envelope, options interface{}) (*pb.Envelope, error) {
	if pmes.Message.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	errorMessage := new(pb.Error)
	err := ptypes.UnmarshalAny(pmes.Message.Payload, errorMessage)
	if err != nil {
		return nil, err
	}

	// TODO

	return nil, nil
}

func unpackMessage(pmes *pb.Envelope) (*pb.SignedThreadBlock, error) {
	if pmes.Message.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(pmes.Message.Payload, signed); err != nil {
		return nil, err
	}
	return signed, nil
}

package service

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"gx/ipfs/Qmej7nf81hi2x2tvjRBF3mcp74sQyuDH4VMYDGd1YtXjb2/go-block-format"
)

const (
	ChatMessageMaxCharacters = 20000
	ChatSubjectMaxCharacters = 500
	//DefaultPointerPrefixLength = 14
)

func (s *TextileService) HandlerForMsgType(t pb.Message_MessageType) func(peer.ID, *pb.Message, interface{}) (*pb.Message, error) {
	switch t {
	case pb.Message_PING:
		return s.handlePing
	case pb.Message_THREAD_BLOCK:
		return s.handleThreadBlock
	case pb.Message_FOLLOW:
		return s.handleFollow
	case pb.Message_UNFOLLOW:
		return s.handleUnFollow
	case pb.Message_OFFLINE_ACK:
		return s.handleOfflineAck
	case pb.Message_OFFLINE_RELAY:
		return s.handleOfflineRelay
	case pb.Message_CHAT:
		return s.handleChat
	//case pb.Message_MODERATOR_ADD:
	//	return s.handleModeratorAdd
	//case pb.Message_MODERATOR_REMOVE:
	//	return s.handleModeratorRemove
	case pb.Message_IPFS_BLOCK:
		return s.handleIPFSBlock
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

func (s *TextileService) handleThreadBlock(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	log.Debug("received THREAD_BLOCK message")

	// unpack it
	signed := new(pb.SignedBlock)
	err := ptypes.UnmarshalAny(pmes.Payload, signed)
	if err != nil {
		return nil, err
	}
	block := new(pb.Block)
	err = proto.Unmarshal(signed.Data, block)
	if err != nil {
		return nil, err
	}
	thrd := s.getThread(signed.ThreadId)

	switch block.Type {
	case pb.Block_INVITE:
		log.Debug("handling Block_INVITE...")
		if thrd != nil {
			return nil, errors.New("thread exists")
		}
		if block.Target != s.self.Pretty() {
			// TODO: should not be error right?
			return nil, errors.New("invalid invite target")
		}
		skb, err := crypto.Decrypt(s.node.PrivateKey, block.TargetKey)
		if err != nil {
			return nil, err
		}
		sk, err := libp2pc.UnmarshalPrivateKey(skb)
		if err != nil {
			return nil, err
		}
		good, err := sk.GetPublic().Verify(signed.Data, signed.Signature)
		if err != nil || !good {
			return nil, errors.New("bad signature")
		}
		// TODO: handle when name leads to conflict (add an int)
		thrd, err = s.addThread(signed.ThreadName, sk)
		if err != nil {
			return nil, err
		}

		// add inviter as local peer
		ppk, err := libp2pc.UnmarshalPublicKey(signed.IssuerPubKey)
		if err != nil {
			return nil, err
		}
		ppkb, err := ppk.Bytes()
		if err != nil {
			return nil, err
		}
		peerId, err := peer.IDFromPublicKey(ppk)
		if err != nil {
			return nil, err
		}
		newPeer := &repo.Peer{
			Row:      ksuid.New().String(),
			Id:       peerId.Pretty(),
			ThreadId: thrd.Id,
			PubKey:   ppkb,
		}
		if err := s.datastore.Peers().Add(newPeer); err != nil {
			return nil, err
		}
	case pb.Block_PHOTO:
		log.Debug("handling Block_PHOTO")
		if thrd == nil {
			return nil, errors.New("thread not found")
		}
		good, err := thrd.Verify(signed.Data, signed.Signature)
		if err != nil || !good {
			return nil, errors.New("bad signature")
		}

	case pb.Block_COMMENT:
		return nil, errors.New("TODO")
	case pb.Block_LIKE:
		return nil, errors.New("TODO")
	}

	// handle block
	return nil, thrd.HandleBlock(signed.Id)
}

func (s *TextileService) handleFollow(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	sd := new(pb.SignedData)
	err := ptypes.UnmarshalAny(pmes.Payload, sd)
	if err != nil {
		return nil, err
	}
	pubkey, err := libp2pc.UnmarshalPublicKey(sd.SenderPubkey)
	if err != nil {
		return nil, err
	}
	id, err := peer.IDFromPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	data := new(pb.SignedData_Command)
	err = proto.Unmarshal(sd.SerializedData, data)
	if err != nil {
		return nil, err
	}
	if data.PeerID != s.node.Identity.Pretty() {
		return nil, errors.New("follow message doesn't include correct peer id")
	}
	if data.Type != pb.Message_FOLLOW {
		return nil, errors.New("data type is not follow")
	}
	good, err := pubkey.Verify(sd.SerializedData, sd.Signature)
	if err != nil || !good {
		return nil, errors.New("bad signature")
	}

	//proof := append(sd.SerializedData, sd.Signature...)
	//err = s.datastore.Followers().Put(id.Pretty(), proof)
	//if err != nil {
	//	return nil, err
	//}
	//n := notifications.FollowNotification{notifications.NewID(), "follow", id.Pretty()}
	//s.broadcast <- n
	//s.datastore.Notifications().Put(n.ID, n, n.Type, time.Now())
	log.Debugf("received FOLLOW message from %s", id.Pretty())
	return nil, nil
}

func (s *TextileService) handleUnFollow(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	sd := new(pb.SignedData)
	err := ptypes.UnmarshalAny(pmes.Payload, sd)
	if err != nil {
		return nil, err
	}
	pubkey, err := libp2pc.UnmarshalPublicKey(sd.SenderPubkey)
	if err != nil {
		return nil, err
	}
	id, err := peer.IDFromPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	data := new(pb.SignedData_Command)
	err = proto.Unmarshal(sd.SerializedData, data)
	if err != nil {
		return nil, err
	}
	if data.PeerID != s.node.Identity.Pretty() {
		return nil, errors.New("unfollow message doesn't include correct peer id")
	}
	if data.Type != pb.Message_UNFOLLOW {
		return nil, errors.New("data type is not unfollow")
	}
	good, err := pubkey.Verify(sd.SerializedData, sd.Signature)
	if err != nil || !good {
		return nil, errors.New("bad signature")
	}
	//err = s.datastore.Followers().Delete(id.Pretty())
	//if err != nil {
	//	return nil, err
	//}
	//n := notifications.UnfollowNotification{notifications.NewID(), "unfollow", id.Pretty()}
	//s.broadcast <- n
	log.Debugf("received UNFOLLOW message from %s", id.Pretty())
	return nil, nil
}

func (s *TextileService) handleOfflineAck(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	_, err := peer.IDB58Decode(string(pmes.Payload.Value))
	if err != nil {
		return nil, err
	}
	//pointer, err := s.datastore.Pointers().Get(pid)
	//if err != nil {
	//	return nil, err
	//}
	//if pointer.CancelID == nil || pointer.CancelID.Pretty() != p.Pretty() {
	//	return nil, errors.New("peer is not authorized to delete pointer")
	//}
	//err = s.datastore.Pointers().Delete(pid)
	//if err != nil {
	//	return nil, err
	//}
	log.Debugf("received OFFLINE_ACK message from %s", pid.Pretty())
	return nil, nil
}

func (s *TextileService) handleOfflineRelay(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	// This acts very similarly to attemptDecrypt&handleMessage in the Offline Message Retreiver
	// However it does not send an ACK, or worry about message ordering

	// Decrypt and unmarshal plaintext
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	var plaintext []byte // FIXME
	//plaintext, err := net.Decrypt(s.node.PrivateKey, pmes.Payload.Value)
	//if err != nil {
	//	return nil, err
	//}

	// Unmarshal plaintext
	env := pb.Envelope{}
	err := proto.Unmarshal(plaintext, &env)
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
	handler := s.HandlerForMsgType(env.Message.MessageType)
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

func (s *TextileService) handleChat(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	// Unmarshall
	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	chat := new(pb.Chat)
	err := ptypes.UnmarshalAny(pmes.Payload, chat)
	if err != nil {
		return nil, err
	}

	if chat.Flag == pb.Chat_TYPING {
		//n := notifications.ChatTyping{
		//	PeerId:  p.Pretty(),
		//	Subject: chat.Subject,
		//}
		//s.broadcast <- notifications.Serialize(n)
		return nil, nil
	}
	if chat.Flag == pb.Chat_READ {
		//n := notifications.ChatRead{
		//	PeerId:    p.Pretty(),
		//	Subject:   chat.Subject,
		//	MessageId: chat.MessageId,
		//}
		//s.broadcast <- n
		//_, _, err = s.datastore.Chat().MarkAsRead(p.Pretty(), chat.Subject, true, chat.MessageId)
		//if err != nil {
		//	return nil, err
		//}
		return nil, nil
	}

	// Validate
	if len(chat.Subject) > ChatSubjectMaxCharacters {
		return nil, errors.New("chat subject over max characters")
	}
	if len(chat.Message) > ChatMessageMaxCharacters {
		return nil, errors.New("chat message over max characters")
	}

	// Use correct timestamp
	//offline, _ := options.(bool)
	//var t time.Time
	//if !offline {
	//	t = time.Now()
	//} else {
	//	if chat.Timestamp == nil {
	//		return nil, errors.New("invalid timestamp")
	//	}
	//	t, err = ptypes.Timestamp(chat.Timestamp)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	// Put to database
	//err = s.datastore.Chat().Put(chat.MessageId, p.Pretty(), chat.Subject, chat.Message, t, false, false)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if chat.Subject != "" {
	//	go func() {
	//		s.datastore.Purchases().MarkAsUnread(chat.Subject)
	//		s.datastore.Sales().MarkAsUnread(chat.Subject)
	//		s.datastore.Cases().MarkAsUnread(chat.Subject)
	//	}()
	//}
	//
	//// Push to websocket
	//n := notifications.ChatMessage{
	//	MessageId: chat.MessageId,
	//	PeerId:    p.Pretty(),
	//	Subject:   chat.Subject,
	//	Message:   chat.Message,
	//	Timestamp: t,
	//}
	//s.broadcast <- n
	log.Debugf("received CHAT message from %s", pid.Pretty())
	return nil, nil
}

func (s *TextileService) handleIPFSBlock(pid peer.ID, pmes *pb.Message, options interface{}) (*pb.Message, error) {
	//// If we aren't accepting store requests then ban this peer
	//if !s.node.AcceptStoreRequests {
	//	s.node.BanManager.AddBlockedId(pid)
	//	return nil, nil
	//}

	if pmes.Payload == nil {
		return nil, errors.New("payload is nil")
	}
	b := new(pb.IPFSBlock)
	err := ptypes.UnmarshalAny(pmes.Payload, b)
	if err != nil {
		return nil, err
	}
	id, err := cid.Decode(b.Cid)
	if err != nil {
		return nil, err
	}
	block, err := blocks.NewBlockWithCid(b.RawData, id)
	if err != nil {
		return nil, err
	}
	if err := s.node.Blocks.AddBlock(block); err != nil {
		return nil, err
	}
	log.Debugf("Received IPFS_BLOCK message from %s", pid.Pretty())
	return nil, nil
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

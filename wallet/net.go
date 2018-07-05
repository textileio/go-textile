package wallet

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/pb"
	"golang.org/x/net/context"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"time"
)

//var OfflineMessageWaitGroup sync.WaitGroup

func (w *Wallet) sendMessage(message *pb.Message, peerId string) error {
	p, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = w.service.SendMessage(ctx, p, message)
	if err != nil {
		//if err := w.SendOfflineMessage(p, k, &message); err != nil {
		return err
		//}
	}
	return nil
}

// Supply of a public key is optional, if nil is instead provided n.EncryptMessage does a lookup
//func (w *Wallet) SendOfflineMessage(p peer.ID, k *libp2pc.PubKey, m *pb.Message) error {
//	pubKeyBytes, err := w.ipfs.PrivateKey.GetPublic().Bytes()
//	if err != nil {
//		return err
//	}
//	ser, err := proto.Marshal(m)
//	if err != nil {
//		return err
//	}
//	sig, err := w.ipfs.PrivateKey.Sign(ser)
//	if err != nil {
//		return err
//	}
//	env := pb.Envelope{Message: m, Pubkey: pubKeyBytes, Signature: sig}
//	messageBytes, merr := proto.Marshal(&env)
//	if merr != nil {
//		return merr
//	}
//	ciphertext, cerr := w.EncryptMessage(p, k, messageBytes)
//	if cerr != nil {
//		return cerr
//	}
//	addr, aerr := n.MessageStorage.Store(p, ciphertext)
//	if aerr != nil {
//		return aerr
//	}
//	mh, mherr := multihash.FromB58String(p.Pretty())
//	if mherr != nil {
//		return mherr
//	}
//	/* TODO: We are just using a default prefix length for now. Eventually we will want to customize this,
//	   but we will need some way to get the recipient's desired prefix length. Likely will be in profile. */
//	pointer, err := ipfs.NewPointer(mh, DefaultPointerPrefixLength, addr, ciphertext)
//	if err != nil {
//		return err
//	}
//	if m.MessageType != pb.Message_OFFLINE_ACK {
//		pointer.Purpose = ipfs.MESSAGE
//		pointer.CancelID = &p
//		err = n.Datastore.Pointers().Put(pointer)
//		if err != nil {
//			return err
//		}
//	}
//	log.Debugf("Sending offline message to: %s, Message Type: %s, PointerID: %s, Location: %s", p.Pretty(), m.MessageType.String(), pointer.Cid.String(), pointer.Value.Addrs[0].String())
//	OfflineMessageWaitGroup.Add(1)
//	go func() {
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//		err := ipfs.PublishPointer(w.ipfs, ctx, pointer)
//		if err != nil {
//			log.Error(err)
//		}
//
//		// Push provider to our push nodes for redundancy
//		for _, p := range n.PushNodes {
//			ctx, cancel := context.WithCancel(context.Background())
//			defer cancel()
//			err := ipfs.PutPointerToPeer(w.ipfs, ctx, p, pointer)
//			if err != nil {
//				log.Error(err)
//			}
//		}
//
//		OfflineMessageWaitGroup.Done()
//	}()
//	return nil
//}

func (w *Wallet) SendOfflineAck(peerId string, pointerID peer.ID) error {
	a := &any.Any{Value: []byte(pointerID.Pretty())}
	m := &pb.Message{
		MessageType: pb.Message_OFFLINE_ACK,
		Payload:     a,
	}
	return w.sendMessage(m, peerId)
}

func (w *Wallet) GetPeerStatus(peerId string) (string, error) {
	p, err := peer.IDB58Decode(peerId)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := pb.Message{MessageType: pb.Message_PING}
	_, err = w.service.SendRequest(ctx, p, &m)
	if err != nil {
		return "offline", nil
	}
	return "online", nil
}

func (w *Wallet) Follow(peerId string) error {
	m := &pb.Message{MessageType: pb.Message_FOLLOW}

	pubkey := w.ipfs.PrivateKey.GetPublic()
	pubkeyBytes, err := pubkey.Bytes()
	if err != nil {
		return err
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return err
	}
	data := &pb.SignedData_Command{
		PeerID:    peerId,
		Type:      pb.Message_FOLLOW,
		Timestamp: ts,
	}
	ser, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	sig, err := w.ipfs.PrivateKey.Sign(ser)
	if err != nil {
		return err
	}
	sd := &pb.SignedData{
		SerializedData: ser,
		SenderPubkey:   pubkeyBytes,
		Signature:      sig,
	}
	payload, err := ptypes.MarshalAny(sd)
	if err != nil {
		return err
	}
	m.Payload = payload

	err = w.sendMessage(m, peerId)
	if err != nil {
		return err
	}
	//err = n.Datastore.Following().Put(peerId)
	//if err != nil {
	//	return err
	//}
	//err = n.UpdateFollow()
	//if err != nil {
	//	return err
	//}
	return nil
}

func (w *Wallet) Unfollow(peerId string) error {
	m := &pb.Message{MessageType: pb.Message_UNFOLLOW}

	pubkey := w.ipfs.PrivateKey.GetPublic()
	pubkeyBytes, err := pubkey.Bytes()
	if err != nil {
		return err
	}
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return err
	}
	data := &pb.SignedData_Command{
		PeerID:    peerId,
		Type:      pb.Message_UNFOLLOW,
		Timestamp: ts,
	}
	ser, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	sig, err := w.ipfs.PrivateKey.Sign(ser)
	if err != nil {
		return err
	}
	sd := &pb.SignedData{
		SerializedData: ser,
		SenderPubkey:   pubkeyBytes,
		Signature:      sig,
	}
	payload, err := ptypes.MarshalAny(sd)
	if err != nil {
		return err
	}
	m.Payload = payload

	err = w.sendMessage(m, peerId)
	if err != nil {
		return err
	}
	//err = n.Datastore.Following().Delete(peerId)
	//if err != nil {
	//	return err
	//}
	//err = n.UpdateFollow()
	//if err != nil {
	//	return err
	//}
	return nil
}

func (w *Wallet) SendError(peerId string, k *libp2pc.PubKey, errorMessage pb.Message) error {
	return w.sendMessage(&errorMessage, peerId)
}

func (w *Wallet) SendChat(peerId string, chatMessage *pb.Chat) error {
	a, err := ptypes.MarshalAny(chatMessage)
	if err != nil {
		return err
	}
	m := pb.Message{
		MessageType: pb.Message_CHAT,
		Payload:     a,
	}

	p, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = w.service.SendMessage(ctx, p, &m)
	if err != nil && chatMessage.Flag != pb.Chat_TYPING {
		//if err := w.SendOfflineMessage(p, nil, &m); err != nil {
		return err
		//}
	}
	return nil
}

func (w *Wallet) SendIPFSBlock(peerId string, id cid.Cid) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	block, err := w.ipfs.Blocks.GetBlock(ctx, &id)
	if err != nil {
		return err
	}

	b := &pb.IPFSBlock{
		Cid:     block.Cid().String(),
		RawData: block.RawData(),
	}
	a, err := ptypes.MarshalAny(b)
	if err != nil {
		return err
	}
	m := pb.Message{
		MessageType: pb.Message_IPFS_BLOCK,
		Payload:     a,
	}

	p, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	return w.service.SendMessage(context.Background(), p, &m)
}

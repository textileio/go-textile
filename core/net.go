package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/net"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	"gx/ipfs/QmZ383TySJVeZWzGnWui6pRcKyYZk9VkKTuW7tmKRWk5au/go-libp2p-routing"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	ds "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore"
	"sync"
	"time"
)

const DefaultPointerPrefixLength = 14

var offlineMessageWaitGroup sync.WaitGroup

func (t *Textile) NewEnvelope(message *pb.Message) (*pb.Envelope, error) {
	serialized, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	authorSig, err := t.ipfs.PrivateKey.Sign(serialized)
	if err != nil {
		return nil, err
	}
	authorPk, err := t.ipfs.PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	return &pb.Envelope{Message: message, Pk: authorPk, Sig: authorSig}, nil
}

func (t *Textile) VerifyEnvelope(env *pb.Envelope) error {
	messageb, err := proto.Marshal(env.Message)
	if err != nil {
		return err
	}
	authorPk, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return err
	}
	return crypto.Verify(authorPk, messageb, env.Sig)
}

func (t *Textile) GetPeerStatus(peerId string) (string, error) {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	message := pb.Message{Type: pb.Message_PING}
	env, err := t.NewEnvelope(&message)
	if err != nil {
		return "", err
	}
	_, err = t.service.SendRequest(ctx, pid, env)
	if err != nil {
		return "offline", nil
	}
	return "online", nil
}

func (t *Textile) SendMessage(env *pb.Envelope, peerId string, hash *string) error {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	var success bool
	go func() {
		err = t.service.SendMessage(ctx, pid, env)
		if err == nil {
			success = true
		}
	}()
	go func() {
		time.Sleep(time.Second * 3)
		if !success {
			log.Debug("TODO send to inbox!")
			//if err := t.sendOfflineMessage(env, pid, hash); err != nil {
			//	log.Debugf("send offline message failed: %s", err)
			//}
		}
		cancel()
	}()
	return nil
}

func (t *Textile) sendOfflineMessage(env *pb.Envelope, pid peer.ID, hash *string) error {
	defer func() {
		if recover() != nil {
			log.Error("recovered from sendOfflineMessage")
		}
	}()
	serialized, err := proto.Marshal(env)
	if err != nil {
		return err
	}

	// if we've already computed the hash, taking that to mean it's already been stored
	var addr ma.Multiaddr
	if hash != nil {
		addr, err = ipfs.MultiaddrFromId(*hash)
	} else {
		addr, err = t.messageStorage.Store(serialized)
	}
	if err != nil {
		return err
	}

	// create a pointer for this peer
	mh, err := multihash.FromB58String(pid.Pretty())
	if err != nil {
		return err
	}
	entropy := ksuid.New().Bytes()
	pointer, err := repo.NewPointer(mh, DefaultPointerPrefixLength, addr, entropy)
	if err != nil {
		return err
	}
	if env.Message.Type != pb.Message_OFFLINE_ACK {
		pointer.Purpose = repo.MESSAGE
		pointer.CancelId = &pid
		err = t.datastore.Pointers().Put(pointer)
		if err != nil {
			return err
		}
	}

	log.Debugf("sending offline message to: %s, type: %s, pointer: %s, addr: %s",
		pid.Pretty(), env.Message.Type.String(), pointer.Cid.String(), pointer.Value.Addrs[0].String())

	offlineMessageWaitGroup.Add(1)
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := repo.PublishPointer(t.ipfs, ctx, pointer)
		if err != nil {
			log.Error(err.Error())
		}
		offlineMessageWaitGroup.Done()
	}()
	return nil
}

func (t *Textile) sendOfflineAck(peerId string, pointerID peer.ID) error {
	payload := &any.Any{Value: []byte(pointerID.Pretty())}
	message := &pb.Message{
		Type:    pb.Message_OFFLINE_ACK,
		Payload: payload,
	}
	env, err := t.NewEnvelope(message)
	if err != nil {
		return err
	}
	return t.SendMessage(env, peerId, nil)
}

func (t *Textile) sendError(peerId string, k *libp2pc.PubKey, errorMessage pb.Envelope) error {
	return t.SendMessage(&errorMessage, peerId, nil)
}

func (t *Textile) sendChat(peerId string, chatMessage *pb.Chat) error {
	oayload, err := ptypes.MarshalAny(chatMessage)
	if err != nil {
		return err
	}
	message := &pb.Message{
		Type:    pb.Message_CHAT,
		Payload: oayload,
	}
	env, err := t.NewEnvelope(message)
	if err != nil {
		return err
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return t.service.SendMessage(ctx, pid, env)
	//if err != nil && chatMessage.Flag != pb.Chat_TYPING {
	//	if err := t.sendOfflineMessage(env, pid, nil); err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (t *Textile) sendStore(peerId string, ids []cid.Cid) error {
	var cids []string
	for _, d := range ids {
		cids = append(cids, d.String())
	}
	cList := new(pb.CidList)
	cList.Cids = cids

	payload, err := ptypes.MarshalAny(cList)
	if err != nil {
		return err
	}

	message := &pb.Message{
		Type:    pb.Message_STORE,
		Payload: payload,
	}
	env, err := t.NewEnvelope(message)
	if err != nil {
		return err
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	pmes, err := t.service.SendRequest(context.Background(), pid, env)
	if err != nil {
		return err
	}
	// TODO: need to disconnect here?
	// defer t.service.DisconnectFromPeer(pid)
	if pmes.Message.Payload == nil {
		return errors.New("peer responded with nil payload")
	}
	if pmes.Message.Type == pb.Message_ERROR {
		err = fmt.Errorf("error response from %s: %s", peerId, string(pmes.Message.Payload.Value))
		log.Errorf(err.Error())
		return err
	}

	resp := new(pb.CidList)
	err = ptypes.UnmarshalAny(pmes.Message.Payload, resp)
	if err != nil {
		return err
	}
	if len(resp.Cids) == 0 {
		log.Debugf("peer %s requested no blocks", peerId)
		return nil
	}
	log.Debugf("sending %d blocks to %s", len(resp.Cids), peerId)
	for _, id := range resp.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			continue
		}
		t.sendBlock(peerId, *decoded)
	}
	return nil
}

func (t *Textile) sendBlock(peerId string, id cid.Cid) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	block, err := t.ipfs.Blocks.GetBlock(ctx, &id)
	if err != nil {
		return err
	}

	pbblock := &pb.Block{
		Cid:     block.Cid().String(),
		RawData: block.RawData(),
	}
	payload, err := ptypes.MarshalAny(pbblock)
	if err != nil {
		return err
	}
	message := &pb.Message{
		Type:    pb.Message_BLOCK,
		Payload: payload,
	}
	env, err := t.NewEnvelope(message)
	if err != nil {
		return err
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	return t.service.SendMessage(context.Background(), pid, env)
}

func (t *Textile) encryptMessage(pid peer.ID, message []byte) (ct []byte, rerr error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pubKey libp2pc.PubKey
	keyval := t.datastore.Peers().GetById(pid.Pretty())
	if keyval != nil {
		var err error
		pubKey, err = libp2pc.UnmarshalPublicKey(keyval.PubKey)
		if err != nil {
			log.Errorf("failed to parse peer public key for %s", pid.Pretty())
			return nil, err
		}
	} else {
		keyval, err := t.ipfs.Repo.Datastore().Get(ds.NewKey(net.KeyCachePrefix + pid.String()))
		if err != nil {
			pubKey, err = routing.GetPublicKey(t.ipfs.Routing, ctx, pid)
			if err != nil {
				log.Errorf("failed to find public key for %s", pid.Pretty())
				return nil, err
			}
		} else {
			pubKey, err = libp2pc.UnmarshalPublicKey(keyval.([]byte))
			if err != nil {
				log.Errorf("failed to find public key for %s", pid.Pretty())
				return nil, err
			}
		}
	}

	if pid.MatchesPublicKey(pubKey) {
		ciphertext, err := crypto.Encrypt(pubKey, message)
		if err != nil {
			return nil, err
		}
		return ciphertext, nil
	} else {
		err := errors.New(fmt.Sprintf("peer public key and id do not match for peer: %s", pid.Pretty()))
		log.Error(err.Error())
		return nil, err
	}
}

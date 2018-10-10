package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"time"
)

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
	_, err = t.threadsService.SendRequest(ctx, pid, env)
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
		err = t.threadsService.SendMessage(ctx, pid, env)
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
	return t.threadsService.SendMessage(ctx, pid, env)
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
	pmes, err := t.threadsService.SendRequest(context.Background(), pid, env)
	if err != nil {
		return err
	}
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
	return t.threadsService.SendMessage(context.Background(), pid, env)
}

func (t *Textile) sendError(peerId string, k *libp2pc.PubKey, errorMessage pb.Envelope) error {
	return t.SendMessage(&errorMessage, peerId, nil)
}

package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmTiWLZ6Fo5j4KcTVutZJ5KWRRJrbxzmxA4td8NfEdrPh7/go-libp2p-routing"
	ds "gx/ipfs/QmXRKBQA4wXP7xWbFiZsR1GP4HV6wMDQ1aWFxZZ4uBcPX9/go-datastore"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"sync"
	"time"
)

const (
	DefaultPointerPrefixLength = 14
)

var offlineMessageWaitGroup sync.WaitGroup

func (w *Wallet) GetPeerStatus(peerId string) (string, error) {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	message := pb.Message{MessageType: pb.Message_PING}
	_, err = w.service.SendRequest(ctx, pid, &message)
	if err != nil {
		return "offline", nil
	}
	return "online", nil
}

func (w *Wallet) SendMessage(message *pb.Message, peerId string) error {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = w.service.SendMessage(ctx, pid, message)
	if err != nil {
		if err := w.sendOfflineMessage(message, pid); err != nil {
			return err
		}
	}
	return nil
}

func (w *Wallet) sendOfflineMessage(message *pb.Message, pid peer.ID) error {
	pubKeyBytes, err := w.ipfs.PrivateKey.GetPublic().Bytes()
	if err != nil {
		return err
	}
	ser, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	sig, err := w.ipfs.PrivateKey.Sign(ser)
	if err != nil {
		return err
	}
	env := pb.Envelope{Message: message, Pubkey: pubKeyBytes, Signature: sig}
	messageBytes, merr := proto.Marshal(&env)
	if merr != nil {
		return merr
	}
	ciphertext, cerr := w.encryptMessage(pid, messageBytes)
	if cerr != nil {
		return cerr
	}
	addr, aerr := w.messageStorage.Store(pid, ciphertext)
	if aerr != nil {
		return aerr
	}
	mh, mherr := multihash.FromB58String(pid.Pretty())
	if mherr != nil {
		return mherr
	}
	pointer, err := repo.NewPointer(mh, DefaultPointerPrefixLength, addr, ciphertext)
	if err != nil {
		return err
	}
	if message.MessageType != pb.Message_OFFLINE_ACK {
		pointer.Purpose = repo.MESSAGE
		pointer.CancelId = &pid
		err = w.datastore.Pointers().Put(pointer)
		if err != nil {
			return err
		}
	}

	log.Debugf("sending offline message to: %s, type: %s, pointer: %s, location: %s",
		pid.Pretty(), message.MessageType.String(), pointer.Cid.String(), pointer.Value.Addrs[0].String())

	offlineMessageWaitGroup.Add(1)
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := repo.PublishPointer(w.ipfs, ctx, pointer)
		if err != nil {
			log.Error(err.Error())
		}
		offlineMessageWaitGroup.Done()
	}()
	return nil
}

func (w *Wallet) sendOfflineAck(peerId string, pointerID peer.ID) error {
	payload := &any.Any{Value: []byte(pointerID.Pretty())}
	message := &pb.Message{
		MessageType: pb.Message_OFFLINE_ACK,
		Payload:     payload,
	}
	return w.SendMessage(message, peerId)
}

func (w *Wallet) sendError(peerId string, k *libp2pc.PubKey, errorMessage pb.Message) error {
	return w.SendMessage(&errorMessage, peerId)
}

func (w *Wallet) sendChat(peerId string, chatMessage *pb.Chat) error {
	oayload, err := ptypes.MarshalAny(chatMessage)
	if err != nil {
		return err
	}
	message := pb.Message{
		MessageType: pb.Message_CHAT,
		Payload:     oayload,
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = w.service.SendMessage(ctx, pid, &message)
	if err != nil && chatMessage.Flag != pb.Chat_TYPING {
		if err := w.sendOfflineMessage(&message, pid); err != nil {
			return err
		}
	}
	return nil
}

func (w *Wallet) sendStore(peerId string, ids []cid.Cid) error {
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

	message := pb.Message{
		MessageType: pb.Message_STORE,
		Payload:     payload,
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	pmes, err := w.service.SendRequest(context.Background(), pid, &message)
	if err != nil {
		return err
	}
	defer w.service.DisconnectFromPeer(pid)
	if pmes.Payload == nil {
		return errors.New("peer responded with nil payload")
	}
	if pmes.MessageType == pb.Message_ERROR {
		log.Errorf("error response from %s: %s", peerId, string(pmes.Payload.Value))
		return errors.New("peer responded with error message")
	}

	resp := new(pb.CidList)
	err = ptypes.UnmarshalAny(pmes.Payload, resp)
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
		w.sendBlock(peerId, *decoded)
	}
	return nil
}

func (w *Wallet) sendBlock(peerId string, id cid.Cid) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	block, err := w.ipfs.Blocks.GetBlock(ctx, &id)
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
	message := pb.Message{
		MessageType: pb.Message_BLOCK,
		Payload:     payload,
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	return w.service.SendMessage(context.Background(), pid, &message)
}

func (w *Wallet) encryptMessage(peerID peer.ID, message []byte) (ct []byte, rerr error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pubKey libp2pc.PubKey
	keyval, err := w.ipfs.Repo.Datastore().Get(ds.NewKey(net.KeyCachePrefix + peerID.String()))
	if err != nil {
		pubKey, err = routing.GetPublicKey(w.ipfs.Routing, ctx, []byte(peerID))
		if err != nil {
			log.Errorf("failed to find public key for %s", peerID.Pretty())
			return nil, err
		}
	} else {
		pubKey, err = libp2pc.UnmarshalPublicKey(keyval.([]byte))
		if err != nil {
			log.Errorf("failed to find public key for %s", peerID.Pretty())
			return nil, err
		}
	}

	if peerID.MatchesPublicKey(pubKey) {
		ciphertext, err := crypto.Encrypt(pubKey, message)
		if err != nil {
			return nil, err
		}
		return ciphertext, nil
	} else {
		err = errors.New(fmt.Sprintf("peer public key and id do not match for peer: %s", peerID.Pretty()))
		log.Error(err.Error())
		return nil, err
	}
}

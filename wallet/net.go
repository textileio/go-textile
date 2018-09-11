package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"gx/ipfs/QmTiWLZ6Fo5j4KcTVutZJ5KWRRJrbxzmxA4td8NfEdrPh7/go-libp2p-routing"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
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

func (w *Wallet) NewEnvelope(message *pb.Message) (*pb.Envelope, error) {
	serialized, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	authorSig, err := w.ipfs.PrivateKey.Sign(serialized)
	if err != nil {
		return nil, err
	}
	authorPk, err := w.ipfs.PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	return &pb.Envelope{Message: message, Pk: authorPk, Sig: authorSig}, nil
}

func (w *Wallet) VerifyEnvelope(env *pb.Envelope) error {
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

func (w *Wallet) GetPeerStatus(peerId string) (string, error) {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	message := pb.Message{Type: pb.Message_PING}
	env, err := w.NewEnvelope(&message)
	if err != nil {
		return "", err
	}
	_, err = w.service.SendRequest(ctx, pid, env)
	if err != nil {
		return "offline", nil
	}
	return "online", nil
}

func (w *Wallet) SendMessage(env *pb.Envelope, peerId string, hash *string) error {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	var success bool
	go func() {
		err = w.service.SendMessage(ctx, pid, env)
		if err == nil {
			success = true
		}
	}()
	go func() {
		time.Sleep(time.Second * 3)
		if !success {
			if err := w.sendOfflineMessage(env, pid, hash); err != nil {
				log.Debugf("send offline message failed: %s", err)
			}
		}
		cancel()
	}()
	return nil
}

func (w *Wallet) sendOfflineMessage(env *pb.Envelope, pid peer.ID, hash *string) error {
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
		addr, err = util.MultiaddrFromId(*hash)
	} else {
		addr, err = w.messageStorage.Store(serialized)
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
		err = w.datastore.Pointers().Put(pointer)
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
		Type:    pb.Message_OFFLINE_ACK,
		Payload: payload,
	}
	env, err := w.NewEnvelope(message)
	if err != nil {
		return err
	}
	return w.SendMessage(env, peerId, nil)
}

func (w *Wallet) sendError(peerId string, k *libp2pc.PubKey, errorMessage pb.Envelope) error {
	return w.SendMessage(&errorMessage, peerId, nil)
}

func (w *Wallet) sendChat(peerId string, chatMessage *pb.Chat) error {
	oayload, err := ptypes.MarshalAny(chatMessage)
	if err != nil {
		return err
	}
	message := &pb.Message{
		Type:    pb.Message_CHAT,
		Payload: oayload,
	}
	env, err := w.NewEnvelope(message)
	if err != nil {
		return err
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = w.service.SendMessage(ctx, pid, env)
	if err != nil && chatMessage.Flag != pb.Chat_TYPING {
		if err := w.sendOfflineMessage(env, pid, nil); err != nil {
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

	message := &pb.Message{
		Type:    pb.Message_STORE,
		Payload: payload,
	}
	env, err := w.NewEnvelope(message)
	if err != nil {
		return err
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	pmes, err := w.service.SendRequest(context.Background(), pid, env)
	if err != nil {
		return err
	}
	// TODO: need to disconnect here?
	// defer w.service.DisconnectFromPeer(pid)
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
	message := &pb.Message{
		Type:    pb.Message_BLOCK,
		Payload: payload,
	}
	env, err := w.NewEnvelope(message)
	if err != nil {
		return err
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	return w.service.SendMessage(context.Background(), pid, env)
}

func (w *Wallet) encryptMessage(pid peer.ID, message []byte) (ct []byte, rerr error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var pubKey libp2pc.PubKey
	keyval := w.datastore.Peers().GetById(pid.Pretty())
	if keyval != nil {
		var err error
		pubKey, err = libp2pc.UnmarshalPublicKey(keyval.PubKey)
		if err != nil {
			log.Errorf("failed to parse peer public key for %s", pid.Pretty())
			return nil, err
		}
	} else {
		keyval, err := w.ipfs.Repo.Datastore().Get(ds.NewKey(net.KeyCachePrefix + pid.String()))
		if err != nil {
			pubKey, err = routing.GetPublicKey(w.ipfs.Routing, ctx, []byte(pid))
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

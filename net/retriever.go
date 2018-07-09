package net

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net/common"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/util"
	routing "gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	ds "gx/ipfs/QmXRKBQA4wXP7xWbFiZsR1GP4HV6wMDQ1aWFxZZ4uBcPX9/go-datastore"
	ps "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"sync"
	"time"
)

const KeyCachePrefix = "PUBKEYCACHE_"

type MRConfig struct {
	Db        repo.Datastore
	Ipfs      *core.IpfsNode
	Service   NetworkService
	PrefixLen int
	SendAck   func(peerId string, pointerID peer.ID) error
	SendError func(peerId string, k *libp2p.PubKey, errorMessage pb.Message) error
}

type MessageRetriever struct {
	db        repo.Datastore
	ipfs      *core.IpfsNode
	service   NetworkService
	prefixLen int
	sendAck   func(peerId string, pointerID peer.ID) error
	sendError func(peerId string, k *libp2p.PubKey, errorMessage pb.Message) error
	queueLock *sync.Mutex
	DoneChan  chan struct{}
	inFlight  chan struct{}
	*sync.WaitGroup
}

type offlineMessage struct {
	addr string
	env  pb.Envelope
}

func NewMessageRetriever(cfg MRConfig) *MessageRetriever {
	mr := MessageRetriever{
		db:        cfg.Db,
		ipfs:      cfg.Ipfs,
		service:   cfg.Service,
		prefixLen: cfg.PrefixLen,
		sendAck:   cfg.SendAck,
		sendError: cfg.SendError,
		queueLock: new(sync.Mutex),
		DoneChan:  make(chan struct{}),
		inFlight:  make(chan struct{}, 5),
		WaitGroup: new(sync.WaitGroup),
	}
	mr.Add(1)
	return &mr
}

func (m *MessageRetriever) Run() {
	dht := time.NewTicker(time.Hour)
	defer dht.Stop()
	go m.fetchPointers(true)
	for {
		select {
		case <-dht.C:
			m.Add(1)
			go m.fetchPointers(true)
		}
	}
}

func (m *MessageRetriever) fetchPointers(useDHT bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := new(sync.WaitGroup)
	downloaded := 0
	mh, _ := multihash.FromB58String(m.ipfs.Identity.Pretty())
	peerOut := make(chan ps.PeerInfo)
	go func(c chan ps.PeerInfo) {
		pwg := new(sync.WaitGroup)
		pwg.Add(1)
		go func(c chan ps.PeerInfo) {
			iout := repo.FindPointersAsync(m.ipfs.Routing.(*routing.IpfsDHT), ctx, mh, m.prefixLen)
			for p := range iout {
				c <- p
			}
			pwg.Done()
		}(c)
		pwg.Wait()
		close(c)
	}(peerOut)

	inFlight := make(map[string]bool)
	// iterate over the pointers, adding 1 to the waitgroup for each pointer found
	for p := range peerOut {
		log.Debugf("retriever found peer info: %s", p.Loggable())
		if len(p.Addrs) > 0 && !m.db.OfflineMessages().Has(p.Addrs[0].String()) && !inFlight[p.Addrs[0].String()] {
			inFlight[p.Addrs[0].String()] = true
			log.Debugf("found pointer with location %s", p.Addrs[0].String())

			// check protocol
			if len(p.Addrs[0].Protocols()) == 1 && p.Addrs[0].Protocols()[0].Code == ma.P_IPFS {
				wg.Add(1)
				downloaded++
				go m.fetch(p.ID, p.Addrs[0], wg)
			}
		}
	}

	// wait for each goroutine to finish then process any remaining messages that needed to be processed last
	wg.Wait()

	m.processQueuedMessages()
	m.Done()
}

// fetchIPFS will attempt to download an encrypted message using IPFS. If the message downloads successfully, we save the
// address to the database to prevent us from wasting bandwidth downloading it again.
func (m *MessageRetriever) fetch(pid peer.ID, addr ma.Multiaddr, wg *sync.WaitGroup) {
	m.inFlight <- struct{}{}
	defer func() {
		wg.Done()
		<-m.inFlight
	}()

	c := make(chan struct{})
	var ciphertext []byte
	var err error

	go func() {
		ciphertext, err = util.GetDataAtPath(m.ipfs, addr.String()+"/msg")
		c <- struct{}{}
	}()

	select {
	case <-c:
		if err != nil {
			log.Errorf("error retrieving offline message from %s, %s", addr.String(), err.Error())
			return
		}
		log.Debugf("successfully downloaded offline message from %s", addr.String())
		m.db.OfflineMessages().Put(addr.String())
		m.attemptDecrypt(ciphertext, pid, addr)
	case <-m.DoneChan:
		return
	}
}

// attemptDecrypt will try to decrypt the message using our identity private key. If it decrypts it will be passed to
// a handler for processing. Not all messages will decrypt. Given the natural of the prefix addressing, we may download
// some messages intended for others. If we can't decrypt it, we can just discard it.
func (m *MessageRetriever) attemptDecrypt(ciphertext []byte, pid peer.ID, addr ma.Multiaddr) {
	// Decrypt and unmarshal plaintext
	plaintext, err := crypto.Decrypt(m.ipfs.PrivateKey, ciphertext)
	if err != nil {
		log.Warning("unable to decrypt offline message from %s: %s", addr.String(), err.Error())
		return
	}

	// unmarshal plaintext
	env := pb.Envelope{}
	err = proto.Unmarshal(plaintext, &env)
	if err != nil {
		log.Warning("unable to decrypt offline message from %s: %s", addr.String(), err.Error())
		return
	}

	// validate the signature
	ser, err := proto.Marshal(env.Message)
	if err != nil {
		log.Warning("unable to decrypt offline message from %s: %s", addr.String(), err.Error())
		return
	}
	pubkey, err := libp2p.UnmarshalPublicKey(env.Pubkey)
	if err != nil {
		log.Warning("unable to decrypt offline message from %s: %s", addr.String(), err.Error())
		return
	}

	valid, err := pubkey.Verify(ser, env.Signature)
	if err != nil || !valid {
		log.Warning("unable to decrypt offline message from %s: %s", addr.String(), err.Error())
		return
	}

	id, err := peer.IDFromPublicKey(pubkey)
	if err != nil {
		log.Warning("unable to decrypt offline message from %s: %s", addr.String(), err.Error())
		return
	}

	m.ipfs.Peerstore.AddPubKey(id, pubkey)
	m.ipfs.Repo.Datastore().Put(ds.NewKey(KeyCachePrefix+id.String()), env.Pubkey)

	// respond with an ACK
	if env.Message.MessageType != pb.Message_OFFLINE_ACK {
		m.sendAck(id.Pretty(), pid)
	}

	// handle
	m.handleMessage(env, addr.String(), nil)
}

// handleMessage loads the hander for this message type and attempts to process the message
func (m *MessageRetriever) handleMessage(env pb.Envelope, addr string, id *peer.ID) error {
	if id == nil {
		// get the peer ID from the public key
		pubkey, err := libp2p.UnmarshalPublicKey(env.Pubkey)
		if err != nil {
			log.Errorf("error processing message %s. type %s: %s", addr, env.Message.MessageType, err.Error())
			return err
		}
		i, err := peer.IDFromPublicKey(pubkey)
		if err != nil {
			log.Errorf("error processing message %s. type %s: %s", addr, env.Message.MessageType, err.Error())
			return err
		}
		id = &i
	}

	// get handler for this message type
	handler := m.service.HandlerForMsgType(env.Message.MessageType)
	if handler == nil {
		err := errors.New(fmt.Sprintf("nil handler for message type %s", env.Message.MessageType))
		log.Error(err.Error())
		return err
	}

	// dispatch handler
	_, err := handler(*id, env.Message, true)
	if err != nil {
		if err == common.OutOfOrderMessage {
			ser, err := proto.Marshal(&env)
			if err == nil {
				err := m.db.OfflineMessages().SetMessage(addr, ser)
				if err != nil {
					log.Errorf("error saving offline message %s to database: %s", addr, err.Error())
				}
			} else {
				log.Errorf("error serializing offline message %s for storage")
			}
		} else {
			log.Errorf("error processing message %s. type %s: %s", addr, env.Message.MessageType, err.Error())
			return err
		}
	}
	return nil
}

var MessageProcessingOrder = []pb.Message_MessageType{
	pb.Message_THREAD_BLOCK,
	pb.Message_CHAT,
	pb.Message_FOLLOW,
	pb.Message_UNFOLLOW,
	pb.Message_MODERATOR_ADD,
	pb.Message_MODERATOR_REMOVE,
	pb.Message_OFFLINE_ACK,
}

// processQueuedMessages loads all the saved messaged from the database for processing. For each message it sorts them into a
// queue based on message type and then processes the queue in order. Any messages that successfully process can then be deleted
// from the databse.
func (m *MessageRetriever) processQueuedMessages() {
	messageQueue := make(map[pb.Message_MessageType][]offlineMessage)
	for _, messageType := range MessageProcessingOrder {
		messageQueue[messageType] = []offlineMessage{}
	}

	// load stored messages from database
	messages, err := m.db.OfflineMessages().GetMessages()
	if err != nil {
		return
	}
	// sort them into the queue by message type
	for url, ser := range messages {
		env := new(pb.Envelope)
		err := proto.Unmarshal(ser, env)
		if err == nil {
			messageQueue[env.Message.MessageType] = append(messageQueue[env.Message.MessageType], offlineMessage{url, *env})
		} else {
			log.Error("error unmarshalling serialized offline message from database")
		}
	}
	var toDelete []string
	// process the queue in order
	for _, messageType := range MessageProcessingOrder {
		queue, ok := messageQueue[messageType]
		if !ok {
			continue
		}
		for _, om := range queue {
			err := m.handleMessage(om.env, om.addr, nil)
			if err == nil {
				toDelete = append(toDelete, om.addr)
			}
		}
	}
	// delete messages that we're successfully processed from the database
	for _, url := range toDelete {
		m.db.OfflineMessages().DeleteMessage(url)
	}
}

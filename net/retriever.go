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
	"github.com/textileio/textile-go/util"
	routing "gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	ds "gx/ipfs/QmXRKBQA4wXP7xWbFiZsR1GP4HV6wMDQ1aWFxZZ4uBcPX9/go-datastore"
	ps "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"sync"
	"time"
)

const KeyCachePrefix = "PUBKEYCACHE_"

const kRetrieveFrequency = time.Minute * 5

type MRConfig struct {
	Datastore repo.Datastore
	Ipfs      *core.IpfsNode
	Service   NetworkService
	PrefixLen int
	SendAck   func(peerId string, pointerID peer.ID) error
	SendError func(peerId string, k *libp2pc.PubKey, errorMessage pb.Envelope) error
}

type MessageRetriever struct {
	datastore repo.Datastore
	ipfs      *core.IpfsNode
	service   NetworkService
	prefixLen int
	sendAck   func(peerId string, pointerID peer.ID) error
	sendError func(peerId string, k *libp2pc.PubKey, errorMessage pb.Envelope) error
	queueLock *sync.Mutex
	DoneChan  chan struct{}
	inFlight  chan struct{}
	*sync.WaitGroup
}

type offlineMessage struct {
	addr string
	env  pb.Envelope
}

func NewMessageRetriever(config MRConfig) *MessageRetriever {
	mr := MessageRetriever{
		datastore: config.Datastore,
		ipfs:      config.Ipfs,
		service:   config.Service,
		prefixLen: config.PrefixLen,
		sendAck:   config.SendAck,
		sendError: config.SendError,
		queueLock: new(sync.Mutex),
		DoneChan:  make(chan struct{}),
		inFlight:  make(chan struct{}, 5),
		WaitGroup: new(sync.WaitGroup),
	}
	mr.Add(1)
	return &mr
}

func (m *MessageRetriever) Run() {
	dht := time.NewTicker(kRetrieveFrequency)
	defer dht.Stop()
	go m.FetchPointers()
	for {
		select {
		case <-dht.C:
			m.Add(1)
			go m.FetchPointers()
		}
	}
}

func (m *MessageRetriever) FetchPointers() {
	log.Debug("fetching pointers...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := new(sync.WaitGroup)
	downloaded := 0
	mh, err := multihash.FromB58String(m.ipfs.Identity.Pretty())
	if err != nil {
		log.Error(err.Error())
	}
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
		if len(p.Addrs) > 0 && !m.datastore.OfflineMessages().Has(p.Addrs[0].String()) && !inFlight[p.Addrs[0].String()] {
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
	addrs := addr.String()
	var payload []byte
	var err error

	go func() {
		payload, err = util.GetDataAtPath(m.ipfs, addrs)
		c <- struct{}{}
	}()

	select {
	case <-c:
		if err != nil {
			log.Errorf("error retrieving offline message from %s, %s", addrs, err)
			return
		}
		log.Debugf("successfully downloaded offline message from %s", addrs)

		// attempt to decrypt and unmarshal
		plaintext, err := crypto.Decrypt(m.ipfs.PrivateKey, payload)
		if err == nil {
			payload = plaintext
		}

		// thread blocks have encrypted contents
		if err := m.unpackMessage(payload, pid, addr); err != nil {
			log.Errorf("unable to unpack offline message from %s: %s", addrs, err)
			return
		}

		// store away
		if err := m.datastore.OfflineMessages().Put(addr.String()); err != nil {
			log.Errorf("put offline message from %s failed: %s", addrs, err)
		}
		return

	case <-m.DoneChan:
		return
	}
}

// unpackMessage unpacks, vefifies, and handles an envelope
func (m *MessageRetriever) unpackMessage(payload []byte, pid peer.ID, addr ma.Multiaddr) error {
	// unmarshal
	env := &pb.Envelope{}
	if err := proto.Unmarshal(payload, env); err != nil {
		return err
	}

	// validate the envelope signature
	ser, err := proto.Marshal(env.Message)
	if err != nil {
		return err
	}
	pk, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return err
	}
	if err := crypto.Verify(pk, ser, env.Sig); err != nil {
		return err
	}
	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return err
	}

	// cache pk, probably should remove this... we already have it in thread peer table
	m.ipfs.Peerstore.AddPubKey(id, pk)
	m.ipfs.Repo.Datastore().Put(ds.NewKey(KeyCachePrefix+id.String()), env.Pk)

	// respond with an ACK
	if env.Message.Type != pb.Message_OFFLINE_ACK {
		m.sendAck(id.Pretty(), pid)
	}

	// handle
	return m.handleMessage(env, addr.String(), &id)
}

// handleMessage loads the hander for this message type and attempts to process the message
func (m *MessageRetriever) handleMessage(env *pb.Envelope, addr string, id *peer.ID) error {
	if id == nil {
		// get the peer ID from the public key
		pubkey, err := libp2pc.UnmarshalPublicKey(env.Pk)
		if err != nil {
			return err
		}
		i, err := peer.IDFromPublicKey(pubkey)
		if err != nil {
			return err
		}
		id = &i
	}

	// get handler for this message type
	handler := m.service.HandlerForMsgType(env.Message.Type)
	if handler == nil {
		return errors.New(fmt.Sprintf("nil handler for message type %s", env.Message.Type))
	}

	// dispatch handler
	_, err := handler(*id, env, true)
	if err != nil {
		if err == common.OutOfOrderMessage {
			ser, err := proto.Marshal(env)
			if err != nil {
				return err
			}
			if err := m.datastore.OfflineMessages().SetMessage(addr, ser); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

var MessageProcessingOrder = []pb.Message_Type{
	pb.Message_THREAD_INVITE,
	pb.Message_THREAD_JOIN,
	pb.Message_THREAD_LEAVE,
	pb.Message_THREAD_DATA,
	pb.Message_THREAD_ANNOTATION,
	pb.Message_THREAD_IGNORE,
	pb.Message_THREAD_MERGE,
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
	messageQueue := make(map[pb.Message_Type][]offlineMessage)
	for _, messageType := range MessageProcessingOrder {
		messageQueue[messageType] = []offlineMessage{}
	}

	// load stored messages from database
	messages, err := m.datastore.OfflineMessages().GetMessages()
	if err != nil {
		return
	}
	// sort them into the queue by message type
	for url, ser := range messages {
		env := new(pb.Envelope)
		err := proto.Unmarshal(ser, env)
		if err == nil {
			messageQueue[env.Message.Type] = append(messageQueue[env.Message.Type], offlineMessage{url, *env})
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
			err := m.handleMessage(&om.env, om.addr, nil)
			if err == nil {
				toDelete = append(toDelete, om.addr)
			}
		}
	}
	// delete messages that we're successfully processed from the database
	for _, url := range toDelete {
		m.datastore.OfflineMessages().DeleteMessage(url)
	}
}

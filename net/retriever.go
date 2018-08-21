package net

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net/common"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	routing "gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	ps "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"sort"
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
	date time.Time
}

type sortedMessages []offlineMessage

func (v sortedMessages) Len() int           { return len(v) }
func (v sortedMessages) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v sortedMessages) Less(i, j int) bool { return v[i].date.Before(v[j].date) }

var messageProcessingOrder = []pb.Message_Type{
	pb.Message_CHAT,
	pb.Message_FOLLOW,
	pb.Message_UNFOLLOW,
	pb.Message_MODERATOR_ADD,
	pb.Message_MODERATOR_REMOVE,
	pb.Message_OFFLINE_ACK,
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

	// find pointers
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

	// iterate over the pointers, adding 1 to the waitgroup for each pointer found
	inFlight := make(map[string]bool)
	for p := range peerOut {
		if len(p.Addrs) > 0 && !m.datastore.OfflineMessages().Has(p.Addrs[0].String()) && !inFlight[p.Addrs[0].String()] {
			inFlight[p.Addrs[0].String()] = true

			// check protocol
			if len(p.Addrs[0].Protocols()) == 1 && p.Addrs[0].Protocols()[0].Code == ma.P_IPFS {
				wg.Add(1)
				downloaded++
				go m.fetch(p.ID, p.Addrs[0], wg)
			}
		}
	}
	wg.Wait()
	// m.processQueuedMessages() // currently not used, message order does not matter
	m.Done()
}

// fetch downloads an message from ipfs
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
			return
		}

		// attempt to decrypt and unmarshal
		plaintext, err := crypto.Decrypt(m.ipfs.PrivateKey, payload)
		if err == nil {
			payload = plaintext
		}

		// thread blocks have encrypted contents
		env, err := m.verifyMessage(payload, pid, addr)
		if err != nil {
			log.Errorf("offline message verification %s failed: %s", addrs, err)
			return
		}

		// respond with an ACK
		if env.Message.Type != pb.Message_OFFLINE_ACK {
			// get sender's id
			id, err := getEnvelopeSenderId(env)
			if err != nil {
				log.Errorf("error getting sender id from env: %s", err)
				return
			}
			m.sendAck(id.Pretty(), pid)
		}

		if err := m.handleMessage(env, addrs); err != nil {
			log.Errorf("error handling offline message: %s", err)
			return
		}

		// store away
		if err := m.datastore.OfflineMessages().Put(addrs); err != nil {
			log.Errorf("put offline message %s failed: %s", addrs, err)
			return
		}
		return

	case <-m.DoneChan:
		return
	}
}

// verifyMessage unpacks, verifies, and handles an envelope
func (m *MessageRetriever) verifyMessage(payload []byte, pid peer.ID, addr ma.Multiaddr) (*pb.Envelope, error) {
	// unmarshal
	env := &pb.Envelope{}
	if err := proto.Unmarshal(payload, env); err != nil {
		return nil, err
	}

	// validate the envelope signature
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

	// handle
	return env, nil
}

// handleMessage loads the hander for this message type and attempts to process the message
func (m *MessageRetriever) handleMessage(env *pb.Envelope, addr string) error {
	// get the peer ID from the public key
	pid, err := getEnvelopeSenderId(env)
	if err != nil {
		return err
	}

	// get handler for this message type
	handler := m.service.HandlerForMsgType(env.Message.Type)
	if handler == nil {
		return errors.New(fmt.Sprintf("nil handler for message type %s", env.Message.Type))
	}

	// dispatch handler
	if _, err := handler(pid, env, true); err != nil {
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

// processQueuedMessages loads all the saved messaged from the database for processing
// - any messages that successfully process can then be deleted from the databse
func (m *MessageRetriever) processQueuedMessages() {
	var threadMessages []offlineMessage
	messageQueue := make(map[pb.Message_Type][]offlineMessage)
	for _, messageType := range messageProcessingOrder {
		messageQueue[messageType] = []offlineMessage{}
	}

	// load stored messages
	messages, err := m.datastore.OfflineMessages().GetMessages()
	if err != nil {
		log.Errorf("error getting offline messages: %s", err)
		return
	}

	// sort them into the queue by message type
	for url, ser := range messages {
		env := new(pb.Envelope)
		if err := proto.Unmarshal(ser, env); err != nil {
			log.Errorf("error unmarshalling offline message: %s", err)
			continue
		}
		switch env.Message.Type {
		case pb.Message_THREAD_MERGE,
			pb.Message_THREAD_INVITE,
			pb.Message_THREAD_EXTERNAL_INVITE,
			pb.Message_THREAD_JOIN,
			pb.Message_THREAD_LEAVE,
			pb.Message_THREAD_DATA,
			pb.Message_THREAD_ANNOTATION,
			pb.Message_THREAD_IGNORE:
			threadMessages = append(threadMessages, offlineMessage{
				addr: url,
				env:  *env,
				date: getThreadEnvelopeDate(env),
			})
		default:
			messageQueue[env.Message.Type] = append(messageQueue[env.Message.Type], offlineMessage{
				addr: url,
				env:  *env,
				date: time.Now(),
			})
		}
	}

	// process the thread list by date ascending
	sort.Sort(sortedMessages(threadMessages))
	var toDelete []string
	for _, om := range threadMessages {
		if err := m.handleMessage(&om.env, om.addr); err != nil {
			log.Errorf("error handling offline thread message: %s", err)
		} else {
			toDelete = append(toDelete, om.addr)
		}
	}

	// process all other messages from queue in order
	for _, messageType := range messageProcessingOrder {
		queue, ok := messageQueue[messageType]
		if !ok {
			continue
		}
		for _, om := range queue {
			if err := m.handleMessage(&om.env, om.addr); err != nil {
				log.Errorf("error handling offline message: %s", err)
			} else {
				toDelete = append(toDelete, om.addr)
			}
		}
	}

	// delete messages that were successfully processed from the database
	for _, url := range toDelete {
		if err := m.datastore.OfflineMessages().DeleteMessage(url); err != nil {
			log.Errorf("error deleting offline message: %s", err)
		}
	}
}

func getThreadEnvelopeDate(env *pb.Envelope) time.Time {
	var date time.Time
	var ts *timestamp.Timestamp
	if env.Message.Payload == nil {
		return date
	}
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(env.Message.Payload, signed); err != nil {
		return date
	}
	threadBlock := new(pb.ThreadHeader)
	if err := proto.Unmarshal(signed.Block, threadBlock); err != nil {
		return date
	}
	if threadBlock.Header != nil {
		ts = threadBlock.Header.Date
	} else {
		// could be merge
		merge := new(pb.ThreadMerge)
		if err := proto.Unmarshal(signed.Block, merge); err != nil {
			return date
		}
		ts = merge.Date
	}
	parsed, err := ptypes.Timestamp(ts)
	if err != nil {
		return date
	}
	return parsed
}

func getEnvelopeSenderId(env *pb.Envelope) (peer.ID, error) {
	pubkey, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return "", err
	}
	return peer.IDFromPublicKey(pubkey)
}

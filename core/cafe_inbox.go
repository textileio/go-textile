package core

import (
	"sync"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// cafeInFlushGroupSize is the size of concurrently processed messages
const cafeInFlushGroupSize = 16

// CafeInbox queues and processes outbound thread messages
type CafeInbox struct {
	service        func() *CafeService
	threadsService func() *ThreadsService
	node           func() *core.IpfsNode
	datastore      repo.Datastore
	mux            sync.Mutex
}

// NewCafeInbox creates a new inbox queue
func NewCafeInbox(
	service func() *CafeService,
	threadsService func() *ThreadsService,
	node func() *core.IpfsNode,
	datastore repo.Datastore,
) *CafeInbox {
	return &CafeInbox{
		service:        service,
		threadsService: threadsService,
		node:           node,
		datastore:      datastore,
	}
}

// CheckMessages asks each active cafe session for new messages
func (q *CafeInbox) CheckMessages() error {
	// get active cafe sessions
	sessions := q.datastore.CafeSessions().List()
	if len(sessions) == 0 {
		return nil
	}

	// check each concurrently
	wg := sync.WaitGroup{}
	var cerr error
	for _, session := range sessions {
		cafe, err := peer.IDB58Decode(session.CafeId)
		if err != nil {
			cerr = err
			continue
		}
		wg.Add(1)
		go func(cafe peer.ID) {
			if err := q.service().CheckMessages(cafe); err != nil {
				cerr = err
			}
			wg.Done()
		}(cafe)
	}
	wg.Wait()
	return cerr
}

// Add adds an inbound message
func (q *CafeInbox) Add(msg *pb.CafeMessage) error {
	date, err := ptypes.Timestamp(msg.Date)
	if err != nil {
		return err
	}
	return q.datastore.CafeMessages().Add(&repo.CafeMessage{
		Id:     msg.Id,
		PeerId: msg.PeerId,
		Date:   date,
	})
}

// Flush processes pending messages
func (q *CafeInbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()

	if q.threadsService() == nil || q.service() == nil {
		return
	}

	if err := q.batch(q.datastore.CafeMessages().List("", cafeInFlushGroupSize)); err != nil {
		log.Errorf("cafe inbox batch error: %s", err)
		return
	}
}

// batch flushes a batch of messages
func (q *CafeInbox) batch(msgs []repo.CafeMessage) error {
	if len(msgs) == 0 {
		return nil
	}

	var berr error
	var toDelete []string
	wg := sync.WaitGroup{}
	for _, msg := range msgs {
		wg.Add(1)
		go func(msg *repo.CafeMessage) {
			if err := q.handle(msg); err != nil {
				berr = err
				return
			}
			toDelete = append(toDelete, msg.Id)
			wg.Done()
		}(&msg)
	}
	wg.Wait()

	// next batch
	offset := msgs[len(msgs)-1].Id
	next := q.datastore.CafeMessages().List(offset, cafeInFlushGroupSize)

	var deleted []string
	for _, id := range toDelete {
		if err := q.datastore.CafeMessages().Delete(id); err != nil {
			log.Errorf("failed to delete cafe message %s: %s", id, err)
			continue
		}
		deleted = append(deleted, id)
	}
	log.Debugf("handled %d cafe messages", len(deleted))

	// keep going unless an error occurred
	if berr == nil {
		return q.batch(next)
	}
	return berr
}

// handle handles a single message
func (q *CafeInbox) handle(msg *repo.CafeMessage) error {
	pid, err := peer.IDB58Decode(msg.PeerId)
	if err != nil {
		return err
	}

	// download the actual message
	ciphertext, err := ipfs.DataAtPath(q.node(), msg.Id)
	if err != nil {
		return err
	}

	envb, err := crypto.Decrypt(q.node().PrivateKey, ciphertext)
	if err != nil {
		return err
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(envb, env); err != nil {
		return err
	}

	if err := q.threadsService().service.VerifyEnvelope(env, pid); err != nil {
		log.Warningf("error verifying cafe message: %s", err)
		return nil
	}

	// pass to thread service for normal handling
	if _, err := q.threadsService().Handle(pid, env); err != nil {
		return err
	}
	return nil
}

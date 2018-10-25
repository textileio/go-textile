package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"sync"
	"time"
)

// how often to run store job (daemon only)
const kCafeInFlushFrequency = time.Minute * 10

// the size of concurrently processed messages
const cafeInFlushGroupSize = 16

// CafeInbox queues and processes outbound thread messages
type CafeInbox struct {
	service   func() *ThreadsService
	node      func() *core.IpfsNode
	datastore repo.Datastore
	mux       sync.Mutex
}

// NewCafeInbox creates a new inbox queue
func NewCafeInbox(service func() *ThreadsService, node func() *core.IpfsNode, datastore repo.Datastore) *CafeInbox {
	return &CafeInbox{
		service:   service,
		node:      node,
		datastore: datastore,
	}
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

// Run starts a job ticker which processes any pending messages
func (q *CafeInbox) Run() {
	tick := time.NewTicker(kCafeInFlushFrequency)
	defer tick.Stop()
	go q.Flush()
	for {
		select {
		case <-tick.C:
			go q.Flush()
		}
	}
}

// Flush processes pending messages
func (q *CafeInbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()

	// check service status
	if q.service() == nil {
		return
	}

	// start at zero offset
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

	// process the batch
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

	// clean up
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
	ciphertext, err := ipfs.GetDataAtPath(q.node(), msg.Id)
	if err != nil {
		return err
	}

	// decrypt and unmarshal to envelope
	envb, err := crypto.Decrypt(q.node().PrivateKey, ciphertext)
	if err != nil {
		return err
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(envb, env); err != nil {
		return err
	}

	// check signature
	if err := q.service().service.VerifyEnvelope(env, pid); err != nil {
		log.Warningf("error verifying cafe message: %s", err)
		return nil
	}

	// finally, pass to thread service for normal handling
	if _, err := q.service().Handle(pid, env); err != nil {
		return err
	}
	return nil
}

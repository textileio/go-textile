package core

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"sync"
	"time"
)

// how often to run store job (daemon only)
const kThreadsFlushFrequency = time.Minute * 10

// the size of concurrently processed requests
// note: msgs from this group are batched to each receiver
const threadsFlushGroupSize = 16

// ThreadsOutbox queues and processes outbound thread messages
type ThreadsOutbox struct {
	service    func() *ThreadsService
	node       func() *core.IpfsNode
	datastore  repo.Datastore
	cafeOutbox *CafeOutbox
	mux        sync.Mutex
}

// NewThreadsOutbox creates a new outbox queue
func NewThreadsOutbox(
	service func() *ThreadsService,
	node func() *core.IpfsNode,
	datastore repo.Datastore,
	cafeOutbox *CafeOutbox,
) *ThreadsOutbox {
	return &ThreadsOutbox{
		service:    service,
		node:       node,
		datastore:  datastore,
		cafeOutbox: cafeOutbox,
	}
}

// Add adds an outbound message
func (q *ThreadsOutbox) Add(pid peer.ID, env *pb.Envelope) {
	msg := &repo.ThreadMessage{
		Id:       ksuid.New().String(),
		PeerId:   pid.Pretty(),
		Envelope: env,
		Date:     time.Now(),
	}
	if err := q.datastore.ThreadMessages().Add(msg); err != nil {
		log.Errorf("error adding thread message for %s: %s", pid.Pretty(), err)
	}

	// try to flush the queue now
	go q.Flush()
}

// Run starts a job ticker which processes any pending messages
func (q *ThreadsOutbox) Run() {
	tick := time.NewTicker(kThreadsFlushFrequency)
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
func (q *ThreadsOutbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()

	// check service status
	if q.service() == nil {
		return
	}

	// start at zero offset
	if err := q.batch(q.datastore.ThreadMessages().List("", threadsFlushGroupSize)); err != nil {
		log.Errorf("thread outbox batch error: %s", err)
		return
	}
}

// batch flushes a batch of requests
func (q *ThreadsOutbox) batch(msgs []repo.ThreadMessage) error {
	if len(msgs) == 0 {
		return nil
	}

	// process the batch
	var berr error
	var toDelete []string
	wg := sync.WaitGroup{}
	for _, msg := range msgs {
		wg.Add(1)
		go func(msg *repo.ThreadMessage) {
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
	next := q.datastore.ThreadMessages().List(offset, threadsFlushGroupSize)

	// clean up
	var deleted []string
	for _, id := range toDelete {
		if err := q.datastore.ThreadMessages().Delete(id); err != nil {
			log.Errorf("failed to delete thread message %s: %s", id, err)
			continue
		}
		deleted = append(deleted, id)
	}
	log.Debugf("handled %d thread messages", len(deleted))

	// keep going unless an error occurred
	if berr == nil {
		return q.batch(next)
	}
	return berr
}

// handle handles a single message
func (q *ThreadsOutbox) handle(msg *repo.ThreadMessage) error {
	pid, err := peer.IDB58Decode(msg.PeerId)
	if err != nil {
		return err
	}
	// first, attempt to send the message directly to the recipient
	if err := q.service().SendMessage(pid, msg.Envelope); err != nil {
		log.Warningf("send thread message direct to %s failed: %s", pid.Pretty(), err)

		// peer is offline, queue an outbound cafe request for the peer's inbox(es)
		contact := q.datastore.Contacts().Get(pid.Pretty())
		if contact != nil && len(contact.Inboxes) > 0 {
			hash, err := q.prepForInbox(pid, msg.Envelope)
			if err != nil {
				return err
			}

			// add an inbox request for message delivery
			q.cafeOutbox.InboxRequest(pid, hash.B58String(), contact.Inboxes)
		}
	}
	return nil
}

// prepForInbox encrypts and pins a message intended for a peer inbox
func (q *ThreadsOutbox) prepForInbox(pid peer.ID, env *pb.Envelope) (mh.Multihash, error) {
	// encrypt envelope w/ recipient's pk
	envb, err := proto.Marshal(env)
	if err != nil {
		return nil, err
	}
	pk, err := pid.ExtractPublicKey()
	if err != nil {
		return nil, err
	}
	ciphertext, err := crypto.Encrypt(pk, envb)
	if err != nil {
		return nil, err
	}

	// pin it
	id, err := ipfs.PinData(q.node(), bytes.NewReader(ciphertext))
	if err != nil {
		return nil, err
	}

	// add a store request for the encrypted message
	q.cafeOutbox.Add(id.Hash().B58String(), repo.CafeStoreRequest)

	return id.Hash(), nil
}

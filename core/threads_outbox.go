package core

import (
	"context"
	"sync"
	"time"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core"

	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/service"
)

// threadsFlushGroupSize is the size of concurrently processed messages
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
func (q *ThreadsOutbox) Add(pid peer.ID, env *pb.Envelope) error {
	log.Debugf("adding thread message for %s", pid.Pretty())
	return q.datastore.ThreadMessages().Add(&repo.ThreadMessage{
		Id:       ksuid.New().String(),
		PeerId:   pid.Pretty(),
		Envelope: env,
		Date:     time.Now(),
	})
}

// Flush processes pending messages
func (q *ThreadsOutbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()
	log.Debug("flushing thread messages")

	if q.service() == nil {
		return
	}

	if err := q.batch(q.datastore.ThreadMessages().List("", threadsFlushGroupSize)); err != nil {
		log.Errorf("thread outbox batch error: %s", err)
		return
	}
}

// batch flushes a batch of messages
func (q *ThreadsOutbox) batch(msgs []repo.ThreadMessage) error {
	log.Debugf("handling %d thread messages", len(msgs))
	if len(msgs) == 0 {
		return nil
	}

	// group by peer id
	groups := make(map[string][]repo.ThreadMessage)
	for _, msg := range msgs {
		groups[msg.PeerId] = append(groups[msg.PeerId], msg)
	}

	var berr error
	var toDelete []string
	wg := sync.WaitGroup{}
	for id, group := range groups {
		pid, err := peer.IDB58Decode(id)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(pid peer.ID, msgs []repo.ThreadMessage) {
			for _, msg := range msgs {
				if err := q.handle(pid, msg); err != nil {
					berr = err
					return
				}
				toDelete = append(toDelete, msg.Id)
			}
			wg.Done()
		}(pid, group)
	}
	wg.Wait()

	// flush the outbox before starting a new batch
	go q.cafeOutbox.Flush()

	// next batch
	offset := msgs[len(msgs)-1].Id
	next := q.datastore.ThreadMessages().List(offset, threadsFlushGroupSize)

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
func (q *ThreadsOutbox) handle(pid peer.ID, msg repo.ThreadMessage) error {
	// first, attempt to send the message directly to the recipient
	ctx, cancel := context.WithTimeout(context.Background(), service.DirectTimeout)
	defer cancel()

	var err error
	if q.service().online {
		err = q.service().SendMessage(ctx, pid, msg.Envelope)
	}
	if !q.service().online || err != nil {
		if err != nil {
			log.Debugf("send thread message direct to %s failed: %s", pid.Pretty(), err)
		}

		// peer is offline, queue an outbound cafe request for the peer's inbox(es)
		contact := q.datastore.Contacts().Get(pid.Pretty())
		if contact != nil && len(contact.Inboxes) > 0 {
			log.Debugf("sending thread message for %s to inbox(es)", pid.Pretty())

			// add an inbox request for message delivery
			if err := q.cafeOutbox.InboxRequest(pid, msg.Envelope, contact.Inboxes); err != nil {
				return err
			}
		}
	}
	return nil
}

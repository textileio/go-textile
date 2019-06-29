package core

import (
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/ipfs/go-ipfs/core"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

// blockFlushGroupSize is the size of concurrently processed messages
// note: msgs from this group are batched to each receiver
const blockFlushGroupSize = 16

// BlockOutbox queues and processes outbound thread messages
type BlockOutbox struct {
	service    func() *ThreadsService
	node       func() *core.IpfsNode
	datastore  repo.Datastore
	cafeOutbox *CafeOutbox
	mux        sync.Mutex
}

// NewBlockOutbox creates a new outbox queue
func NewBlockOutbox(service func() *ThreadsService, node func() *core.IpfsNode, datastore repo.Datastore, cafeOutbox *CafeOutbox) *BlockOutbox {
	return &BlockOutbox{
		service:    service,
		node:       node,
		datastore:  datastore,
		cafeOutbox: cafeOutbox,
	}
}

// Add adds an outbound message
func (q *BlockOutbox) Add(peerId string, env *pb.Envelope) error {
	log.Debugf("adding block message for %s", peerId)
	return q.datastore.BlockMessages().Add(&pb.BlockMessage{
		Id:   ksuid.New().String(),
		Peer: peerId,
		Env:  env,
		Date: ptypes.TimestampNow(),
	})
}

// Flush processes pending messages
func (q *BlockOutbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()
	log.Debug("flushing block messages")

	if q.service() == nil {
		return
	}

	q.batch(q.datastore.BlockMessages().List("", blockFlushGroupSize))
}

// batch flushes a batch of messages
func (q *BlockOutbox) batch(msgs []pb.BlockMessage) {
	log.Debugf("handling %d block messages", len(msgs))
	if len(msgs) == 0 {
		return
	}

	// group by peer id
	groups := make(map[string][]pb.BlockMessage)
	for _, msg := range msgs {
		groups[msg.Peer] = append(groups[msg.Peer], msg)
	}

	var toDelete []string
	wg := sync.WaitGroup{}
	for id, group := range groups {
		wg.Add(1)
		go func(id string, msgs []pb.BlockMessage) {
			for _, msg := range msgs {
				if err := q.handle(msg); err != nil {
					log.Warningf("error handling block message %s: %s", msg.Id, err)
					continue
				}
				toDelete = append(toDelete, msg.Id)
			}
			wg.Done()
		}(id, group)
	}
	wg.Wait()

	// next batch
	offset := msgs[len(msgs)-1].Id
	next := q.datastore.BlockMessages().List(offset, blockFlushGroupSize)

	var deleted []string
	for _, id := range toDelete {
		if err := q.datastore.BlockMessages().Delete(id); err != nil {
			log.Errorf("failed to delete block message %s: %s", id, err)
			continue
		}
		deleted = append(deleted, id)
	}
	log.Debugf("handled %d block messages", len(deleted))

	q.batch(next)
}

// handle handles a single message
func (q *BlockOutbox) handle(msg pb.BlockMessage) error {
	err := q.service().SendMessage(msg)
	if err != nil {
		log.Debugf("send block message direct to %s failed: %s", msg.Peer, err)

		// peer is offline, queue an outbound cafe request for the peer's inbox(es)
		contact := q.datastore.Peers().Get(msg.Peer)
		if contact != nil && len(contact.Inboxes) > 0 {
			log.Debugf("sending block message for %s to inbox(es)", msg.Peer)

			// add an inbox request for message delivery
			err = q.cafeOutbox.AddForInbox(msg.Peer, msg.Env, contact.Inboxes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

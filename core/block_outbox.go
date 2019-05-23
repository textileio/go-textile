package core

import (
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/ipfs/go-ipfs/core"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/ipfs"
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
func (q *BlockOutbox) Add(pid peer.ID, env *pb.Envelope) error {
	log.Debugf("adding block message for %s", pid.Pretty())
	return q.datastore.BlockMessages().Add(&pb.BlockMessage{
		Id:   ksuid.New().String(),
		Peer: pid.Pretty(),
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

	// exclude blocks that have an incomplete sync group
	// TODO: also exclude newer
	syncing := q.datastore.CafeRequests().ListIncompleteSyncGroups()

	err := q.batch(q.datastore.BlockMessages().List("", blockFlushGroupSize, syncing))
	if err != nil {
		log.Errorf("block outbox batch error: %s", err)
		return
	}
}

// batch flushes a batch of messages
func (q *BlockOutbox) batch(msgs []pb.BlockMessage) error {
	log.Debugf("handling %d block messages", len(msgs))
	if len(msgs) == 0 {
		return nil
	}

	// group by peer id
	groups := make(map[string][]pb.BlockMessage)
	for _, msg := range msgs {
		groups[msg.Peer] = append(groups[msg.Peer], msg)
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
		go func(pid peer.ID, msgs []pb.BlockMessage) {
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

	// next batch
	offset := msgs[len(msgs)-1].Id
	syncing := q.datastore.CafeRequests().ListIncompleteSyncGroups()
	next := q.datastore.BlockMessages().List(offset, blockFlushGroupSize, syncing)

	var deleted []string
	for _, id := range toDelete {
		if err := q.datastore.BlockMessages().Delete(id); err != nil {
			log.Errorf("failed to delete block message %s: %s", id, err)
			continue
		}
		deleted = append(deleted, id)
	}
	log.Debugf("handled %d block messages", len(deleted))

	// keep going unless an error occurred
	if berr == nil {
		return q.batch(next)
	}
	return berr
}

// handle handles a single message
func (q *BlockOutbox) handle(pid peer.ID, msg pb.BlockMessage) error {
	// first, attempt to send the message directly to the recipient
	sendable := q.service().online
	if sendable {
		connected, err := ipfs.SwarmConnected(q.node(), pid)
		if err != nil {
			return err
		}
		if !connected {
			sendable = false
		}
	}
	var err error
	if sendable {
		err = q.service().SendMessage(nil, pid, msg.Env)
	}
	if !sendable || err != nil {
		if err != nil {
			log.Debugf("send block message direct to %s failed: %s", pid.Pretty(), err)
		}

		// peer is offline, queue an outbound cafe request for the peer's inbox(es)
		contact := q.datastore.Peers().Get(pid.Pretty())
		if contact != nil && len(contact.Inboxes) > 0 {
			log.Debugf("sending block message for %s to inbox(es)", pid.Pretty())

			// add an inbox request for message delivery
			err = q.cafeOutbox.AddForInbox(pid, msg.Env, contact.Inboxes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

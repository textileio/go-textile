package core

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-ipfs/core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

// cafeInFlushGroupSize is the size of concurrently processed messages
const cafeInFlushGroupSize = 16

// cafeInMaxDownloadAttempts is the number of times a message can fail to download before being deleted
const cafeInMaxDownloadAttempts = 5

// CafeInbox queues and processes downloaded cafe messages
type CafeInbox struct {
	service        func() *CafeService
	threadsService func() *ThreadsService
	node           func() *core.IpfsNode
	datastore      repo.Datastore
	checking       bool
	lock           sync.Mutex
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
	if q.checking {
		return nil
	}
	q.checking = true
	defer func() {
		q.checking = false
	}()

	// get active cafe sessions
	sessions := q.datastore.CafeSessions().List().Items
	if len(sessions) == 0 {
		return nil
	}

	// check each concurrently
	wg := sync.WaitGroup{}
	var cerr error
	for _, session := range sessions {
		wg.Add(1)
		go func(cafeId string) {
			if err := q.service().CheckMessages(cafeId); err != nil {
				cerr = err
			}
			wg.Done()
		}(session.Id)
	}
	wg.Wait()
	return cerr
}

// Add adds an inbound message
func (q *CafeInbox) Add(msg *pb.CafeMessage) error {
	log.Debugf("received cafe message from %s: %s", ipfs.ShortenID(msg.Peer), msg.Id)

	return q.datastore.CafeMessages().Add(msg)
}

// Flush processes pending messages
func (q *CafeInbox) Flush() {
	q.lock.Lock()
	defer q.lock.Unlock()
	log.Debug("flushing cafe inbox")

	if q.threadsService() == nil || !q.threadsService().online || q.service() == nil {
		return
	}

	q.batch(q.datastore.CafeMessages().List("", cafeInFlushGroupSize))
}

// batch flushes a batch of messages
func (q *CafeInbox) batch(msgs []pb.CafeMessage) {
	log.Debugf("handling %d cafe messages", len(msgs))
	if len(msgs) == 0 {
		return
	}

	for _, msg := range msgs {
		go func(msg pb.CafeMessage) {
			err := q.handle(msg)
			if err != nil {
				log.Warningf("handle attempt failed for cafe message %s: %s", msg.Id, err)
				return
			}
			err = q.datastore.CafeMessages().Delete(msg.Id)
			if err != nil {
				log.Errorf("failed to delete cafe message %s: %s", msg.Id, err)
			} else {
				log.Debugf("handled cafe message %s", msg.Id)
			}
		}(msg)
	}

	// next batch
	offset := msgs[len(msgs)-1].Id
	next := q.datastore.CafeMessages().List(offset, cafeInFlushGroupSize)

	// keep going
	q.batch(next)
}

// handle handles a single message
func (q *CafeInbox) handle(msg pb.CafeMessage) error {
	pid, err := peer.IDB58Decode(msg.Peer)
	if err != nil {
		return q.handleErr(fmt.Errorf("error decoding msg peer: %s", err), msg)
	}

	envb, err := ipfs.DataAtPath(q.node(), msg.Id)
	if err != nil {
		return q.handleErr(fmt.Errorf("error getting msg data: %s", err), msg)
	}
	env := new(pb.Envelope)
	err = proto.Unmarshal(envb, env)
	if err != nil {
		return q.handleErr(err, msg)
	}

	err = q.threadsService().service.VerifyEnvelope(env, pid)
	if err != nil {
		return q.handleErr(err, msg)
	}

	// pass to thread service for normal handling
	_, err = q.threadsService().Handle(env, pid)
	if err != nil {
		return q.handleErr(err, msg)
	}
	return nil
}

// handleErr deletes or adds an attempt to a message processing error
func (q *CafeInbox) handleErr(herr error, msg pb.CafeMessage) error {
	var err error
	if msg.Attempts+1 >= cafeInMaxDownloadAttempts {
		err = q.datastore.CafeMessages().Delete(msg.Id)
	} else {
		err = q.datastore.CafeMessages().AddAttempt(msg.Id)
	}
	if err != nil {
		return err
	}
	return herr
}

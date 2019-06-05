package core

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-ipfs/core"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

// cafeInFlushGroupSize is the size of concurrently processed messages
const cafeInFlushGroupSize = 16

// maxDownloadAttempts is the number of times a message can fail to download before being deleted
const maxDownloadAttempts = 5

// CafeInbox queues and processes outbound thread messages
type CafeInbox struct {
	service        func() *CafeService
	threadsService func() *ThreadsService
	node           func() *core.IpfsNode
	datastore      repo.Datastore
	checking       bool
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
		cafe, err := peer.IDB58Decode(session.Id)
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
	log.Debugf("received cafe message from %s: %s", ipfs.ShortenID(msg.Peer), msg.Id)

	// reset attempts
	return q.datastore.CafeMessages().Add(&pb.CafeMessage{
		Id:   msg.Id,
		Peer: msg.Peer,
		Date: msg.Date,
	})
}

// Flush processes pending messages
func (q *CafeInbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()
	log.Debug("flushing cafe inbox")

	if q.threadsService() == nil || !q.threadsService().online || q.service() == nil {
		return
	}

	if err := q.batch(q.datastore.CafeMessages().List("", cafeInFlushGroupSize)); err != nil {
		log.Errorf("cafe inbox batch error: %s", err)
		return
	}
}

// batch flushes a batch of messages
func (q *CafeInbox) batch(msgs []pb.CafeMessage) error {
	log.Debugf("handling %d cafe messages", len(msgs))
	if len(msgs) == 0 {
		return nil
	}

	for _, msg := range msgs {
		go func(msg pb.CafeMessage) {
			if err := q.handle(msg); err != nil {
				log.Warningf("handle attempt failed for cafe message %s: %s", msg.Id, err)
				return
			}
			if err := q.datastore.CafeMessages().Delete(msg.Id); err != nil {
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
	return q.batch(next)
}

// handle handles a single message
func (q *CafeInbox) handle(msg pb.CafeMessage) error {
	pid, err := peer.IDB58Decode(msg.Peer)
	if err != nil {
		return q.handleErr(fmt.Errorf("error decoding msg peer: %s", err), msg)
	}

	// download the actual message
	ciphertext, err := ipfs.DataAtPath(q.node(), msg.Id)
	if err != nil {
		return q.handleErr(fmt.Errorf("error getting msg data: %s", err), msg)
	}

	envb, err := crypto.Decrypt(q.node().PrivateKey, ciphertext)
	if err != nil {
		return q.handleErr(fmt.Errorf("error decrypting msg: %s", err), msg)
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(envb, env); err != nil {
		return q.handleErr(fmt.Errorf("error unmarshaling env: %s", err), msg)
	}

	if err := q.threadsService().service.VerifyEnvelope(env, pid); err != nil {
		return q.handleErr(fmt.Errorf("error verifying env: %s", err), msg)
	}

	// pass to thread service for normal handling
	if _, err := q.threadsService().Handle(pid, env); err != nil {
		return q.handleErr(fmt.Errorf("error handling msg: %s", err), msg)
	}
	return nil
}

// handleErr deletes or adds an attempt to a message processing error
func (q *CafeInbox) handleErr(herr error, msg pb.CafeMessage) error {
	if msg.Attempts+1 >= maxDownloadAttempts {
		if err := q.datastore.CafeMessages().Delete(msg.Id); err != nil {
			return err
		}
	} else {
		if err := q.datastore.CafeMessages().AddAttempt(msg.Id); err != nil {
			return err
		}
	}
	return herr
}

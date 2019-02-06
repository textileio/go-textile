package core

import (
	"bytes"
	"errors"
	"sync"
	"time"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core"

	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// cafeOutFlushGroupSize is the size of concurrently processed requests
// note: reqs from this group are batched to each cafe
const cafeOutFlushGroupSize = 16

// CafeOutbox queues and processes outbound cafe requests
type CafeOutbox struct {
	service   func() *CafeService
	node      func() *core.IpfsNode
	datastore repo.Datastore
	mux       sync.Mutex
}

// NewCafeOutbox creates a new outbox queue
func NewCafeOutbox(service func() *CafeService, node func() *core.IpfsNode, datastore repo.Datastore) *CafeOutbox {
	return &CafeOutbox{
		service:   service,
		node:      node,
		datastore: datastore,
	}
}

// Add adds a request for each active cafe session
func (q *CafeOutbox) Add(target string, rtype repo.CafeRequestType) error {
	if rtype == repo.CafePeerInboxRequest {
		return errors.New("inbox request to own inbox, aborting")
	}

	sessions := q.datastore.CafeSessions().List()
	if len(sessions) == 0 {
		return nil
	}

	// for each session, add a req
	for _, session := range sessions {
		// all possible request types are for our own peer
		if err := q.add(q.node().Identity, target, *session.Cafe, rtype); err != nil {
			return err
		}
	}
	return nil
}

// InboxRequest adds a request for a peer's inbox(es)
func (q *CafeOutbox) InboxRequest(pid peer.ID, env *pb.Envelope, inboxes []repo.Cafe) error {
	if len(inboxes) == 0 {
		return nil
	}

	hash, err := q.prepForInbox(pid, env)
	if err != nil {
		return err
	}

	for _, inbox := range inboxes {
		cafe := repoCafeToProto(inbox)
		if err := q.add(pid, hash.B58String(), *cafe, repo.CafePeerInboxRequest); err != nil {
			return err
		}
	}
	return nil
}

// Flush processes pending requests
func (q *CafeOutbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()
	log.Debug("flushing cafe outbox")

	if q.service() == nil {
		return
	}

	if err := q.batch(q.datastore.CafeRequests().List("", cafeOutFlushGroupSize)); err != nil {
		log.Errorf("cafe outbox batch error: %s", err)
		return
	}
}

// add queues a single request
func (q *CafeOutbox) add(pid peer.ID, target string, cafe pb.Cafe, rtype repo.CafeRequestType) error {
	log.Debugf("adding cafe %s request for %s to %s: %s",
		rtype.Description(), ipfs.ShortenID(pid.Pretty()), ipfs.ShortenID(cafe.Peer), target)
	return q.datastore.CafeRequests().Add(&repo.CafeRequest{
		Id:       ksuid.New().String(),
		PeerId:   pid.Pretty(),
		TargetId: target,
		Cafe:     cafe,
		Type:     rtype,
		Date:     time.Now(),
	})
}

// batch flushes a batch of requests
func (q *CafeOutbox) batch(reqs []repo.CafeRequest) error {
	log.Debugf("handling %d cafe requests", len(reqs))
	if len(reqs) == 0 {
		return nil
	}

	// group reqs by cafe
	groups := make(map[string][]repo.CafeRequest)
	for _, req := range reqs {
		groups[req.Cafe.Peer] = append(groups[req.Cafe.Peer], req)
	}

	// process each cafe group concurrently
	var berr error
	var toDelete []string
	wg := sync.WaitGroup{}
	for cafeId, group := range groups {
		cafe, err := peer.IDB58Decode(cafeId)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(cafe peer.ID, reqs []repo.CafeRequest) {
			// group by type
			types := make(map[repo.CafeRequestType][]repo.CafeRequest)
			for _, req := range reqs {
				types[req.Type] = append(types[req.Type], req)
			}
			for t, group := range types {
				handled, err := q.handle(group, t, cafe)
				if err != nil {
					berr = err
				}
				for _, id := range handled {
					toDelete = append(toDelete, id)
				}
			}
			wg.Done()
		}(cafe, group)
	}
	wg.Wait()

	// next batch
	offset := reqs[len(reqs)-1].Id
	next := q.datastore.CafeRequests().List(offset, cafeOutFlushGroupSize)

	var deleted []string
	for _, id := range toDelete {
		if err := q.datastore.CafeRequests().Delete(id); err != nil {
			log.Errorf("failed to delete cafe request %s: %s", id, err)
			continue
		}
		deleted = append(deleted, id)
	}
	log.Debugf("handled %d cafe requests", len(deleted))

	// keep going unless an error occurred
	if berr == nil {
		return q.batch(next)
	}
	return berr
}

// handle handles a group of requests for a single cafe
func (q *CafeOutbox) handle(reqs []repo.CafeRequest, rtype repo.CafeRequestType, cafe peer.ID) ([]string, error) {
	var handled []string
	var herr error
	switch rtype {

	// store requests are handled in bulk
	case repo.CafeStoreRequest:
		var cids []string
		for _, req := range reqs {
			cids = append(cids, req.TargetId)
		}

		stored, err := q.service().Store(cids, cafe)
		for _, s := range stored {
			for _, r := range reqs {
				if r.TargetId == s {
					handled = append(handled, r.Id)
				}
			}
		}
		if err != nil {
			log.Errorf("cafe %s request to %s failed: %s", rtype.Description(), cafe.Pretty(), err)
			herr = err
		}

	case repo.CafeStoreThreadRequest:
		for _, req := range reqs {
			thrd := q.datastore.Threads().Get(req.TargetId)
			if thrd == nil {
				log.Warningf("could not find thread: %s", req.TargetId)
				handled = append(handled, req.Id)
				continue
			}

			if err := q.service().StoreThread(thrd, cafe); err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.Description(), cafe.Pretty(), err)
				herr = err
				continue
			}
			handled = append(handled, req.Id)
		}

	case repo.CafePeerInboxRequest:
		for _, req := range reqs {
			pid, err := peer.IDB58Decode(req.PeerId)
			if err != nil {
				herr = err
				continue
			}

			if err := q.service().DeliverMessage(req.TargetId, pid, protoCafeToRepo(&req.Cafe)); err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.Description(), cafe.Pretty(), err)
				herr = err
				continue
			}
			handled = append(handled, req.Id)
		}

	}
	return handled, herr
}

// prepForInbox encrypts and pins a message intended for a peer inbox
func (q *CafeOutbox) prepForInbox(pid peer.ID, env *pb.Envelope) (mh.Multihash, error) {
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

	id, err := ipfs.AddData(q.node(), bytes.NewReader(ciphertext), true)
	if err != nil {
		return nil, err
	}

	if err := q.Add(id.Hash().B58String(), repo.CafeStoreRequest); err != nil {
		return nil, err
	}

	return id.Hash(), nil
}

package core

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
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
	// get active cafe sessions
	sessions := q.datastore.CafeSessions().List()
	if len(sessions) == 0 {
		return nil
	}

	// for each session, add a req
	for _, session := range sessions {
		// all possible request types are for our own peer
		if err := q.add(q.node().Identity, target, session.CafeId, rtype); err != nil {
			return err
		}
	}

	// flush the queue now
	go q.Flush()

	return nil
}

// InboxRequest adds a request for a peer's inbox(es)
func (q *CafeOutbox) InboxRequest(pid peer.ID, env *pb.Envelope, inboxes []string) error {
	if len(inboxes) == 0 {
		return nil
	}

	// encrypt for peer
	hash, err := q.prepForInbox(pid, env)
	if err != nil {
		return err
	}

	// for each inbox, add a req
	for _, inbox := range inboxes {
		q.add(pid, hash.B58String(), inbox, repo.CafePeerInboxRequest)
	}
	return nil
}

// Flush processes pending requests
func (q *CafeOutbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()

	// check service status
	if q.service() == nil {
		return
	}

	// start at zero offset
	if err := q.batch(q.datastore.CafeRequests().List("", cafeOutFlushGroupSize)); err != nil {
		log.Errorf("cafe outbox batch error: %s", err)
		return
	}
}

// add queues a single request
func (q *CafeOutbox) add(pid peer.ID, target string, cafeId string, rtype repo.CafeRequestType) error {
	return q.datastore.CafeRequests().Add(&repo.CafeRequest{
		Id:       ksuid.New().String(),
		PeerId:   pid.Pretty(),
		TargetId: target,
		CafeId:   cafeId,
		Type:     rtype,
		Date:     time.Now(),
	})
}

// batch flushes a batch of requests
func (q *CafeOutbox) batch(reqs []repo.CafeRequest) error {
	if len(reqs) == 0 {
		return nil
	}

	// group reqs by cafe
	groups := make(map[string][]repo.CafeRequest)
	for _, req := range reqs {
		groups[req.CafeId] = append(groups[req.CafeId], req)
	}

	// process each cafe group concurrently
	var berr error
	var toDelete []string
	wg := sync.WaitGroup{}
	for cafeId, group := range groups {
		cafe, err := peer.IDB58Decode(cafeId)
		if err != nil {
			berr = err
			continue
		}
		// group reqs by type
		types := make(map[repo.CafeRequestType][]repo.CafeRequest)
		for _, req := range group {
			types[req.Type] = append(types[req.Type], req)
		}
		wg.Add(1)
		go func(types map[repo.CafeRequestType][]repo.CafeRequest, cafe peer.ID) {
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
		}(types, cafe)
	}
	wg.Wait()

	// next batch
	offset := reqs[len(reqs)-1].Id
	next := q.datastore.CafeRequests().List(offset, cafeOutFlushGroupSize)

	// clean up
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
				err := errors.New(fmt.Sprintf("could not find thread: %s", req.TargetId))
				log.Error(err.Error())
				herr = err
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
			cid, err := peer.IDB58Decode(req.CafeId)
			if err != nil {
				herr = err
				continue
			}
			if err := q.service().DeliverMessage(req.TargetId, pid, cid); err != nil {
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

	// pin it
	id, err := ipfs.PinData(q.node(), bytes.NewReader(ciphertext))
	if err != nil {
		return nil, err
	}

	// add a store request for the encrypted message
	q.Add(id.Hash().B58String(), repo.CafeStoreRequest)

	return id.Hash(), nil
}

//func Store(node *core.IpfsNode, id string, session *repo.CafeSession) error {
//	// load local content
//	cType := "application/octet-stream"
//	var reader io.Reader
//	data, err := ipfsutil.GetDataAtPath(node, id)
//	if err != nil {
//		if err == iface.ErrIsDir {
//			reader, err = ipfsutil.GetArchiveAtPath(node, id)
//			if err != nil {
//				return err
//			}
//			cType = "application/gzip"
//		} else {
//			return err
//		}
//	} else {
//		reader = bytes.NewReader(data)
//	}
//}

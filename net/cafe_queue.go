package net

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"sync"
	"time"
)

// how often to run store job (daemon only)
const kFrequency = time.Minute * 10

// the size of concurrently processed requests
// note: reqs from this group are sent as a group to each cafe
const groupSize = 16

// CafeRequestQueue holds pending cafe store requests
type CafeRequestQueue struct {
	service   func() *CafeService
	node      func() *core.IpfsNode
	datastore repo.Datastore
	mux       sync.Mutex
}

// NewCafeRequestQueue creates a new queue
func NewCafeRequestQueue(service func() *CafeService, node func() *core.IpfsNode, datastore repo.Datastore) *CafeRequestQueue {
	return &CafeRequestQueue{
		service:   service,
		node:      node,
		datastore: datastore,
	}
}

// Put adds a request for each active cafe session
func (q *CafeRequestQueue) Put(target string, rtype repo.CafeRequestType) {
	// get active cafe sessions
	sessions := q.datastore.CafeSessions().List()
	if len(sessions) == 0 {
		return
	}

	// for each session, add a req
	for _, session := range sessions {
		req := &repo.CafeRequest{
			Id:       ksuid.New().String(),
			TargetId: target,
			CafeId:   session.CafeId,
			Type:     rtype,
			Date:     time.Now(),
		}
		if err := q.datastore.CafeRequests().Put(req); err != nil {
			log.Errorf("error adding cafe request %s: %s", target, err)
		}
	}

	// try to flush the queue now
	go q.Flush()
}

// Run starts a job ticker which processes any pending requests
func (q *CafeRequestQueue) Run() {
	tick := time.NewTicker(kFrequency)
	defer tick.Stop()
	go q.Flush()
	for {
		select {
		case <-tick.C:
			go q.Flush()
		}
	}
}

// Flush processes pending requests
func (q *CafeRequestQueue) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()

	// check service status
	if q.service() == nil {
		return
	}

	// start at zero offset
	if err := q.batch(q.datastore.CafeRequests().List("", groupSize)); err != nil {
		log.Errorf("cafe queue batch error: %s", err)
		return
	}
}

// batch flushes a batch of requests
func (q *CafeRequestQueue) batch(reqs []repo.CafeRequest) error {
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
	next := q.datastore.CafeRequests().List(offset, groupSize)

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

// handle handles a single request
func (q *CafeRequestQueue) handle(reqs []repo.CafeRequest, rtype repo.CafeRequestType, cafe peer.ID) ([]string, error) {
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
	}
	return handled, herr
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

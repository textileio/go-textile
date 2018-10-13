package net

import (
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
const groupSize = 20

// CafeStoreRequestQueue holds pending cafe store requests
type CafeStoreRequestQueue struct {
	service   func() *CafeService
	node      func() *core.IpfsNode
	datastore repo.Datastore
	mux       sync.Mutex
}

// NewCafeStoreRequestQueue creates a new queue
func NewCafeStoreRequestQueue(service func() *CafeService, node func() *core.IpfsNode, datastore repo.Datastore) *CafeStoreRequestQueue {
	return &CafeStoreRequestQueue{
		service:   service,
		node:      node,
		datastore: datastore,
	}
}

// Put adds a request for each active cafe session
func (p *CafeStoreRequestQueue) Put(id string) {
	// get active cafe sessions
	sessions := p.datastore.CafeSessions().List()
	if len(sessions) == 0 {
		return
	}

	// for each session, add a req
	for _, session := range sessions {
		req := &repo.CafeStoreRequest{
			Id:     id,
			CafeId: session.CafeId,
			Date:   time.Now(),
		}
		if err := p.datastore.CafeStoreRequests().Put(req); err != nil {
			log.Warningf("store request %s exists for cafe %s", id, session.CafeId)
		}
	}

	// run it now
	go p.Store()
}

// Run starts a job ticker which processes any pending requests
func (p *CafeStoreRequestQueue) Run() {
	tick := time.NewTicker(kFrequency)
	defer tick.Stop()
	go p.Store()
	for {
		select {
		case <-tick.C:
			go p.Store()
		}
	}
}

// Store retrieves and processes pending requests
func (p *CafeStoreRequestQueue) Store() {
	p.mux.Lock()
	defer p.mux.Unlock()

	// check service status
	if p.service() == nil {
		return
	}

	// start at zero offset
	if err := p.handle(p.datastore.CafeStoreRequests().List("", groupSize)); err != nil {
		log.Errorf("store queue handle error: %s", err)
		return
	}
}

// handle handles a batch of requests
func (p *CafeStoreRequestQueue) handle(reqs []repo.CafeStoreRequest) error {
	if len(reqs) == 0 {
		return nil
	}

	// group reqs by cafe
	grps := make(map[string][]string)
	for _, req := range reqs {
		grps[req.CafeId] = append(grps[req.CafeId], req.Id)
	}

	// process concurrently
	var toDelete []string
	wg := sync.WaitGroup{}
	for cafeId, cids := range grps {
		cafe, err := peer.IDB58Decode(cafeId)
		if err != nil {
			continue
		}
		wg.Add(1)
		go func(cids []string, cafe peer.ID) {
			stored, err := p.service().Store(cids, cafe)
			if err != nil {
				log.Errorf("store request to cafe %s failed: %s", cafe.Pretty(), err)
			}
			for _, s := range stored {
				toDelete = append(toDelete, s)
			}
			wg.Done()
		}(cids, cafe)
	}
	wg.Wait()

	// next batch
	offset := reqs[len(reqs)-1].Id
	next := p.datastore.CafeStoreRequests().List(offset, groupSize)

	// clean up
	var deleted []string
	for _, id := range toDelete {
		if err := p.datastore.CafeStoreRequests().Delete(id); err != nil {
			log.Errorf("failed to delete store request %s: %s", id, err)
			continue
		}
		deleted = append(deleted, id)
	}
	log.Debugf("handled %d store requests", len(deleted))

	// keep going
	return p.handle(next)
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

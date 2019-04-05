package core

import (
	"bytes"
	"fmt"
	"sync"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
	mh "gx/ipfs/QmerPMzPk1mJVowm8KgmoknWa4yCYvvugMPsgWmDNUvDLW/go-multihash"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

// cafeOutFlushGroupSize is the size of concurrently processed requests
// note: reqs from this group are batched to each cafe
const cafeOutFlushGroupSize = 32

// cafeReqOpt is an instance helper for creating request options
var cafeReqOpt CafeRequestOption

// RequestHandler is fullfilled by the layer responsible for cafe network requests
//   Desktop and Server => CafeService over libp2p
//   Mobile => Objc and Java SDKs
type RequestHandler interface {
	Flush()
}

// CafeRequestSettings for a request
type CafeRequestSettings struct {
	Size  int
	Group string
}

// CafeRequestOption returns a request setting from an option
type CafeRequestOption func(*CafeRequestSettings)

// Group sets the request's group field
func (CafeRequestOption) Group(val string) CafeRequestOption {
	return func(settings *CafeRequestSettings) {
		settings.Group = val
	}
}

// Size sets the request's size in bytes
func (CafeRequestOption) Size(val int) CafeRequestOption {
	return func(settings *CafeRequestSettings) {
		settings.Size = val
	}
}

// CafeRequestOptions returns request settings from options
func CafeRequestOptions(opts ...CafeRequestOption) *CafeRequestSettings {
	options := &CafeRequestSettings{
		Group: "",
	}

	for _, opt := range opts {
		opt(options)
	}
	return options
}

// CafeRequestGroupStat reports the status of a request group
type CafeRequestGroupStat struct {
	NumTotal    int
	NumComplete int
	SizeTotal   int
	SizComplete int
}

// CafeOutbox queues and processes outbound cafe requests
type CafeOutbox struct {
	service   func() *CafeService
	node      func() *core.IpfsNode
	datastore repo.Datastore
	mux       sync.Mutex
}

// NewCafeOutbox creates a new outbox queue
func NewCafeOutbox(
	service func() *CafeService,
	node func() *core.IpfsNode,
	datastore repo.Datastore,
) *CafeOutbox {
	return &CafeOutbox{
		service:   service,
		node:      node,
		datastore: datastore,
	}
}

// Add adds a request for each active cafe session
func (q *CafeOutbox) Add(target string, rtype pb.CafeRequest_Type, opts ...CafeRequestOption) error {
	pid := q.node().Identity
	settings := CafeRequestOptions(opts...)

	switch rtype {
	case pb.CafeRequest_INBOX:
		return fmt.Errorf("inbox request to own inbox, aborting")
	case pb.CafeRequest_STORE, pb.CafeRequest_UNSTORE:
		if settings.Size == 0 {
			stat, err := ipfs.StatObjectAtPath(q.node(), target)
			if err != nil {
				return err
			}
			settings.Size = stat.BlockSize
		}
	}

	// add a request for each session
	sessions := q.datastore.CafeSessions().List().Items
	for _, session := range sessions {
		// all possible request types are for our own peer
		if err := q.add(pid, target, session.Cafe, rtype, settings); err != nil {
			return err
		}
	}
	return nil
}

// AddForInbox adds a request for a peer's inbox(es)
func (q *CafeOutbox) AddForInbox(pid peer.ID, env *pb.Envelope, inboxes []*pb.Cafe) error {
	if len(inboxes) == 0 {
		return nil
	}

	hash, err := q.prepForInbox(pid, env)
	if err != nil {
		return err
	}

	target := hash.B58String()
	settings := &CafeRequestSettings{
		Group: target,
	}
	for _, inbox := range inboxes {
		if err := q.add(pid, target, inbox, pb.CafeRequest_INBOX, settings); err != nil {
			return err
		}
	}
	return nil
}

// List returns a batch of not complete requests
func (q *CafeOutbox) List(offset string, limit int) *pb.CafeRequestList {
	return q.datastore.CafeRequests().List(offset, limit)
}

// Complete marks a single request as complete, deleting the group if all from its group are complete
func (q *CafeOutbox) Complete(requestId string) error {
	req := q.datastore.CafeRequests().Get(requestId)
	if req == nil {
		return nil
	}
	if err := q.datastore.CafeRequests().Complete(requestId); err != nil {
		return err
	}

	// see if the group can be removed yet
	if q.datastore.CafeRequests().CountByGroup(req.Group) == 0 {
		return q.datastore.CafeRequests().DeleteByGroup(req.Group)
	}
	return nil
}

// StatRequestGroup returns the status of a request group
func (q *CafeOutbox) StatRequestGroup(group string) *pb.CafeRequestGroupStats {
	return q.datastore.CafeRequests().StatGroup(group)
}

// Flush processes pending requests
func (q *CafeOutbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()
	log.Debug("flushing cafe outbox")

	if q.service() == nil {
		return
	}

	if err := q.batch(q.List("", cafeOutFlushGroupSize)); err != nil {
		log.Errorf("cafe outbox batch error: %s", err)
		return
	}
}

// add queues a single request
func (q *CafeOutbox) add(pid peer.ID, target string, cafe *pb.Cafe, rtype pb.CafeRequest_Type, settings *CafeRequestSettings) error {
	log.Debugf("adding cafe %s request for %s to %s: %s",
		rtype.String(), ipfs.ShortenID(pid.Pretty()), ipfs.ShortenID(cafe.Peer), target)

	return q.datastore.CafeRequests().Add(&pb.CafeRequest{
		Id:     ksuid.New().String(),
		Peer:   pid.Pretty(),
		Target: target,
		Cafe:   cafe,
		Type:   rtype,
		Size:   int64(settings.Size),
		Group:  settings.Group,
		Date:   ptypes.TimestampNow(),
	})
}

// batch flushes a batch of requests
func (q *CafeOutbox) batch(reqs *pb.CafeRequestList) error {
	log.Debugf("handling %d cafe requests", len(reqs.Items))
	if len(reqs.Items) == 0 {
		return nil
	}

	// group reqs by cafe
	groups := make(map[string][]*pb.CafeRequest)
	for _, req := range reqs.Items {
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
		go func(cafe peer.ID, reqs []*pb.CafeRequest) {
			// group by type
			types := make(map[pb.CafeRequest_Type][]*pb.CafeRequest)
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
	offset := reqs.Items[len(reqs.Items)-1].Id
	next := q.List(offset, cafeOutFlushGroupSize)

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
func (q *CafeOutbox) handle(reqs []*pb.CafeRequest, rtype pb.CafeRequest_Type, cafe peer.ID) ([]string, error) {
	var handled []string
	var herr error
	switch rtype {

	// store requests are handled in bulk
	case pb.CafeRequest_STORE:
		var cids []string
		for _, req := range reqs {
			cids = append(cids, req.Target)
		}

		stored, err := q.service().Store(cids, cafe)
		for _, s := range stored {
			for _, r := range reqs {
				if r.Target == s {
					handled = append(handled, r.Id)
				}
			}
		}
		if err != nil {
			log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafe.Pretty(), err)
			herr = err
		}

	case pb.CafeRequest_UNSTORE:
		var cids []string
		for _, req := range reqs {
			cids = append(cids, req.Target)
		}

		unstored, err := q.service().Unstore(cids, cafe)
		for _, u := range unstored {
			for _, r := range reqs {
				if r.Target == u {
					handled = append(handled, r.Id)
				}
			}
		}
		if err != nil {
			log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafe.Pretty(), err)
			herr = err
		}

	case pb.CafeRequest_STORE_THREAD:
		for _, req := range reqs {
			thrd := q.datastore.Threads().Get(req.Target)
			if thrd == nil {
				log.Warningf("could not find thread: %s", req.Target)
				handled = append(handled, req.Id)
				continue
			}

			if err := q.service().StoreThread(thrd, cafe); err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafe.Pretty(), err)
				herr = err
				continue
			}
			handled = append(handled, req.Id)
		}

	case pb.CafeRequest_UNSTORE_THREAD:
		for _, req := range reqs {
			if err := q.service().UnstoreThread(req.Target, cafe); err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafe.Pretty(), err)
				herr = err
				continue
			}
			handled = append(handled, req.Id)
		}

	case pb.CafeRequest_INBOX:
		for _, req := range reqs {
			pid, err := peer.IDB58Decode(req.Peer)
			if err != nil {
				herr = err
				continue
			}

			if err := q.service().DeliverMessage(req.Target, pid, req.Cafe); err != nil {
				log.Errorf("cafe %s request to %s failed: %s", rtype.String(), cafe.Pretty(), err)
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
	hash := id.Hash().B58String()

	if err := q.Add(hash, pb.CafeRequest_STORE, cafeReqOpt.Group(hash)); err != nil {
		return nil, err
	}

	return id.Hash(), nil
}

package core

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ipfs/go-ipfs/core"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

// cafeReqOpt is an instance helper for creating request options
var cafeReqOpt CafeRequestOption

// CafeOutboxHandler is fullfilled by the layer responsible for cafe network requests
//   Desktop and Server => CafeService over libp2p
//   Mobile => Objc and Java SDKs
type CafeOutboxHandler interface {
	Flush()
}

// CafeRequestSettings for a request
type CafeRequestSettings struct {
	Size      int
	Group     string
	SyncGroup string
	Cafe      string
}

// Options converts settings back to options
func (s *CafeRequestSettings) Options() []CafeRequestOption {
	return []CafeRequestOption{
		cafeReqOpt.Size(s.Size),
		cafeReqOpt.Group(s.Group),
		cafeReqOpt.SyncGroup(s.SyncGroup),
		cafeReqOpt.Cafe(s.Cafe),
	}
}

// CafeRequestOption returns a request setting from an option
type CafeRequestOption func(*CafeRequestSettings)

// Group sets the request's group field
func (CafeRequestOption) Group(val string) CafeRequestOption {
	return func(settings *CafeRequestSettings) {
		settings.Group = val
	}
}

// SyncGroup sets the request's sync group field
func (CafeRequestOption) SyncGroup(val string) CafeRequestOption {
	return func(settings *CafeRequestSettings) {
		settings.SyncGroup = val
	}
}

// Size sets the request's size in bytes
func (CafeRequestOption) Size(val int) CafeRequestOption {
	return func(settings *CafeRequestSettings) {
		settings.Size = val
	}
}

// Cafe limits the request to a single cafe
func (CafeRequestOption) Cafe(val string) CafeRequestOption {
	return func(settings *CafeRequestSettings) {
		settings.Cafe = val
	}
}

// CafeRequestOptions returns request settings from options
func CafeRequestOptions(opts ...CafeRequestOption) *CafeRequestSettings {
	options := &CafeRequestSettings{
		Group:     ksuid.New().String(),
		SyncGroup: ksuid.New().String(),
	}

	for _, opt := range opts {
		opt(options)
	}
	return options
}

// CafeOutbox queues and processes outbound cafe requests
type CafeOutbox struct {
	node        func() *core.IpfsNode
	datastore   repo.Datastore
	handler     CafeOutboxHandler
	flushBlocks func()
	lock        sync.Mutex
}

// NewCafeOutbox creates a new outbox queue
func NewCafeOutbox(
	node func() *core.IpfsNode,
	datastore repo.Datastore,
	handler CafeOutboxHandler,
	flushBlocks func()) *CafeOutbox {
	return &CafeOutbox{
		node:        node,
		datastore:   datastore,
		handler:     handler,
		flushBlocks: flushBlocks,
	}
}

// Add adds a request for each active cafe session
func (q *CafeOutbox) Add(target string, rtype pb.CafeRequest_Type, opts ...CafeRequestOption) error {
	peerId := q.node().Identity.Pretty()
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
	group := settings.Group
	var err error
	for _, session := range sessions {
		if settings.Cafe != "" && settings.Cafe != session.Id {
			continue
		}
		// all possible request types are for our own peer
		settings.Group = fmt.Sprintf("%s-%s", group, session.Id)
		err = q.add(peerId, target, session.Cafe, rtype, settings)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddForInbox adds a request for a peer's inbox(es)
func (q *CafeOutbox) AddForInbox(peerId string, env *pb.Envelope, inboxes []*pb.Cafe) error {
	if len(inboxes) == 0 {
		return nil
	}

	envb, err := proto.Marshal(env)
	if err != nil {
		return err
	}
	id, err := ipfs.AddData(q.node(), bytes.NewReader(envb), true, false)
	if err != nil {
		return err
	}

	target := id.Hash().B58String()
	settings := CafeRequestOptions()
	for _, inbox := range inboxes {
		err = q.add(peerId, target, inbox, pb.CafeRequest_INBOX, settings)
		if err != nil {
			return err
		}
	}
	return nil
}

// Flush processes pending requests
func (q *CafeOutbox) Flush(skipBlocks bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	log.Debug("flushing cafe outbox")

	if q.handler == nil {
		return
	}
	q.handler.Flush()
	if !skipBlocks {
		q.flushBlocks()
	}

	// clean up
	err := q.datastore.CafeRequests().DeleteCompleteSyncGroups()
	if err != nil {
		log.Error(err.Error())
		return
	}
}

// add queues a single request
func (q *CafeOutbox) add(peerId string, target string, cafe *pb.Cafe, rtype pb.CafeRequest_Type, settings *CafeRequestSettings) error {
	log.Debugf("adding cafe %s request: %s", rtype.String(), target)

	return q.datastore.CafeRequests().Add(&pb.CafeRequest{
		Id:        ksuid.New().String(),
		Peer:      peerId,
		Target:    target,
		Cafe:      cafe,
		Group:     settings.Group,
		SyncGroup: settings.SyncGroup,
		Type:      rtype,
		Date:      ptypes.TimestampNow(),
		Size:      int64(settings.Size),
		Status:    pb.CafeRequest_NEW,
	})
}

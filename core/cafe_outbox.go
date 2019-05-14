package core

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ipfs/go-ipfs/core"
	peer "github.com/libp2p/go-libp2p-peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/crypto"
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

// CafeOutbox queues and processes outbound cafe requests
type CafeOutbox struct {
	node      func() *core.IpfsNode
	datastore repo.Datastore
	handler   CafeOutboxHandler
	mux       sync.Mutex
}

// NewCafeOutbox creates a new outbox queue
func NewCafeOutbox(node func() *core.IpfsNode, datastore repo.Datastore, handler CafeOutboxHandler) *CafeOutbox {
	return &CafeOutbox{
		node:      node,
		datastore: datastore,
		handler:   handler,
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

// Flush processes pending requests
func (q *CafeOutbox) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()
	log.Debug("flushing cafe outbox")

	if q.handler == nil {
		return
	}
	q.handler.Flush()
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

	// TODO: remove pin after req is handled
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

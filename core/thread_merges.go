package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"sort"
	"time"
)

// Merge adds a merge block, which are kept local until subsequent updates, avoiding possibly endless echoes
func (t *Thread) Merge(head mh.Multihash) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build custom merge header
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	// delete author since we want these identical across peers
	header.Author = ""
	// add a second parent
	header.Parents = append(header.Parents, head.B58String())
	// sort to ensure a deterministic (the order may be reversed on other peers)
	sort.Strings(header.Parents)
	// choose newest to use for date
	p1b := t.datastore.Blocks().Get(header.Parents[0])
	if p1b == nil {
		return nil, errors.New("first merge parent not found")
	}
	p2b := t.datastore.Blocks().Get(header.Parents[1])
	if p2b == nil {
		return nil, errors.New("second merge parent not found")
	}
	var date time.Time
	if p1b.Date.Before(p2b.Date) {
		date = p2b.Date
	} else {
		date = p1b.Date
	}
	// add a small amount to date to keep it ahead of both parents
	date = date.Add(time.Millisecond)
	// set content
	pdate, err := ptypes.TimestampProto(date)
	if err != nil {
		return nil, err
	}
	header.Date = pdate

	// build block
	block := &pb.ThreadBlock{
		Header: header,
		Type:   pb.ThreadBlock_MERGE,
	}
	plaintext, err := proto.Marshal(block)
	if err != nil {
		return nil, err
	}

	// add plaintext to ipfs
	hash, err := t.addBlock(plaintext)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(&commitResult{hash: hash, header: header}, repo.MergeBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(hash); err != nil {
		return nil, err
	}

	log.Debugf("adding MERGE to %s: %s", t.Id, hash.B58String())

	// all done
	return hash, nil
}

// HandleMergeBlock handles an incoming merge block
func (t *Thread) HandleMergeBlock(from *peer.ID, hash mh.Multihash, block *pb.ThreadBlock, following bool) error {
	// index it locally
	if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.MergeBlock, nil); err != nil {
		return err
	}

	// back prop
	newPeers, err := t.FollowParents(block.Header.Parents, from)
	if err != nil {
		return err
	}

	// handle HEAD
	if following {
		return nil
	}
	if _, err := t.handleHead(hash, block.Header.Parents); err != nil {
		return err
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return err
		}
	}
	return nil
}

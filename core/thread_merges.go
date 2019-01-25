package core

import (
	"errors"
	"sort"
	"time"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// merge adds a merge block, which are kept local until subsequent updates, avoiding possibly endless echoes
func (t *Thread) merge(head mh.Multihash) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

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
	pdate, err := ptypes.TimestampProto(date)
	if err != nil {
		return nil, err
	}
	header.Date = pdate

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

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: header,
	}, repo.MergeBlock, "", ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(hash); err != nil {
		return nil, err
	}

	log.Debugf("adding MERGE to %s: %s", t.Id, hash.B58String())

	return hash, nil
}

// handleMergeBlock handles an incoming merge block
func (t *Thread) handleMergeBlock(hash mh.Multihash, block *pb.ThreadBlock) error {
	if !t.readable(t.config.Account.Address) {
		return ErrNotReadable
	}
	if !t.readable(block.Header.Address) {
		return ErrNotReadable
	}

	return t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.MergeBlock, "", "")
}

package core

import (
	"fmt"
	"strings"

	"github.com/textileio/textile-go/ipfs"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
)

// AddIgnore adds an outgoing ignore block targeted at another block to ignore
func (t *Thread) AddIgnore(block string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.annotatable(t.config.Account.Address) {
		return nil, ErrNotAnnotatable
	}

	// adding an ignore specific prefix here to ensure future flexibility
	target := fmt.Sprintf("ignore-%s", block)

	msg := &pb.ThreadIgnore{
		Target: target,
	}

	res, err := t.commitBlock(msg, pb.Block_IGNORE, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, pb.Block_IGNORE, target, ""); err != nil {
		return nil, err
	}

	rblock := t.datastore.Blocks().Get(block)
	if err := t.ignoreBlockTarget(rblock); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	// cleanup
	if err := t.datastore.Notifications().DeleteByBlock(block); err != nil {
		return nil, err
	}

	log.Debugf("added IGNORE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleIgnoreBlock handles an incoming ignore block
func (t *Thread) handleIgnoreBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadIgnore, error) {
	msg := new(pb.ThreadIgnore)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}
	if !t.annotatable(block.Header.Address) {
		return nil, ErrNotAnnotatable
	}

	// cleanup
	blockId := strings.Replace(msg.Target, "ignore-", "", 1)
	if err := t.datastore.Notifications().DeleteByBlock(blockId); err != nil {
		return nil, err
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_IGNORE, msg.Target, ""); err != nil {
		return nil, err
	}

	rblock := t.datastore.Blocks().Get(blockId)
	if err := t.ignoreBlockTarget(rblock); err != nil {
		return nil, err
	}

	return msg, nil
}

// ignoreBlockTarget conditionally removes block target and files
func (t *Thread) ignoreBlockTarget(block *pb.Block) error {
	if block == nil || block.Target == "" {
		return nil
	}

	switch block.Type {
	case pb.Block_FILES:
		node, err := ipfs.NodeAtPath(t.node(), block.Target)
		if err != nil {
			return err
		}

		return t.removeFiles(node)
	default:
		return nil
	}
}

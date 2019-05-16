package core

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
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

	res, err := t.commitBlock(msg, pb.Block_IGNORE, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(res, pb.Block_IGNORE, target, "")
	if err != nil {
		return nil, err
	}

	rblock := t.datastore.Blocks().Get(block)
	err = t.ignoreBlockTarget(rblock)
	if err != nil {
		return nil, err
	}

	err = t.updateHead(res.hash)
	if err != nil {
		return nil, err
	}

	err = t.post(res, t.Peers())
	if err != nil {
		return nil, err
	}

	// cleanup
	err = t.datastore.Notifications().DeleteByBlock(block)
	if err != nil {
		return nil, err
	}

	log.Debugf("added IGNORE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleIgnoreBlock handles an incoming ignore block
func (t *Thread) handleIgnoreBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadIgnore, error) {
	msg := new(pb.ThreadIgnore)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
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
	err = t.datastore.Notifications().DeleteByBlock(blockId)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_IGNORE, msg.Target, "")
	if err != nil {
		return nil, err
	}

	rblock := t.datastore.Blocks().Get(blockId)
	err = t.ignoreBlockTarget(rblock)
	if err != nil {
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

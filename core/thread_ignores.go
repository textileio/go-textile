package core

import (
	"strings"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

// AddIgnore adds an outgoing ignore block targeted at another block to ignore
func (t *Thread) AddIgnore(block string) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.annotatable(t.config.Account.Address) {
		return nil, ErrNotAnnotatable
	}

	res, err := t.commitBlock(nil, pb.Block_IGNORE, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_IGNORE,
		Date:   res.header.Date,
		Target: block,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	err = t.ignoreBlockTarget(t.datastore.Blocks().Get(block))
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
func (t *Thread) handleIgnoreBlock(bnode *blockNode, block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	msg := new(pb.ThreadIgnore)
	if block.Payload != nil {
		err := ptypes.UnmarshalAny(block.Payload, msg)
		if err != nil {
			return res, err
		}
	}

	if !t.readable(t.config.Account.Address) {
		return res, ErrNotReadable
	}
	if !t.annotatable(block.Header.Address) {
		return res, ErrNotAnnotatable
	}

	var target string
	if msg.Target != "" {
		target = msg.Target
	} else {
		target = bnode.target
	}

	// cleanup
	target = strings.Replace(target, "ignore-", "", 1)
	err := t.datastore.Notifications().DeleteByBlock(target)
	if err != nil {
		return res, err
	}

	err = t.ignoreBlockTarget(t.datastore.Blocks().Get(target))
	if err != nil {
		return res, err
	}

	res.oldTarget = target
	return res, err
}

// ignoreBlockTarget conditionally ignore the given block
func (t *Thread) ignoreBlockTarget(block *pb.Block) error {
	if block == nil {
		return nil
	}

	switch block.Type {
	case pb.Block_FILES:
		if block.Data == "" {
			return nil
		}

		node, err := ipfs.NodeAtPath(t.node(), block.Data, ipfs.CatTimeout)
		if err != nil {
			return err
		}

		return t.removeFiles(node)
	default:
		return nil
	}
}

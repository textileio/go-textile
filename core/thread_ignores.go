package core

import (
	"fmt"
	"strings"

	"github.com/textileio/textile-go/ipfs"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// AddIgnore adds an outgoing ignore block targeted at another block to ignore
func (t *Thread) AddIgnore(block string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// adding an ignore specific prefix here to ensure future flexibility
	target := fmt.Sprintf("ignore-%s", block)

	msg := &pb.ThreadIgnore{
		Target: target,
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_IGNORE, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, repo.IgnoreBlock, target, ""); err != nil {
		return nil, err
	}

	// unpin
	t.unpinTarget(block)

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

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

	// delete notifications
	blockId := strings.Replace(msg.Target, "ignore-", "", 1)
	if err := t.datastore.Notifications().DeleteByBlock(blockId); err != nil {
		return nil, err
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.IgnoreBlock, msg.Target, ""); err != nil {
		return nil, err
	}

	// unpin
	t.unpinTarget(blockId)

	return msg, nil
}

// unpinTarget unpins block target if present and not part of another thread
func (t *Thread) unpinTarget(blockId string) {
	block := t.datastore.Blocks().Get(blockId)
	if block != nil && block.Target != "" {
		blocks := t.datastore.Blocks().List("", -1, "target='"+block.Target+"'")
		if len(blocks) == 1 {
			// safe to unpin

			switch block.Type {
			case repo.FilesBlock:
				if err := ipfs.UnpinPath(t.node(), block.Target); err != nil {
					log.Warningf("failed to unpin %s: %s", block.Target, err)
				}
			}
		}
	}
}

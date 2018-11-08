package core

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"strings"
)

// Ignore adds an outgoing ignore block targeted at another block to ignore
func (t *Thread) Ignore(blockId string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// adding an ignore specific prefix here to ensure future flexibility
	dataId := fmt.Sprintf("ignore-%s", blockId)

	// build block
	msg := &pb.ThreadIgnore{
		Data: dataId,
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_IGNORE, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	dconf := &repo.FilesBlockConfig{
		DataId: msg.Data,
	}
	if err := t.indexBlock(res, repo.IgnoreBlock, dconf); err != nil {
		return nil, err
	}

	// unpin
	t.unpinBlockData(blockId)

	// update head
	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	// post it
	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	// delete notifications
	if err := t.datastore.Notifications().DeleteByBlock(blockId); err != nil {
		return nil, err
	}

	log.Debugf("added IGNORE to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// handleIgnoreBlock handles an incoming ignore block
func (t *Thread) handleIgnoreBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadIgnore, error) {
	msg := new(pb.ThreadIgnore)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// delete notifications
	blockId := strings.Replace(msg.Data, "ignore-", "", 1)
	if err := t.datastore.Notifications().DeleteByBlock(blockId); err != nil {
		return nil, err
	}

	// index it locally
	dconf := &repo.FilesBlockConfig{
		DataId: msg.Data,
	}
	if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.IgnoreBlock, dconf); err != nil {
		return nil, err
	}

	// unpin
	t.unpinBlockData(blockId)

	return msg, nil
}

// unpinBlockData unpins block data if present and not part of another thread
func (t *Thread) unpinBlockData(blockId string) {
	block := t.datastore.Blocks().Get(blockId)
	if block != nil && block.DataId != "" {
		blocks := t.datastore.Blocks().List("", -1, "dataId='"+block.DataId+"'")
		if len(blocks) == 1 {
			// safe to unpin

			switch block.Type {
			case repo.FilesBlock:
				// unpin image paths
				path := fmt.Sprintf("%s/thumb", block.DataId)
				if err := ipfs.UnpinPath(t.node(), path); err != nil {
					log.Warningf("failed to unpin %s: %s", path, err)
				}
			}
		}
	}
}

package core

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
)

// Flag adds an outgoing flag block targeted at another block to flag
func (t *Thread) Flag(blockId string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// adding a flag specific prefix here to ensure future flexibility
	dataId := fmt.Sprintf("flag-%s", blockId)

	// build block
	msg := &pb.ThreadFlag{
		Data: dataId,
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_FLAG, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId: msg.Data,
	}
	if err := t.indexBlock(res, repo.FlagBlock, dconf); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	// post it
	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added FLAG to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// handleFlagBlock handles an incoming flag block
func (t *Thread) handleFlagBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadFlag, error) {
	msg := new(pb.ThreadFlag)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// TODO: how do we want to handle flags? making visible to UIs would be a good start

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId: msg.Data,
	}
	if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.FlagBlock, dconf); err != nil {
		return nil, err
	}

	return msg, nil
}

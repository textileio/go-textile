package core

import (
	"fmt"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// AddFlag adds an outgoing flag block targeted at another block to flag
func (t *Thread) AddFlag(block string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.annotatable(t.config.Account.Address) {
		return nil, ErrNotAnnotatable
	}

	// adding a flag specific prefix here to ensure future flexibility
	target := fmt.Sprintf("flag-%s", block)

	msg := &pb.ThreadFlag{
		Target: target,
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_FLAG, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, repo.FlagBlock, target, ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added FLAG to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleFlagBlock handles an incoming flag block
func (t *Thread) handleFlagBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadFlag, error) {
	msg := new(pb.ThreadFlag)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}
	if !t.annotatable(block.Header.Address) {
		return nil, ErrNotAnnotatable
	}

	// TODO: how do we want to handle flags? making visible to UIs would be a good start

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.FlagBlock, msg.Target, ""); err != nil {
		return nil, err
	}

	return msg, nil
}

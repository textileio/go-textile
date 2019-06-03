package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
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

	res, err := t.commitBlock(msg, pb.Block_FLAG, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:      res.hash.B58String(),
		Thread:  t.Id,
		Author:  res.header.Author,
		Type:    pb.Block_FLAG,
		Date:    res.header.Date,
		Parents: res.parents,
		Target:  target,
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("added FLAG to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleFlagBlock handles an incoming flag block
func (t *Thread) handleFlagBlock(hash mh.Multihash, block *pb.ThreadBlock) error {
	msg := new(pb.ThreadFlag)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return err
	}

	if !t.readable(t.config.Account.Address) {
		return ErrNotReadable
	}
	if !t.annotatable(block.Header.Address) {
		return ErrNotAnnotatable
	}

	return nil
}

package core

import (
	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// AddLike adds an outgoing like block
func (t *Thread) AddLike(target string) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.annotatable(t.config.Account.Address) {
		return nil, ErrNotAnnotatable
	}

	res, err := t.commitBlock(nil, pb.Block_LIKE, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_LIKE,
		Date:   res.header.Date,
		Target: target,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	log.Debugf("added LIKE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleLikeBlock handles an incoming like block
func (t *Thread) handleLikeBlock(block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	msg := new(pb.ThreadLike)
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

	res.oldTarget = msg.Target
	return res, nil
}

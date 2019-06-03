package core

import (
	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// AddLike adds an outgoing like block
func (t *Thread) AddLike(target string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.annotatable(t.config.Account.Address) {
		return nil, ErrNotAnnotatable
	}

	msg := &pb.ThreadLike{
		Target: target,
	}

	res, err := t.commitBlock(msg, pb.Block_LIKE, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:      res.hash.B58String(),
		Thread:  t.Id,
		Author:  res.header.Author,
		Type:    pb.Block_LIKE,
		Date:    res.header.Date,
		Parents: res.parents,
		Target:  target,
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("added LIKE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleLikeBlock handles an incoming like block
func (t *Thread) handleLikeBlock(hash mh.Multihash, block *pb.ThreadBlock) (string, error) {
	msg := new(pb.ThreadLike)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return "", err
	}

	if !t.readable(t.config.Account.Address) {
		return "", ErrNotReadable
	}
	if !t.annotatable(block.Header.Address) {
		return "", ErrNotAnnotatable
	}

	return msg.Target, nil
}

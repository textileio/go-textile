package core

import (
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
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

	res, err := t.commitBlock(msg, pb.Block_LIKE, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, pb.Block_LIKE, target, ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added LIKE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleLikeBlock handles an incoming like block
func (t *Thread) handleLikeBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadLike, error) {
	msg := new(pb.ThreadLike)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}
	if !t.annotatable(block.Header.Address) {
		return nil, ErrNotAnnotatable
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_LIKE, msg.Target, ""); err != nil {
		return nil, err
	}
	return msg, nil
}

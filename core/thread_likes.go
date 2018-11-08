package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
)

// AddLike adds an outgoing like block
func (t *Thread) AddLike(target string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	msg := &pb.ThreadLike{
		Target: target,
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_LIKE, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.LikeBlock, target, ""); err != nil {
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

	log.Debugf("added LIKE to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// handleLikeBlock handles an incoming like block
func (t *Thread) handleLikeBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadLike, error) {
	msg := new(pb.ThreadLike)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.LikeBlock, msg.Target, ""); err != nil {
		return nil, err
	}
	return msg, nil
}
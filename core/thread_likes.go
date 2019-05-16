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

	err = t.indexBlock(res, pb.Block_LIKE, target, "")
	if err != nil {
		return nil, err
	}

	err = t.updateHead(res.hash)
	if err != nil {
		return nil, err
	}

	err = t.post(res, t.Peers())
	if err != nil {
		return nil, err
	}

	log.Debugf("added LIKE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleLikeBlock handles an incoming like block
func (t *Thread) handleLikeBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadLike, error) {
	msg := new(pb.ThreadLike)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return nil, err
	}

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}
	if !t.annotatable(block.Header.Address) {
		return nil, ErrNotAnnotatable
	}

	err = t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_LIKE, msg.Target, "")
	if err != nil {
		return nil, err
	}
	return msg, nil
}

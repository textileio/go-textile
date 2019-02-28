package core

import (
	"strings"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
)

// AddComment adds an outgoing comment block
func (t *Thread) AddComment(target string, body string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.annotatable(t.config.Account.Address) {
		return nil, ErrNotAnnotatable
	}

	body = strings.TrimSpace(body)
	msg := &pb.ThreadComment{
		Target: target,
		Body:   body,
	}

	res, err := t.commitBlock(msg, pb.Block_COMMENT, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, pb.Block_COMMENT, target, body); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added COMMENT to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleCommentBlock handles an incoming comment block
func (t *Thread) handleCommentBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadComment, error) {
	msg := new(pb.ThreadComment)
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
	}, pb.Block_COMMENT, msg.Target, msg.Body); err != nil {
		return nil, err
	}
	return msg, nil
}

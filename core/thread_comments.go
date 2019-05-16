package core

import (
	"strings"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
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

	res, err := t.commitBlock(msg, pb.Block_COMMENT, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(res, pb.Block_COMMENT, target, body)
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

	log.Debugf("added COMMENT to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleCommentBlock handles an incoming comment block
func (t *Thread) handleCommentBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadComment, error) {
	msg := new(pb.ThreadComment)
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
	}, pb.Block_COMMENT, msg.Target, msg.Body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

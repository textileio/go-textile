package core

import (
	"strings"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// AddComment adds an outgoing comment block
func (t *Thread) AddComment(target string, body string) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

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

	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_COMMENT,
		Date:   res.header.Date,
		Target: target,
		Body:   body,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	log.Debugf("added COMMENT to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleCommentBlock handles an incoming comment block
func (t *Thread) handleCommentBlock(block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	msg := new(pb.ThreadComment)
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
	res.body = msg.Body
	return res, nil
}

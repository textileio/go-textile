package core

import (
	"strings"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// AddMessage adds an outgoing message block
func (t *Thread) AddMessage(target string, body string) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.writable(t.config.Account.Address) {
		return nil, ErrNotWritable
	}

	body = strings.TrimSpace(body)
	msg := &pb.ThreadMessage{
		Body: body,
	}

	res, err := t.commitBlock(msg, pb.Block_TEXT, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_TEXT,
		Date:   res.header.Date,
		Target: target,
		Body:   msg.Body,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	log.Debugf("added MESSAGE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleMessageBlock handles an incoming message block
func (t *Thread) handleMessageBlock(block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	msg := new(pb.ThreadMessage)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return res, err
	}

	if !t.readable(t.config.Account.Address) {
		return res, ErrNotReadable
	}
	if !t.writable(block.Header.Address) {
		return res, ErrNotWritable
	}

	res.body = msg.Body
	return res, nil
}

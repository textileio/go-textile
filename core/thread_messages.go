package core

import (
	"strings"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// AddMessage adds an outgoing message block
func (t *Thread) AddMessage(body string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

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

	err = t.indexBlock(res, pb.Block_TEXT, "", body)
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

	log.Debugf("added MESSAGE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleMessageBlock handles an incoming message block
func (t *Thread) handleMessageBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadMessage, error) {
	msg := new(pb.ThreadMessage)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return nil, err
	}

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}
	if !t.writable(block.Header.Address) {
		return nil, ErrNotWritable
	}

	err = t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_TEXT, "", msg.Body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

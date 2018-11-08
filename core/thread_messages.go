package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
)

// AddMessage adds an outgoing message block
func (t *Thread) AddMessage(body string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	msg := &pb.ThreadMessage{
		Body: body,
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_MESSAGE, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.MessageBlock, "", body); err != nil {
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

	log.Debugf("added MESSAGE to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// handleMessageBlock handles an incoming message block
func (t *Thread) handleMessageBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadMessage, error) {
	msg := new(pb.ThreadMessage)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.MessageBlock, "", msg.Body); err != nil {
		return nil, err
	}
	return msg, nil
}

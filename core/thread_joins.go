package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// join creates an outgoing join block
func (t *Thread) join(inviter string) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	self := t.datastore.Peers().Get(t.node().Identity.Pretty())
	if self == nil {
		return nil, fmt.Errorf("unable to join, no peer for self")
	}

	res, err := t.commitBlock(&pb.ThreadJoin{
		Inviter: inviter,
		Peer:    self,
	}, pb.Block_JOIN, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_JOIN,
		Date:   res.header.Date,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleJoinBlock handles an incoming join block
func (t *Thread) handleJoinBlock(block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	msg := new(pb.ThreadJoin)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return res, err
	}

	if !t.readable(t.config.Account.Address) {
		return res, ErrNotReadable
	}
	if !t.readable(block.Header.Address) {
		return res, ErrNotReadable
	}

	// join's peer _must_ match the sender
	if msg.Peer.Id != block.Header.Author {
		return res, ErrInvalidThreadBlock
	}

	// collect author as an unwelcomed peer
	if msg.Peer != nil {
		err = t.addOrUpdatePeer(msg.Peer, false)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

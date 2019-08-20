package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// announce creates an outgoing announce block
func (t *Thread) Annouce(msg *pb.ThreadAnnounce) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	if msg == nil {
		msg = &pb.ThreadAnnounce{}
	}
	if msg.Peer == nil {
		peer := t.datastore.Peers().Get(t.node().Identity.Pretty())
		if peer == nil {
			return nil, fmt.Errorf("unable to announce, no peer for self")
		}
		msg.Peer = peer
	}

	// do not annouce for other account peers
	if msg.Peer.Address == t.account.Address() && msg.Peer.Id != t.node().Identity.Pretty() {
		return nil, nil
	}

	res, err := t.commitBlock(msg, pb.Block_ANNOUNCE, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_ANNOUNCE,
		Date:   res.header.Date,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	log.Debugf("added ANNOUNCE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleAnnounceBlock handles an incoming announce block
func (t *Thread) handleAnnounceBlock(block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	msg := new(pb.ThreadAnnounce)
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

	// unless this is our account thread, announce's peer _must_ match the sender
	if msg.Peer != nil {
		if t.Id != t.config.Account.Thread && msg.Peer.Id != block.Header.Author {
			return res, ErrInvalidThreadBlock
		}
	}

	// only initiators can change a thread's name
	if msg.Name != "" {
		if t.initiator != block.Header.Address {
			return res, ErrInvalidThreadBlock
		}
	}

	// update author info
	if msg.Peer != nil && msg.Peer.Id != t.node().Identity.Pretty() {
		if t.Id == t.config.Account.Thread && msg.Peer.Id != block.Header.Author {
			err = t.addPeer(msg.Peer)
		} else {
			err = t.addOrUpdatePeer(msg.Peer, false)
		}
		if err != nil {
			return res, err
		}
	}

	// update thread name
	if msg.Name != "" {
		t.Name = msg.Name
		err = t.datastore.Threads().UpdateName(t.Id, msg.Name)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

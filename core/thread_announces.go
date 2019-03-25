package core

import (
	"fmt"

	mh "gx/ipfs/QmerPMzPk1mJVowm8KgmoknWa4yCYvvugMPsgWmDNUvDLW/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
)

// announce creates an outgoing announce block
func (t *Thread) annouce(msg *pb.ThreadAnnounce) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	if msg == nil {
		var err error
		msg, err = t.buildAnnounce()
		if err != nil {
			return nil, err
		}
	}

	// do not annouce for other account peers
	if msg.Peer.Address == t.account.Address() && msg.Peer.Id != t.node().Identity.Pretty() {
		return nil, nil
	}

	res, err := t.commitBlock(msg, pb.Block_ANNOUNCE, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, pb.Block_ANNOUNCE, "", ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added ANNOUNCE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleAnnounceBlock handles an incoming announce block
func (t *Thread) handleAnnounceBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadAnnounce, error) {
	msg := new(pb.ThreadAnnounce)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}
	if !t.readable(block.Header.Address) {
		return nil, ErrNotReadable
	}

	// unless this is our account thread, announce's peer _must_ match the sender
	if msg.Peer != nil {
		if t.Id != t.config.Account.Thread && msg.Peer.Id != block.Header.Author {
			return nil, ErrInvalidThreadBlock
		}
	}

	// only initiators can change a thread's name
	if msg.Name != "" {
		if t.initiator != block.Header.Address {
			return nil, ErrInvalidThreadBlock
		}
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_ANNOUNCE, "", ""); err != nil {
		return nil, err
	}

	// update author info
	if msg.Peer != nil && msg.Peer.Id != t.node().Identity.Pretty() {
		if t.Id == t.config.Account.Thread && msg.Peer.Id != block.Header.Author {
			if err := t.addPeer(msg.Peer); err != nil {
				return nil, err
			}
		} else {
			if err := t.addOrUpdatePeer(msg.Peer); err != nil {
				return nil, err
			}
		}
	}

	// update thread name
	if msg.Name != "" {
		t.Name = msg.Name
		if err := t.datastore.Threads().UpdateName(t.Id, msg.Name); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

// buildAnnounce builds up a Announce block
func (t *Thread) buildAnnounce() (*pb.ThreadAnnounce, error) {
	msg := &pb.ThreadAnnounce{}
	peer := t.datastore.Peers().Get(t.node().Identity.Pretty())
	if peer == nil {
		return nil, fmt.Errorf("unable to announce, no peer for self")
	}
	msg.Peer = peer
	return msg, nil
}

package core

import (
	"fmt"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// announce creates an outgoing announce block
func (t *Thread) annouce() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	msg, err := t.buildAnnounce()
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_ANNOUNCE, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, repo.AnnounceBlock, "", ""); err != nil {
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

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.AnnounceBlock, "", ""); err != nil {
		return nil, err
	}

	// update author info
	if msg.Contact != nil {
		if err := t.addOrUpdateContact(protoContactToRepo(msg.Contact)); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

// buildAnnounce builds up a Announce block
func (t *Thread) buildAnnounce() (*pb.ThreadAnnounce, error) {
	msg := &pb.ThreadAnnounce{}
	contact := t.datastore.Contacts().Get(t.node().Identity.Pretty())
	if contact == nil {
		return nil, fmt.Errorf("unable to announce, no contact for self")
	}
	msg.Contact = repoContactToProto(contact)
	return msg, nil
}

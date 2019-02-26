package core

import (
	"fmt"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
)

// joinInitial creates an outgoing join block for an emtpy thread
func (t *Thread) joinInitial() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	msg, err := t.buildJoin(t.node().Identity.Pretty())
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.Block_JOIN, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, pb.Block_JOIN, "", ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// join creates an outgoing join block
func (t *Thread) join(inviterId peer.ID) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	msg, err := t.buildJoin(inviterId.Pretty())
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.Block_JOIN, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, pb.Block_JOIN, "", ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleJoinBlock handles an incoming join block
func (t *Thread) handleJoinBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadJoin, error) {
	msg := new(pb.ThreadJoin)
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
	}, pb.Block_JOIN, "", ""); err != nil {
		return nil, err
	}

	// collect author as an unwelcomed peer
	if msg.Contact != nil {
		if cjson, err := pbMarshaler.MarshalToString(msg.Contact); err == nil {
			log.Debugf("found contact: %s", cjson)
		}
		if err := t.addOrUpdateContact(msg.Contact); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

// buildJoin builds up a join block
func (t *Thread) buildJoin(inviterId string) (*pb.ThreadJoin, error) {
	msg := &pb.ThreadJoin{
		Inviter: inviterId,
	}
	contact := t.datastore.Contacts().Get(t.node().Identity.Pretty())
	if contact == nil {
		return nil, fmt.Errorf("unable to join, no contact for self")
	}
	msg.Contact = contact
	return msg, nil
}

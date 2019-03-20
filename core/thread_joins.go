package core

import (
	"fmt"

	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
	mh "gx/ipfs/QmerPMzPk1mJVowm8KgmoknWa4yCYvvugMPsgWmDNUvDLW/go-multihash"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
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

	// join's peer _must_ match the sender
	if msg.Peer.Id != block.Header.Author {
		return nil, ErrInvalidThreadBlock
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_JOIN, "", ""); err != nil {
		return nil, err
	}

	// collect author as an unwelcomed peer
	if msg.Peer != nil {
		if cjson, err := pbMarshaler.MarshalToString(msg.Peer); err == nil {
			log.Debugf("found peer: %s", cjson)
		}
		if err := t.addOrUpdatePeer(msg.Peer); err != nil {
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
	p := t.datastore.Peers().Get(t.node().Identity.Pretty())
	if p == nil {
		return nil, fmt.Errorf("unable to join, no peer for self")
	}
	msg.Peer = p
	return msg, nil
}

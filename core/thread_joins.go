package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// JoinInitial creates an outgoing join block for an emtpy thread
func (t *Thread) JoinInitial() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	msg, err := t.buildJoin(t.node().Identity.Pretty(), "")
	if err != nil {
		return nil, err
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_JOIN, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// Join creates an outgoing join block
func (t *Thread) Join(inviterId peer.ID, inviteId string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	msg, err := t.buildJoin(inviterId.Pretty(), inviteId)
	if err != nil {
		return nil, err
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_JOIN, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	// add new peer
	if inviterId.Pretty() != t.node().Identity.Pretty() {
		newPeer := &repo.ThreadPeer{
			Id:       inviterId.Pretty(),
			ThreadId: t.Id,
		}
		if err := t.datastore.ThreadPeers().Add(newPeer); err != nil {
			log.Errorf("error adding peer: %s", err)
		}
	}

	// post it
	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// HandleJoinBlock handles an incoming join block
func (t *Thread) HandleJoinBlock(from *peer.ID, hash mh.Multihash, block *pb.ThreadBlock, joined *repo.ThreadPeer, following bool) (*pb.ThreadJoin, error) {
	msg := new(pb.ThreadJoin)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	newPeers, err := t.FollowParents(block.Header.Parents, from)
	if err != nil {
		return nil, err
	}

	// short circuit if we're traversing history as a new peer
	if following {
		// if a new peer is discovered during back prop, we'll need to send a welcome
		// but not until _after_ HEAD has been updated at the update entry point, where
		// the new peers will be collected
		// NOTE: if from == nil, we've started with an invite, in which case there is
		// no need to handle new peers in this manner (they're sent OUR join)
		if joined != nil && from != nil && joined.Id != from.Pretty() {
			return msg, nil
		}
		return msg, nil
	}

	// send latest direct to this peer if they could use a merge, i.e., we have newer updates
	head, err := t.Head()
	if err != nil {
		return nil, err
	}
	if joined != nil && head != msg.Invite {
		if err := t.sendWelcome(*joined); err != nil {
			return nil, err
		}
	}

	// handle HEAD
	if _, err := t.handleHead(hash, block.Header.Parents); err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

// buildJoin builds up a join block
func (t *Thread) buildJoin(inviterId string, inviteId string) (*pb.ThreadJoin, error) {
	msg := &pb.ThreadJoin{
		Invite:  inviteId,
		Inviter: inviterId,
	}
	username, err := t.datastore.Profile().GetUsername()
	if err != nil {
		return nil, err
	}
	if username != nil {
		msg.Username = *username
	}
	for _, ses := range t.datastore.CafeSessions().List() {
		msg.Inboxes = append(msg.Inboxes, ses.CafeId)
	}
	return msg, nil
}

// welcome sends the latest HEAD block
func (t *Thread) sendWelcome(joined repo.ThreadPeer) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	// get current HEAD
	head, err := t.Head()
	if err != nil {
		return err
	}
	if head == "" {
		return nil
	}
	hash, err := mh.FromB58String(head)
	if err != nil {
		return err
	}

	// download it
	ciphertext, err := ipfs.GetDataAtPath(t.node(), hash.B58String())
	if err != nil {
		return err
	}

	log.Debugf("WELCOME sent to %s at %s", joined.Id, hash.B58String())

	// post it
	return t.post(&commitResult{hash: hash, ciphertext: ciphertext}, []repo.ThreadPeer{joined})
}

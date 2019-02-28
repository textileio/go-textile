package core

import (
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/pb"
)

// AddInvite creates an outgoing invite block, which is sent directly to the recipient
// and does not become part of the hash chain
func (t *Thread) AddInvite(inviteeId peer.ID) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	contact := t.datastore.Contacts().Get(inviteeId.Pretty())
	if contact == nil {
		return nil, ErrContactNotFound
	}

	if !t.shareable(t.config.Account.Address, contact.Address) {
		return nil, ErrNotShareable
	}

	self := t.datastore.Contacts().Get(t.node().Identity.Pretty())
	msg := &pb.ThreadInvite{
		Thread:  t.datastore.Threads().Get(t.Id),
		Inviter: self,
	}

	inviteePk, err := inviteeId.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.Block_INVITE, func(plaintext []byte) ([]byte, error) {
		return crypto.Encrypt(inviteePk, plaintext)
	})
	if err != nil {
		return nil, err
	}

	// create new peer for posting (it will get added if+when they accept)
	target := pb.ThreadPeer{Id: contact.Id}

	if err := t.post(res, []pb.ThreadPeer{target}); err != nil {
		return nil, err
	}

	log.Debugf("sent INVITE to %s for %s", contact.Id, t.Id)

	return res.hash, nil
}

// AddExternalInvite creates an external invite, which can be retrieved by any peer
// and does not become part of the hash chain
func (t *Thread) AddExternalInvite() (mh.Multihash, []byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.shareable(t.config.Account.Address, "") {
		return nil, nil, ErrNotShareable
	}

	self := t.datastore.Contacts().Get(t.node().Identity.Pretty())
	msg := &pb.ThreadInvite{
		Thread:  t.datastore.Threads().Get(t.Id),
		Inviter: self,
	}

	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, nil, err
	}

	res, err := t.commitBlock(msg, pb.Block_INVITE, func(plaintext []byte) ([]byte, error) {
		return crypto.EncryptAES(plaintext, key)
	})
	if err != nil {
		return nil, nil, err
	}

	go t.cafeOutbox.Flush()

	log.Debugf("created external INVITE for %s", t.Id)

	return res.hash, key, nil
}

// handleInviteMessage handles an incoming invite.
// This happens right before a join. The invite is not kept on-chain,
// so we only need to follow parents and update HEAD.
func (t *Thread) handleInviteMessage(block *pb.ThreadBlock) error {
	if err := t.followParents(block.Header.Parents); err != nil {
		return err
	}

	// update HEAD if parents of the invite are actual updates
	if len(block.Header.Parents) > 0 {
		hash, err := mh.FromB58String(block.Header.Parents[0])
		if err != nil {
			return err
		}
		if err := t.updateHead(hash); err != nil {
			return err
		}
	}
	return nil
}

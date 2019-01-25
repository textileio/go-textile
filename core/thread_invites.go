package core

import (
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
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

	threadSk, err := t.privKey.Bytes()
	if err != nil {
		return nil, err
	}
	self := t.datastore.Contacts().Get(t.node().Identity.Pretty())
	msg := &pb.ThreadInvite{
		Sk:        threadSk,
		Name:      t.Name,
		Schema:    t.schemaId,
		Initiator: t.initiator,
		Contact:   repoContactToProto(self),
		Type:      int32(t.ttype),
		Sharing:   int32(t.sharing),
		Members:   t.members,
	}

	inviteePk, err := inviteeId.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_INVITE, func(plaintext []byte) ([]byte, error) {
		return crypto.Encrypt(inviteePk, plaintext)
	})
	if err != nil {
		return nil, err
	}

	// create new peer for posting (it will get added if+when they accept)
	target := repo.ThreadPeer{Id: inviteeId.Pretty()}

	if err := t.post(res, []repo.ThreadPeer{target}); err != nil {
		return nil, err
	}

	log.Debugf("sent INVITE to %s for %s", inviteeId.Pretty(), t.Id)

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

	threadSk, err := t.privKey.Bytes()
	if err != nil {
		return nil, nil, err
	}
	contact := t.datastore.Contacts().Get(t.node().Identity.Pretty())
	msg := &pb.ThreadInvite{
		Sk:        threadSk,
		Name:      t.Name,
		Schema:    t.schemaId,
		Initiator: t.initiator,
		Contact:   repoContactToProto(contact),
		Type:      int32(t.ttype),
		Sharing:   int32(t.sharing),
		Members:   t.members,
	}

	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, nil, err
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_INVITE, func(plaintext []byte) ([]byte, error) {
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

package core

import (
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// AddInvite creates an outgoing invite block, which is sent directly to the recipient
// and does not become part of the hash chain
func (t *Thread) AddInvite(inviteeId peer.ID) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	threadSk, err := t.privKey.Bytes()
	if err != nil {
		return nil, err
	}
	msg := &pb.ThreadInvite{
		Sk:   threadSk,
		Name: t.Name,
	}

	// get the peer pub key from the id
	inviteePk, err := inviteeId.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_INVITE, func(plaintext []byte) ([]byte, error) {
		return crypto.Encrypt(inviteePk, plaintext)
	})
	if err != nil {
		return nil, err
	}

	// create new peer for posting (it will get added if+when they accept)
	target := repo.ThreadPeer{Id: inviteeId.Pretty()}

	// post it
	if err := t.post(res, []repo.ThreadPeer{target}); err != nil {
		return nil, err
	}

	log.Debugf("sent INVITE to %s for %s", inviteeId.Pretty(), t.Id)

	// all done
	return res.hash, nil
}

// AddExternalInvite creates an external invite, which can be retrieved by any peer
// and does not become part of the hash chain
func (t *Thread) AddExternalInvite() (mh.Multihash, []byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	threadSk, err := t.privKey.Bytes()
	if err != nil {
		return nil, nil, err
	}
	msg := &pb.ThreadInvite{
		Sk:   threadSk,
		Name: t.Name,
	}

	// generate an aes key
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, nil, err
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_INVITE, func(plaintext []byte) ([]byte, error) {
		return crypto.EncryptAES(plaintext, key)
	})
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("created external INVITE for %s", t.Id)

	// all done
	return res.hash, key, nil
}

// HandleInviteMessage handles an incoming invite
// - this happens right before a join
// - the invite is not kept on-chain, so we only need to follow parents and update HEAD
func (t *Thread) HandleInviteMessage(block *pb.ThreadBlock) error {
	// back prop
	if _, err := t.FollowParents(block.Header.Parents, nil); err != nil {
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

package thread

import (
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// AddInvite creates an outgoing invite block, which is sent directly to the recipient
// and does not become part of the hash chain
func (t *Thread) AddInvite(inviteePk libp2pc.PubKey) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// get the peer id from the pub key
	inviteeId, err := peer.IDFromPublicKey(inviteePk)
	if err != nil {
		return nil, err
	}

	// encypt thread secret with the recipient's public key
	threadSk, err := t.PrivKey.Bytes()
	if err != nil {
		return nil, err
	}
	threadSkCipher, err := crypto.Encrypt(inviteePk, threadSk)
	if err != nil {
		return nil, err
	}

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadInvite{
		Header:        header,
		SkCipher:      threadSkCipher,
		SuggestedName: t.Name,
		InviteeId:     inviteeId.Pretty(),
	}

	// commit to ipfs
	message, addr, err := t.commitBlock(content, pb.Message_THREAD_INVITE)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// create new peer for posting, but don't add it yet. it will get added if+when they accept.
	inviteePkb, err := inviteePk.Bytes()
	if err != nil {
		return nil, err
	}
	target := repo.Peer{
		Id:     inviteeId.Pretty(),
		PubKey: inviteePkb,
	}

	// post it
	t.post(message, id, []repo.Peer{target})

	log.Debugf("sent INVITE to %s for %s", inviteeId.Pretty(), t.Id)

	// all done
	return addr, nil
}

// HandleInviteMessage handles an incoming invite
// - this happens right before a join
// - the invite is not kept on-chain, so we only need to follow parents and update HEAD
func (t *Thread) HandleInviteMessage(content *pb.ThreadInvite) error {
	// back prop
	if _, err := t.FollowParents(content.Header.Parents, nil); err != nil {
		return err
	}

	// update HEAD if parents of the invite are actual updates
	if len(content.Header.Parents) > 0 {
		if err := t.updateHead(content.Header.Parents[0]); err != nil {
			return err
		}
	}

	return nil
}

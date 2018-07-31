package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// AddInvite creates an outgoing invite block
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

	// index it locally
	if err := t.indexBlock(id, header, repo.InviteBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// create new peer for posting, but don't add it yet. it will get added if+when they accept.
	inviteePkb, err := inviteePk.Bytes()
	if err != nil {
		return nil, err
	}
	peers := []repo.Peer{{
		Id:     inviteeId.Pretty(),
		PubKey: inviteePkb,
	}}
	for _, p := range t.Peers() {
		if p.Id != inviteeId.Pretty() {
			peers = append(peers, p)
		}
	}

	// post it
	t.post(message, id, peers)

	log.Debugf("added invite to %s for %s: %s", t.Id, inviteeId.Pretty(), id)

	// all done
	return addr, nil
}

// HandleInviteBlock handles an incoming invite block
func (t *Thread) HandleInviteBlock(message *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadInvite, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadInvite)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, err
		}
	}

	// add to ipfs
	addr, err := t.addBlock(message)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// check if we aleady have this block indexed
	// (should only happen if a misbehaving peer keeps sending the same block)
	index := t.blocks().Get(id)
	if index != nil {
		return nil, nil
	}

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.InviteBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	if err := t.FollowParents(content.Header.Parents); err != nil {
		return nil, err
	}

	// handle HEAD
	if following {
		return addr, nil
	}
	if _, err := t.handleHead(id, content.Header.Parents, false); err != nil {
		return nil, err
	}

	return addr, nil
}

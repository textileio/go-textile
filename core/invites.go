package core

import (
	"errors"
	"fmt"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/util"
)

// ErrThreadInviteNotFound indicates thread invite is not found
var ErrThreadInviteNotFound = errors.New("thread invite not found")

// InviteView gets a pending invite as a view object, which does not include the block payload
func (t *Textile) InviteView(invite *pb.Invite) *pb.InviteView {
	if invite == nil {
		return nil
	}

	view := &pb.InviteView{
		Id:   invite.Id,
		Name: invite.Name,
		Date: invite.Date,
	}

	ex := t.datastore.Contacts().Get(invite.Inviter.Id)
	if ex != nil && (invite.Inviter == nil || util.ProtoTsIsNewer(ex.Updated, invite.Inviter.Updated)) {
		view.Inviter = t.User(ex.Id)
	} else if invite.Inviter != nil {
		view.Inviter = &pb.User{
			Address: invite.Inviter.Address,
			Name:    invite.Inviter.Username,
			Avatar:  invite.Inviter.Avatar,
		}
	}

	return view
}

// Invites lists info on all pending invites
func (t *Textile) Invites() *pb.InviteViewList {
	list := &pb.InviteViewList{Items: make([]*pb.InviteView, 0)}

	for _, invite := range t.datastore.Invites().List().Items {
		view := t.InviteView(invite)
		list.Items = append(list.Items, view)
	}

	return list
}

// AcceptInvite adds a new thread, and notifies the inviter of the join
func (t *Textile) AcceptInvite(inviteId string) (mh.Multihash, error) {
	invite := t.datastore.Invites().Get(inviteId)
	if invite == nil {
		return nil, ErrThreadInviteNotFound
	}

	hash, err := t.handleThreadInvite(invite.Block)
	if err != nil {
		return nil, err
	}

	if err := t.IgnoreInvite(inviteId); err != nil {
		return nil, err
	}

	return hash, nil
}

// AcceptExternalInvite attemps to download an encrypted thread key from an external invite,
// adds a new thread, and notifies the inviter of the join
func (t *Textile) AcceptExternalInvite(inviteId string, key []byte) (mh.Multihash, error) {
	ciphertext, err := ipfs.DataAtPath(t.node, fmt.Sprintf("%s", inviteId))
	if err != nil {
		return nil, err
	}

	// attempt decrypt w/ key
	plaintext, err := crypto.DecryptAES(ciphertext, key)
	if err != nil {
		return nil, ErrInvalidThreadBlock
	}
	return t.handleThreadInvite(plaintext)
}

// IgnoreInvite deletes the invite and removes the associated notification.
func (t *Textile) IgnoreInvite(inviteId string) error {
	if err := t.datastore.Invites().Delete(inviteId); err != nil {
		return err
	}
	return t.datastore.Notifications().DeleteByBlock(inviteId)
}

// handleThreadInvite uses an invite block to join a thread
func (t *Textile) handleThreadInvite(plaintext []byte) (mh.Multihash, error) {
	block := new(pb.ThreadBlock)
	if err := proto.Unmarshal(plaintext, block); err != nil {
		return nil, err
	}
	if block.Type != pb.Block_INVITE {
		return nil, ErrInvalidThreadBlock
	}
	msg := new(pb.ThreadInvite)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}
	if msg.Thread == nil || msg.Inviter == nil {
		return nil, ErrInvalidThreadBlock
	}

	// check if we're allowed to get an invite
	// Note: just using a dummy thread here because having these access+sharing
	// methods on Thread is very nice elsewhere.
	dummy := &Thread{
		initiator: msg.Thread.Initiator,
		ttype:     msg.Thread.Type,
		sharing:   msg.Thread.Sharing,
		members:   msg.Thread.Members,
	}
	if !dummy.shareable(msg.Inviter.Address, t.config.Account.Address) {
		return nil, ErrNotShareable
	}

	sk, err := libp2pc.UnmarshalPrivateKey(msg.Thread.Sk)
	if err != nil {
		return nil, err
	}

	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	if thrd := t.Thread(id.Pretty()); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	config := pb.AddThreadConfig{
		Key:  msg.Thread.Key,
		Name: msg.Thread.Name,
		Schema: &pb.AddThreadConfig_Schema{
			Id: msg.Thread.Schema,
		},
		Type:    msg.Thread.Type,
		Sharing: msg.Thread.Sharing,
		Members: msg.Thread.Members,
		Force:   true,
	}
	thrd, err := t.AddThread(config, sk, msg.Thread.Initiator, false)
	if err != nil {
		return nil, err
	}

	if err := thrd.addOrUpdateContact(msg.Inviter); err != nil {
		return nil, err
	}

	// follow parents, update head
	if err := thrd.handleInviteMessage(block); err != nil {
		return nil, err
	}

	// mark any discovered peers as welcomed
	// there's no need to send a welcome because we're about to send a join message
	if err := t.datastore.ThreadPeers().WelcomeByThread(thrd.Id); err != nil {
		return nil, err
	}

	// join the thread
	author, err := peer.IDB58Decode(block.Header.Author)
	if err != nil {
		return nil, err
	}
	hash, err := thrd.join(author)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

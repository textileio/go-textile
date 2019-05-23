package core

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/mr-tron/base58/base58"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

// ErrContactNotFound indicates a local contact was not found
var ErrContactNotFound = fmt.Errorf("contact not found")

// ErrThreadInviteNotFound indicates thread invite is not found
var ErrThreadInviteNotFound = fmt.Errorf("thread invite not found")

// AddInvite creates an invite for each of the target address's peers
func (t *Textile) AddInvite(threadId string, address string) error {
	thrd := t.Thread(threadId)
	if thrd == nil {
		return ErrThreadNotFound
	}

	peers := t.datastore.Peers().List(fmt.Sprintf("address='%s'", address))
	if len(peers) == 0 {
		return ErrContactNotFound
	}

	var err error
	for _, p := range peers {
		_, err = thrd.AddInvite(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddExternalInvite generates a new external invite link to a thread
func (t *Textile) AddExternalInvite(threadId string) (*pb.ExternalInvite, error) {
	thrd := t.Thread(threadId)
	if thrd == nil {
		return nil, ErrThreadNotFound
	}

	hash, key, err := thrd.AddExternalInvite()
	if err != nil {
		return nil, err
	}

	return &pb.ExternalInvite{
		Id:      hash.B58String(),
		Key:     base58.FastBase58Encoding(key),
		Inviter: t.account.Address(),
	}, nil
}

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

	x := t.datastore.Peers().Get(invite.Inviter.Id)
	if x != nil && (invite.Inviter == nil || util.ProtoTsIsNewer(x.Updated, invite.Inviter.Updated)) {
		view.Inviter = t.PeerUser(x.Id)
	} else if invite.Inviter != nil {
		view.Inviter = &pb.User{
			Address: invite.Inviter.Address,
			Name:    invite.Inviter.Name,
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
func (t *Textile) AcceptInvite(id string) (mh.Multihash, error) {
	invite := t.datastore.Invites().Get(id)
	if invite == nil {
		return nil, ErrThreadInviteNotFound
	}

	hash, err := t.handleThreadAdd(invite.Block)
	if err != nil {
		return nil, err
	}

	if err := t.IgnoreInvite(id); err != nil {
		return nil, err
	}

	return hash, nil
}

// AcceptExternalInvite attemps to download an encrypted thread key from an external invite,
// adds a new thread, and notifies the inviter of the join
func (t *Textile) AcceptExternalInvite(id string, key []byte) (mh.Multihash, error) {
	ciphertext, err := ipfs.DataAtPath(t.node, fmt.Sprintf("%s", id))
	if err != nil {
		return nil, err
	}

	// attempt decrypt w/ key
	plaintext, err := crypto.DecryptAES(ciphertext, key)
	if err != nil {
		return nil, ErrInvalidThreadBlock
	}
	return t.handleThreadAdd(plaintext)
}

// IgnoreInvite deletes the invite and removes the associated notification.
func (t *Textile) IgnoreInvite(id string) error {
	err := t.datastore.Invites().Delete(id)
	if err != nil {
		return err
	}
	return t.datastore.Notifications().DeleteByBlock(id)
}

// handleThreadAdd uses an add block to join a thread
func (t *Textile) handleThreadAdd(plaintext []byte) (mh.Multihash, error) {
	block := new(pb.ThreadBlock)
	err := proto.Unmarshal(plaintext, block)
	if err != nil {
		return nil, err
	}
	if block.Type != pb.Block_ADD {
		return nil, ErrInvalidThreadBlock
	}
	msg := new(pb.ThreadAdd)
	err = ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
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
		whitelist: msg.Thread.Whitelist,
	}
	if !dummy.shareable(msg.Inviter.Address, t.config.Account.Address) {
		return nil, ErrNotShareable
	}

	sk, err := ipfs.UnmarshalPrivateKey(msg.Thread.Sk)
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
		Type:      msg.Thread.Type,
		Sharing:   msg.Thread.Sharing,
		Whitelist: msg.Thread.Whitelist,
		Force:     true,
	}
	thrd, err := t.AddThread(config, sk, msg.Thread.Initiator, false, !t.isAccountPeer(msg.Inviter.Id))
	if err != nil {
		return nil, err
	}

	err = thrd.addOrUpdatePeer(msg.Inviter)
	if err != nil {
		return nil, err
	}

	// follow parents, update head
	err = thrd.handleAddBlock(block)
	if err != nil {
		return nil, err
	}

	// mark any discovered peers as welcomed
	// there's no need to send a welcome because we're about to send a join message
	err = t.datastore.ThreadPeers().WelcomeByThread(thrd.Id)
	if err != nil {
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

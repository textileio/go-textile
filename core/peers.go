package core

import (
	"fmt"

	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

// PeerUser returns a user object with the most recently updated contact for the given id
// Note: If no underlying contact is found, this will return a blank object w/ a
// generic user name for display-only purposes.
func (t *Textile) PeerUser(id string) *pb.User {
	user := t.datastore.Peers().GetBestUser(id)
	if user != nil {
		return user
	}
	return &pb.User{
		Name: ipfs.ShortenID(id),
	}
}

// AddPeer adds or updates a peer
func (t *Textile) AddPeer(peer *pb.Peer) error {
	x := t.datastore.Peers().Get(peer.Id)
	if x != nil && (peer.Updated == nil || util.ProtoTsIsNewer(x.Updated, peer.Updated)) {
		return nil
	}

	// peer is new / newer, update
	err := t.datastore.Peers().AddOrUpdate(peer)
	if err != nil {
		return err
	}

	if x == nil && peer.Address == t.account.Address() {
		t.sendUpdate(&pb.AccountUpdate{
			Id:   peer.Id,
			Type: pb.AccountUpdate_ACCOUNT_PEER_ADDED,
		})
	}

	// ensure new update is actually different before announcing to account
	if x != nil {
		if peersEqual(x, peer) {
			return nil
		}
	}

	thrd := t.AccountThread()
	if thrd == nil {
		return fmt.Errorf("account thread not found")
	}

	_, err = thrd.Annouce(&pb.ThreadAnnounce{Peer: peer})
	if err != nil {
		return err
	}
	return nil
}

// PublishPeer publishes this peer's info to the cafe network
func (t *Textile) PublishPeer() error {
	self := t.datastore.Peers().Get(t.node.Identity.Pretty())
	if self == nil {
		return nil
	}

	sessions := t.datastore.CafeSessions().List().Items
	if len(sessions) == 0 {
		return nil
	}
	for _, session := range sessions {
		if err := t.cafe.PublishPeer(self, session.Id); err != nil {
			return err
		}
	}
	return nil
}

// UpdatePeerInboxes sets own peer inboxes from the current cafe sessions
func (t *Textile) UpdatePeerInboxes() error {
	var inboxes []*pb.Cafe
	for _, session := range t.datastore.CafeSessions().List().Items {
		inboxes = append(inboxes, session.Cafe)
	}
	return t.datastore.Peers().UpdateInboxes(t.node.Identity.Pretty(), inboxes)
}

// peersEqual returns whether or not the two peers are identical
// Note: this does not consider Peer.Created or Peer.Updated
func peersEqual(a *pb.Peer, b *pb.Peer) bool {
	if a.Id != b.Id {
		return false
	}
	if a.Address != b.Address {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Avatar != b.Avatar {
		return false
	}
	if len(a.Inboxes) != len(b.Inboxes) {
		return false
	}
	ac := make(map[string]*pb.Cafe)
	for _, c := range a.Inboxes {
		ac[c.Peer] = c
	}
	for _, j := range b.Inboxes {
		i, ok := ac[j.Peer]
		if !ok {
			return false
		}
		if !cafesEqual(i, j) {
			return false
		}
	}
	return true
}

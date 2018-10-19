package core

import (
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/thread"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"time"
)

// AccountPeers lists all peers
func (t *Textile) AccountPeers() []repo.AccountPeer {
	return t.datastore.AccountPeers().List("")
}

// AddAccountPeer creates an invite for every current and future thread
func (t *Textile) AddAccountPeer(pid peer.ID, name string) error {
	if !t.Online() {
		return ErrOffline
	}

	// index a new peer
	mod := &repo.AccountPeer{
		Id:   pid.Pretty(),
		Name: name,
	}
	if err := t.datastore.AccountPeers().Add(mod); err != nil {
		return err
	}
	log.Infof("added account peer '%s'", pid.Pretty())

	// invite peer to existing threads
	for _, thrd := range t.threads {
		if _, err := thrd.AddInvite(pid); err != nil {
			return err
		}
	}

	// notify listeners
	t.sendUpdate(Update{Id: mod.Id, Name: mod.Name, Type: AccountPeerAdded})

	// send notification
	id, err := t.Id()
	if err != nil {
		return err
	}
	notification := &repo.Notification{
		Id:            ksuid.New().String(),
		Date:          time.Now(),
		ActorId:       id.Pretty(),
		ActorUsername: "You",
		Subject:       mod.Name,
		SubjectId:     mod.Id,
		Type:          repo.AccountPeerAddedNotification,
		Body:          "paired account with a new peer",
	}
	return t.sendNotification(notification)
}

// InviteAccountPeers sends a thread invite to all peers
func (t *Textile) InviteAccountPeers(thrd *thread.Thread) error {
	for _, ap := range t.AccountPeers() {
		id, err := peer.IDB58Decode(ap.Id)
		if err != nil {
			return err
		}
		if _, err := thrd.AddInvite(id); err != nil {
			return err
		}
	}
	return nil
}

// RemoveAccountPeer removes a peer
func (t *Textile) RemoveAccountPeer(id string) error {
	if !t.Online() {
		return ErrOffline
	}

	// delete model
	mod := t.datastore.AccountPeers().Get(id)
	if mod == nil {
		return errors.New("peer not found")
	}
	if err := t.datastore.AccountPeers().Delete(id); err != nil {
		return err
	}

	// delete notifications
	if err := t.datastore.Notifications().DeleteBySubject(mod.Id); err != nil {
		return err
	}

	log.Infof("removed peer '%s'", id)

	// notify listeners
	t.sendUpdate(Update{Id: mod.Id, Name: mod.Name, Type: AccountPeerRemoved})

	return nil
}

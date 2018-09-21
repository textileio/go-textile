package core

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
)

// GetNotifications lists notifications
func (t *Textile) GetNotifications(offset string, limit int) []repo.Notification {
	return t.datastore.Notifications().List(offset, limit, "")
}

// CountUnreadNotifications counts unread notifications
func (t *Textile) CountUnreadNotifications() int {
	return t.datastore.Notifications().CountUnread()
}

// ReadNotification marks a notification as read
func (t *Textile) ReadNotification(id string) error {
	return t.datastore.Notifications().Read(id)
}

// ReadAllNotifications marks all notification as read
func (t *Textile) ReadAllNotifications() error {
	return t.datastore.Notifications().ReadAll()
}

// AcceptThreadInviteViaNotification uses an invite notification to accept an invite to a thread
func (t *Textile) AcceptThreadInviteViaNotification(id string) (mh.Multihash, error) {
	// look up notification
	notification := t.datastore.Notifications().Get(id)
	if notification == nil {
		return nil, errors.New(fmt.Sprintf("could not find notification: %s", id))
	}
	if notification.Type != repo.ReceivedInviteNotification {
		return nil, errors.New(fmt.Sprintf("notification not invite type"))
	}

	// block is the invite's block id
	return t.AcceptThreadInvite(notification.BlockId)
}

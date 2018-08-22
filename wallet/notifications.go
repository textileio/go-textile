package wallet

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/repo"
)

// GetNotifications lists notifications
func (w *Wallet) GetNotifications(offset string, limit int) []repo.Notification {
	return w.datastore.Notifications().List(offset, limit, "")
}

// CountUnreadNotifications counts unread notifications
func (w *Wallet) CountUnreadNotifications() int {
	return w.datastore.Notifications().CountUnread()
}

// ReadNotification marks a notification as read
func (w *Wallet) ReadNotification(id string) error {
	return w.datastore.Notifications().Read(id)
}

// ReadAllNotifications marks all notification as read
func (w *Wallet) ReadAllNotifications() error {
	return w.datastore.Notifications().ReadAll()
}

// AcceptThreadInviteViaNotification uses an invite notification to accept an invite to a thread
func (w *Wallet) AcceptThreadInviteViaNotification(id string) (*string, error) {
	// look up notification
	notification := w.datastore.Notifications().Get(id)
	if notification == nil {
		return nil, errors.New(fmt.Sprintf("could not find notification: %s", id))
	}
	if notification.Type != repo.ReceivedInviteNotification {
		return nil, errors.New(fmt.Sprintf("notification not invite type"))
	}

	// target id is the invite's block id
	return w.AcceptThreadInvite(notification.TargetId)
}

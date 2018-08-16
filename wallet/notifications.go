package wallet

import "github.com/textileio/textile-go/repo"

// GetNotifications lists notifications
func (w *Wallet) GetNotifications(offset string, limit int, query string) []repo.Notification {
	return w.datastore.Notifications().List(offset, limit, query)
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

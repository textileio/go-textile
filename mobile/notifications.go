package mobile

import (
	"github.com/textileio/textile-go/repo"
)

// Notifications call core Notifications
func (m *Mobile) Notifications(offset string, limit int) (string, error) {
	notes := make([]repo.Notification, 0)
	notes = m.node.Notifications(offset, limit)
	return toJSON(notes)
}

// CountUnreadNotifications calls core CountUnreadNotifications
func (m *Mobile) CountUnreadNotifications() int {
	return m.node.CountUnreadNotifications()
}

// ReadNotification calls core ReadNotification
func (m *Mobile) ReadNotification(id string) error {
	return m.node.ReadNotification(id)
}

// ReadAllNotifications calls core ReadAllNotifications
func (m *Mobile) ReadAllNotifications() error {
	return m.node.ReadAllNotifications()
}

// AcceptThreadInviteViaNotification call core AcceptThreadInviteViaNotification
func (m *Mobile) AcceptThreadInviteViaNotification(id string) (string, error) {
	addr, err := m.node.AcceptThreadInviteViaNotification(id)
	if err != nil {
		return "", err
	}
	return addr.B58String(), nil
}

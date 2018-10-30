package mobile

import (
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
)

// Notifications is a wrapper around a list of Notifications
type Notifications struct {
	Items []repo.Notification `json:"items"`
}

// Notifications call core Notifications
func (m *Mobile) Notifications(offset string, limit int) (string, error) {
	notes := Notifications{Items: make([]repo.Notification, 0)}
	fetched := core.Node.Notifications(offset, limit)
	if len(fetched) > 0 {
		notes.Items = fetched
	}
	return toJSON(notes)
}

// CountUnreadNotifications calls core CountUnreadNotifications
func (m *Mobile) CountUnreadNotifications() int {
	return core.Node.CountUnreadNotifications()
}

// ReadNotification calls core ReadNotification
func (m *Mobile) ReadNotification(id string) error {
	return core.Node.ReadNotification(id)
}

// ReadAllNotifications calls core ReadAllNotifications
func (m *Mobile) ReadAllNotifications() error {
	return core.Node.ReadAllNotifications()
}

// AcceptThreadInviteViaNotification call core AcceptThreadInviteViaNotification
func (m *Mobile) AcceptThreadInviteViaNotification(id string) (string, error) {
	addr, err := core.Node.AcceptThreadInviteViaNotification(id)
	if err != nil {
		return "", err
	}
	return addr.B58String(), nil
}

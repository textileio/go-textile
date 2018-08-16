package mobile

import (
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
)

// Notifications is a wrapper around a list of Notifications
type Notifications struct {
	Items []repo.Notification `json:"items"`
}

// GetNotifications call core GetNotifications
func (m *Mobile) GetNotifications(offset string, limit int) (string, error) {
	notes := Notifications{Items: make([]repo.Notification, 0)}
	fetched := core.Node.Wallet.GetNotifications(offset, limit)
	if len(fetched) > 0 {
		notes.Items = fetched
	}
	return toJSON(notes)
}

// CountUnreadNotifications calls core CountUnreadNotifications
func (m *Mobile) CountUnreadNotifications() int {
	return core.Node.Wallet.CountUnreadNotifications()
}

// ReadNotification calls core ReadNotification
func (m *Mobile) ReadNotification(id string) error {
	return core.Node.Wallet.ReadNotification(id)
}

// ReadAllNotifications calls core ReadAllNotifications
func (m *Mobile) ReadAllNotifications() error {
	return core.Node.Wallet.ReadAllNotifications()
}

// AcceptThreadInviteViaNotification call core AcceptThreadInviteViaNotification
func (m *Mobile) AcceptThreadInviteViaNotification(id string) (string, error) {
	addr, err := core.Node.Wallet.AcceptThreadInviteViaNotification(id)
	if err != nil {
		return "", err
	}
	return addr.B58String(), nil
}

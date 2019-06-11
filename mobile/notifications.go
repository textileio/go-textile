package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
)

// Notifications call core Notifications
func (m *Mobile) Notifications(offset string, limit int) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.Notifications(offset, limit))
}

// CountUnreadNotifications calls core CountUnreadNotifications
func (m *Mobile) CountUnreadNotifications() int {
	if !m.node.Started() {
		return 0
	}

	return m.node.CountUnreadNotifications()
}

// ReadNotification calls core ReadNotification
func (m *Mobile) ReadNotification(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.ReadNotification(id)
}

// ReadAllNotifications calls core ReadAllNotifications
func (m *Mobile) ReadAllNotifications() error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.ReadAllNotifications()
}

// AcceptInviteViaNotification call core AcceptInviteViaNotification
func (m *Mobile) AcceptInviteViaNotification(id string) (string, error) {
	if !m.node.Online() {
		return "", core.ErrOffline
	}

	addr, err := m.node.AcceptInviteViaNotification(id)
	if err != nil {
		return "", err
	}

	m.node.FlushCafes()

	return addr.B58String(), nil
}

// IgnoreInviteViaNotification call core IgnoreInviteViaNotification
func (m *Mobile) IgnoreInviteViaNotification(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.IgnoreInviteViaNotification(id)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

package mobile

import (
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
)

// RegisterCafe calls core RegisterCafe
func (m *Mobile) RegisterCafe(peerId string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	if _, err := m.node.RegisterCafe(peerId); err != nil {
		return err
	}
	return nil
}

// CafeSessions calls core CafeSessions
func (m *Mobile) CafeSessions() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	items, err := m.node.CafeSessions()
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		items = make([]repo.CafeSession, 0)
	}
	return toJSON(items)
}

// CafeSession calls core CafeSession
func (m *Mobile) CafeSession(peerId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	session, err := m.node.CafeSession(peerId)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", nil
	}
	return toJSON(session)
}

// RefreshCafeSession calls core RefreshCafeSession
func (m *Mobile) RefreshCafeSession(peerId string) (string, error) {
	if !m.node.Online() {
		return "", core.ErrOffline
	}

	session, err := m.node.RefreshCafeSession(peerId)
	if err != nil {
		return "", err
	}
	return toJSON(session)
}

// DeegisterCafe calls core DeregisterCafe
func (m *Mobile) DeregisterCafe(peerId string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.DeregisterCafe(peerId)
}

// CheckCafeMessages calls core CheckCafeMessages
func (m *Mobile) CheckCafeMessages() error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	return m.node.CheckCafeMessages()
}

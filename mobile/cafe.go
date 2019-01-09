package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
)

// RegisterCafe calls core RegisterCafe
func (m *Mobile) RegisterCafe(host string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	if _, err := m.node.RegisterCafe(host); err != nil {
		return err
	}
	return nil
}

// CafeSessions calls core CafeSessions
func (m *Mobile) CafeSessions() ([]byte, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	items, err := m.node.CafeSessions()
	if err != nil {
		return [], err
	}
	return proto.Marshal(items)
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
	if !m.node.Started() {
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
	if !m.node.Started() {
		return core.ErrOffline
	}

	return m.node.CheckCafeMessages()
}

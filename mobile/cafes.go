package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
)

// RegisterCafe calls core RegisterCafe
func (m *Mobile) RegisterCafe(host string, token string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	if _, err := m.node.RegisterCafe(host, token); err != nil {
		return err
	}
	return nil
}

// CafeSession calls core CafeSession
func (m *Mobile) CafeSession(peerId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	session, err := m.node.CafeSession(peerId)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(session)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// CafeSessions calls core CafeSessions
func (m *Mobile) CafeSessions() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	bytes, err := proto.Marshal(m.node.CafeSessions())
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// RefreshCafeSession calls core RefreshCafeSession
func (m *Mobile) RefreshCafeSession(peerId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	session, err := m.node.RefreshCafeSession(peerId)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(session)
	if err != nil {
		return nil, err
	}
	return bytes, nil
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

// ListCafeRequests calls core ListCafeRequests
func (m *Mobile) ListCafeRequests(offset string, limit int) ([]byte, error) {
	return proto.Marshal(m.node.ListCafeRequests(offset, limit))
}

// MarkCafeRequestPending calls core MarkCafeRequestPending
func (m *Mobile) MarkCafeRequestPending(requestId string) error {
	return m.node.MarkCafeRequestPending(requestId)
}

// MarkCafeRequestComplete calls core MarkCafeRequestComplete
func (m *Mobile) MarkCafeRequestComplete(requestId string) error {
	return m.node.MarkCafeRequestComplete(requestId)
}

// CafeRequestGroupStatus calls core CafeRequestGroupStatus
func (m *Mobile) CafeRequestGroupStatus(group string) ([]byte, error) {
	return proto.Marshal(m.node.CafeRequestGroupStatus(group))
}

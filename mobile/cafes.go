package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
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
func (m *Mobile) CafeSession(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	session, err := m.node.CafeSession(id)
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
func (m *Mobile) RefreshCafeSession(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	session, err := m.node.RefreshCafeSession(id)
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
func (m *Mobile) DeregisterCafe(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.DeregisterCafe(id)
}

// CheckCafeMessages calls core CheckCafeMessages
func (m *Mobile) CheckCafeMessages() error {
	if !m.node.Started() {
		return core.ErrOffline
	}

	return m.node.CheckCafeMessages()
}

// CafeRequests calls core ListCafeRequests
func (m *Mobile) CafeRequests(offset string, limit int) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.CafeRequests(offset, limit))
}

// SetCafeRequestPending marks a request as pending
func (m *Mobile) SetCafeRequestPending(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.UpdateCafeRequestStatus(id, pb.CafeRequest_PENDING)
}

// SetCafeRequestComplete marks a request as complete
func (m *Mobile) SetCafeRequestComplete(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.UpdateCafeRequestStatus(id, pb.CafeRequest_COMPLETE)
}

// WriteCafeHTTPRequest calls core WriteCafeHTTPRequest
func (m *Mobile) WriteCafeHTTPRequest(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	req, err := m.node.WriteCafeHTTPRequest(id)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(req)
}

// CafeRequestSyncGroupStatus calls core CafeRequestSyncGroupStatus
func (m *Mobile) CafeRequestSyncGroupStatus(group string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.CafeRequestSyncGroupStatus(group))
}

// CleanupCafeRequests calls core CleanupCafeRequests
func (m *Mobile) CleanupCafeRequests() error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.CleanupCafeRequests()
}

package mobile

import (
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
)

// CafeSessions is a wrapper around a list of sessions
type CafeSessions struct {
	Items []repo.CafeSession `json:"items"`
}

// RegisterCafe calls core RegisterCafe
func (m *Mobile) RegisterCafe(peerId string) error {
	return core.Node.RegisterCafe(peerId)
}

// DeegisterCafe calls core DeregisterCafe
func (m *Mobile) DeregisterCafe(peerId string) error {
	return core.Node.DeregisterCafe(peerId)
}

// ListCafeSessions calls core ListCafeSessions
func (m *Mobile) ListCafeSessions() (string, error) {
	items, err := core.Node.ListCafeSessions()
	if err != nil {
		return "", err
	}
	sessions := &CafeSessions{Items: make([]repo.CafeSession, 0)}
	if len(items) > 0 {
		sessions.Items = items
	}
	return toJSON(sessions)
}

// RefreshCafeSession calls core RefreshCafeSession
func (m *Mobile) RefreshCafeSession(cafeId string) (string, error) {
	session, err := core.Node.RefreshCafeSession(cafeId)
	if err != nil {
		return "", err
	}
	return toJSON(session)
}

// CheckCafeMessages calls core CheckCafeMessages
func (m *Mobile) CheckCafeMessages() error {
	return core.Node.CheckCafeMessages()
}

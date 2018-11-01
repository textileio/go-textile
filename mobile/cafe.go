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
	if _, err := core.Node.RegisterCafe(peerId); err != nil {
		return err
	}
	return nil
}

// CafeSessions calls core ListCafeSessions
func (m *Mobile) CafeSessions() (string, error) {
	items, err := core.Node.CafeSessions()
	if err != nil {
		return "", err
	}
	sessions := &CafeSessions{Items: make([]repo.CafeSession, 0)}
	if len(items) > 0 {
		sessions.Items = items
	}
	return toJSON(sessions)
}

// CafeSession calls core CafeSession
func (m *Mobile) CafeSession(peerId string) (string, error) {
	session, err := core.Node.CafeSession(peerId)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", nil
	}
	return toJSON(session)
}

// RefreshCafeSession calls core RefreshCafeSession
func (m *Mobile) RefreshCafeSession(cafeId string) (string, error) {
	session, err := core.Node.RefreshCafeSession(cafeId)
	if err != nil {
		return "", err
	}
	return toJSON(session)
}

// DeegisterCafe calls core DeregisterCafe
func (m *Mobile) DeregisterCafe(peerId string) error {
	return core.Node.DeregisterCafe(peerId)
}

// CheckCafeMail calls core CheckCafeMessages
func (m *Mobile) CheckCafeMail() error {
	return core.Node.CheckCafeMail()
}

package core

import (
	"errors"

	"github.com/textileio/textile-go/pb"
)

// RegisterCafe registers this account with another peer (the "cafe"),
// which provides a session token for the service
func (t *Textile) RegisterCafe(host string, token string) (*pb.CafeSession, error) {
	session, err := t.cafe.Register(host, token)
	if err != nil {
		return nil, err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(); err != nil {
			return nil, err
		}
	}

	if err := t.UpdateContactInboxes(); err != nil {
		return nil, err
	}

	if err := t.PublishContact(); err != nil {
		return nil, err
	}

	return session, nil
}

// CafeSessions lists active cafe sessions
func (t *Textile) CafeSessions() *pb.CafeSessionList {
	return t.datastore.CafeSessions().List()
}

// CafeSession returns an active session by id
func (t *Textile) CafeSession(peerId string) (*pb.CafeSession, error) {
	return t.datastore.CafeSessions().Get(peerId), nil
}

// RefreshCafeSession attempts to refresh a token with a cafe
func (t *Textile) RefreshCafeSession(peerId string) (*pb.CafeSession, error) {
	session := t.datastore.CafeSessions().Get(peerId)
	if session == nil {
		return nil, errors.New("session not found")
	}
	return t.cafe.refresh(session)
}

// DeregisterCafe removes the session associated with the given cafe
func (t *Textile) DeregisterCafe(peerId string) error {
	session := t.datastore.CafeSessions().Get(peerId)
	if session == nil {
		return nil
	}

	// clean up
	if err := t.datastore.CafeRequests().DeleteByCafe(session.Id); err != nil {
		return err
	}
	if err := t.datastore.CafeSessions().Delete(peerId); err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(); err != nil {
			return err
		}
	}

	if err := t.UpdateContactInboxes(); err != nil {
		return err
	}

	return t.PublishContact()
}

// CheckCafeMessages fetches new messages from registered cafes
func (t *Textile) CheckCafeMessages() error {
	return t.cafeInbox.CheckMessages()
}

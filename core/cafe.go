package core

import (
	"errors"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// RegisterCafe registers this account with another peer (the "cafe"),
// which provides a session token for the service
func (t *Textile) RegisterCafe(peerId string) (*repo.CafeSession, error) {
	if !t.Online() {
		return nil, ErrOffline
	}

	// call up the peer, see if they're offering a cafe
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return nil, err
	}
	if err := t.cafeService.Register(pid); err != nil {
		return nil, err
	}

	// publish profile w/ updated inboxes
	if err := t.PublishProfile(); err != nil {
		return nil, err
	}

	return t.datastore.CafeSessions().Get(pid.Pretty()), nil
}

// CafeSessions lists active cafe sessions
func (t *Textile) CafeSessions() ([]repo.CafeSession, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.CafeSessions().List(), nil
}

// CafeSession returns an active session by id
func (t *Textile) CafeSession(peerId string) (*repo.CafeSession, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.CafeSessions().Get(peerId), nil
}

// RefreshCafeSession attempts to refresh a token with a cafe
func (t *Textile) RefreshCafeSession(cafeId string) (*repo.CafeSession, error) {
	if !t.Online() {
		return nil, ErrOffline
	}
	session := t.datastore.CafeSessions().Get(cafeId)
	if session == nil {
		return nil, errors.New("session not found")
	}
	return t.cafeService.refresh(session)
}

// DeregisterCafe removes the session associated with the given cafe
func (t *Textile) DeregisterCafe(peerId string) error {
	if err := t.touchDatastore(); err != nil {
		return err
	}
	if err := t.datastore.CafeSessions().Delete(peerId); err != nil {
		return err
	}

	// publish profile w/ updated inboxes
	return t.PublishProfile()
}

// CheckCafeMail fetches new messages from registered cafes
func (t *Textile) CheckCafeMail() error {
	if err := t.touchDatastore(); err != nil {
		return err
	}
	if !t.Online() {
		return ErrOffline
	}
	return t.cafeInbox.CheckMessages()
}

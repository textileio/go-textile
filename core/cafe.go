package core

import (
	"errors"
	"strings"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/textileio/textile-go/repo"
)

// RegisterCafe registers this account with another peer (the "cafe"),
// which provides a session token for the service
func (t *Textile) RegisterCafe(peerId string) (*repo.CafeSession, error) {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return nil, err
	}
	if err := t.cafeService.Register(pid); err != nil {
		return nil, err
	}

	// add to bootstrap
	session := t.datastore.CafeSessions().Get(pid.Pretty())
	if session != nil {
		var peers []string
		for _, s := range session.SwarmAddrs {
			if !strings.Contains(s, "/ws/") {
				peers = append(peers, s+"/ipfs/"+session.CafeId)
			}
		}
		if err := updateBootstrapConfig(t.repoPath, peers, []string{}); err != nil {
			return nil, err
		}
	}

	for _, thrd := range t.threads {
		if _, err := thrd.annouce(); err != nil {
			return nil, err
		}
	}

	if err := t.PublishProfile(); err != nil {
		return nil, err
	}

	return session, nil
}

// CafeSessions lists active cafe sessions
func (t *Textile) CafeSessions() ([]repo.CafeSession, error) {
	return t.datastore.CafeSessions().List(), nil
}

// CafeSession returns an active session by id
func (t *Textile) CafeSession(peerId string) (*repo.CafeSession, error) {
	return t.datastore.CafeSessions().Get(peerId), nil
}

// RefreshCafeSession attempts to refresh a token with a cafe
func (t *Textile) RefreshCafeSession(peerId string) (*repo.CafeSession, error) {
	session := t.datastore.CafeSessions().Get(peerId)
	if session == nil {
		return nil, errors.New("session not found")
	}
	return t.cafeService.refresh(session)
}

// DeregisterCafe removes the session associated with the given cafe
func (t *Textile) DeregisterCafe(peerId string) error {
	session := t.datastore.CafeSessions().Get(peerId)
	if session == nil {
		return nil
	}

	// remove from bootstrap
	var peers []string
	for _, s := range session.SwarmAddrs {
		peers = append(peers, s+"/ipfs/"+session.CafeId)
	}
	if err := updateBootstrapConfig(t.repoPath, []string{}, peers); err != nil {
		return err
	}

	// clean up
	if err := t.datastore.CafeRequests().DeleteByCafe(session.CafeId); err != nil {
		return err
	}
	if err := t.datastore.CafeSessions().Delete(peerId); err != nil {
		return err
	}

	for _, thrd := range t.threads {
		if _, err := thrd.annouce(); err != nil {
			return err
		}
	}

	return t.PublishProfile()
}

// CheckCafeMessages fetches new messages from registered cafes
func (t *Textile) CheckCafeMessages() error {
	return t.cafeInbox.CheckMessages()
}

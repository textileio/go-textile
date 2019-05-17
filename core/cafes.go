package core

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/textileio/go-textile/pb"
)

// RegisterCafe registers this account with another peer (the "cafe"),
// which provides a session token for the service
func (t *Textile) RegisterCafe(host string, token string) (*pb.CafeSession, error) {
	session, err := t.cafe.Register(host, token)
	if err != nil {
		return nil, err
	}

	err = t.updatePeerInboxes()
	if err != nil {
		return nil, err
	}

	for _, thrd := range t.loadedThreads {
		_, err = thrd.annouce(nil)
		if err != nil {
			return nil, err
		}
	}

	err = t.publishPeer()
	if err != nil {
		return nil, err
	}

	err = t.SnapshotThreads()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// CafeSession returns an active session by id
func (t *Textile) CafeSession(id string) (*pb.CafeSession, error) {
	return t.datastore.CafeSessions().Get(id), nil
}

// CafeSessions lists active cafe sessions
func (t *Textile) CafeSessions() *pb.CafeSessionList {
	return t.datastore.CafeSessions().List()
}

// RefreshCafeSession attempts to refresh a token with a cafe
func (t *Textile) RefreshCafeSession(id string) (*pb.CafeSession, error) {
	session := t.datastore.CafeSessions().Get(id)
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}
	return t.cafe.refresh(session)
}

// DeregisterCafe removes the session associated with the given cafe
func (t *Textile) DeregisterCafe(id string) error {
	cafe, err := peer.IDB58Decode(id)
	if err != nil {
		return err
	}
	err = t.cafe.Deregister(cafe)
	if err != nil {
		return err
	}

	err = t.updatePeerInboxes()
	if err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		_, err := thrd.annouce(nil)
		if err != nil {
			return err
		}
	}

	return t.publishPeer()
}

// CheckCafeMessages fetches new messages from registered cafes
func (t *Textile) CheckCafeMessages() error {
	return t.cafeInbox.CheckMessages()
}

// cafesEqual returns whether or not the two cafes are identical
// Note: swarms are allowed to be in different order and still be "equal"
func cafesEqual(a *pb.Cafe, b *pb.Cafe) bool {
	if a.Peer != b.Peer {
		return false
	}
	if a.Address != b.Address {
		return false
	}
	if a.Api != b.Api {
		return false
	}
	if a.Protocol != b.Protocol {
		return false
	}
	if a.Node != b.Node {
		return false
	}
	if a.Url != b.Url {
		return false
	}
	return true
}

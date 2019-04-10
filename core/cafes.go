package core

import (
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/proto"

	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"

	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

// RegisterCafe registers this account with another peer (the "cafe"),
// which provides a session token for the service
func (t *Textile) RegisterCafe(host string, token string) (*pb.CafeSession, error) {
	session, err := t.cafe.Register(host, token)
	if err != nil {
		return nil, err
	}

	if err := t.updatePeerInboxes(); err != nil {
		return nil, err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(nil); err != nil {
			return nil, err
		}
	}

	if err := t.publishPeer(); err != nil {
		return nil, err
	}

	if err := t.SnapshotThreads(); err != nil {
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
	if err := t.cafe.Deregister(cafe); err != nil {
		return err
	}

	if err := t.updatePeerInboxes(); err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(nil); err != nil {
			return err
		}
	}

	return t.publishPeer()
}

// CheckCafeMessages fetches new messages from registered cafes
func (t *Textile) CheckCafeMessages() error {
	return t.cafeInbox.CheckMessages()
}

// CafeRequests returns a batch of requests
func (t *Textile) CafeRequests(offset string, limit int) *pb.CafeRequestList {
	return t.datastore.CafeRequests().List(offset, limit)
}

// UpdateCafeRequestStatus updates a request status
func (t *Textile) UpdateCafeRequestStatus(id string, status pb.CafeRequest_Status) error {
	return t.datastore.CafeRequests().UpdateStatus(id, status)
}

// CafeHTTPRequest returns the type, path, headers, body and token for an HTTP cafe request
// - store: PUT /store/:cid, body => raw object data
// - unstore: DELETE /store/:cid, body => none
// - store thread: PUT /threads/:id, body => encrypted thread object (snapshot)
// - unstore thread: DELETE /threads/:id, body => none
// - deliver message: POST /inbox/:pid, body => encrypted message
func (t *Textile) CafeHTTPRequest(id string) (*pb.CafeHTTPRequest, error) {
	req := t.datastore.CafeRequests().Get(id)
	if req == nil {
		return nil, fmt.Errorf("request not found")
	}

	session := t.datastore.CafeSessions().Get(req.Cafe.Peer)
	if session == nil {
		return nil, fmt.Errorf("session for cafe %s not found", req.Cafe.Peer)
	}

	hreq := &pb.CafeHTTPRequest{
		Path: "/api/v1",
		Headers: map[string]string{
			"Authorization":  "Basic " + session.Access,
			"X-Textile-Peer": t.node.Identity.Pretty(),
		},
	}

	switch req.Type {
	case pb.CafeRequest_STORE:
		hreq.Type = pb.CafeHTTPRequest_PUT
		hreq.Path += "/store/" + req.Target

		data, err := ipfs.DataAtPath(t.node, req.Target)
		if err != nil {
			if err == iface.ErrIsDir {
				data, err := ipfs.GetObjectAtPath(t.node, req.Target)
				if err != nil {
					return nil, err
				}
				hreq.Headers["X-Textile-Store-Type"] = "object"
				hreq.Body = data
			} else {
				return nil, err
			}
		} else {
			hreq.Headers["X-Textile-Store-Type"] = "data"
			hreq.Body = data
		}

	case pb.CafeRequest_UNSTORE:
		hreq.Type = pb.CafeHTTPRequest_DELETE
		hreq.Path += "/store/" + req.Target

	case pb.CafeRequest_STORE_THREAD:
		hreq.Type = pb.CafeHTTPRequest_PUT
		hreq.Path += "/threads/" + req.Target

		thrd := t.datastore.Threads().Get(req.Target)
		if thrd == nil {
			return nil, ErrThreadNotFound
		}
		plaintext, err := proto.Marshal(thrd)
		if err != nil {
			return nil, err
		}
		ciphertext, err := t.Encrypt(plaintext)
		if err != nil {
			return nil, err
		}
		hreq.Body = ciphertext

	case pb.CafeRequest_UNSTORE_THREAD:
		hreq.Type = pb.CafeHTTPRequest_DELETE
		hreq.Path += "/threads/" + req.Target

	case pb.CafeRequest_INBOX:
		hreq.Type = pb.CafeHTTPRequest_POST
		hreq.Path += "/inbox/" + req.Peer
		hreq.Body = []byte(req.Target)
	}

	if hreq.Body != nil {
		sig, err := t.node.PrivateKey.Sign(hreq.Body)
		if err != nil {
			return nil, err
		}
		hreq.Headers["X-Textile-Peer-Sig"] = hex.EncodeToString(sig)
	}

	return hreq, nil
}

// CafeRequestGroupStatus returns the status of a request group
func (t *Textile) CafeRequestGroupStatus(group string) *pb.CafeRequestGroupStatus {
	return t.datastore.CafeRequests().GroupStatus(group)
}

// CleanupCafeRequests deletes request groups that are completely completed
func (t *Textile) CleanupCafeRequests() error {
	for _, group := range t.datastore.CafeRequests().ListCompletedGroups() {
		if err := t.datastore.CafeRequests().DeleteByGroup(group); err != nil {
			return err
		}
	}
	return nil
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
	if len(a.Swarm) != len(b.Swarm) {
		return false
	}
	as := make(map[string]struct{})
	for _, s := range a.Swarm {
		as[s] = struct{}{}
	}
	for _, s := range b.Swarm {
		if _, ok := as[s]; !ok {
			return false
		}
	}
	return true
}

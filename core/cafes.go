package core

import (
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	iface "github.com/ipfs/interface-go-ipfs-core"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
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

// WriteCafeHTTPRequest returns the type, url, headers, and body path for an HTTP cafe request
// - store: PUT /store/:cid, body => raw object data
// - unstore: DELETE /store/:cid, body => none
// - store thread: PUT /threads/:id, body => encrypted thread object (snapshot)
// - unstore thread: DELETE /threads/:id, body => none
// - deliver message: POST /inbox/:pid, body => encrypted message
func (t *Textile) WriteCafeHTTPRequest(id string) (*pb.CafeHTTPRequest, error) {
	req := t.datastore.CafeRequests().Get(id)
	if req == nil {
		return nil, fmt.Errorf("request not found")
	}

	session := t.datastore.CafeSessions().Get(req.Cafe.Peer)
	if session == nil {
		return nil, fmt.Errorf("session for cafe %s not found", req.Cafe.Peer)
	}

	hreq := &pb.CafeHTTPRequest{
		Url: session.Cafe.Url + "/api/v1",
		Headers: map[string]string{
			"Authorization":  "Basic " + session.Access,
			"X-Textile-Peer": t.node.Identity.Pretty(),
		},
	}

	var body []byte
	switch req.Type {
	case pb.CafeRequest_STORE:
		hreq.Type = pb.CafeHTTPRequest_PUT
		hreq.Url += "/store/" + req.Target

		data, err := ipfs.DataAtPath(t.node, req.Target)
		if err != nil {
			if err == iface.ErrIsDir {
				data, err := ipfs.GetObjectAtPath(t.node, req.Target)
				if err != nil {
					return nil, err
				}
				hreq.Headers["X-Textile-Store-Type"] = "object"
				body = data
			} else {
				return nil, err
			}
		} else {
			hreq.Headers["X-Textile-Store-Type"] = "data"
			body = data
		}

	case pb.CafeRequest_UNSTORE:
		hreq.Type = pb.CafeHTTPRequest_DELETE
		hreq.Url += "/store/" + req.Target

	case pb.CafeRequest_STORE_THREAD:
		hreq.Type = pb.CafeHTTPRequest_PUT
		hreq.Url += "/threads/" + req.Target

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
		body = ciphertext

	case pb.CafeRequest_UNSTORE_THREAD:
		hreq.Type = pb.CafeHTTPRequest_DELETE
		hreq.Url += "/threads/" + req.Target

	case pb.CafeRequest_INBOX:
		hreq.Type = pb.CafeHTTPRequest_POST
		hreq.Url += "/inbox/" + req.Peer
		body = []byte(req.Target)
	}

	if body != nil {
		sig, err := t.node.PrivateKey.Sign(body)
		if err != nil {
			return nil, err
		}
		hreq.Headers["X-Textile-Peer-Sig"] = hex.EncodeToString(sig)

		hreq.Path = filepath.Join(t.repoPath, "tmp", id)
		err = util.WriteFileByPath(hreq.Path, body)
		if err != nil {
			return nil, err
		}
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
	return true
}

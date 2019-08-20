package core

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

// RegisterCafe registers this account with another peer (the "cafe"),
// which provides a session token for the service
func (t *Textile) RegisterCafe(id string, token string) (*pb.CafeSession, error) {
	// ensure id is a peer id
	_, err := peer.IDB58Decode(id)
	if err != nil {
		return nil, fmt.Errorf("not a valid peerID: %s", id)
	}

	// return existing session
	if x := t.datastore.CafeSessions().Get(id); x != nil {
		return x, nil
	}

	session, err := t.cafe.Register(id, token)
	if err != nil {
		return nil, err
	}

	err = t.UpdatePeerInboxes()
	if err != nil {
		return nil, err
	}

	// sync all blocks and files target
	err = t.CafeRequestThreadsContent(session.Id)
	if err != nil {
		return nil, err
	}

	for _, thrd := range t.loadedThreads {
		_, err = thrd.Annouce(nil)
		if err != nil {
			return nil, err
		}
	}

	err = t.PublishPeer()
	if err != nil {
		return nil, err
	}

	err = t.SnapshotThreads()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// DeregisterCafe removes the session associated with the given cafe
func (t *Textile) DeregisterCafe(id string) error {
	// ensure id is a peer id
	_, err := peer.IDB58Decode(id)
	if err != nil {
		return fmt.Errorf("not a valid peerID: %s", id)
	}

	err = t.cafe.Deregister(id)
	if err != nil {
		return err
	}

	err = t.UpdatePeerInboxes()
	if err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		_, err := thrd.Annouce(nil)
		if err != nil {
			return err
		}
	}

	return t.PublishPeer()
}

// RefreshCafeSession attempts to refresh a token with a cafe
func (t *Textile) RefreshCafeSession(id string) (*pb.CafeSession, error) {
	session := t.datastore.CafeSessions().Get(id)
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}
	return t.cafe.refresh(session)
}

// CheckCafeMessages fetches new messages from registered cafes
func (t *Textile) CheckCafeMessages() error {
	return t.cafeInbox.CheckMessages()
}

// CafeSession returns an active session by id
func (t *Textile) CafeSession(id string) (*pb.CafeSession, error) {
	return t.datastore.CafeSessions().Get(id), nil
}

// CafeSessions lists active cafe sessions
func (t *Textile) CafeSessions() *pb.CafeSessionList {
	return t.datastore.CafeSessions().List()
}

// CafeRequestThreadContent sync the entire thread conents (blocks and files) to the given cafe
func (t *Textile) CafeRequestThreadsContent(cafe string) error {
	for _, thrd := range t.loadedThreads {
		blocks := t.Blocks("", -1, fmt.Sprintf("threadId='%s'", thrd.Id))
		for _, b := range blocks.Items {

			// store the block itself
			err := t.cafeOutbox.Add(b.Id, pb.CafeRequest_STORE, cafeReqOpt.SyncGroup(b.Id), cafeReqOpt.Cafe(cafe))
			if err != nil {
				return err
			}

			// store the files DAGs
			if b.Type == pb.Block_FILES {
				dec, err := cid.Decode(b.Data)
				if err != nil {
					return err
				}
				node, err := ipfs.NodeAtCid(t.node, dec)
				if err != nil {
					return err
				}
				err = thrd.cafeReqFileData(node, b.Id, cafe)
				if err != nil {
					return err
				}
			}
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

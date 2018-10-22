package core

import (
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// Leave creates an outgoing leave block
func (t *Thread) Leave() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// commit to ipfs
	res, err := t.commitBlock(nil, pb.ThreadBlock_LEAVE, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.LeaveBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	// post it
	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	// delete blocks
	if err := t.datastore.Blocks().DeleteByThread(t.Id); err != nil {
		return nil, err
	}

	// delete peers
	if err := t.datastore.ThreadPeers().DeleteByThread(t.Id); err != nil {
		return nil, err
	}

	// delete notifications
	if err := t.datastore.Notifications().DeleteBySubject(t.Id); err != nil {
		return nil, err
	}

	log.Debugf("added LEAVE to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// HandleLeaveBlock handles an incoming leave block
func (t *Thread) HandleLeaveBlock(from *peer.ID, hash mh.Multihash, block *pb.ThreadBlock, following bool) error {
	// remove peer
	if err := t.datastore.ThreadPeers().Delete(block.Header.Author, t.Id); err != nil {
		return err
	}
	if err := t.datastore.Notifications().DeleteByActor(block.Header.Author); err != nil {
		return err
	}

	// index it locally
	if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.LeaveBlock, nil); err != nil {
		return err
	}

	// back prop
	newPeers, err := t.FollowParents(block.Header.Parents, from)
	if err != nil {
		return err
	}

	// handle HEAD
	if following {
		return nil
	}
	if _, err := t.handleHead(hash, block.Header.Parents); err != nil {
		return nil
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return nil
		}
	}
	return nil
}

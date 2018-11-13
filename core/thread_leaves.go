package core

import (
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
)

// leave creates an outgoing leave block
func (t *Thread) leave() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	res, err := t.commitBlock(nil, pb.ThreadBlock_LEAVE, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, repo.LeaveBlock, "", ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	// cleanup
	// TODO: delete files
	if err := t.datastore.Blocks().DeleteByThread(t.Id); err != nil {
		return nil, err
	}
	if err := t.datastore.ThreadPeers().DeleteByThread(t.Id); err != nil {
		return nil, err
	}
	if err := t.datastore.Notifications().DeleteBySubject(t.Id); err != nil {
		return nil, err
	}

	log.Debugf("added LEAVE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleLeaveBlock handles an incoming leave block
func (t *Thread) handleLeaveBlock(hash mh.Multihash, block *pb.ThreadBlock) error {
	// remove peer
	if err := t.datastore.ThreadPeers().Delete(block.Header.Author, t.Id); err != nil {
		return err
	}
	if err := t.datastore.Notifications().DeleteByActor(block.Header.Author); err != nil {
		return err
	}

	return t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.LeaveBlock, "", "")
}

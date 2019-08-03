package core

import (
	"fmt"

	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// leave creates an outgoing leave block
func (t *Thread) leave() (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	res, err := t.commitBlock(nil, pb.Block_LEAVE, true, nil)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_LEAVE,
		Date:   res.header.Date,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	// cleanup
	query := fmt.Sprintf("threadId='%s'", t.Id)
	for _, block := range t.datastore.Blocks().List("", -1, query).Items {
		err = t.ignoreBlockTarget(block)
		if err != nil {
			return nil, err
		}
	}
	err = t.datastore.Blocks().DeleteByThread(t.Id)
	if err != nil {
		return nil, err
	}
	err = t.datastore.ThreadPeers().DeleteByThread(t.Id)
	if err != nil {
		return nil, err
	}
	err = t.datastore.Notifications().DeleteBySubject(t.Id)
	if err != nil {
		return nil, err
	}

	log.Debugf("added LEAVE to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleLeaveBlock handles an incoming leave block
func (t *Thread) handleLeaveBlock(block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	if !t.readable(t.config.Account.Address) {
		return res, ErrNotReadable
	}
	if !t.readable(block.Header.Address) {
		return res, ErrNotReadable
	}

	err := t.datastore.ThreadPeers().Delete(block.Header.Author, t.Id)
	if err != nil {
		return res, err
	}
	err = t.datastore.Notifications().DeleteByActor(block.Header.Author)
	if err != nil {
		return res, err
	}

	return res, nil
}

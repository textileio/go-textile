package core

import (
	"fmt"

	"github.com/ipfs/go-ipfs/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

// downloadsFlushGroupSize is the size of concurrently processed downloads
const downloadsFlushGroupSize = 16

// maxDownloadAttempts is the number of times a download can fail to download before being deleted
const maxDownloadAttempts = 5

// BlockDownloads manages a queue of pending downloads
type BlockDownloads struct {
	node      func() *core.IpfsNode
	datastore repo.Datastore
	getThread func(id string) *Thread
	flushing  bool
}

// NewBlockDownloads creates a new download queue
func NewBlockDownloads(node func() *core.IpfsNode, datastore repo.Datastore, getThread func(id string) *Thread) *BlockDownloads {
	return &BlockDownloads{
		node:      node,
		datastore: datastore,
		getThread: getThread,
	}
}

// Add queues a download, starting it if flush is not active
func (q *BlockDownloads) Add(download *pb.Block) error {
	err := q.datastore.Blocks().Add(download)
	if err == nil {
		go q.Flush()
	}
	return err
}

// Flush processes pending messages
func (q *BlockDownloads) Flush() {
	if q.flushing {
		return
	}
	q.flushing = true
	defer func() {
		q.flushing = false
	}()
	log.Debug("flushing downloads")

	query := fmt.Sprintf("status=%d", pb.Block_PENDING)
	q.batch(q.datastore.Blocks().List("", downloadsFlushGroupSize, query).Items)
}

// batch flushes a batch of downloads
func (q *BlockDownloads) batch(downloads []*pb.Block) {
	log.Debugf("handling %d downloads", len(downloads))
	if len(downloads) == 0 {
		return
	}

	for _, dl := range downloads {
		go func(dl *pb.Block) {
			err := q.handle(dl)
			if err != nil {
				log.Warningf("handle attempt failed for download %s: %s", dl.Id, err)
				return
			}
			log.Debugf("handled download %s", dl.Id)
		}(dl)
	}

	// next batch
	offset := downloads[len(downloads)-1].Id
	query := fmt.Sprintf("status=%d", pb.Block_PENDING)
	q.batch(q.datastore.Blocks().List(offset, downloadsFlushGroupSize, query).Items)
}

// handle handles a single message
func (q *BlockDownloads) handle(dl *pb.Block) error {
	fail := func(reason string) error {
		log.Warningf("download %s failed: %s", dl.Id, reason)
		return q.datastore.Blocks().Delete(dl.Id)
	}

	ciphertext, err := ipfs.DataAtPath(q.node(), dl.Id)
	if err != nil {
		return q.handleErr(err, dl)
	}

	thread := q.getThread(dl.Thread)
	if thread == nil {
		return fail("thread not found")
	}

	_, err = thread.handle(&blockNode{hash: dl.Id,
		ciphertext: ciphertext,
		parents:    dl.Parents,
		target:     dl.Target,
		data:       dl.Data,
	}, true)
	if err != nil {
		return fail(err.Error())
	}
	return nil
}

// handleErr deletes or adds an attempt to a download processing error
func (q *BlockDownloads) handleErr(herr error, dl *pb.Block) error {
	var err error
	if dl.Attempts+1 >= maxDownloadAttempts {
		err = q.datastore.Blocks().Delete(dl.Id)
	} else {
		err = q.datastore.Blocks().AddAttempt(dl.Id)
	}
	if err != nil {
		return err
	}
	return herr
}

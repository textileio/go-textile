package core

import (
	"sync"

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
	mux       sync.Mutex
}

// NewBlockDownloads creates a new download queue
func NewBlockDownloads(node func() *core.IpfsNode, datastore repo.Datastore) *BlockDownloads {
	return &BlockDownloads{
		node:      node,
		datastore: datastore,
	}
}

// Add adds a download to the queue
func (q *BlockDownloads) Add(dl *pb.BlockDownload) error {
	log.Debugf("adding download for %s: %s", ipfs.ShortenID(dl.Thread), dl.Id)

	return q.datastore.BlockDownloads().Add(dl)
}

// Flush processes pending messages
func (q *BlockDownloads) Flush() {
	q.mux.Lock()
	defer q.mux.Unlock()
	log.Debug("flushing downloads")

	err := q.batch(q.datastore.BlockDownloads().List("", downloadsFlushGroupSize))
	if err != nil {
		log.Errorf("downloads batch error: %s", err)
		return
	}
}

// batch flushes a batch of downloads
func (q *BlockDownloads) batch(downloads []pb.BlockDownload) error {
	log.Debugf("handling %d downloads", len(downloads))
	if len(downloads) == 0 {
		return nil
	}

	for _, dl := range downloads {
		go func(dl pb.BlockDownload) {
			err := q.handle(dl)
			if err != nil {
				log.Warningf("handle attempt failed for download %s: %s", dl.Id, err)
				return
			}
			err = q.datastore.BlockDownloads().Delete(dl.Id)
			if err != nil {
				log.Errorf("failed to delete download %s: %s", dl.Id, err)
			} else {
				log.Debugf("handled download %s", dl.Id)
			}
		}(dl)
	}

	// next batch
	offset := downloads[len(downloads)-1].Id
	next := q.datastore.BlockDownloads().List(offset, downloadsFlushGroupSize)

	// keep going
	return q.batch(next)
}

// handle handles a single message
func (q *BlockDownloads) handle(dl pb.BlockDownload) error {
	fail := func() error {
		return q.datastore.BlockDownloads().Delete(dl.Id)
	}

	thread := q.getThread(dl.Thread)
	if thread == nil {
		return fail()
	}

	ciphertext, err := ipfs.DataAtPath(q.node(), dl.Id)
	if err != nil {
		return q.handleErr(err, dl)
	}

	block, err := thread.handleBlock(dl.Id, ciphertext)
	if err != nil {
		if err == ErrBlockExists {
			// exists, abort
			log.Debugf("%s exists, aborting", dl.Id)
			return fail()
		}
		return fail()
	}

	err = thread.indexBlock(&pb.Block{
		Id:     dl.Id,
		Thread: thread.Id,
		Author: block.Header.Author,
		Type:   block.Type,
		Date:   block.Header.Date,
		//Parents: commit.parents,
		//Target:  target,
		//Body:    block.,
	})

	return nil
}

// handleErr deletes or adds an attempt to a download processing error
func (q *BlockDownloads) handleErr(herr error, dl pb.BlockDownload) error {
	var err error
	if dl.Attempts+1 >= maxDownloadAttempts {
		err = q.datastore.BlockDownloads().Delete(dl.Id)
	} else {
		err = q.datastore.BlockDownloads().AddAttempt(dl.Id)
	}
	if err != nil {
		return err
	}
	return herr
}

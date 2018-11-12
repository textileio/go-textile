package core

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
)

// AddFile adds an outgoing files block
func (t *Thread) AddFiles(node ipld.Node, caption string, keys Keys) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.schema == nil {
		return nil, ErrThreadSchemaRequired
	}

	target := node.Cid().Hash().B58String()

	// each link should point to a dag described by the thread schema
	for _, link := range node.Links() {
		nd, err := ipfs.LinkNode(t.node(), link)
		if err != nil {
			return nil, err
		}
		if err := t.process(t.schema, nd, false); err != nil {
			return nil, err
		}
	}

	t.cafeOutbox.Add(target, repo.CafeStoreRequest)

	msg := &pb.ThreadFiles{
		Target: node.Cid().Hash().B58String(),
		Body:   caption,
		Keys:   keys,
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_FILES, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, repo.FilesBlock, msg.Target, msg.Body); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added FILES to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleFilesBlock handles an incoming files block
func (t *Thread) handleFilesBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadFiles, error) {
	msg := new(pb.ThreadFiles)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	if t.schema == nil {
		return nil, ErrThreadSchemaRequired
	}

	var ignore bool
	ignored := t.datastore.Blocks().GetByTarget(fmt.Sprintf("ignore-%s", hash.B58String()))
	if ignored != nil {
		date, err := ptypes.Timestamp(block.Header.Date)
		if err != nil {
			return nil, err
		}
		// ignore if the ignore block came after (could happen during back prop)
		if ignored.Date.After(date) {
			ignore = true
		}
	}
	if !ignore {

		id, err := cid.Parse(hash)
		if err != nil {
			return nil, err
		}
		node, err := ipfs.CidNode(t.node(), id)
		if err != nil {
			return nil, err
		}
		// each link should point to a dag described by the thread schema
		for _, link := range node.Links() {
			nd, err := ipfs.LinkNode(t.node(), link)
			if err != nil {
				return nil, err
			}
			if err := t.process(t.schema, nd, true); err != nil {
				return nil, err
			}
		}

		t.cafeOutbox.Add(hash.B58String(), repo.CafeStoreRequest)

		for hash, key := range msg.Keys {
			if err := t.datastore.ThreadFileKeys().Add(&repo.ThreadFileKey{
				Hash: hash,
				Key:  key,
			}); err != nil {
				return nil, err
			}
		}
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.FilesBlock, msg.Target, msg.Body); err != nil {
		return nil, err
	}

	return msg, nil
}

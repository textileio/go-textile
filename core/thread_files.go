package core

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
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
		nd, err := ipfs.NodeAtLink(t.node(), link)
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

		target, err := cid.Parse(msg.Target)
		if err != nil {
			return nil, err
		}
		node, err := ipfs.NodeAtCid(t.node(), target)
		if err != nil {
			return nil, err
		}
		if err := ipfs.PinNode(t.node(), node, false); err != nil {
			return nil, err
		}

		// each link should point to a dag described by the thread schema
		for _, link := range node.Links() {
			nd, err := ipfs.NodeAtLink(t.node(), link)
			if err != nil {
				return nil, err
			}
			if err := t.process(t.schema, nd, true); err != nil {
				return nil, err
			}
		}

		t.cafeOutbox.Add(msg.Target, repo.CafeStoreRequest)

		// use msg keys to decrypt each file
		for pth, key := range msg.Keys {
			nd, err := ipfs.NodeAtPath(t.node(), path.Path(msg.Target+pth+FileLinkName))
			if err != nil {
				return nil, err
			}

			keyb, err := base58.Decode(key)
			if err != nil {
				return nil, err
			}
			plaintext, err := crypto.DecryptAES(nd.RawData(), keyb)
			if err != nil {
				return nil, err
			}

			var file repo.File
			if err := json.Unmarshal(plaintext, &file); err != nil {
				return nil, err
			}
			if err := t.datastore.Files().Add(&file); err != nil {
				log.Debugf("received file already exists: %s", file.Hash)
			} else {
				log.Debugf("received file: %s", file.Hash)
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

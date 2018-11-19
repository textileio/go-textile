package core

import (
	"encoding/json"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"

	"github.com/golang/protobuf/ptypes"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// AddFile adds an outgoing files block
func (t *Thread) AddFiles(node ipld.Node, caption string, keys Keys) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.Schema == nil {
		return nil, ErrThreadSchemaRequired
	}
	if node == nil {
		return nil, ErrInvalidFileNode
	}

	target := node.Cid().Hash().B58String()

	// each link should point to a dag described by the thread schema
	for _, link := range node.Links() {
		nd, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return nil, err
		}
		if err := t.processNode(t.Schema, nd, false); err != nil {
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

	if t.Schema == nil {
		return nil, ErrThreadSchemaRequired
	}

	var ignore bool
	ignored := t.datastore.Blocks().List("", -1, "target='ignore-"+hash.B58String()+"'")
	if len(ignored) > 0 {
		date, err := ptypes.Timestamp(block.Header.Date)
		if err != nil {
			return nil, err
		}
		// ignore if the first (latest) ignore came after (could happen during back prop)
		if ignored[0].Date.After(date) {
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
			if err := t.processNode(t.Schema, nd, true); err != nil {
				return nil, err
			}
		}

		t.cafeOutbox.Add(msg.Target, repo.CafeStoreRequest)

		// use msg keys to decrypt each file
		for pth, key := range msg.Keys {
			fd, err := ipfs.DataAtPath(t.node(), msg.Target+pth+FileLinkName)
			if err != nil {
				return nil, err
			}

			keyb, err := base58.Decode(key)
			if err != nil {
				return nil, err
			}
			plaintext, err := crypto.DecryptAES(fd, keyb)
			if err != nil {
				return nil, err
			}

			var file repo.File
			if err := json.Unmarshal(plaintext, &file); err != nil {
				return nil, err
			}
			log.Debugf("received file: %s", file.Hash)
			if err := t.datastore.Files().Add(&file); err != nil {
				log.Debugf("file exists: %s", file.Hash)
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

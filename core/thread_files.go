package core

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
)

// AddFile adds an outgoing files block
func (t *Thread) AddFiles(node ipld.Node, caption string, keys map[string]string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// each link should point to a dag described by the thread schema
	//links := node.Links()
	//
	//for sl, s := range t.schema.Nodes {
	//	s.
	//}
	//
	//for _, link := range node.Links() {
	//
	//	n, err := ipfs.GetNode(t.node(), link)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	// build block
	msg := &pb.ThreadFiles{
		Target: node.Cid().Hash().B58String(),
		Body:   caption,
		Keys:   keys,
	}

	// TODO: verify files exist? schema matched? pin

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_FILES, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.FilesBlock, msg.Target, msg.Body); err != nil {
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

	log.Debugf("added FILES to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// handleFilesBlock handles an incoming files block
func (t *Thread) handleFilesBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadFiles, error) {
	msg := new(pb.ThreadFiles)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// check if this block has been ignored
	var ignore bool
	ignored := t.datastore.Blocks().GetByTarget(fmt.Sprintf("ignore-%s", hash.B58String()))
	if ignored != nil {
		date, err := ptypes.Timestamp(block.Header.Date)
		if err != nil {
			return nil, err
		}
		// ignore if the ignore block came after (could happen when bootstrapping the thread during back prop)
		if ignored.Date.After(date) {
			ignore = true
		}
	}
	if !ignore {
		// pin the schema (it's likely already pinned)

		// pin data first (it may not be available)
		// TODO: this shouldn't block, may need to queue somewhere
		if err := ipfs.PinPath(t.node(), fmt.Sprintf("%s/thumb", msg.Data), false); err != nil {
			return nil, err
		}

		// get metadata
		// TODO: same here
		meta, err := getMetadata(t.node(), msg.Data, msg.Key)
		if err != nil {
			return nil, err
		}
		dconf.DataMetadata = meta
	}

	// index
	if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.FilesBlock, dconf); err != nil {
		return nil, err
	}

	return msg, nil
}

// getMetadata downloads and decrypts metadata
//func getMetadata(node *core.IpfsNode, dataId string, key []byte) (*images.Metadata, error) {
//	metacipher, err := ipfs.DataAtPath(node, fmt.Sprintf("%s/meta", dataId))
//	if err != nil {
//		return nil, err
//	}
//	metaplain, err := crypto.DecryptAES(metacipher, key)
//	if err != nil {
//		return nil, err
//	}
//	var meta *images.Metadata
//	if metaplain != nil {
//		if err := json.Unmarshal(metaplain, &meta); err != nil {
//			return nil, err
//		}
//	}
//	return meta, nil
//}

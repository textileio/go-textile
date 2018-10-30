package core

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/images"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
)

// AddPhoto adds an outgoing photo block
func (t *Thread) AddPhoto(dataId string, caption string, key []byte) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	msg := &pb.ThreadData{
		Type:    pb.ThreadData_PHOTO,
		Data:    dataId,
		Key:     key,
		Caption: caption,
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_DATA, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	meta, err := getMetadata(t.node(), dataId, key)
	if err != nil {
		return nil, err
	}
	dconf := &repo.DataBlockConfig{
		DataId:       dataId,
		DataKey:      key,
		DataCaption:  caption,
		DataMetadata: meta,
	}
	if err := t.indexBlock(res, repo.PhotoBlock, dconf); err != nil {
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

	log.Debugf("added DATA to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// handleDataBlock handles an incoming data block
func (t *Thread) handleDataBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadData, error) {
	msg := new(pb.ThreadData)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId:      msg.Data,
		DataKey:     msg.Key,
		DataCaption: msg.Caption,
	}
	switch msg.Type {
	case pb.ThreadData_PHOTO:
		// check if this block has been ignored, if so, don't pin locally, but we still want the block
		var ignore bool
		ignored := t.datastore.Blocks().GetByData(fmt.Sprintf("ignore-%s", hash.B58String()))
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
		if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.PhotoBlock, dconf); err != nil {
			return nil, err
		}
	case pb.ThreadData_TEXT:
		// TODO: chat
		break
	}
	return msg, nil
}

// getMetadata downloads and decrypts metadata
func getMetadata(node *core.IpfsNode, dataId string, key []byte) (*images.Metadata, error) {
	metacipher, err := ipfs.DataAtPath(node, fmt.Sprintf("%s/meta", dataId))
	if err != nil {
		return nil, err
	}
	metaplain, err := crypto.DecryptAES(metacipher, key)
	if err != nil {
		return nil, err
	}
	var meta *images.Metadata
	if metaplain != nil {
		if err := json.Unmarshal(metaplain, &meta); err != nil {
			return nil, err
		}
	}
	return meta, nil
}

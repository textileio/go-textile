package core

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/photo"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
)

// AddPhoto adds an outgoing photo block
func (t *Thread) AddPhoto(dataId string, caption string, key string) (mh.Multihash, error) {
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

// HandleDataBlock handles an incoming data block
func (t *Thread) HandleDataBlock(from *peer.ID, hash mh.Multihash, block *pb.ThreadBlock, following bool) (*pb.ThreadData, error) {
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
			if err := ipfs.PinPath(t.node(), fmt.Sprintf("%s/thumb", msg.Data), false); err != nil {
				return nil, err
			}
			if err := ipfs.PinPath(t.node(), fmt.Sprintf("%s/small", msg.Data), false); err != nil {
				log.Warningf("photo set missing small size")
			}

			// get metadata
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

	// back prop
	newPeers, err := t.FollowParents(block.Header.Parents, from)
	if err != nil {
		return nil, err
	}

	// handle HEAD
	if following {
		return msg, nil
	}
	if _, err := t.handleHead(hash, block.Header.Parents); err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

// getMetadata downloads and decrypts metadata
func getMetadata(node *core.IpfsNode, dataId string, key string) (*photo.Metadata, error) {
	metacipher, err := ipfs.GetDataAtPath(node, fmt.Sprintf("%s/meta", dataId))
	if err != nil {
		return nil, err
	}
	metaplain, err := crypto.DecryptAES(metacipher, []byte(key))
	if err != nil {
		return nil, err
	}
	var metadata *photo.Metadata
	if err := json.Unmarshal(metaplain, &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

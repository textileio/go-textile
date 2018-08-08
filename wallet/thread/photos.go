package thread

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// AddPhoto adds an outgoing photo block
func (t *Thread) AddPhoto(dataId string, caption string, key []byte) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// encrypt AES key with thread pk
	keyCipher, err := t.Encrypt(key)
	if err != nil {
		return nil, err
	}

	// encrypt caption with thread pk
	var captionCipher []byte
	if caption != "" {
		captionCipher, err = t.Encrypt([]byte(caption))
		if err != nil {
			return nil, err
		}
	}

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadData{
		Header:        header,
		Type:          pb.ThreadData_PHOTO,
		DataId:        dataId,
		KeyCipher:     keyCipher,
		CaptionCipher: captionCipher,
	}

	// commit to ipfs
	message, addr, err := t.commitBlock(content, pb.Message_THREAD_DATA)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId:            dataId,
		DataKeyCipher:     keyCipher,
		DataCaptionCipher: captionCipher,
	}
	if err := t.indexBlock(id, header, repo.PhotoBlock, dconf); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// post it
	t.post(message, id, t.Peers())

	log.Debugf("added photo to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// HandleDataBlock handles an incoming data block
func (t *Thread) HandleDataBlock(from *peer.ID, message *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadData, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadData)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, err
		}
	}

	// add to ipfs
	addr, err := t.addBlock(message)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// check if we aleady have this block indexed
	// (should only happen if a misbehaving peer keeps sending the same block)
	index := t.blocks().Get(id)
	if index != nil {
		return nil, nil
	}

	// get the author id
	authorPk, err := libp2pc.UnmarshalPublicKey(content.Header.AuthorPk)
	if err != nil {
		return nil, err
	}
	authorId, err := peer.IDFromPublicKey(authorPk)
	if err != nil {
		return nil, err
	}

	// add author as a new local peer, just in case we haven't found this peer yet.
	// double-check not self in case we're re-discovering the thread
	self := authorId.Pretty() == t.ipfs().Identity.Pretty()
	if !self {
		newPeer := &repo.Peer{
			Row:      ksuid.New().String(),
			Id:       authorId.Pretty(),
			ThreadId: libp2pc.ConfigEncodeKey(content.Header.ThreadPk),
			PubKey:   content.Header.AuthorPk,
		}
		if err := t.peers().Add(newPeer); err != nil {
			// TODO: #202 (Properly handle database/sql errors)
			log.Warningf("peer with id %s already exists in thread %s", newPeer.Id, t.Id)
		}
	}

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId:            content.DataId,
		DataKeyCipher:     content.KeyCipher,
		DataCaptionCipher: content.CaptionCipher,
	}
	switch content.Type {
	case pb.ThreadData_PHOTO:
		// pin data first (it may not be available)
		if err := util.PinPath(t.ipfs(), fmt.Sprintf("%s/thumb", content.DataId), false); err != nil {
			return nil, err
		}
		if err := util.PinPath(t.ipfs(), fmt.Sprintf("%s/meta", content.DataId), false); err != nil {
			return nil, err
		}
		if err := util.PinPath(t.ipfs(), fmt.Sprintf("%s/pk", content.DataId), false); err != nil {
			return nil, err
		}
		if err := t.indexBlock(id, content.Header, repo.PhotoBlock, dconf); err != nil {
			return nil, err
		}
	case pb.ThreadData_TEXT:
		// TODO: comments
		break
	}

	// back prop
	newPeers, err := t.FollowParents(content.Header.Parents, from)
	if err != nil {
		return nil, err
	}

	// handle HEAD
	if following {
		return addr, nil
	}
	if _, err := t.handleHead(id, content.Header.Parents); err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return nil, err
		}
	}

	return addr, nil
}

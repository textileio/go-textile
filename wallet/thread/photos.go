package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// AddPhoto adds an outgoing photo block
func (t *Thread) AddPhoto(sourceId string, caption string, key []byte) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// encrypt AES key with thread pk
	keyCipher, err := t.Encrypt(key)
	if err != nil {
		return nil, err
	}

	// encrypt caption with thread pk
	captionCipher, err := t.Encrypt([]byte(caption))
	if err != nil {
		return nil, err
	}

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadData{
		Header:        header,
		Type:          pb.ThreadData_PHOTO,
		DataId:        sourceId,
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
	if err := t.indexBlock(id, header, repo.PhotoBlock, &sourceId, keyCipher); err != nil {
		return nil, err
	}

	// post it
	t.post(message, id, t.Peers())

	// all done
	return addr, nil
}

// HandleDataBlock handles an incoming data block
func (t *Thread) HandleDataBlock(message *pb.Message, signed *pb.SignedThreadBlock, content *pb.ThreadData) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadData)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, err
		}
	}

	// verify author sig
	if err := t.verifyAuthor(signed, content.Header); err != nil {
		return nil, err
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
		return nil, err
	}

	// remove peer
	authorPk, err := libp2pc.UnmarshalPublicKey(content.Header.AuthorPk)
	if err != nil {
		return nil, err
	}
	authorId, err := peer.IDFromPublicKey(authorPk)
	if err != nil {
		return nil, err
	}
	if err := t.peers().Delete(authorId.Pretty(), t.Id); err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.PhotoBlock, &content.DataId, content.KeyCipher); err != nil {
		return nil, err
	}

	// back prop
	if err := t.followParents(content.Header.Parents); err != nil {
		return nil, err
	}

	return addr, nil
}

package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	"time"
)

// AddExternalInvite creates an outgoing external invite
func (t *Thread) AddExternalInvite() (mh.Multihash, []byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// generate an aes key
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, nil, err
	}

	// encypt thread secret with the key
	threadSk, err := t.PrivKey.Bytes()
	if err != nil {
		return nil, nil, err
	}
	threadSkCipher, err := crypto.EncryptAES(threadSk, key)
	if err != nil {
		return nil, nil, err
	}

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, nil, err
	}
	content := &pb.ThreadExternalInvite{
		Header:        header,
		SkCipher:      threadSkCipher,
		SuggestedName: t.Name,
	}

	// commit to ipfs
	message, addr, err := t.commitBlock(content, pb.Message_THREAD_EXTERNAL_INVITE)
	if err != nil {
		return nil, nil, err
	}
	id := addr.B58String()

	// index it locally
	if err := t.indexBlock(id, header, repo.ExternalInviteBlock, nil); err != nil {
		return nil, nil, err
	}

	// post it
	t.post(message, id, t.Peers())

	log.Debugf("added external invite to %s: %s", t.Id, id)

	// all done
	return addr, key, nil
}

// HandleExternalInviteBlock handles an incoming external invite block
func (t *Thread) HandleExternalInviteBlock(message *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadExternalInvite) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadExternalInvite)
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

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.ExternalInviteBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	if err := t.FollowParents(content.Header.Parents); err != nil {
		return nil, err
	}

	return addr, nil
}

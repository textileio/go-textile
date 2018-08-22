package thread

import (
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
)

// AddExternalInvite creates an external invite, which can be retrieved by any peer
// and does not become part of the hash chain
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
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, nil, err
	}
	content := &pb.ThreadExternalInvite{
		Header:        header,
		SkCipher:      threadSkCipher,
		SuggestedName: t.Name,
	}

	// commit to ipfs
	_, addr, err := t.commitBlock(content, pb.Message_THREAD_EXTERNAL_INVITE)
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("created EXTERNAL_INVITE for %s", t.Id)

	// all done
	return addr, key, nil
}

// HandleExternalInviteMessage handles an incoming external invite
// - this happens right before a join
// - the invite is not kept on-chain, so we only need to follow parents and update HEAD
func (t *Thread) HandleExternalInviteMessage(content *pb.ThreadExternalInvite) error {
	// back prop
	if _, err := t.FollowParents(content.Header.Parents, nil); err != nil {
		return err
	}

	// update HEAD if parents of the invite are actual updates
	if len(content.Header.Parents) > 0 {
		if err := t.updateHead(content.Header.Parents[0]); err != nil {
			return err
		}
	}

	return nil
}

package core

import (
	peer "github.com/libp2p/go-libp2p-peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/pb"
)

// AddInvite creates an outgoing add block, which is sent directly to the recipient
// and does not become part of the hash chain
func (t *Thread) AddInvite(p *pb.Peer) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.shareable(t.config.Account.Address, p.Address) {
		return nil, ErrNotShareable
	}

	self := t.datastore.Peers().Get(t.node().Identity.Pretty())
	msg := &pb.ThreadAdd{
		Thread:  t.datastore.Threads().Get(t.Id),
		Inviter: self,
		Invitee: p.Id,
	}

	pid, err := peer.IDB58Decode(p.Id)
	if err != nil {
		return nil, err
	}
	pk, err := pid.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.Block_ADD, true, func(plaintext []byte) ([]byte, error) {
		return crypto.Encrypt(pk, plaintext)
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("created ADD to %s for %s", p.Id, t.Id)

	return res.hash, nil
}

// AddExternalInvite creates an add block, which can be retrieved by any peer
// and does not become part of the hash chain
func (t *Thread) AddExternalInvite() (mh.Multihash, []byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	self := t.datastore.Peers().Get(t.node().Identity.Pretty())
	msg := &pb.ThreadAdd{
		Thread:  t.datastore.Threads().Get(t.Id),
		Inviter: self,
	}

	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, nil, err
	}

	res, err := t.commitBlock(msg, pb.Block_ADD, true, func(plaintext []byte) ([]byte, error) {
		return crypto.EncryptAES(plaintext, key)
	})
	if err != nil {
		return nil, nil, err
	}
	nhash, err := t.commitNode(res.hash.B58String(), nil, false)
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("created external ADD for %s", t.Id)

	return nhash, key, nil
}

// handleAddBlock handles an incoming add.
// This happens right before a join. The invite is not kept on-chain,
// so we only need to follow parents and update HEAD.
func (t *Thread) handleAddBlock(parents []string) error {
	_, err := t.followParents(parents)
	if err != nil {
		return err
	}

	// update HEAD if parents of the invite are actual updates
	if len(parents) > 0 {
		err = t.updateHead(parents)
		if err != nil {
			return err
		}
	}
	return nil
}

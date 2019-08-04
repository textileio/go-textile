package core

import (
	peer "github.com/libp2p/go-libp2p-core/peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/pb"
)

// AddInvite creates an outgoing add block, which is sent directly to the recipient
// and does not become part of the hash chain
func (t *Thread) AddInvite(p *pb.Peer) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

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

	// add directly, no need for an update event which happens w/ indexBlock
	// Note: this will be deleted after posting
	err = t.datastore.Blocks().Add(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_ADD,
		Date:   res.header.Date,
		Body:   msg.Invitee, // ugly and tmp way to retain invitee address when posting
		Status: pb.Block_QUEUED,
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("added ADD to %s for %s", p.Id, t.Id)

	return res.hash, nil
}

// AddExternalInvite creates an add block, which can be retrieved by any peer
// and does not become part of the hash chain
func (t *Thread) AddExternalInvite() (mh.Multihash, []byte, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

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
	nhash, err := t.commitNode(&pb.Block{Id: res.hash.B58String()}, nil, false)
	if err != nil {
		return nil, nil, err
	}

	// add directly, no need for an update event which happens w/ indexBlock
	// Note: this will be deleted after posting (technically not needed at all because their won't
	// be any recipients, but since this logic is tied to sync group management, creating a dummy
	// block here is the cleanest solution at this point).
	err = t.datastore.Blocks().Add(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_ADD,
		Date:   res.header.Date,
		Status: pb.Block_QUEUED,
	})
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("added external ADD for %s", t.Id)

	return nhash, key, nil
}

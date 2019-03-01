package mobile

import (
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
)

// AddInvite adds a new invite to a thread
func (m *Mobile) AddInvite(threadId string, inviteeId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	pid, err := peer.IDB58Decode(inviteeId)
	if err != nil {
		return "", err
	}

	hash, err := thrd.AddInvite(pid)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

// AddExternalInvite generates a new external invite link to a thread
func (m *Mobile) AddExternalInvite(threadId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return nil, core.ErrThreadNotFound
	}

	hash, key, err := thrd.AddExternalInvite()
	if err != nil {
		return nil, err
	}

	username, _ := m.Username()
	invite := &pb.NewInvite{
		Id:      hash.B58String(),
		Key:     base58.FastBase58Encoding(key),
		Inviter: username,
	}

	return proto.Marshal(invite)
}

// AcceptExternalInvite notifies the thread of a join
func (m *Mobile) AcceptExternalInvite(id string, key string) (string, error) {
	if !m.node.Online() {
		return "", core.ErrOffline
	}

	keyb, err := base58.Decode(key)
	if err != nil {
		return "", err
	}

	hash, err := m.node.AcceptExternalInvite(id, keyb)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

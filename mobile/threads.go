package mobile

import (
	"crypto/rand"

	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
)

// AddThread adds a new thread with the given name
func (m *Mobile) AddThread(config []byte) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	conf := new(pb.AddThreadConfig)
	if err := proto.Unmarshal(config, conf); err != nil {
		return nil, err
	}

	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	thrd, err := m.node.AddThread(*conf, sk, m.node.Account().Address(), true)
	if err != nil {
		return nil, err
	}

	view, err := m.node.ThreadView(thrd.Id)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(view)
}

// AddOrUpdateThread calls core AddOrUpdateThread
func (m *Mobile) AddOrUpdateThread(thrd []byte) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	mthrd := new(pb.Thread)
	if err := proto.Unmarshal(thrd, mthrd); err != nil {
		return err
	}

	return m.node.AddOrUpdateThread(mthrd)
}

// RemoveThread call core RemoveThread
func (m *Mobile) RemoveThread(id string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	hash, err := m.node.RemoveThread(id)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

// Threads lists all threads
func (m *Mobile) Threads() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	views := &pb.ThreadList{
		Items: make([]*pb.Thread, 0),
	}
	for _, thrd := range m.node.Threads() {
		view, err := m.node.ThreadView(thrd.Id)
		if err != nil {
			return nil, err
		}
		views.Items = append(views.Items, view)
	}

	return proto.Marshal(views)
}

// Thread calls core Thread
func (m *Mobile) Thread(threadId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	view, err := m.node.ThreadView(threadId)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(view)
}

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

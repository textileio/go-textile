package mobile

import (
	"crypto/rand"

	libp2pc "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"

	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
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

// RenameThread call core RenameThread
func (m *Mobile) RenameThread(id string, name string) error {
	return m.node.RenameThread(id, name)
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

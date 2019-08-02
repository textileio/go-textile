package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/go-textile/core"
)

// AddInvite call core AddInvite
func (m *Mobile) AddInvite(threadId string, address string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.AddInvite(threadId, address)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// AddExternalInvite generates a new external invite link to a thread
func (m *Mobile) AddExternalInvite(threadId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	invite, err := m.node.AddExternalInvite(threadId)
	if err != nil {
		return nil, err
	}

	m.node.FlushCafes()

	return proto.Marshal(invite)
}

// Invites calls core Invites
func (m *Mobile) Invites() ([]byte, error) {
	return proto.Marshal(m.node.Invites())
}

// AcceptInvite calls core AcceptInvite
func (m *Mobile) AcceptInvite(id string) (string, error) {
	if !m.node.Online() {
		return "", core.ErrOffline
	}

	hash, err := m.node.AcceptInvite(id)
	if err != nil {
		return "", err
	}

	m.node.FlushCafes()

	return hash.B58String(), nil
}

// AcceptExternalInvite calls core AcceptExternalInvite
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

	m.node.FlushCafes()

	return hash.B58String(), nil
}

// IgnoreInvite calls core IgnoreInvite
func (m *Mobile) IgnoreInvite(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.IgnoreInvite(id)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

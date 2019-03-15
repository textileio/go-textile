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

	return m.node.AddInvite(threadId, address)
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

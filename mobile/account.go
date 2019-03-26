package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
)

// Address returns account address
func (m *Mobile) Address() string {
	if !m.node.Started() {
		return ""
	}
	return m.node.Account().Address()
}

// Seed returns account seed
func (m *Mobile) Seed() string {
	if !m.node.Started() {
		return ""
	}
	return m.node.Account().Seed()
}

// Encrypt calls core Encrypt
func (m *Mobile) Encrypt(input []byte) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}
	return m.node.Encrypt(input)
}

// Decrypt call core Decrypt
func (m *Mobile) Decrypt(input []byte) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}
	return m.node.Decrypt(input)
}

// AccountContact calls core AccountContact
func (m *Mobile) AccountContact() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.AccountContact())
}

// SyncAccount calls core SyncAccount
func (m *Mobile) SyncAccount(options []byte) (*SearchHandle, error) {
	if !m.node.Online() {
		return nil, core.ErrOffline
	}

	moptions := new(pb.QueryOptions)
	if err := proto.Unmarshal(options, moptions); err != nil {
		return nil, err
	}

	cancel, err := m.node.SyncAccount(moptions)
	if err != nil {
		return nil, err
	}

	return &SearchHandle{
		Id:     ksuid.New().String(),
		cancel: cancel,
		done:   func() {},
	}, nil
}

package mobile

import (
	"fmt"

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

// Sign calls core Sign
func (m *Mobile) Sign(input []byte) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}
	return m.node.Sign(input)
}

// Verify calls core verify
func (m *Mobile) Verify(input []byte, sig []byte) error {
	if !m.node.Started() {
		return core.ErrStopped
	}
	return m.node.Verify(input, sig)
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

// AccountThread calls core AccountThread
func (m *Mobile) AccountThread() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrd := m.node.AccountThread()
	if thrd == nil {
		return nil, fmt.Errorf("account thread not found")
	}
	view, err := m.node.ThreadView(thrd.Id)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(view)
}

// AccountContact calls core AccountContact
func (m *Mobile) AccountContact() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	contact := m.node.AccountContact()
	if contact == nil {
		return nil, fmt.Errorf("self contact not found")
	}

	return proto.Marshal(contact)
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

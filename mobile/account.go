package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
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

// AccountPeers calls core AccountPeers
func (m *Mobile) AccountPeers(input []byte) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.AccountPeers())
}

// FindThreadBackups calls core FindThreadBackups
func (m *Mobile) FindThreadBackups(query []byte, options []byte) (*SearchHandle, error) {
	if !m.node.Online() {
		return nil, core.ErrOffline
	}

	mquery := new(pb.ThreadBackupQuery)
	if err := proto.Unmarshal(query, mquery); err != nil {
		return nil, err
	}
	moptions := new(pb.QueryOptions)
	if err := proto.Unmarshal(options, moptions); err != nil {
		return nil, err
	}

	resCh, errCh, cancel, err := m.node.FindThreadBackups(mquery, moptions)
	if err != nil {
		return nil, err
	}

	return m.handleSearchStream(resCh, errCh, cancel)
}

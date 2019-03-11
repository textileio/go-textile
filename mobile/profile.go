package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
)

// Profile calls core Profile
func (m *Mobile) Profile() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	self := m.node.Profile()
	if self == nil {
		return nil, nil
	}

	return proto.Marshal(self)
}

// Name calls core Name
func (m *Mobile) Name() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	return m.node.Name(), nil
}

// SetName calls core SetName
func (m *Mobile) SetName(username string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	return m.node.SetName(username)
}

// Avatar calls core Avatar
func (m *Mobile) Avatar() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	return m.node.Avatar(), nil
}

// SetAvatar calls core SetAvatar
func (m *Mobile) SetAvatar(hash string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	return m.node.SetAvatar(hash)
}

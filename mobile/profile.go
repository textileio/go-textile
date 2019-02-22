package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
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

// Username calls core Username
func (m *Mobile) Username() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	return m.node.Username(), nil
}

// SetUsername calls core SetUsername
func (m *Mobile) SetUsername(username string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	return m.node.SetUsername(username)
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

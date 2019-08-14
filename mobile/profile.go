package mobile

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/core"
)

// Profile calls core Profile
func (m *Mobile) Profile() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	self := m.node.Profile()
	if self == nil {
		return nil, fmt.Errorf("profile not found")
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

	err := m.node.SetName(username)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// Avatar calls core Avatar
func (m *Mobile) Avatar() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	return m.node.Avatar(), nil
}

// SetAvatar adds the image at pth to the account thread and calls core SetAvatar
func (m *Mobile) SetAvatar(pth string, cb ProtoCallback) {
	m.node.WaitAdd(1, "Mobile.SetAvatar")
	go func() {
		defer m.node.WaitDone("Mobile.SetAvatar")

		hash, err := m.setAvatar(pth)
		if err != nil {
			cb.Call(nil, err)
			return
		}

		cb.Call(m.blockView(hash))
	}()
}

func (m *Mobile) setAvatar(pth string) (mh.Multihash, error) {
	if !m.node.Online() {
		return nil, core.ErrOffline
	}

	thrd := m.node.AccountThread()
	if thrd == nil {
		return nil, fmt.Errorf("account thread not found")
	}

	hash, err := m.addFiles([]string{pth}, thrd.Id, "")
	if err != nil {
		return nil, err
	}

	err = m.node.SetAvatar()
	if err != nil {
		return nil, err
	}

	m.node.FlushCafes()

	return hash, nil
}

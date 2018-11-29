package mobile

import (
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/textileio/textile-go/core"
)

// SetUsername calls core SetUsername
func (m *Mobile) SetUsername(username string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	return m.node.SetUsername(username)
}

// Username calls core Username
func (m *Mobile) Username() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	username, err := m.node.Username()
	if err != nil {
		return "", err
	}
	if username == nil {
		return "", nil
	}
	return *username, nil
}

// SetAvatar calls core SetAvatar
func (m *Mobile) SetAvatar(hash string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	return m.node.SetAvatar(hash)
}

// Profile returns the local profile
func (m *Mobile) Profile() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	id, err := m.node.PeerId()
	if err != nil {
		return "", err
	}
	prof, err := m.node.Profile(id)
	if err != nil {
		return "", err
	}
	return toJSON(prof)
}

// PeerProfile looks up a profile by id
func (m *Mobile) PeerProfile(peerId string) (string, error) {
	if !m.node.Online() {
		return "", core.ErrOffline
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return "", err
	}
	prof, err := m.node.Profile(pid)
	if err != nil {
		return "", err
	}
	return toJSON(prof)
}

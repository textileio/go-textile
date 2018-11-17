package mobile

import (
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// SetUsername calls core SetUsername
func (m *Mobile) SetUsername(username string) error {
	return m.node.SetUsername(username)
}

// Username calls core Username
func (m *Mobile) Username() (string, error) {
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
func (m *Mobile) SetAvatar(id string) error {
	return m.node.SetAvatar(id)
}

// Profile returns the local profile
func (m *Mobile) Profile() (string, error) {
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

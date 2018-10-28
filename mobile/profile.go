package mobile

import (
	"github.com/textileio/textile-go/core"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// SetUsername calls core SetUsername
func (m *Mobile) SetUsername(username string) error {
	return core.Node.SetUsername(username)
}

// GetUsername calls core GetUsername
func (m *Mobile) GetUsername() (string, error) {
	username, err := core.Node.GetUsername()
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
	return core.Node.SetAvatar(id)
}

// GetProfile returns the local profile
func (m *Mobile) GetProfile() (string, error) {
	id, err := core.Node.PeerId()
	if err != nil {
		return "", err
	}
	prof, err := core.Node.GetProfile(id)
	if err != nil {
		return "", err
	}
	return toJSON(prof)
}

// GetPeerProfile looks up a profile by id
func (m *Mobile) GetPeerProfile(peerId string) (string, error) {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return "", err
	}
	prof, err := core.Node.GetProfile(pid)
	if err != nil {
		return "", err
	}
	return toJSON(prof)
}

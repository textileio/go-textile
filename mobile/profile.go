package mobile

import (
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/util"
)

// GetId calls core GetId
func (m *Mobile) GetId() (string, error) {
	id, err := core.Node.Wallet.GetId()
	if err != nil {
		return "", err
	}
	if id == nil {
		return "", nil
	}
	return id.Pretty(), nil
}

// GetPubKey returns the profile public key string
func (m *Mobile) GetPubKey() (string, error) {
	key, err := core.Node.Wallet.GetKey()
	if err != nil {
		return "", err
	}
	if key == nil {
		return "", nil
	}
	return util.EncodeKey(key.GetPublic())
}

// SetUsername calls core SetUsername
func (m *Mobile) SetUsername(username string) error {
	return core.Node.Wallet.SetUsername(username)
}

// GetUsername calls core GetUsername
func (m *Mobile) GetUsername() (string, error) {
	username, err := core.Node.Wallet.GetUsername()
	if err != nil {
		return "", err
	}
	if username == nil {
		return "", nil
	}
	return *username, nil
}

// SetAvatarId calls core SetAvatarId
func (m *Mobile) SetAvatarId(id string) error {
	return core.Node.Wallet.SetAvatarId(id)
}

// GetProfile returns the local profile
func (m *Mobile) GetProfile() (string, error) {
	id, err := core.Node.Wallet.GetId()
	if err != nil {
		log.Errorf("error getting id %s: %s", id, err)
		return "", err
	}
	if id == nil {
		return "", errors.New("profile does not exist")
	}
	prof, err := core.Node.Wallet.GetProfile(id.Pretty())
	if err != nil {
		log.Errorf("error getting profile %s: %s", id, err)
		return "", err
	}
	return toJSON(prof)
}

// GetOtherProfile looks up a profile by id
func (m *Mobile) GetOtherProfile(peerId string) (string, error) {
	prof, err := core.Node.Wallet.GetProfile(peerId)
	if err != nil {
		log.Errorf("error getting profile %s: %s", peerId, err)
		return "", err
	}
	return toJSON(prof)
}

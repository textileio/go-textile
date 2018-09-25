package mobile

import (
	"github.com/textileio/textile-go/core"
)

// GetID calls core GetID
func (m *Mobile) GetID() (string, error) {
	id, err := core.Node.GetID()
	if err != nil {
		return "", err
	}
	return id.Pretty(), nil
}

// GetAddress returns account address
func (m *Mobile) GetAddress() (string, error) {
	accnt, err := core.Node.GetAccount()
	if err != nil {
		return "", err
	}
	if accnt == nil {
		return "", nil
	}
	return accnt.Address(), nil
}

// GetSeed returns account seed
func (m *Mobile) GetSeed() (string, error) {
	accnt, err := core.Node.GetAccount()
	if err != nil {
		return "", err
	}
	if accnt == nil {
		return "", nil
	}
	return accnt.Seed(), nil
}

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

// SetAvatarId calls core SetAvatarId
func (m *Mobile) SetAvatarId(id string) error {
	return core.Node.SetAvatarId(id)
}

// GetProfile returns the local profile
func (m *Mobile) GetProfile() (string, error) {
	id, err := core.Node.GetID()
	if err != nil {
		log.Errorf("error getting profile (get id): %s", err)
		return "", err
	}
	prof, err := core.Node.GetProfile(id.Pretty())
	if err != nil {
		log.Errorf("error getting profile %s: %s", id, err)
		return "", err
	}
	return toJSON(prof)
}

// GetOtherProfile looks up a profile by id
func (m *Mobile) GetOtherProfile(peerId string) (string, error) {
	prof, err := core.Node.GetProfile(peerId)
	if err != nil {
		log.Errorf("error getting profile %s: %s", peerId, err)
		return "", err
	}
	return toJSON(prof)
}

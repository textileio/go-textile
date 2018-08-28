package mobile

import (
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core"
)

// SignUpWithEmail creates an email based registration and calls core signup
func (m *Mobile) SignUpWithEmail(email string, username string, password string, referral string) error {
	// build registration
	reg := &models.Registration{
		Username: username,
		Password: password,
		Identity: &models.Identity{
			Type:  models.EmailAddress,
			Value: email,
		},
		Referral: referral,
	}
	return core.Node.Wallet.SignUp(reg)
}

// SignIn build credentials and calls core SignIn
func (m *Mobile) SignIn(username string, password string) error {
	// build creds
	creds := &models.Credentials{
		Username: username,
		Password: password,
	}
	return core.Node.Wallet.SignIn(creds)
}

// SignOut calls core SignOut
func (m *Mobile) SignOut() error {
	return core.Node.Wallet.SignOut()
}

// IsSignedIn calls core IsSignedIn
func (m *Mobile) IsSignedIn() bool {
	si, _ := core.Node.Wallet.IsSignedIn()
	return si
}

// GetId calls core GetId
func (m *Mobile) GetId() (string, error) {
	return core.Node.Wallet.GetId()
}

// GetPubKey calls core GetPubKeyString
func (m *Mobile) GetPubKey() (string, error) {
	return core.Node.Wallet.GetPubKeyString()
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

// GetTokens calls core GetTokens
func (m *Mobile) GetTokens() (string, error) {
	tokens, err := core.Node.Wallet.GetTokens()
	if err != nil {
		return "", err
	}
	if tokens == nil {
		return "", nil
	}
	return toJSON(tokens)
}

// SetAvatarId calls core SetAvatarId
func (m *Mobile) SetAvatarId(id string) error {
	return core.Node.Wallet.SetAvatarId(id)
}

// GetProfile returns this peer's profile
func (m *Mobile) GetProfile() (string, error) {
	id, err := core.Node.Wallet.GetId()
	if err != nil {
		log.Errorf("error getting id %s: %s", id, err)
		return "", err
	}
	prof, err := core.Node.Wallet.GetProfile(id)
	if err != nil {
		log.Errorf("error getting profile %s: %s", id, err)
		return "", err
	}
	return toJSON(prof)
}

// GetPeerProfile uses a peer id to look up a profile
func (m *Mobile) GetPeerProfile(peerId string) (string, error) {
	prof, err := core.Node.Wallet.GetProfile(peerId)
	if err != nil {
		log.Errorf("error getting profile %s: %s", peerId, err)
		return "", err
	}
	return toJSON(prof)
}

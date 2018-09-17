package mobile

import (
	"github.com/textileio/textile-go/core"
)

// CafeRegister calls core CafeRegister
func (m *Mobile) CafeRegister(referral string) error {
	return core.Node.Wallet.CafeRegister(referral)
}

// CafeLogin calls core CafeLogin
func (m *Mobile) CafeLogin() error {
	return core.Node.Wallet.CafeLogin()
}

// GetCafeTokens calls core GetCafeTokens
func (m *Mobile) GetCafeTokens(forceRefresh bool) (string, error) {
	tokens, err := core.Node.Wallet.GetCafeTokens(forceRefresh)
	if err != nil {
		return "", err
	}
	if tokens == nil {
		return "", nil
	}
	return toJSON(tokens)
}

// CafeLogout calls core CafeLogout
func (m *Mobile) CafeLogout() error {
	return core.Node.Wallet.CafeLogout()
}

// CafeLoggedIn calls core CafeLoggedIn
func (m *Mobile) CafeLoggedIn() bool {
	loggedIn, _ := core.Node.Wallet.CafeLoggedIn()
	return loggedIn
}

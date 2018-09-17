package mobile

import (
	"github.com/textileio/textile-go/core"
)

// CafeRegister calls core CafeRegister with the current profile key
func (m *Mobile) CafeRegister(referral string) error {
	key, err := core.Node.Wallet.GetKey()
	if err != nil {
		return err
	}
	return core.Node.Wallet.CafeRegister(key, referral)
}

// CafeLogin calls core CafeLogin with the current profile key
func (m *Mobile) CafeLogin() error {
	key, err := core.Node.Wallet.GetKey()
	if err != nil {
		return err
	}
	return core.Node.Wallet.CafeLogin(key)
}

// CafeLogout calls core CafeLogout
func (m *Mobile) CafeLogout() error {
	return core.Node.Wallet.CafeLogout()
}

// IsSignedIn calls core IsSignedIn
func (m *Mobile) IsSignedIn() bool {
	loggedIn, _ := core.Node.Wallet.CafeLoggedIn()
	return loggedIn
}

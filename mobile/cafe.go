package mobile

import "github.com/textileio/textile-go/core"

// CafeRegister calls core CafeRegister
func (m *Mobile) CafeRegister(referral string) error {
	return core.Node.CafeRegister(referral)
}

// CafeLogin calls core CafeLogin
func (m *Mobile) CafeLogin() error {
	return core.Node.CafeLogin()
}

// GetCafeTokens calls core GetCafeTokens
func (m *Mobile) GetCafeTokens(forceRefresh bool) (string, error) {
	tokens, err := core.Node.GetCafeTokens(forceRefresh)
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
	return core.Node.CafeLogout()
}

// CafeLoggedIn calls core CafeLoggedIn
func (m *Mobile) CafeLoggedIn() bool {
	loggedIn, _ := core.Node.CafeLoggedIn()
	return loggedIn
}

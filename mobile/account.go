package mobile

import (
	"github.com/textileio/textile-go/core"
)

// Address returns account address
func (m *Mobile) Address() (string, error) {
	accnt, err := core.Node.Account()
	if err != nil {
		return "", err
	}
	if accnt == nil {
		return "", nil
	}
	return accnt.Address(), nil
}

// Seed returns account seed
func (m *Mobile) Seed() (string, error) {
	accnt, err := core.Node.Account()
	if err != nil {
		return "", err
	}
	if accnt == nil {
		return "", nil
	}
	return accnt.Seed(), nil
}

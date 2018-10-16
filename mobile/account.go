package mobile

import (
	"github.com/textileio/textile-go/core"
)

// Id calls core Id
func (m *Mobile) GetID() (string, error) {
	id, err := core.Node.Id()
	if err != nil {
		return "", err
	}
	return id.Pretty(), nil
}

// GetAddress returns account address
func (m *Mobile) GetAddress() (string, error) {
	accnt, err := core.Node.Account()
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
	accnt, err := core.Node.Account()
	if err != nil {
		return "", err
	}
	if accnt == nil {
		return "", nil
	}
	return accnt.Seed(), nil
}

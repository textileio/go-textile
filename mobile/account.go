package mobile

import (
	"github.com/textileio/textile-go/core"
)

// Address returns account address
func (m *Mobile) Address() string {
	return core.Node.Account().Address()
}

// Seed returns account seed
func (m *Mobile) Seed() string {
	return core.Node.Account().Seed()
}

package mobile

// Address returns account address
func (m *Mobile) Address() string {
	return m.node.Account().Address()
}

// Seed returns account seed
func (m *Mobile) Seed() string {
	return m.node.Account().Seed()
}

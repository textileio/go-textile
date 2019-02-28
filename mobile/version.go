package mobile

import "github.com/textileio/go-textile/common"

// Version returns common Version
func (m *Mobile) Version() string {
	return "v" + common.Version
}

// GitSummary returns common GitSummary
func (m *Mobile) GitSummary() string {
	return common.GitSummary
}

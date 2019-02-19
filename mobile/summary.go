package mobile

import "github.com/textileio/textile-go/core"

// Summary calls core Summary
func (m *Mobile) Summary() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	stats, err := m.node.Summary()
	if err != nil {
		return "", err
	}
	return toJSON(stats)
}

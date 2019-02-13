package mobile

import "github.com/textileio/textile-go/core"

// Feed calls core Feed
func (m *Mobile) Feed(offset string, limit int, threadId string, annotated bool) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	items, err := m.node.Feed(offset, limit, threadId, annotated)
	if err != nil {
		return "", err
	}

	return toJSON(items)
}

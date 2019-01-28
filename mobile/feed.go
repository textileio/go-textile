package mobile

import "github.com/textileio/textile-go/core"

// ThreadFeed calls core ThreadFeed
func (m *Mobile) ThreadFeed(offset string, limit int, threadId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	items, err := m.node.ThreadFeed(offset, limit, threadId)
	if err != nil {
		return "", err
	}

	return toJSON(items)
}

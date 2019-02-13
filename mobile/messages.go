package mobile

import "github.com/textileio/textile-go/core"

// AddThreadMessage adds a message to a thread
func (m *Mobile) AddThreadMessage(threadId string, body string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	hash, err := thrd.AddMessage(body)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

// Messages calls core Messages
func (m *Mobile) Messages(offset string, limit int, threadId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	msgs, err := m.node.Messages(offset, limit, threadId)
	if err != nil {
		return "", err
	}

	return toJSON(msgs)
}

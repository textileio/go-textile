package mobile

import "github.com/textileio/go-textile/core"

// AddComment adds a comment targeted at the given block
func (m *Mobile) AddComment(blockId string, body string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	block, err := m.node.Block(blockId)
	if err != nil {
		return "", err
	}

	thrd := m.node.Thread(block.Thread)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	hash, err := thrd.AddComment(block.Id, body)
	if err != nil {
		return "", err
	}

	m.node.FlushCafes()

	return hash.B58String(), nil
}

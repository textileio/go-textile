package mobile

import "github.com/textileio/textile-go/core"

// AddLike adds a like targeted at the given block
func (m *Mobile) AddLike(blockId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	block, err := m.node.Block(blockId)
	if err != nil {
		return "", err
	}

	thrd := m.node.Thread(block.ThreadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	hash, err := thrd.AddLike(block.Id)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

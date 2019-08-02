package mobile

import "github.com/textileio/go-textile/core"

// AddIgnore adds an ignore targeted at the given block and unpins any associated target data
func (m *Mobile) AddIgnore(blockId string) (string, error) {
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

	hash, err := thrd.AddIgnore(block.Id)
	if err != nil {
		return "", err
	}

	m.node.FlushCafes()

	return hash.B58String(), nil
}

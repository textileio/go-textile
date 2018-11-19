package mobile

import "github.com/textileio/textile-go/core"

// AddThreadIgnore adds an ignore targeted at the given block and unpins any associated target data
func (m *Mobile) AddThreadIgnore(blockId string) (string, error) {
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

	hash, err := thrd.AddIgnore(block.Id)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

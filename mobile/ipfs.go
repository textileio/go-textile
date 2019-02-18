package mobile

import "github.com/textileio/textile-go/core"

// PeerId returns the ipfs peer id
func (m *Mobile) PeerId() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	pid, err := m.node.PeerId()
	if err != nil {
		return "", err
	}
	return pid.Pretty(), nil
}

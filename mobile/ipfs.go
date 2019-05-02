package mobile

import (
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/textileio/go-textile/core"
)

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

// DataAtPath calls core DataAtPath
func (m *Mobile) DataAtPath(pth string) ([]byte, error) {
	if !m.node.Online() {
		return nil, core.ErrOffline
	}

	data, err := m.node.DataAtPath(pth)
	if err != nil {
		if err == ipld.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

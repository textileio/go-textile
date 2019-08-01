package mobile

import (
	"bytes"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/textileio/go-textile/core"
)

// PeerId returns the ipfs peer id
func (m *Mobile) PeerId() (string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.node.Started() {
		return "", core.ErrStopped
	}

	pid, err := m.node.PeerId()
	if err != nil {
		return "", err
	}
	return pid.Pretty(), nil
}

// DataAtPath is the async version of dataAtPath
func (m *Mobile) DataAtPath(pth string, cb DataCallback) {
	go func() {
		m.mux.Lock()
		defer m.mux.Unlock()

		cb.Call(m.dataAtPath(pth))
	}()
}

// dataAtPath calls core DataAtPath
func (m *Mobile) dataAtPath(pth string) ([]byte, string, error) {
	if !m.node.Started() {
		return nil, "", core.ErrStopped
	}

	data, err := m.node.DataAtPath(pth)
	if err != nil {
		if err == ipld.ErrNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	media, err := m.node.GetMedia(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	return data, media, nil
}

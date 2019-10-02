package mobile

import (
	"bytes"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
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

// SwarmConnect opens a new direct connection to a peer using an IPFS multiaddr
func (m *Mobile) SwarmConnect(address string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	results, err := ipfs.SwarmConnect(m.node.Ipfs(), []string{address})
	if err != nil {
		return "", err
	}

	return results[0], nil
}

// DataAtPath is the async version of dataAtPath
func (m *Mobile) DataAtPath(pth string, cb DataCallback) {
	m.node.WaitAdd(1, "Mobile.DataAtPath")
	go func() {
		defer m.node.WaitDone("Mobile.DataAtPath")
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

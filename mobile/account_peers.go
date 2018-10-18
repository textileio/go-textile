package mobile

import (
	"github.com/textileio/textile-go/core"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// AccountPeer is a simple meta data wrapper around an AccountPeer
type AccountPeer struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// AccountPeers is a wrapper around a list of AccountPeers
type AccountPeers struct {
	Items []AccountPeer `json:"items"`
}

// AccountPeers lists all devices
func (m *Mobile) AccountPeers() (string, error) {
	peers := AccountPeers{Items: make([]AccountPeer, 0)}
	for _, dev := range core.Node.AccountPeers() {
		item := AccountPeer{Id: dev.Id, Name: dev.Name}
		peers.Items = append(peers.Items, item)
	}
	return toJSON(peers)
}

// AddAccountPeer calls core AddAccountPeer
func (m *Mobile) AddAccountPeer(id string, name string) error {
	m.waitForOnline()
	pid, err := peer.IDB58Decode(id)
	if err != nil {
		return err
	}
	return core.Node.AddAccountPeer(pid, name)
}

// RemoveAccountPeer call core RemoveAccountPeer
func (m *Mobile) RemoveAccountPeer(id string) error {
	return core.Node.RemoveAccountPeer(id)
}

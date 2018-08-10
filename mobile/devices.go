package mobile

import (
	"github.com/textileio/textile-go/core"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// Device is a simple meta data wrapper around a Device
type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Devices is a wrapper around a list of Devices
type Devices struct {
	Items []Device `json:"items"`
}

// Devices lists all devices
func (m *Mobile) Devices() (string, error) {
	devices := Devices{Items: make([]Device, 0)}
	for _, dev := range core.Node.Wallet.Devices() {
		item := Device{Id: dev.Id, Name: dev.Name}
		devices.Items = append(devices.Items, item)
	}
	return toJSON(devices)
}

// AddDevice calls core AddDevice
func (m *Mobile) AddDevice(name string, pubKey string) error {
	m.waitForOnline()
	pkb, err := libp2pc.ConfigDecodeKey(pubKey)
	if err != nil {
		return err
	}
	pk, err := libp2pc.UnmarshalPublicKey(pkb)
	if err != nil {
		return err
	}
	return core.Node.Wallet.AddDevice(name, pk)
}

// RemoveDevice call core RemoveDevice
func (m *Mobile) RemoveDevice(id string) error {
	return core.Node.Wallet.RemoveDevice(id)
}

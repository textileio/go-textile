package mobile

import "github.com/textileio/textile-go/core"

// RegisterCafe calls core RegisterCafe
func (m *Mobile) RegisterCafe(peerId string) error {
	return core.Node.RegisterCafe(peerId)
}

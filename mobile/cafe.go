package mobile

import "github.com/textileio/textile-go/core"

// CafeRegister calls core CafeRegister
func (m *Mobile) CafeRegister(peerId string) error {
	return core.Node.CafeRegister(peerId)
}

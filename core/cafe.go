package core

import (
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// CafeRegister registers a public key w/ a cafe, requests a session token, and saves it locally
func (t *Textile) CafeRegister(peerId string) error {
	if err := t.touchDatastore(); err != nil {
		return err
	}
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}
	return t.cafeService.Register(pid)
}

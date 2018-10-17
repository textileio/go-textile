package core

import (
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

// GetPeerId returns peer id
func (t *Textile) GetPeerId() (peer.ID, error) {
	if !t.started {
		return "", ErrStopped
	}
	return t.ipfs.Identity, nil
}

// GetPrivKey returns the current peer private key
func (t *Textile) GetPeerPrivKey() (libp2pc.PrivKey, error) {
	if !t.started {
		return nil, ErrStopped
	}
	if t.ipfs.PrivateKey == nil {
		if err := t.ipfs.LoadPrivateKey(); err != nil {
			return nil, err
		}
	}
	return t.ipfs.PrivateKey, nil
}

// GetPeerPubKey returns the current peer public key
func (t *Textile) GetPeerPubKey() (libp2pc.PubKey, error) {
	sk, err := t.GetPeerPrivKey()
	if err != nil {
		return nil, err
	}
	return sk.GetPublic(), nil
}

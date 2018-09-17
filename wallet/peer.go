package wallet

import (
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// GetPeerId returns peer id
func (w *Wallet) GetPeerId() (string, error) {
	if !w.started {
		return "", ErrStopped
	}
	return w.ipfs.Identity.Pretty(), nil
}

// GetPrivKey returns the current peer private key
func (w *Wallet) GetPeerPrivKey() (libp2pc.PrivKey, error) {
	if !w.started {
		return nil, ErrStopped
	}
	if w.ipfs.PrivateKey == nil {
		if err := w.ipfs.LoadPrivateKey(); err != nil {
			return nil, err
		}
	}
	return w.ipfs.PrivateKey, nil
}

// GetPeerPubKey returns the current peer public key
func (w *Wallet) GetPeerPubKey() (libp2pc.PubKey, error) {
	secret, err := w.GetPeerPrivKey()
	if err != nil {
		return nil, err
	}
	return secret.GetPublic(), nil
}

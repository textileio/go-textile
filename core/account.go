package core

import (
	"github.com/textileio/textile-go/keypair"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

// Account returns account keypair
func (t *Textile) Account() (*keypair.Full, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.Config().GetAccount()
}

// Address returns account address
func (t *Textile) Address() (string, error) {
	accnt, err := t.Account()
	if err != nil {
		return "", err
	}
	return accnt.Address(), nil
}

// ID returns account id
func (t *Textile) ID() (*peer.ID, error) {
	accnt, err := t.Account()
	if err != nil {
		return nil, err
	}
	id, err := accnt.PeerID()
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// Sign signs input with account seed
func (t *Textile) Sign(input []byte) ([]byte, error) {
	accnt, err := t.Account()
	if err != nil {
		return nil, err
	}
	return accnt.Sign(input)
}

// Verify verifies input with account address
func (t *Textile) Verify(input []byte, sig []byte) error {
	accnt, err := t.Account()
	if err != nil {
		return err
	}
	return accnt.Verify(input, sig)
}

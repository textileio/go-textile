package core

import (
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/keypair"
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

// Encrypt encrypts input with account address
func (t *Textile) Encrypt(input []byte) ([]byte, error) {
	accnt, err := t.Account()
	if err != nil {
		return nil, err
	}
	pk, err := accnt.LibP2PPubKey()
	if err != nil {
		return nil, err
	}
	return crypto.Encrypt(pk, input)
}

// Decrypt decrypts input with account address
func (t *Textile) Decrypt(input []byte) ([]byte, error) {
	accnt, err := t.Account()
	if err != nil {
		return nil, err
	}
	sk, err := accnt.LibP2PPrivKey()
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(sk, input)
}

package core

import (
	"fmt"

	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

// Account returns account keypair
func (t *Textile) Account() *keypair.Full {
	return t.account
}

// Sign signs input with account seed
func (t *Textile) Sign(input []byte) ([]byte, error) {
	return t.account.Sign(input)
}

// Verify verifies input with account address
func (t *Textile) Verify(input []byte, sig []byte) error {
	return t.account.Verify(input, sig)
}

// Encrypt encrypts input with account address
func (t *Textile) Encrypt(input []byte) ([]byte, error) {
	pk, err := t.account.LibP2PPubKey()
	if err != nil {
		return nil, err
	}
	return crypto.Encrypt(pk, input)
}

// Decrypt decrypts input with account address
func (t *Textile) Decrypt(input []byte) ([]byte, error) {
	sk, err := t.account.LibP2PPrivKey()
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(sk, input)
}

// AccountThread returns the account private thread
func (t *Textile) AccountThread() *Thread {
	return t.ThreadByKey(t.config.Account.Address)
}

// AccountContact returns a contact for this account
func (t *Textile) AccountContact() *pb.Contact {
	return t.contact(t.account.Address(), false)
}

// SyncAccount performs a thread backup search and applies the result
//func (t *Textile) SyncAccount() error {
//	//t.FindThreadBackups()
//}

// accountPeers returns all known account peers
func (t *Textile) accountPeers() []*pb.Peer {
	query := fmt.Sprintf("address='%s' and id!='%s'", t.account.Address(), t.node.Identity.Pretty())
	return t.datastore.Peers().List(query)
}

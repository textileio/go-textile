package core

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
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

// AccountPeers returns all known account peers
func (t *Textile) AccountPeers() *pb.ContactList {
	peerId := t.node.Identity.Pretty()
	address := t.account.Address()

	query := fmt.Sprintf("address='%s' and id!='%s'", address, peerId)
	list := t.datastore.Contacts().List(query)
	for i, model := range list.Items {
		list.Items[i] = t.contactView(model, false)
	}

	return list
}

// FindThreadBackups searches the network for backups
func (t *Textile) FindThreadBackups(query *pb.ThreadBackupQuery, options *pb.QueryOptions) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, nil, nil, err
	}

	options.Filter = pb.QueryOptions_NO_FILTER

	resCh, errCh, cancel := t.search(&pb.Query{
		Type:    pb.Query_THREAD_BACKUPS,
		Options: options,
		Payload: &any.Any{
			TypeUrl: "/ThreadBackupQuery",
			Value:   payload,
		},
	})

	// transform and filter results by decrypting
	backups := make(map[string]struct{})
	tresCh := make(chan *pb.QueryResult)
	terrCh := make(chan error)
	go func() {
		for {
			select {
			case res, ok := <-resCh:
				if !ok {
					close(tresCh)
					return
				}

				backup := new(pb.CafeClientThread)
				if err := ptypes.UnmarshalAny(res.Value, backup); err != nil {
					terrCh <- err
					break
				}

				plaintext, err := t.account.Decrypt(backup.Ciphertext)
				if err != nil {
					terrCh <- err
					break
				}

				thrd := new(pb.Thread)
				if err := proto.Unmarshal(plaintext, thrd); err != nil {
					terrCh <- err
					break
				}

				res.Id += ":" + thrd.Head
				if _, ok := backups[res.Id]; ok {
					continue
				}
				backups[res.Id] = struct{}{}

				res.Value = &any.Any{
					TypeUrl: "/Thread",
					Value:   plaintext,
				}
				tresCh <- res

			case err := <-errCh:
				terrCh <- err
			}
		}
	}()

	return tresCh, terrCh, cancel, nil
}

package core

import (
	"fmt"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
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
func (t *Textile) AccountPeers() ([]ContactInfo, error) {
	peers := make([]ContactInfo, 0)

	peerId := t.node.Identity.Pretty()
	address := t.account.Address()
	query := fmt.Sprintf("address='%s' and id!='%s'", address, peerId)
	for _, model := range t.datastore.Contacts().List(query) {
		info := t.contactInfo(t.datastore.Contacts().Get(model.Id), false)
		if info != nil {
			peers = append(peers, *info)
		}
	}

	return peers, nil
}

// FindThreadBackups searches the network for backups
func (t *Textile) FindThreadBackups(query *pb.ThreadBackupQuery, options *pb.QueryOptions) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, nil, nil, err
	}

	options.Filter = pb.QueryOptions_NO_FILTER

	resCh, errCh, cancel := t.search(&pb.Query{
		Type:    pb.QueryType_THREAD_BACKUPS,
		Options: options,
		Payload: &any.Any{
			TypeUrl: "/ThreadBackupQuery",
			Value:   payload,
		},
	})

	// transform results by decrypting
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

// ApplyThreadBackup dencrypts and adds a thread from a backup
func (t *Textile) ApplyThreadBackup(backup *pb.Thread) error {
	// check if we're allowed to get an invite
	// Note: just using a dummy thread here because having these access+sharing
	// methods on Thread is very nice elsewhere.
	dummy := &Thread{
		initiator: backup.Initiator,
		ttype:     repo.ThreadType(backup.Type),
		sharing:   repo.ThreadSharing(backup.Sharing),
		members:   backup.Members,
	}
	if !dummy.shareable(t.config.Account.Address, t.config.Account.Address) {
		return ErrNotShareable
	}

	sk, err := libp2pc.UnmarshalPrivateKey(backup.Sk)
	if err != nil {
		return err
	}

	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return err
	}
	if thrd := t.Thread(id.Pretty()); thrd != nil {
		// thread exists, aborting
		return nil
	}

	var sch mh.Multihash
	if backup.Schema != "" {
		sch, err = mh.FromB58String(backup.Schema)
		if err != nil {
			return err
		}

	}
	config := AddThreadConfig{
		Key:       ksuid.New().String(),
		Name:      backup.Name,
		Schema:    sch,
		Initiator: backup.Initiator,
		Type:      repo.ThreadType(backup.Type),
		Sharing:   repo.ThreadSharing(backup.Sharing),
		Members:   backup.Members,
		Join:      false,
	}
	thrd, err := t.AddThread(sk, config)
	if err != nil {
		return err
	}

	if err := thrd.followParents([]string{backup.Head}); err != nil {
		return err
	}
	hash, err := mh.FromB58String(backup.Head)
	if err != nil {
		return err
	}

	return thrd.updateHead(hash)
}

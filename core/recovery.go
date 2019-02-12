package core

import (
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

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
					TypeUrl: "/CafeThread",
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
func (t *Textile) ApplyThreadBackup(backup *pb.CafeThread) error {
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

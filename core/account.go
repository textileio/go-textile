package core

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/broadcast"
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
	return t.account.Encrypt(input)
}

// Decrypt decrypts input with account address
func (t *Textile) Decrypt(input []byte) ([]byte, error) {
	return t.account.Decrypt(input)
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
func (t *Textile) SyncAccount(options *pb.QueryOptions) (*broadcast.Broadcaster, error) {
	query := &pb.ThreadSnapshotQuery{
		Address: t.account.Address(),
	}

	resCh, errCh, cancel, err := t.SearchThreadSnapshots(query, options)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case res, ok := <-resCh:
				if !ok {
					return
				}
				err := t.applySnapshot(res)
				if err != nil {
					log.Errorf("error applying snap %s: %s", res.Id, err)
				}

			case err := <-errCh:
				log.Errorf("error during account sync: %s", err)
			}
		}
	}()

	return cancel, err
}

// maybeSyncAccount runs SyncAccount if it has not been run in the last kSyncAccountFreq
func (t *Textile) maybeSyncAccount() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.cancelSync != nil {
		t.cancelSync.Close()
		t.cancelSync = nil
	}

	daily, err := t.datastore.Config().GetLastDaily()
	if err != nil {
		log.Errorf("error get last daily: %s", err)
		return
	}

	if daily.Add(kSyncAccountFreq).Before(time.Now()) {
		var err error
		t.cancelSync, err = t.SyncAccount(&pb.QueryOptions{
			Wait: 10,
		})
		if err != nil {
			log.Errorf("error sync account: %s", err)
			return
		}

		err = t.datastore.Config().SetLastDaily()
		if err != nil {
			log.Errorf("error set last daily: %s", err)
		}
	}
}

// accountPeers returns all known account peers
func (t *Textile) accountPeers() []*pb.Peer {
	query := fmt.Sprintf("address='%s' and id!='%s'", t.account.Address(), t.node.Identity.Pretty())
	return t.datastore.Peers().List(query)
}

// isAccountPeer returns whether or not the given id is an account peer
func (t *Textile) isAccountPeer(id string) bool {
	query := fmt.Sprintf("address='%s' and id='%s'", t.account.Address(), id)
	return len(t.datastore.Peers().List(query)) > 0
}

// applySnapshot unmarshals and adds an unencrypted thread snapshot from a search result
func (t *Textile) applySnapshot(result *pb.QueryResult) error {
	snap := new(pb.Thread)
	if err := ptypes.UnmarshalAny(result.Value, snap); err != nil {
		return err
	}

	log.Debugf("applying snapshot %s", snap.Id)

	return t.AddOrUpdateThread(snap)
}

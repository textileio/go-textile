package core

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

// AddThread adds a thread with a given name and secret key
func (t *Textile) AddThread(name string, sk libp2pc.PrivKey, join bool) (*Thread, error) {
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	skb, err := sk.Bytes()
	if err != nil {
		return nil, err
	}
	threadModel := &repo.Thread{
		Id:      id.Pretty(),
		Name:    name,
		PrivKey: skb,
	}
	if err := t.datastore.Threads().Add(threadModel); err != nil {
		return nil, err
	}

	// load as active thread
	thrd, err := t.loadThread(threadModel)
	if err != nil {
		return nil, err
	}

	// we join here if we're the creator
	if join {
		if _, err := thrd.joinInitial(); err != nil {
			return nil, err
		}
	}

	// notify listeners
	t.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadAdded})

	log.Debugf("added a new thread %s with name %s", thrd.Id, name)

	return thrd, nil
}

// RemoveThread removes a thread
func (t *Textile) RemoveThread(id string) (mh.Multihash, error) {
	if !t.Online() {
		return nil, ErrOffline
	}

	// get the loaded thread
	i, thrd := t.GetThread(id)
	if thrd == nil {
		return nil, errors.New("thread not found")
	}

	// notify peers
	addr, err := thrd.Leave()
	if err != nil {
		return nil, err
	}

	// remove model from db
	if err := t.datastore.Threads().Delete(thrd.Id); err != nil {
		return nil, err
	}

	// clean up
	copy(t.threads[*i:], t.threads[*i+1:])
	t.threads[len(t.threads)-1] = nil
	t.threads = t.threads[:len(t.threads)-1]

	// notify listeners
	t.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadRemoved})

	log.Infof("removed thread %s with name %s", thrd.Id, thrd.Name)

	return addr, nil
}

// AcceptThreadInvite attemps to download an encrypted thread key from an internal invite,
// add the thread, and notify the inviter of the join
func (t *Textile) AcceptThreadInvite(inviteId string) (mh.Multihash, error) {
	if !t.Online() {
		return nil, ErrOffline
	}

	// download
	ciphertext, err := ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s", inviteId))
	if err != nil {
		return nil, err
	}

	// attempt decrypt w/ own keys
	plaintext, err := crypto.Decrypt(t.ipfs.PrivateKey, ciphertext)
	if err != nil {
		return nil, ErrInvalidThreadBlock
	}
	return t.handleThreadInvite(plaintext)
}

// AcceptExternalThreadInvite attemps to download an encrypted thread key from an external invite,
// add the thread, and notify the inviter of the join
func (t *Textile) AcceptExternalThreadInvite(inviteId string, key []byte) (mh.Multihash, error) {
	if !t.Online() {
		return nil, ErrOffline
	}

	// download
	ciphertext, err := ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s", inviteId))
	if err != nil {
		return nil, err
	}

	// attempt decrypt w/ key
	plaintext, err := crypto.DecryptAES(ciphertext, key)
	if err != nil {
		return nil, ErrInvalidThreadBlock
	}
	return t.handleThreadInvite(plaintext)
}

// Threads lists loaded threads
func (t *Textile) Threads() []*Thread {
	return t.threads
}

// GetThread get a thread by id from loaded threads
func (t *Textile) GetThread(id string) (*int, *Thread) {
	for i, thrd := range t.threads {
		if thrd.Id == id {
			return &i, thrd
		}
	}
	return nil, nil
}

// ThreadInfo gets thread info
func (t *Textile) ThreadInfo(id string) (*ThreadInfo, error) {
	_, thrd := t.GetThread(id)
	if thrd == nil {
		return nil, errors.New(fmt.Sprintf("cound not find thread: %s", id))
	}
	return thrd.Info()
}

// handleThreadInvite
func (t *Textile) handleThreadInvite(plaintext []byte) (mh.Multihash, error) {
	block := new(pb.ThreadBlock)
	if err := proto.Unmarshal(plaintext, block); err != nil {
		return nil, err
	}
	if block.Type != pb.ThreadBlock_INVITE {
		return nil, ErrInvalidThreadBlock
	}
	msg := new(pb.ThreadInvite)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// unpack thread secret
	sk, err := libp2pc.UnmarshalPrivateKey(msg.Sk)
	if err != nil {
		return nil, err
	}

	// ensure we dont already have this thread
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	if _, thrd := t.GetThread(id.Pretty()); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	// add it
	thrd, err := t.AddThread(msg.Name, sk, false)
	if err != nil {
		return nil, err
	}

	// follow parents, update head
	if err := thrd.handleInviteMessage(block); err != nil {
		return nil, err
	}

	// mark any discovered peers as welcomed
	// there's no need to send a welcome because we're about to send a join message
	if err := t.datastore.ThreadPeers().WelcomeByThread(thrd.Id); err != nil {
		return nil, err
	}

	// join the thread
	author, err := peer.IDB58Decode(block.Header.Author)
	if err != nil {
		return nil, err
	}
	hash, err := thrd.join(author)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

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
	"github.com/textileio/textile-go/thread"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

// AddThread adds a thread with a given name and secret key
func (t *Textile) AddThread(name string, secret libp2pc.PrivKey, join bool) (*thread.Thread, error) {
	skb, err := secret.Bytes()
	if err != nil {
		return nil, err
	}
	pkb, err := secret.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	pk := libp2pc.ConfigEncodeKey(pkb)
	threadModel := &repo.Thread{
		Id:      pk,
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
		if _, err := thrd.JoinInitial(); err != nil {
			return nil, err
		}
	}

	// notify listeners
	t.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadAdded})

	// add cafe update request
	t.cafeRequestQueue.Put(thrd.Id, repo.CafeAddThreadRequest)

	log.Debugf("added a new thread %s with name %s", thrd.Id, name)

	return thrd, nil
}

// RemoveThread removes a thread
func (t *Textile) RemoveThread(id string) (mh.Multihash, error) {
	if !t.IsOnline() {
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

	// add cafe update request
	t.cafeRequestQueue.Put(thrd.Id, repo.CafeRemoveThreadRequest)

	log.Infof("removed thread %s with name %s", thrd.Id, thrd.Name)

	return addr, nil
}

// AcceptThreadInvite attemps to download an encrypted thread key from an internal invite,
// add the thread, and notify the inviter of the join
func (t *Textile) AcceptThreadInvite(blockId string) (mh.Multihash, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}

	// download
	messageb, err := ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s", blockId))
	if err != nil {
		return nil, err
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(messageb, env); err != nil {
		return nil, err
	}

	// verify author sig
	authorPk, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return nil, err
	}
	if err := t.threadsService.VerifyEnvelope(env); err != nil {
		return nil, err
	}

	// unpack invite
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(env.Message.Payload, signed); err != nil {
		return nil, err
	}
	invite := new(pb.ThreadInvite)
	if err := proto.Unmarshal(signed.Block, invite); err != nil {
		return nil, err
	}

	// verify invitee
	if invite.InviteeId != t.ipfs.Identity.Pretty() {
		return nil, errors.New("invalid invitee")
	}

	// decrypt thread key with private key
	skb, err := crypto.Decrypt(t.ipfs.PrivateKey, invite.SkCipher)
	if err != nil {
		return nil, err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}
	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// ensure we dont already have this thread
	id := libp2pc.ConfigEncodeKey(pkb)
	if _, thrd := t.GetThread(id); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	// verify thread sig
	if err := crypto.Verify(sk.GetPublic(), signed.Block, signed.ThreadSig); err != nil {
		return nil, err
	}

	// add it
	thrd, err := t.AddThread(invite.SuggestedName, sk, false)
	if err != nil {
		return nil, err
	}

	// follow parents, update head
	if err := thrd.HandleInviteMessage(invite); err != nil {
		return nil, err
	}

	// join the thread
	addr, err := thrd.Join(authorPk, blockId)
	if err != nil {
		return nil, err
	}

	// invite devices
	if err := t.InviteDevices(thrd); err != nil {
		return nil, err
	}

	return addr, nil
}

// AcceptExternalThreadInvite attemps to download an encrypted thread key from an external invite,
// add the thread, and notify the inviter of the join
func (t *Textile) AcceptExternalThreadInvite(blockId string, key []byte) (mh.Multihash, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}

	// download
	envb, err := ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s", blockId))
	if err != nil {
		return nil, err
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(envb, env); err != nil {
		return nil, err
	}

	// verify author sig
	authorPk, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return nil, err
	}
	if err := t.threadsService.VerifyEnvelope(env); err != nil {
		return nil, err
	}

	// unpack invite
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(env.Message.Payload, signed); err != nil {
		return nil, err
	}
	invite := new(pb.ThreadExternalInvite)
	if err := proto.Unmarshal(signed.Block, invite); err != nil {
		return nil, err
	}

	// decrypt thread key
	skb, err := crypto.DecryptAES(invite.SkCipher, key)
	if err != nil {
		return nil, err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}
	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// ensure we dont already have this thread
	id := libp2pc.ConfigEncodeKey(pkb)
	if _, thrd := t.GetThread(id); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	// verify thread sig
	if err := crypto.Verify(sk.GetPublic(), signed.Block, signed.ThreadSig); err != nil {
		return nil, err
	}

	// add it
	thrd, err := t.AddThread(invite.SuggestedName, sk, false)
	if err != nil {
		return nil, err
	}

	// follow parents, update head
	if err := thrd.HandleExternalInviteMessage(invite); err != nil {
		return nil, err
	}

	// join the thread
	addr, err := thrd.Join(authorPk, blockId)
	if err != nil {
		return nil, err
	}

	// invite devices
	if err := t.InviteDevices(thrd); err != nil {
		return nil, err
	}

	return addr, nil
}

// Threads lists loaded threads
func (t *Textile) Threads() []*thread.Thread {
	return t.threads
}

// GetThread get a thread by id from loaded threads
func (t *Textile) GetThread(id string) (*int, *thread.Thread) {
	for i, thrd := range t.threads {
		if thrd.Id == id {
			return &i, thrd
		}
	}
	return nil, nil
}

// ThreadInfo gets thread info
func (t *Textile) ThreadInfo(id string) (*thread.Info, error) {
	_, thrd := t.GetThread(id)
	if thrd == nil {
		return nil, errors.New(fmt.Sprintf("cound not find thread: %s", id))
	}
	return thrd.Info()
}

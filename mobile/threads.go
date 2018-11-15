package mobile

import (
	"crypto/rand"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"

	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
)

// Thread is a simple meta data wrapper around a Thread
type Thread struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Peers int    `json:"peers"`
}

// Threads is a wrapper around a list of Threads
type Threads struct {
	Items []Thread `json:"items"`
}

// ExternalInvite is a wrapper around an invite id and key
type ExternalInvite struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Inviter string `json:"inviter"`
}

// Threads lists all threads
func (m *Mobile) Threads() (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range core.Node.Threads() {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

// AddThread adds a new thread with the given name
func (m *Mobile) AddThread(key string, name string, schema string) (string, error) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return "", err
	}
	shash, err := mh.FromB58String(schema)
	if err != nil {
		return "", err
	}
	pid, err := core.Node.PeerId()
	if err != nil {
		return "", err
	}
	config := core.NewThreadConfig{
		Key:       key,
		Name:      name,
		Schema:    shash,
		Initiator: pid.Pretty(),
		Type:      repo.OpenThread,
		Join:      true,
	}
	thrd, err := core.Node.AddThread(sk, config)
	if err != nil {
		return "", err
	}

	peers := thrd.Peers()
	item := Thread{
		Id:    thrd.Id,
		Name:  thrd.Name,
		Peers: len(peers),
	}
	return toJSON(item)
}

// ThreadInfo calls core ThreadInfo
func (m *Mobile) ThreadInfo(threadId string) (string, error) {
	info, err := core.Node.ThreadInfo(threadId)
	if err != nil {
		return "", err
	}
	return toJSON(info)
}

// AddThreadInvite adds a new invite to a thread
func (m *Mobile) AddThreadInvite(threadId string, inviteeId string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	thrd := core.Node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	pid, err := peer.IDB58Decode(inviteeId)
	if err != nil {
		return "", err
	}

	hash, err := thrd.AddInvite(pid)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

// AddExternalThreadInvite generates a new external invite link to a thread
func (m *Mobile) AddExternalThreadInvite(threadId string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	thrd := core.Node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	hash, key, err := thrd.AddExternalInvite()
	if err != nil {
		return "", err
	}

	username, _ := m.Username()
	invite := ExternalInvite{
		Id:      hash.B58String(),
		Key:     string(key),
		Inviter: username,
	}

	return toJSON(invite)
}

// AcceptExternalThreadInvite notifies the thread of a join
func (m *Mobile) AcceptExternalThreadInvite(id string, key string) (string, error) {
	hash, err := core.Node.AcceptExternalThreadInvite(id, []byte(key))
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// RemoveThread call core RemoveThread
func (m *Mobile) RemoveThread(id string) (string, error) {
	hash, err := core.Node.RemoveThread(id)
	if err != nil {
		return "", err
	}
	return hash.B58String(), err
}

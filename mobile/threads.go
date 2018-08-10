package mobile

import (
	"errors"
	"fmt"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet/thread"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
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
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range core.Node.Wallet.Threads() {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

// AddThread adds a new thread with the given name
func (m *Mobile) AddThread(name string, mnemonic string) (string, error) {
	var mnem *string
	if mnemonic != "" {
		mnem = &mnemonic
	}
	thrd, _, err := core.Node.Wallet.AddThreadWithMnemonic(name, mnem)
	if err != nil {
		return "", err
	}

	// build json
	peers := thrd.Peers()
	item := Thread{
		Id:    thrd.Id,
		Name:  thrd.Name,
		Peers: len(peers),
	}
	return toJSON(item)
}

// AddThreadInvite adds a new invite to a thread
func (m *Mobile) AddThreadInvite(threadId string, inviteePk string) (string, error) {
	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread: %s", threadId))
	}

	// decode pubkey
	ikb, err := libp2pc.ConfigDecodeKey(inviteePk)
	if err != nil {
		return "", err
	}
	ipk, err := libp2pc.UnmarshalPublicKey(ikb)
	if err != nil {
		return "", err
	}

	// add it
	addr, err := thrd.AddInvite(ipk)
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// AddExternalThreadInvite generates a new external invite link to a thread
func (m *Mobile) AddExternalThreadInvite(threadId string) (string, error) {
	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread: %s", threadId))
	}

	// add it
	addr, key, err := thrd.AddExternalInvite()
	if err != nil {
		return "", err
	}

	// create a structured invite
	username, _ := m.GetUsername()
	invite := ExternalInvite{
		Id:      addr.B58String(),
		Key:     string(key),
		Inviter: username,
	}

	return toJSON(invite)
}

// AcceptExternalThreadInvite notifies the thread of a join
func (m *Mobile) AcceptExternalThreadInvite(id string, key string) (string, error) {
	m.waitForOnline()
	addr, err := core.Node.Wallet.AcceptExternalThreadInvite(id, []byte(key))
	if err != nil {
		return "", err
	}
	return addr.B58String(), nil
}

// RemoveThread call core RemoveDevice
func (m *Mobile) RemoveThread(id string) (string, error) {
	addr, err := core.Node.Wallet.RemoveThread(id)
	if err != nil {
		return "", err
	}
	return addr.B58String(), err
}

// subscribe to thread and pass updates to messenger
func (m *Mobile) subscribe(thrd *thread.Thread) {
	for {
		select {
		case update, ok := <-thrd.Updates():
			if !ok {
				return
			}
			payload, err := toJSON(update)
			if err == nil {
				m.messenger.Notify(&Event{Name: "onThreadUpdate", Payload: payload})
			}
		}
	}
}

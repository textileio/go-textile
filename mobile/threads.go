package mobile

import (
	"crypto/rand"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema/textile"
)

// ExternalInvite is a wrapper around an invite id and key
type ExternalInvite struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Inviter string `json:"inviter"`
}

// MobileThreadConfig is the mobile-client specific config for creating a new thread
type MobileThreadConfig struct {
	Key      	string `json:"key"`
	Name  	 	string `json:"name"`
	Type  	 	string `json:"type"`
	Sharing  	string `json:"sharing"`
	Members   	[]string `json:"members"`
	Schema	 	string `json:"schema"`
	Media	 	bool `json:"media"`
	CameraRoll	bool `json:"cameraRoll"`
}

// Threads lists all threads
func (m *Mobile) Threads() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	infos := make([]core.ThreadInfo, 0)
	for _, thrd := range m.node.Threads() {
		info, err := thrd.Info()
		if err != nil {
			return "", err
		}
		infos = append(infos, *info)
	}

	return toJSON(infos)
}

// AddThread adds a new thread with the given name
func (m *Mobile) AddThread(config *MobileThreadConfig) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	threadType, err := repo.ThreadTypeFromString(config.Type)
	if err != nil {
		return "", err
	}

	sharingType, err := repo.ThreadSharingFromString(config.Sharing)
	if err != nil {
		return "", err
	}

	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return "", err
	}

	// tmp use the built-in schemas for all mobile threads
	// until we're ready to let the app define its own schemas.
	var sch string
	if config.Media {
		sch = textile.Media
	} else if config.CameraRoll {
		sch = textile.CameraRoll
	} else {
		sch = config.Schema
	}

	schema, err := m.addSchema(sch)
	if err != nil {
		return "", err
	}
	shash, err := mh.FromB58String(schema.Hash)
	if err != nil {
		return "", err
	}

	conf := core.AddThreadConfig{
		Key:       config.Key,
		Name:      config.Name,
		Schema:    shash,
		Initiator: m.node.Account().Address(),
		Type:      threadType,
		Sharing:   sharingType,
		Members:   config.Members,
		Join:      true,
	}

	thrd, err := m.node.AddThread(sk, conf)
	if err != nil {
		return "", err
	}

	info, err := thrd.Info()
	if err != nil {
		return "", err
	}

	return toJSON(info)
}

// ThreadInfo calls core ThreadInfo
func (m *Mobile) ThreadInfo(threadId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	info, err := m.node.ThreadInfo(threadId)
	if err != nil {
		return "", err
	}
	return toJSON(info)
}

// AddThreadInvite adds a new invite to a thread
func (m *Mobile) AddThreadInvite(threadId string, inviteeId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
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
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
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
		Key:     base58.FastBase58Encoding(key),
		Inviter: username,
	}

	return toJSON(invite)
}

// AcceptExternalThreadInvite notifies the thread of a join
func (m *Mobile) AcceptExternalThreadInvite(id string, key string) (string, error) {
	if !m.node.Online() {
		return "", core.ErrOffline
	}

	keyb, err := base58.Decode(key)
	if err != nil {
		return "", err
	}

	hash, err := m.node.AcceptExternalThreadInvite(id, keyb)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

// RemoveThread call core RemoveThread
func (m *Mobile) RemoveThread(id string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	hash, err := m.node.RemoveThread(id)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

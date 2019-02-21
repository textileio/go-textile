package mobile

import (
	"crypto/rand"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema/textile"
	"github.com/textileio/textile-go/util"
)

// ExternalInvite is a wrapper around an invite id and key
type ExternalInvite struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Inviter string `json:"inviter"`
}

// AddThreadConfig is the mobile-client specific config for creating a new thread
type AddThreadConfig struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Sharing    string `json:"sharing"`
	Members    string `json:"members"`
	Schema     string `json:"schema"`
	Media      bool   `json:"media"`
	CameraRoll bool   `json:"cameraRoll"`
}

// AddThread adds a new thread with the given name
func (m *Mobile) AddThread(config *AddThreadConfig) (string, error) {
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
		Members:   util.SplitString(config.Members, ","),
		Join:      true,
	}

	thrd, err := m.node.AddThread(sk, conf)
	if err != nil {
		return "", err
	}

	info, err := thrd.View()
	if err != nil {
		return "", err
	}

	return toJSON(info)
}

// AddOrUpdateThread calls core AddOrUpdateThread
func (m *Mobile) AddOrUpdateThread(thrd []byte) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	mthrd := new(pb.Thread)
	if err := proto.Unmarshal(thrd, mthrd); err != nil {
		return err
	}

	return m.node.AddOrUpdateThread(mthrd)
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

// Threads lists all threads
func (m *Mobile) Threads() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	infos := make([]core.ThreadInfo, 0)
	for _, thrd := range m.node.Threads() {
		info, err := thrd.View()
		if err != nil {
			return "", err
		}
		infos = append(infos, *info)
	}

	return toJSON(infos)
}

// Thread calls core Thread
func (m *Mobile) Thread(threadId string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	info, err := m.node.ThreadView(threadId)
	if err != nil {
		return "", err
	}
	return toJSON(info)
}

// AddInvite adds a new invite to a thread
func (m *Mobile) AddInvite(threadId string, inviteeId string) (string, error) {
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

// AddExternalInvite generates a new external invite link to a thread
func (m *Mobile) AddExternalInvite(threadId string) (string, error) {
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

// AcceptExternalInvite notifies the thread of a join
func (m *Mobile) AcceptExternalInvite(id string, key string) (string, error) {
	if !m.node.Online() {
		return "", core.ErrOffline
	}

	keyb, err := base58.Decode(key)
	if err != nil {
		return "", err
	}

	hash, err := m.node.AcceptExternalInvite(id, keyb)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

package core

import (
	"time"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
)

// ContactInfo display info about a contact
type ContactInfo struct {
	Id        string      `json:"id"`
	Address   string      `json:"address"`
	Username  string      `json:"username"`
	Inboxes   []repo.Cafe `json:"inboxes"`
	Added     time.Time   `json:"added"`
	ThreadIds []string    `json:"thread_ids"`
}

// AddContact adds a contact for the first time
// Note: Existing contacts will not be overwritten
func (t *Textile) AddContact(id string, address string, username string) error {
	return t.datastore.Contacts().Add(&repo.Contact{
		Id:       id,
		Address:  address,
		Username: username,
		Added:    time.Now(),
	})
}

// Contact looks up a contact by peer id
func (t *Textile) Contact(id string) *ContactInfo {
	model := t.datastore.Contacts().Get(id)
	if model == nil {
		return nil
	}

	threads := make([]string, 0)
	for _, peer := range t.datastore.ThreadPeers().ListById(id) {
		threads = append(threads, peer.ThreadId)
	}

	return &ContactInfo{
		Id:        model.Id,
		Address:   model.Address,
		Username:  toUsername(model),
		Inboxes:   model.Inboxes,
		Added:     model.Added,
		ThreadIds: threads,
	}
}

// Contacts returns all contacts this peer has interacted with
func (t *Textile) Contacts() ([]ContactInfo, error) {
	contacts := make([]ContactInfo, 0)

	for _, model := range t.datastore.Contacts().List() {

		threads := make([]string, 0)
		for _, peer := range t.datastore.ThreadPeers().ListById(model.Id) {
			threads = append(threads, peer.ThreadId)
		}

		contacts = append(contacts, ContactInfo{
			Id:        model.Id,
			Address:   model.Address,
			Username:  toUsername(&model),
			Inboxes:   model.Inboxes,
			Added:     model.Added,
			ThreadIds: threads,
		})
	}

	return contacts, nil
}

// ContactUsername returns the username for the peer id if known
func (t *Textile) ContactUsername(id string) string {
	if id == t.node.Identity.Pretty() {
		username, err := t.Username()
		if err == nil && username != nil && *username != "" {
			return *username
		}
		return ipfs.ShortenID(id)
	}
	contact := t.datastore.Contacts().Get(id)
	if contact == nil {
		return ipfs.ShortenID(id)
	}
	return toUsername(contact)
}

// ContactThreads returns all threads with the given peer
func (t *Textile) ContactThreads(id string) ([]ThreadInfo, error) {
	peers := t.datastore.ThreadPeers().ListById(id)
	if len(peers) == 0 {
		return nil, nil
	}

	var infos []ThreadInfo
	for _, peer := range peers {
		info, err := t.ThreadInfo(peer.ThreadId)
		if err != nil {
			return nil, err
		}
		infos = append(infos, *info)
	}

	return infos, nil
}

// toUsername returns a contact's username or trimmed peer id
func toUsername(contact *repo.Contact) string {
	if contact == nil || contact.Id == "" {
		return ""
	}
	if contact.Username != "" {
		return contact.Username
	}
	if len(contact.Id) >= 7 {
		return contact.Id[len(contact.Id)-7:]
	}
	return ""
}

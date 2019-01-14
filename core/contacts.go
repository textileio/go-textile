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
	Username  string      `json:"username,omitempty"`
	Avatar    string      `json:"avatar,omitempty"`
	Inboxes   []repo.Cafe `json:"inboxes,omitempty"`
	Created   time.Time   `json:"created"`
	Updated   time.Time   `json:"updated"`
	ThreadIds []string    `json:"thread_ids,omitempty"`
}

// AddContact adds a contact for the first time
// Note: Existing contacts will not be overwritten
func (t *Textile) AddContact(contact *repo.Contact) error {
	return t.datastore.Contacts().Add(contact)
}

// Contact looks up a contact by peer id
func (t *Textile) Contact(id string) *ContactInfo {
	return t.contactInfo(t.datastore.Contacts().Get(id), true)
}

// Contacts returns all contacts this peer has interacted with
func (t *Textile) Contacts() ([]ContactInfo, error) {
	contacts := make([]ContactInfo, 0)

	self := t.node.Identity.Pretty()
	for _, model := range t.datastore.Contacts().List() {
		if model.Id == self {
			continue
		}
		info := t.contactInfo(t.datastore.Contacts().Get(model.Id), true)
		if info != nil {
			contacts = append(contacts, *info)
		}
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

// contactInfo expands a contact into a more detailed view
func (t *Textile) contactInfo(model *repo.Contact, addThreads bool) *ContactInfo {
	if model == nil {
		return nil
	}

	var threads []string
	if addThreads {
		threads = make([]string, 0)
		for _, p := range t.datastore.ThreadPeers().ListById(model.Id) {
			threads = append(threads, p.ThreadId)
		}
	}

	return &ContactInfo{
		Id:        model.Id,
		Address:   model.Address,
		Username:  toUsername(model),
		Avatar:    model.Avatar,
		Inboxes:   model.Inboxes,
		Created:   model.Created,
		Updated:   model.Updated,
		ThreadIds: threads,
	}
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

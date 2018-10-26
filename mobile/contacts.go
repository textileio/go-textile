package mobile

import (
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/core"
)

// Contacts is a wrapper around a list of Contacts
type Contacts struct {
	Items []*core.Contact `json:"items"`
}

// Contact calls core Contact
func (m *Mobile) Contact(id string) (string, error) {
	contact := core.Node.Contact(id)
	if contact != nil {
		return toJSON(contact)
	}
	return "", errors.New("contact not found")
}

// Contacts calls core Contacts
func (m *Mobile) Contacts() (string, error) {
	contacts := Contacts{Items: make([]*core.Contact, 0)}
	items := core.Node.Contacts()
	if items != nil {
		contacts.Items = items
	}
	return toJSON(contacts)
}

// ContactThreads calls core ContactThreads
// - id is a contact's peer id
func (m *Mobile) ContactThreads(id string) (string, error) {
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range core.Node.ContactThreads(id) {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

// getUsername returns a contact's username by peer id if known
func getUsername(id string) string {
	contact := core.Node.Contact(id)
	if contact != nil {
		return contact.Username
	}
	return id[:8]
}

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
	contact := m.node.Contact(id)
	if contact != nil {
		return toJSON(contact)
	}
	return "", errors.New("contact not found")
}

// Contacts calls core Contacts
func (m *Mobile) Contacts() (string, error) {
	contacts := Contacts{Items: make([]*core.Contact, 0)}
	items := m.node.Contacts()
	if items != nil {
		contacts.Items = items
	}
	return toJSON(contacts)
}

// ContactUsername calls core ContactUsername
func (m *Mobile) ContactUsername(id string) string {
	return m.node.ContactUsername(id)
}

// ContactThreads calls core ContactThreads
func (m *Mobile) ContactThreads(id string) (string, error) {
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range m.node.ContactThreads(id) {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

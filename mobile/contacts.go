package mobile

import (
	"encoding/json"
	"errors"

	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// AddContact calls core AddContact
func (m *Mobile) AddContact(contact string) error {
	var model *repo.Contact
	if err := json.Unmarshal([]byte(contact), &model); err != nil {
		return err
	}

	return m.node.AddContact(model)
}

// Contact calls core Contact
func (m *Mobile) Contact(id string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	contact := m.node.Contact(id)
	if contact != nil {
		return toJSON(contact)
	}
	return "", errors.New("contact not found")
}

// Contacts calls core Contacts
func (m *Mobile) Contacts() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	contacts, err := m.node.Contacts()
	if err != nil {
		return "", err
	}
	if len(contacts) == 0 {
		contacts = make([]core.ContactInfo, 0)
	}
	return toJSON(contacts)
}

// RemoveContact calls core RemoveContact
func (m *Mobile) RemoveContact(id string) error {
	return m.node.RemoveContact(id)
}

// ContactThreads calls core ContactThreads
func (m *Mobile) ContactThreads(id string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	infos, err := m.node.ContactThreads(id)
	if err != nil {
		return "", err
	}
	if len(infos) == 0 {
		infos = make([]core.ThreadInfo, 0)
	}
	return toJSON(infos)
}

// SearchContacts calls core SearchContacts
func (m *Mobile) SearchContacts(query *pb.ContactQuery, options *pb.QueryOptions, cb Callback) (func(), error) {
	resCh, errCh, cancel, err := m.node.SearchContacts(query, options)
	if err != nil {
		return nil, err
	}

	return handleSearch(resCh, errCh, cancel, cb)
}

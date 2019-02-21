package mobile

import (
	"encoding/json"
	"errors"

	"github.com/golang/protobuf/proto"

	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// AddContact calls core AddContact
func (m *Mobile) AddContact(contact string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

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
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.RemoveContact(id)
}

// ContactThreads calls core ContactThreads
func (m *Mobile) ContactThreads(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrds, err := m.node.ContactThreads(id)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(thrds)
}

// SearchContacts calls core SearchContacts
func (m *Mobile) SearchContacts(query []byte, options []byte, cb Callback) (*CancelFn, error) {
	if !m.node.Online() {
		return nil, core.ErrOffline
	}

	mquery := new(pb.ContactQuery)
	if err := proto.Unmarshal(query, mquery); err != nil {
		return nil, err
	}
	moptions := new(pb.QueryOptions)
	if err := proto.Unmarshal(options, moptions); err != nil {
		return nil, err
	}

	resCh, errCh, cancel, err := m.node.SearchContacts(mquery, moptions)
	if err != nil {
		return nil, err
	}

	return handleSearchStream(resCh, errCh, cancel, cb)
}

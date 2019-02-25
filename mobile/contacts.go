package mobile

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
)

// AddContact calls core AddContact
func (m *Mobile) AddContact(contact []byte) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	model := new(pb.Contact)
	if err := proto.Unmarshal(contact, model); err != nil {
		return err
	}

	return m.node.AddContact(model)
}

// Contact calls core Contact
func (m *Mobile) Contact(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	contact := m.node.Contact(id)
	if contact == nil {
		return nil, errors.New("contact not found")
	}

	return proto.Marshal(contact)
}

// Contacts calls core Contacts
func (m *Mobile) Contacts() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.Contacts())
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
func (m *Mobile) SearchContacts(query []byte, options []byte) (*SearchHandle, error) {
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

	return m.handleSearchStream(resCh, errCh, cancel)
}

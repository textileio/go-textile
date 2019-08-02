package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
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

	err := m.node.AddContact(model)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// Contact calls core Contact
func (m *Mobile) Contact(address string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	contact := m.node.Contact(address)
	if contact == nil {
		return nil, nil
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
func (m *Mobile) RemoveContact(address string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.RemoveContact(address)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// ContactThreads calls core ContactThreads
func (m *Mobile) ContactThreads(address string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrds, err := m.node.ContactThreads(address)
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

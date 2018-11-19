package mobile

import (
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/core"
)

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
	contacts := make([]core.Contact, 0)
	contacts = m.node.Contacts()
	return toJSON(contacts)
}

// ContactUsername calls core ContactUsername
func (m *Mobile) ContactUsername(id string) string {
	return m.node.ContactUsername(id)
}

// ContactThreads calls core ContactThreads
func (m *Mobile) ContactThreads(id string) (string, error) {
	infos := make([]core.ThreadInfo, 0)
	var err error
	infos, err = m.node.ContactThreads(id)
	if err != nil {
		return "", err
	}
	return toJSON(infos)
}

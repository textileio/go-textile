package mobile

import (
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet"
)

// Contacts is a wrapper around a list of Contacts
type Contacts struct {
	Items []*wallet.Contact `json:"items"`
}

// Contacts calls core Contacts
func (m *Mobile) Contacts() (string, error) {
	contacts := Contacts{Items: make([]*wallet.Contact, 0)}
	items := core.Node.Wallet.Contacts()
	if items != nil {
		contacts.Items = items
	}
	return toJSON(contacts)
}

// ContactThreads calls core ContactThreads
// - id is a contact's peer id
func (m *Mobile) ContactThreads(id string) (string, error) {
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range core.Node.Wallet.ContactThreads(id) {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

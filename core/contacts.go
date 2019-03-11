package core

import (
	"fmt"

	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

// addContact adds or updates a contact
func (t *Textile) addContact(contact *pb.Contact) error {
	ex := t.datastore.Contacts().Get(contact.Id)
	if ex != nil && (contact.Updated == nil || util.ProtoTsIsNewer(ex.Updated, contact.Updated)) {
		return nil
	}

	// contact is new / newer, update
	if err := t.datastore.Contacts().AddOrUpdate(contact); err != nil {
		return err
	}

	// ensure new update is actually different before announcing to account
	if ex != nil {
		if contactsEqual(ex, contact) {
			return nil
		}
	}

	thrd := t.AccountThread()
	if thrd == nil {
		return fmt.Errorf("account thread not found")
	}

	if _, err := thrd.annouce(&pb.ThreadAnnounce{Contact: contact}); err != nil {
		return err
	}
	return nil
}

// publishContact publishes this peer's contact info to the cafe network
func (t *Textile) publishContact() error {
	self := t.datastore.Contacts().Get(t.node.Identity.Pretty())
	if self == nil {
		return nil
	}

	sessions := t.datastore.CafeSessions().List().Items
	if len(sessions) == 0 {
		return nil
	}
	for _, session := range sessions {
		pid, err := peer.IDB58Decode(session.Id)
		if err != nil {
			return err
		}
		if err := t.cafe.PublishContact(self, pid); err != nil {
			return err
		}
	}
	return nil
}

// updateContactInboxes sets this node's own contact's inboxes from the current cafe sessions
func (t *Textile) updateContactInboxes() error {
	var inboxes []*pb.Cafe
	for _, session := range t.datastore.CafeSessions().List().Items {
		inboxes = append(inboxes, session.Cafe)
	}
	return t.datastore.Contacts().UpdateInboxes(t.node.Identity.Pretty(), inboxes)
}

// contactView adds view info fields to a contact
func (t *Textile) contactView(model *pb.Contact, addThreads bool) *pb.Contact {
	if model == nil {
		return nil
	}

	if addThreads {
		model.Threads = make([]string, 0)
		for _, p := range t.datastore.ThreadPeers().ListById(model.Id) {
			model.Threads = append(model.Threads, p.Thread)
		}
	}

	return model
}

// toUserName returns a contact's name or trimmed address
func toUserName(contact *pb.Contact) string {
	if contact == nil || contact.Address == "" {
		return ""
	}
	if contact.Name != "" {
		return contact.Name
	}
	if len(contact.Address) >= 7 {
		return contact.Address[:7]
	}
	return ""
}

// contactsEqual returns whether or not the two contacts are identical
// Note: this does not consider Contact.Created or Contact.Updated
func contactsEqual(a *pb.Contact, b *pb.Contact) bool {
	if a.Id != b.Id {
		return false
	}
	if a.Address != b.Address {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Avatar != b.Avatar {
		return false
	}
	if len(a.Inboxes) != len(b.Inboxes) {
		return false
	}
	ac := make(map[string]*pb.Cafe)
	for _, c := range a.Inboxes {
		ac[c.Peer] = c
	}
	for _, j := range b.Inboxes {
		i, ok := ac[j.Peer]
		if !ok {
			return false
		}
		if !cafesEqual(i, j) {
			return false
		}
	}
	return true
}

// AddContact adds or updates a card
func (t *Textile) AddContact(card *pb.ContactCard) error {
	for _, contact := range card.Contacts {
		if err := t.addContact(contact); err != nil {
			return err
		}
	}
	return nil
}

// Contact looks up a card by address
func (t *Textile) Contact(address string) *pb.ContactCard {
	return t.contact(address, true)
}

// Contacts returns all known contacts as cards (grouped by address)
func (t *Textile) Contacts() *pb.ContactCardList {
	return t.contacts(fmt.Sprintf("address!='%s'", t.account.Address()), true)
}

// RemoveContact removes all contacts that share the given address
func (t *Textile) RemoveContact(address string) error {
	return t.datastore.Contacts().DeleteByAddress(address)
}

// SearchContacts searches the network for contacts
func (t *Textile) SearchContacts(query *pb.ContactQuery, options *pb.QueryOptions) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, nil, nil, err
	}

	options.Filter = pb.QueryOptions_HIDE_OLDER

	self := t.Profile()
	if self != nil {
		options.Exclude = append(options.Exclude, self.Id)
	}

	resCh, errCh, cancel := t.search(&pb.Query{
		Type:    pb.Query_CONTACTS,
		Options: options,
		Payload: &any.Any{
			TypeUrl: "/ContactQuery",
			Value:   payload,
		},
	})
	return resCh, errCh, cancel, nil
}

// contact looks up a card by address, optionally listing threads the address is part of
func (t *Textile) contact(address string, addThreads bool) *pb.ContactCard {
	list := t.datastore.Contacts().List(fmt.Sprintf("address='%s'", address))
	if len(list) == 0 {
		return nil
	}

	card := &pb.ContactCard{
		User: &pb.User{
			Address: address,
		},
		Contacts: make([]*pb.Contact, len(list)),
	}
	for i := 0; i < len(card.Contacts); i++ {
		if i == 0 {
			card.User.Name = list[i].Name
			card.User.Avatar = list[i].Avatar
		}
		card.Contacts[i] = t.contactView(list[i], addThreads)
	}

	return card
}

// contacts returns a list of contacts matching the given query
func (t *Textile) contacts(query string, addThreads bool) *pb.ContactCardList {
	groups := make(map[string]*pb.ContactCard)
	for _, model := range t.datastore.Contacts().List(query) {
		if groups[model.Address] == nil {
			groups[model.Address] = &pb.ContactCard{
				User: &pb.User{
					Address: model.Address,
					Name:    model.Name,
					Avatar:  model.Avatar,
				},
			}
		}
		groups[model.Address].Contacts = append(
			groups[model.Address].Contacts, t.contactView(model, addThreads))
	}

	cards := &pb.ContactCardList{
		Items: make([]*pb.ContactCard, 0),
	}
	for _, card := range groups {
		cards.Items = append(cards.Items, card)
	}

	return cards
}

// User returns a user object with the most recently updated contact for the given id
// Note: If no underlying contact is found, this will return an blank object w/ a
// generic user name for display-only purposes.
func (t *Textile) User(id string) *pb.User {
	contact := t.datastore.Contacts().GetBest(id)
	if contact == nil {
		return &pb.User{
			Name: ipfs.ShortenID(id),
		}
	}
	return &pb.User{
		Address: contact.Address,
		Name:    toUserName(contact),
		Avatar:  contact.Avatar,
	}
}

// UserThreads returns all threads with the given address
func (t *Textile) UserThreads(address string) (*pb.ThreadList, error) {
	threads := make(map[string]struct{})

	list := &pb.ThreadList{Items: make([]*pb.Thread, 0)}
	for _, contact := range t.datastore.Contacts().List(fmt.Sprintf("address='%s'", address)) {
		peers := t.datastore.ThreadPeers().ListById(contact.Id)
		for _, p := range peers {
			if _, ok := threads[p.Thread]; ok {
				continue
			}
			view, err := t.ThreadView(p.Thread)
			if err != nil {
				return nil, err
			}
			list.Items = append(list.Items, view)
			threads[p.Thread] = struct{}{}
		}
	}

	return list, nil
}

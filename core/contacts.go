package core

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/pb"
)

// AddContact adds or updates a card
func (t *Textile) AddContact(card *pb.Contact) error {
	var err error
	for _, peer := range card.Peers {
		err = t.AddPeer(peer)
		if err != nil {
			return err
		}
	}

	return nil
}

// Contact looks up a contact by address
func (t *Textile) Contact(address string) *pb.Contact {
	return t.contact(address, true)
}

// Contacts returns all known contacts, excluding self
func (t *Textile) Contacts() *pb.ContactList {
	return t.contacts(true)
}

// RemoveContact removes all contacts that share the given address
// @todo Add ignore to account thread targeted at the join
func (t *Textile) RemoveContact(address string) error {
	if address == t.account.Address() {
		return fmt.Errorf("cannot remove own contact")
	}

	return t.datastore.Peers().DeleteByAddress(address)
}

// ContactThreads returns all threads with the given address
func (t *Textile) ContactThreads(address string) (*pb.ThreadList, error) {
	threads := make(map[string]struct{})
	list := &pb.ThreadList{Items: make([]*pb.Thread, 0)}
	for _, p := range t.datastore.Peers().List(fmt.Sprintf("address='%s'", address)) {
		peers := t.datastore.ThreadPeers().ListById(p.Id)
		for _, tp := range peers {
			if _, ok := threads[tp.Thread]; ok {
				continue
			}
			view, err := t.ThreadView(tp.Thread)
			if err != nil {
				return nil, err
			}
			list.Items = append(list.Items, view)
			threads[tp.Thread] = struct{}{}
		}
	}

	return list, nil
}

// SearchContacts searches the network for peers and returns contacts
func (t *Textile) SearchContacts(query *pb.ContactQuery, options *pb.QueryOptions) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, nil, nil, err
	}

	// settings required for contacts
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

	// transform and results into contacts
	contacts := make(map[string]*pb.Contact)
	tresCh := make(chan *pb.QueryResult)
	terrCh := make(chan error)
	go func() {
		for {
			select {
			case res, ok := <-resCh:
				if !ok {
					close(tresCh)
					return
				}

				peer := new(pb.Peer)
				err := ptypes.UnmarshalAny(res.Value, peer)
				if err != nil {
					terrCh <- err
					break
				}

				if _, ok := contacts[peer.Address]; !ok {
					contacts[peer.Address] = &pb.Contact{
						Address: peer.Address,
						Name:    peer.Name,
						Avatar:  peer.Avatar,
					}
				}
				contacts[peer.Address].Peers = append(contacts[peer.Address].Peers, peer)

				value, err := proto.Marshal(ensureContactUser(contacts[peer.Address]))
				if err != nil {
					terrCh <- err
					break
				}

				res.Id = peer.Address
				res.Value = &any.Any{
					TypeUrl: "/Contact",
					Value:   value,
				}
				tresCh <- res

			case err := <-errCh:
				terrCh <- err
			}
		}
	}()

	return tresCh, terrCh, cancel, nil
}

// contact returns all peers with the given address as a contact
func (t *Textile) contact(address string, addThreads bool) *pb.Contact {
	list := t.datastore.Peers().List(fmt.Sprintf("address='%s'", address))
	if len(list) == 0 {
		return nil
	}

	contact := &pb.Contact{
		Address: address,
		Name:    list[0].Name,
		Avatar:  list[0].Avatar,
		Peers:   list,
	}
	return t.contactView(contact, addThreads)
}

// contacts returns a list of contacts for the given address
func (t *Textile) contacts(addThreads bool) *pb.ContactList {
	groups := make(map[string]*pb.Contact)
	for _, p := range t.datastore.Peers().List(fmt.Sprintf("address!='%s'", t.account.Address())) {
		if groups[p.Address] == nil {
			groups[p.Address] = &pb.Contact{
				Address: p.Address,
				Name:    p.Name,
				Avatar:  p.Avatar,
			}
		}
		groups[p.Address].Peers = append(groups[p.Address].Peers, p)
	}

	contacts := &pb.ContactList{
		Items: make([]*pb.Contact, 0),
	}
	for _, contact := range groups {
		contacts.Items = append(contacts.Items, t.contactView(contact, addThreads))
	}

	return contacts
}

// contactView adds view info fields to a contact
func (t *Textile) contactView(contact *pb.Contact, addThreads bool) *pb.Contact {
	if contact == nil {
		return nil
	}

	if addThreads {
		threads := make(map[string]struct{})
		for _, p := range contact.Peers {
			for _, tp := range t.datastore.ThreadPeers().ListById(p.Id) {
				if _, ok := threads[tp.Thread]; ok {
					continue
				}
				threads[tp.Thread] = struct{}{}
				contact.Threads = append(contact.Threads, tp.Thread)
			}
		}
	}

	return ensureContactUser(contact)
}

// ensureContactUser chooses the "best" user info (name, avatar) from peers to display
func ensureContactUser(contact *pb.Contact) *pb.Contact {
	if contact.Name != "" && contact.Avatar != "" {
		return contact
	}
	for _, p := range contact.Peers {
		if p.Name != "" && contact.Name == "" {
			contact.Name = p.Name
		}
		if p.Avatar != "" && contact.Avatar == "" {
			contact.Avatar = p.Avatar
		}
		if contact.Name != "" && contact.Avatar != "" {
			return contact
		}
	}
	return contact
}

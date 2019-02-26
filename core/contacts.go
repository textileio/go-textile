package core

import (
	"fmt"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/util"
)

// AddContact adds or updates a contact
func (t *Textile) AddContact(contact *pb.Contact) error {
	ex := t.datastore.Contacts().Get(contact.Id)
	if ex != nil && (contact.Updated == nil || util.ProtoTsIsNewer(ex.Updated, contact.Updated)) {
		return nil
	}

	return t.datastore.Contacts().AddOrUpdate(contact)
}

// Contact looks up a contact by peer id
func (t *Textile) Contact(id string) *pb.Contact {
	return t.contactView(t.datastore.Contacts().Get(id), true)
}

// Contacts returns all contacts this peer has interacted with
func (t *Textile) Contacts() *pb.ContactList {
	self := t.node.Identity.Pretty()
	list := t.datastore.Contacts().List(fmt.Sprintf("id!='%s'", self))

	for i, model := range list.Items {
		list.Items[i] = t.contactView(model, true)
	}

	return list
}

// RemoveContact removes a contact
func (t *Textile) RemoveContact(id string) error {
	return t.datastore.Contacts().Delete(id)
}

// User returns a user object by finding the most recently updated contact for the given id
func (t *Textile) User(id string) *pb.User {
	contact := t.datastore.Contacts().GetBest(id)
	if contact == nil {
		return &pb.User{
			Name: ipfs.ShortenID(id),
		}
	}
	return &pb.User{
		Address: contact.Address,
		Name:    toName(contact),
		Avatar:  contact.Avatar,
	}
}

// ContactThreads returns all threads with the given peer
func (t *Textile) ContactThreads(id string) (*pb.ThreadList, error) {
	peers := t.datastore.ThreadPeers().ListById(id)
	if len(peers) == 0 {
		return nil, nil
	}

	list := &pb.ThreadList{Items: make([]*pb.Thread, 0)}
	for _, p := range peers {
		view, err := t.ThreadView(p.Thread)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, view)
	}

	return list, nil
}

// PublishContact publishes this peer's contact info to the cafe network
func (t *Textile) PublishContact() error {
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

// UpdateContactInboxes sets this node's own contact's inboxes from the current cafe sessions
func (t *Textile) UpdateContactInboxes() error {
	var inboxes []*pb.Cafe
	for _, session := range t.datastore.CafeSessions().List().Items {
		inboxes = append(inboxes, session.Cafe)
	}
	return t.datastore.Contacts().UpdateInboxes(t.node.Identity.Pretty(), inboxes)
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

// toName returns a contact's name or trimmed address
func toName(contact *pb.Contact) string {
	if contact == nil || contact.Address == "" {
		return ""
	}
	if contact.Username != "" {
		return contact.Username
	}
	if len(contact.Address) >= 7 {
		return contact.Address[:7]
	}
	return ""
}

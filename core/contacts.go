package core

import (
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
)

// ContactInfo display info about a contact
type ContactInfo struct {
	Id        string      `json:"id"`
	Address   string      `json:"address"`
	Username  string      `json:"username,omitempty"`
	Avatar    string      `json:"avatar,omitempty"`
	Inboxes   []repo.Cafe `json:"inboxes,omitempty"`
	Created   time.Time   `json:"created"`
	Updated   time.Time   `json:"updated"`
	ThreadIds []string    `json:"thread_ids,omitempty"`
}

// AddContact adds a contact for the first time
// Note: Existing contacts will not be overwritten
func (t *Textile) AddContact(contact *repo.Contact) error {
	return t.datastore.Contacts().Add(contact)
}

// Contact looks up a contact by peer id
func (t *Textile) Contact(id string) (*ContactInfo, error) {
	self := t.node.Identity.Pretty()
	if id == self {
		prof, err := t.Profile(t.node.Identity)
		if err != nil || prof == nil {
			return nil, err
		}

		return &ContactInfo{
			Id:       self,
			Address:  prof.Address,
			Username: prof.Username,
			Avatar:   strings.Replace(prof.AvatarUri, "/ipfs/", "", 1),
			Inboxes:  prof.Inboxes,
			Created:  time.Now(),
			Updated:  time.Now(),
		}, nil
	}

	return t.contactInfo(t.datastore.Contacts().Get(id), true), nil
}

// Contacts returns all contacts this peer has interacted with
func (t *Textile) Contacts() ([]ContactInfo, error) {
	contacts := make([]ContactInfo, 0)

	self := t.node.Identity.Pretty()
	for _, model := range t.datastore.Contacts().List() {
		if model.Id == self {
			continue
		}
		info := t.contactInfo(t.datastore.Contacts().Get(model.Id), true)
		if info != nil {
			contacts = append(contacts, *info)
		}
	}

	return contacts, nil
}

// ContactUsername returns the username for the peer id if known
func (t *Textile) ContactUsername(id string) string {
	if id == t.node.Identity.Pretty() {
		username, err := t.Username()
		if err == nil && username != nil && *username != "" {
			return *username
		}
		return ipfs.ShortenID(id)
	}
	contact, err := t.contact(id)
	if contact == nil || err != nil {
		return ipfs.ShortenID(id)
	}
	return toUsername(contact)
}

// ContactThreads returns all threads with the given peer
func (t *Textile) ContactThreads(id string) ([]ThreadInfo, error) {
	peers := t.datastore.ThreadPeers().ListById(id)
	if len(peers) == 0 {
		return nil, nil
	}

	var infos []ThreadInfo
	for _, peer := range peers {
		info, err := t.ThreadInfo(peer.ThreadId)
		if err != nil {
			return nil, err
		}
		infos = append(infos, *info)
	}

	return infos, nil
}

// contactInfo expands a contact into a more detailed view
func (t *Textile) contactInfo(model *repo.Contact, addThreads bool) *ContactInfo {
	if model == nil {
		return nil
	}

	var threads []string
	if addThreads {
		threads = make([]string, 0)
		for _, p := range t.datastore.ThreadPeers().ListById(model.Id) {
			threads = append(threads, p.ThreadId)
		}
	}

	return &ContactInfo{
		Id:        model.Id,
		Address:   model.Address,
		Username:  toUsername(model),
		Avatar:    model.Avatar,
		Inboxes:   model.Inboxes,
		Created:   model.Created,
		Updated:   model.Updated,
		ThreadIds: threads,
	}
}

// toUsername returns a contact's username or trimmed peer id
func toUsername(contact *repo.Contact) string {
	if contact == nil || contact.Id == "" {
		return ""
	}
	if contact.Username != "" {
		return contact.Username
	}
	if len(contact.Id) >= 7 {
		return contact.Id[len(contact.Id)-7:]
	}
	return ""
}

// protoContactToModel is a tmp method just converting proto contact to the repo version
func protoContactToModel(pro *pb.Contact) *repo.Contact {
	if pro == nil {
		return nil
	}
	var inboxes []repo.Cafe
	for _, i := range pro.Inboxes {
		if i != nil {
			inboxes = append(inboxes, protoCafeToModel(i))
		}
	}
	created, err := ptypes.Timestamp(pro.Created)
	if err != nil {
		created = time.Now()
	}
	updated, err := ptypes.Timestamp(pro.Updated)
	if err != nil {
		updated = time.Now()
	}
	return &repo.Contact{
		Id:       pro.Id,
		Address:  pro.Address,
		Username: pro.Username,
		Avatar:   pro.Avatar,
		Inboxes:  inboxes,
		Created:  created,
		Updated:  updated,
	}
}

// repoContactToProto is a tmp method just converting repo contact to the proto version
func repoContactToProto(rep *repo.Contact) *pb.Contact {
	if rep == nil {
		return nil
	}
	var inboxes []*pb.Cafe
	for _, i := range rep.Inboxes {
		inboxes = append(inboxes, repoCafeToProto(i))
	}
	created, err := ptypes.TimestampProto(rep.Created)
	if err != nil {
		created = ptypes.TimestampNow()
	}
	updated, err := ptypes.TimestampProto(rep.Updated)
	if err != nil {
		updated = ptypes.TimestampNow()
	}
	return &pb.Contact{
		Id:       rep.Id,
		Address:  rep.Address,
		Username: rep.Username,
		Avatar:   rep.Avatar,
		Inboxes:  inboxes,
		Created:  created,
		Updated:  updated,
	}
}

// tmp
func (t *Textile) contact(id string) (*repo.Contact, error) {
	self := t.node.Identity.Pretty()
	if id == self {
		prof, err := t.Profile(t.node.Identity)
		if err != nil || prof == nil {
			return nil, err
		}

		return &repo.Contact{
			Id:       self,
			Address:  prof.Address,
			Username: prof.Username,
			Avatar:   strings.Replace(prof.AvatarUri, "/ipfs/", "", 1),
			Inboxes:  prof.Inboxes,
			Created:  time.Now(),
			Updated:  time.Now(),
		}, nil
	}

	return t.datastore.Contacts().Get(id), nil
}

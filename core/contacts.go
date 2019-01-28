package core

import (
	"fmt"
	"sync"
	"time"

	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// ContactInfo displays info about a contact
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

// ContactInfoQuery describes a contact search query
type ContactInfoQuery struct {
	Id       string
	Address  string
	Username string
	Local    bool
	Limit    int
	Wait     int
	Add      bool
}

// ContactInfoQueryResult displays info about a contact search result
type ContactInfoQueryResult struct {
	Local  []ContactInfo `json:"local,omitempty"`
	Remote []ContactInfo `json:"remote,omitempty"`
}

// AddContact adds a contact for the first time
// Note: Existing contacts will not be overwritten
func (t *Textile) AddContact(contact *repo.Contact) error {
	return t.datastore.Contacts().Add(contact)
}

// Contact looks up a contact by peer id
func (t *Textile) Contact(id string) *ContactInfo {
	return t.contactInfo(t.datastore.Contacts().Get(id), true)
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

// RemoveContact removes a contact
func (t *Textile) RemoveContact(id string) error {
	return t.datastore.Contacts().Delete(id)
}

// ContactDisplayInfo returns the username and avatar for the peer id if known
func (t *Textile) ContactDisplayInfo(id string) (string, string) {
	contact := t.datastore.Contacts().Get(id)
	if contact == nil {
		return ipfs.ShortenID(id), ""
	}
	return toUsername(contact), contact.Avatar
}

// ContactThreads returns all threads with the given peer
func (t *Textile) ContactThreads(id string) ([]ThreadInfo, error) {
	peers := t.datastore.ThreadPeers().ListById(id)
	if len(peers) == 0 {
		return nil, nil
	}

	var infos []ThreadInfo
	for _, p := range peers {
		info, err := t.ThreadInfo(p.ThreadId)
		if err != nil {
			return nil, err
		}
		infos = append(infos, *info)
	}

	return infos, nil
}

// PublishContact publishes this peer's contact info to the cafe network
func (t *Textile) PublishContact() error {
	self := t.datastore.Contacts().Get(t.node.Identity.Pretty())
	if self == nil {
		return nil
	}

	sessions := t.datastore.CafeSessions().List()
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
	var inboxes []repo.Cafe
	for _, session := range t.datastore.CafeSessions().List() {
		inboxes = append(inboxes, protoCafeToRepo(session.Cafe))
	}
	return t.datastore.Contacts().UpdateInboxes(t.node.Identity.Pretty(), inboxes)
}

// FindContact searches the network for contacts via cafe sessions
func (t *Textile) FindContact(query *ContactInfoQuery) (*ContactInfoQueryResult, error) {
	if query.Add && query.Username != "" {
		return nil, fmt.Errorf("add not allowed when searching by username")
	}

	result := &ContactInfoQueryResult{
		Local:  make([]ContactInfo, 0),
		Remote: make([]ContactInfo, 0),
	}

	// search local
	for _, c := range t.datastore.Contacts().Find(query.Id, query.Address, query.Username) {
		i := t.contactInfo(&c, true)
		if i != nil {
			result.Local = append(result.Local, *i)
		}
	}

	if query.Local || len(result.Local) >= query.Limit {
		return result, nil
	}

	// search the network via cafe if present
	var inbound []*pb.Contact
	sessions := t.datastore.CafeSessions().List()
	if len(sessions) == 0 {

		// search via pubsub directly
		res, err := t.cafe.FindContactPubSub(&pb.CafeContactQuery{
			FindId:       query.Id,
			FindAddress:  query.Address,
			FindUsername: query.Username,
			Limit:        int32(query.Limit),
			Wait:         int32(query.Wait),
		}, false)
		if err != nil {
			return nil, err
		}
		inbound = res.Contacts

	} else {
		groups := make([][]*pb.Contact, len(sessions))
		wg := sync.WaitGroup{}
		var ferr error
		for i, session := range sessions {
			cafe, err := peer.IDB58Decode(session.Id)
			if err != nil {
				return result, err
			}
			wg.Add(1)
			go func(i int, cafe peer.ID) {
				defer wg.Done()
				res, err := t.cafe.FindContact(query, cafe)
				if err != nil {
					ferr = err
					return
				}
				groups[i] = res
			}(i, cafe)
		}
		wg.Wait()
		if ferr != nil {
			return result, ferr
		}

		for _, g := range groups {
			inbound = deduplicateContactResults(append(inbound, g...))
		}
	}

	for _, i := range inbound {
		contact := protoContactToRepo(i)
		if query.Add {
			ex := t.datastore.Contacts().Get(contact.Id)
			if ex == nil || ex.Updated.UnixNano() < contact.Updated.UnixNano() {
				if err := t.datastore.Contacts().AddOrUpdate(contact); err != nil {
					return nil, err
				}
			}
		}

		info := t.contactInfo(contact, false)
		if info != nil {
			result.Remote = append(result.Remote, *info)
		}
	}

	return result, nil
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

// protoContactToRepo is a tmp method just converting proto contact to the repo version
func protoContactToRepo(pro *pb.Contact) *repo.Contact {
	if pro == nil {
		return nil
	}
	var inboxes []repo.Cafe
	for _, i := range pro.Inboxes {
		if i != nil {
			inboxes = append(inboxes, protoCafeToRepo(i))
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

package core

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
)

// ContactInfo displays info about a contact
type ContactInfo struct {
	Id        string      `json:"id"`
	Address   string      `json:"address"`
	Username  string      `json:"username,omitempty"`
	Avatar    string      `json:"avatar,omitempty"`
	Inboxes   []repo.Cafe `json:"inboxes,omitempty"`
	Added     time.Time   `json:"added"`
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
}

// ContactInfoQueryResult displays info about a contact search result
type ContactInfoQueryResult struct {
	Local  []ContactInfo `json:"local,omitempty"`
	Remote []ContactInfo `json:"remote,omitempty"`
}

// AddContact adds a contact for the first time
// Note: Existing contacts will not be overwritten
func (t *Textile) AddContact(id string, address string, username string) error {
	return t.datastore.Contacts().Add(&repo.Contact{
		Id:       id,
		Address:  address,
		Username: username,
		Added:    time.Now(),
	})
}

// Contact looks up a contact by peer id
func (t *Textile) Contact(id string) *ContactInfo {
	return t.contactInfo(t.datastore.Contacts().Get(id))
}

// Contacts returns all contacts this peer has interacted with
func (t *Textile) Contacts() ([]ContactInfo, error) {
	contacts := make([]ContactInfo, 0)

	self := t.node.Identity.Pretty()
	for _, model := range t.datastore.Contacts().List() {
		if model.Id == self {
			continue
		}
		info := t.contactInfo(t.datastore.Contacts().Get(model.Id))
		if info != nil {
			contacts = append(contacts, *info)
		}
	}

	return contacts, nil
}

// ContactUsername returns the username for the peer id if known
func (t *Textile) ContactUsername(id string) string {
	contact := t.datastore.Contacts().Get(id)
	if contact == nil {
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
		inboxes = append(inboxes, session.Cafe)
	}
	return t.datastore.Contacts().UpdateInboxes(t.node.Identity.Pretty(), inboxes)
}

// FindContact searches the network for contacts
func (t *Textile) FindContact(query *ContactInfoQuery) (*ContactInfoQueryResult, error) {
	sessions := t.datastore.CafeSessions().List()
	if len(sessions) == 0 {
		return nil, nil
	}

	result := &ContactInfoQueryResult{
		Local:  make([]ContactInfo, 0),
		Remote: make([]ContactInfo, 0),
	}

	// find local contacts
	for _, c := range t.datastore.Contacts().ListByUsername(query.Username) {
		i := t.contactInfo(&c)
		if i != nil {
			result.Local = append(result.Local, *i)
		}
	}

	// search the network
	if !query.Local {
		for _, session := range sessions {
			pid, err := peer.IDB58Decode(session.Id)
			if err != nil {
				return result, err
			}
			res, err := t.cafe.FindContact(query, pid)
			if err != nil {
				return result, err
			}
			for _, c := range res {
				i := t.contactInfo(&c)
				if i != nil {
					result.Remote = append(result.Remote, *i)
				}
			}
		}
	}

	return result, nil
}

// contactInfo expands a contact into a more detailed view
func (t *Textile) contactInfo(model *repo.Contact) *ContactInfo {
	if model == nil {
		return nil
	}

	threads := make([]string, 0)
	for _, p := range t.datastore.ThreadPeers().ListById(model.Id) {
		threads = append(threads, p.ThreadId)
	}

	return &ContactInfo{
		Id:        model.Id,
		Address:   model.Address,
		Username:  toUsername(model),
		Avatar:    model.Avatar,
		Inboxes:   model.Inboxes,
		Added:     model.Added,
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
	added, _ := ptypes.Timestamp(pro.Added)
	return &repo.Contact{
		Id:       pro.Id,
		Address:  pro.Address,
		Username: pro.Username,
		Avatar:   pro.Avatar,
		Inboxes:  inboxes,
		Added:    added,
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
	added, _ := ptypes.TimestampProto(rep.Added)
	return &pb.Contact{
		Id:       rep.Id,
		Address:  rep.Address,
		Username: rep.Username,
		Avatar:   rep.Avatar,
		Inboxes:  inboxes,
		Added:    added,
	}
}

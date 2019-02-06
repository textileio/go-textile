package core

import (
	"sync"
	"time"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/broadcast"
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

// ContactQuery describes a contact search query
type ContactQuery struct {
	Id       string
	Address  string
	Username string
	Local    bool
	Limit    int
	Wait     int
}

// ContactQueryResult displays info about a contact search result
type ContactQueryResult struct {
	Local   bool        `json:"local"`
	Contact ContactInfo `json:"contact"`
}

// contactSet holds a unique set of contact search results
type contactSet struct {
	items map[string]*pb.Contact
	mux   sync.Mutex
}

func newContactSet() *contactSet {
	return &contactSet{
		items: make(map[string]*pb.Contact, 0),
	}
}

// Add only adds a contact to the set if it's newer than last
func (s *contactSet) Add(items ...*pb.Contact) []*pb.Contact {
	s.mux.Lock()
	defer s.mux.Unlock()

	var added []*pb.Contact
	for _, contact := range items {
		last := s.items[contact.Id]
		if last == nil || protoTimeToNano(contact.Updated) > protoTimeToNano(last.Updated) {
			s.items[contact.Id] = contact
			added = append(added, contact)
		}
	}

	return added
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

// FindContacts searches the network for contacts
func (t *Textile) FindContacts(query *ContactQuery) (<-chan *ContactQueryResult, <-chan error, *broadcast.Broadcaster) {
	set := newContactSet()
	var searchChs []chan *ContactQueryResult

	// local results channel
	localCh := make(chan *ContactQueryResult)
	searchChs = append(searchChs, localCh)

	// remote results channel(s)
	var cafeChs []chan *ContactQueryResult
	clientCh := make(chan *ContactQueryResult)
	sessions := t.datastore.CafeSessions().List()
	if len(sessions) > 0 {
		for range sessions {
			cafeCh := make(chan *ContactQueryResult)
			cafeChs = append(cafeChs, cafeCh)
			searchChs = append(searchChs, cafeCh)
		}
	} else {
		searchChs = append(searchChs, clientCh)
	}

	resultCh := mergeContactQueryResults(searchChs)
	errCh := make(chan error)
	cancel := broadcast.NewBroadcaster(0)

	go func() {
		defer func() {
			for _, ch := range searchChs {
				close(ch)
			}
		}()

		// search local
		for _, c := range t.datastore.Contacts().Find(query.Id, query.Address, query.Username) {
			added := set.Add(repoContactToProto(&c))
			if len(added) == 0 {
				continue
			}
			info := t.contactInfo(&c, true)
			if info != nil {
				localCh <- &ContactQueryResult{Local: true, Contact: *info}
			}
		}

		if query.Local || len(set.items) >= query.Limit {
			return
		}

		// search the network via cafe if present
		if len(sessions) == 0 {

			// search via pubsub directly
			canceler := cancel.Listen()
			defer canceler.Close()
			if err := t.cafe.FindContactPubSub(&pb.CafeContactQuery{
				FindId:       query.Id,
				FindAddress:  query.Address,
				FindUsername: query.Username,
				Limit:        int32(query.Limit),
				Wait:         int32(query.Wait),
			}, set, func(res *pb.CafeContactQueryResult) {
				for _, c := range res.Contacts {
					contact := t.contactInfo(protoContactToRepo(c), false)
					if contact != nil {
						clientCh <- &ContactQueryResult{Contact: *contact}
					}
				}
			}, canceler.Ch, false); err != nil {
				errCh <- err
				return
			}

		} else {

			// search via cafes
			wg := sync.WaitGroup{}
			for i, session := range sessions {
				cafe, err := peer.IDB58Decode(session.Id)
				if err != nil {
					errCh <- err
					return
				}
				canceler := cancel.Listen()

				wg.Add(1)
				go func(i int, cafe peer.ID) {
					defer func() {
						canceler.Close()
						wg.Done()
					}()

					if err := t.cafe.FindContact(query, cafe, func(res *pb.CafeContactQueryResult) {
						added := set.Add(res.Contacts...)
						for _, c := range added {
							contact := t.contactInfo(protoContactToRepo(c), false)
							if contact != nil {
								cafeChs[i] <- &ContactQueryResult{Contact: *contact}
							}
						}
					}, canceler.Ch); err != nil {
						errCh <- err
						return
					}
				}(i, cafe)
			}

			wg.Wait()
		}
	}()

	return resultCh, errCh, cancel
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

// mergeContactQueryResults merges results from mulitple queries
func mergeContactQueryResults(cs []chan *ContactQueryResult) chan *ContactQueryResult {
	out := make(chan *ContactQueryResult)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c chan *ContactQueryResult) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
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

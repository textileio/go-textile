package core

import "time"

// Contact is wrapper around Peer, with thread info
type Contact struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	ThreadIds []string  `json:"thread_ids"`
	Added     time.Time `json:"added"`
}

// Contact looks up a contact by peer id
func (t *Textile) Contact(id string) *Contact {
	model := t.datastore.Contacts().Get(id)
	if model == nil {
		return nil
	}

	peers := t.datastore.ThreadPeers().ListById(id)
	if len(peers) == 0 {
		return nil
	}

	var threads []string
	for _, peer := range peers {
		threads = append(threads, peer.ThreadId)
	}

	return &Contact{
		Id:        model.Id,
		ThreadIds: threads,
		Username:  model.Username,
		Added:     model.Added,
	}
}

// Contacts returns all contacts this peer has interacted with
func (t *Textile) Contacts() []Contact {
	var contacts []Contact
	set := make(map[string]Contact)

	for _, peer := range t.datastore.ThreadPeers().List() {
		c, ok := set[peer.Id]
		if ok {
			c.ThreadIds = append(set[peer.Id].ThreadIds, peer.ThreadId)
		} else {
			if peer.Id == "" {
				continue
			}
			contact := Contact{
				Id:        peer.Id,
				ThreadIds: []string{peer.ThreadId},
			}

			model := t.datastore.Contacts().Get(peer.Id)
			if model != nil && model.Username != "" {
				contact.Username = model.Username
				contact.Added = model.Added
			} else {
				if len(peer.Id) >= 7 {
					contact.Username = peer.Id[len(peer.Id)-7:]
				}
			}

			set[peer.Id] = contact
			contacts = append(contacts, set[peer.Id])
		}
	}

	return contacts
}

// ContactUsername returns the username for the peer id if known
func (t *Textile) ContactUsername(id string) string {
	var username string
	if id == "" {
		return ""
	}
	if len(id) >= 7 {
		username = id[len(id)-7:]
	}

	contact := t.datastore.Contacts().Get(id)
	if contact != nil && contact.Username != "" {
		username = contact.Username
	}

	return username
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

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
func (t *Textile) Contacts() []*Contact {
	var contacts []*Contact
	set := make(map[string]*Contact)

	for _, peer := range t.datastore.ThreadPeers().List() {
		c, ok := set[peer.Id]
		if ok {
			c.ThreadIds = append(set[peer.Id].ThreadIds, peer.ThreadId)
		} else {
			username := peer.Id[:8]
			contact := t.datastore.Contacts().Get(peer.Id)
			if contact != nil {
				username = contact.Username
			}
			set[peer.Id] = &Contact{
				Id:        contact.Id,
				ThreadIds: []string{peer.ThreadId},
				Username:  username,
				Added:     contact.Added,
			}
			contacts = append(contacts, set[peer.Id])
		}
	}

	return contacts
}

// ContactUsername returns the username for the peer id if known
func (t *Textile) ContactUsername(id string) string {
	username := id[len(id)-7:]

	contact := t.datastore.Contacts().Get(id)
	if contact != nil && contact.Username != "" {
		username = contact.Username
	}

	return username
}

// ContactThreads returns all threads with the given peer
func (t *Textile) ContactThreads(id string) []*Thread {
	peers := t.datastore.ThreadPeers().ListById(id)
	if len(peers) == 0 {
		return nil
	}

	var threads []*Thread
	for _, peer := range peers {
		if thrd := t.Thread(peer.ThreadId); thrd != nil {
			threads = append(threads, thrd)
		}
	}

	return threads
}

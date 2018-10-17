package core

import (
	"github.com/textileio/textile-go/thread"
)

// Contact is wrapper around Peer, with thread info
type Contact struct {
	Id        string   `json:"id"`
	ThreadIds []string `json:"thread_ids"`
}

// GetContacts returns all contacts this peer has interacted with
func (t *Textile) Contacts() []*Contact {
	var contacts []*Contact
	set := make(map[string]*Contact)
	for _, peer := range t.datastore.Peers().List(-1, "") {
		c, ok := set[peer.Id]
		if ok {
			c.ThreadIds = append(set[peer.Id].ThreadIds, peer.ThreadId)
		} else {
			set[peer.Id] = &Contact{
				Id:        peer.Id,
				ThreadIds: []string{peer.ThreadId},
			}
			contacts = append(contacts, set[peer.Id])
		}
	}
	return contacts
}

// ContactThreads returns all threads with the given peer
func (t *Textile) ContactThreads(id string) []*thread.Thread {
	peers := t.datastore.Peers().List(-1, "id='"+id+"'")
	if len(peers) == 0 {
		return nil
	}
	var threads []*thread.Thread
	for _, peer := range peers {
		if _, thrd := t.GetThread(peer.ThreadId); thrd != nil {
			threads = append(threads, thrd)
		}
	}
	return threads
}

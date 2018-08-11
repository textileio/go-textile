package wallet

import (
	"github.com/textileio/textile-go/wallet/thread"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// Contact is wrapper around Peer, with thread info
type Contact struct {
	Id        string   `json:"id"`
	Pk        string   `json:"pk"`
	ThreadIds []string `json:"thread_ids"`
}

// GetContacts returns all contacts this peer has interacted with
func (w *Wallet) Contacts() []*Contact {
	var contacts []*Contact
	set := make(map[string]*Contact)
	for _, peer := range w.datastore.Peers().List("", -1, "") {
		c, ok := set[peer.Id]
		if ok {
			c.ThreadIds = append(set[peer.Id].ThreadIds, peer.ThreadId)
		} else {
			set[peer.Id] = &Contact{
				Id:        peer.Id,
				Pk:        libp2pc.ConfigEncodeKey(peer.PubKey),
				ThreadIds: []string{peer.ThreadId},
			}
			contacts = append(contacts, set[peer.Id])
		}
	}
	return contacts
}

// ContactThreads returns all threads with the given peer
func (w *Wallet) ContactThreads(id string) []*thread.Thread {
	peers := w.datastore.Peers().List("", -1, "id='"+id+"'")
	if len(peers) == 0 {
		return nil
	}
	var threads []*thread.Thread
	for _, peer := range peers {
		if _, thrd := w.GetThread(peer.ThreadId); thrd != nil {
			threads = append(threads, thrd)
		}
	}
	return threads
}

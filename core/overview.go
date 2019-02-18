package core

import (
	"fmt"

	"github.com/textileio/textile-go/pb"
)

// Summary is a wallet summary object
type Summary struct {
	AccountPeerCount int `json:"account_peer_cnt"`
	ThreadCount      int `json:"thread_cnt"`
	FileCount        int `json:"file_cnt"`
	ContactCount     int `json:"contact_cnt"`
}

// Summary returns a summary of node data
func (t *Textile) Summary() (*Summary, error) {
	selfId := t.node.Identity.Pretty()
	selfAddress := t.account.Address()

	peers := t.datastore.Contacts().Count(fmt.Sprintf("address!='%s'", selfAddress))
	threads := t.datastore.Threads().Count()
	files := t.datastore.Blocks().Count(fmt.Sprintf("type=%d", pb.Block_FILES))
	contacts := t.datastore.Contacts().Count(fmt.Sprintf("id!='%s'", selfId))

	return &Summary{
		AccountPeerCount: peers,
		ThreadCount:      threads,
		FileCount:        files,
		ContactCount:     contacts,
	}, nil
}

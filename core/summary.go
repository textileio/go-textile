package core

import (
	"fmt"

	"github.com/textileio/textile-go/pb"
)

// Summary returns a summary of node data
func (t *Textile) Summary() *pb.Summary {
	selfId := t.node.Identity.Pretty()
	selfAddress := t.account.Address()

	peers := t.datastore.Contacts().Count(fmt.Sprintf("address!='%s'", selfAddress))
	threads := t.datastore.Threads().Count()
	files := t.datastore.Blocks().Count(fmt.Sprintf("type=%d", pb.Block_FILES))
	contacts := t.datastore.Contacts().Count(fmt.Sprintf("id!='%s'", selfId))

	return &pb.Summary{
		AccountPeerCount: int32(peers),
		ThreadCount:      int32(threads),
		FileCount:        int32(files),
		ContactCount:     int32(contacts),
	}
}

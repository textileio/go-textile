package core

import (
	"fmt"

	"github.com/textileio/go-textile/pb"
)

// Summary returns a summary of node data
func (t *Textile) Summary() *pb.Summary {
	peers := t.datastore.Peers().Count(fmt.Sprintf("address='%s'", t.account.Address()))
	threads := t.datastore.Threads().Count()
	files := t.datastore.Blocks().Count(fmt.Sprintf("type=%d", pb.Block_FILES))
	contacts := len(t.Contacts().Items)

	return &pb.Summary{
		Id:               t.node.Identity.Pretty(),
		Address:          t.account.Address(),
		AccountPeerCount: int32(peers) - 1,
		ThreadCount:      int32(threads),
		FilesCount:       int32(files),
		ContactCount:     int32(contacts),
	}
}

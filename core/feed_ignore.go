package core

import (
	"strings"

	"github.com/textileio/go-textile/pb"
)

func (t *Textile) ignore(block *pb.Block, opts feedItemOpts) (*pb.Ignore, error) {
	if block.Type != pb.Block_IGNORE {
		return nil, ErrBlockWrongType
	}

	targetId := strings.TrimPrefix(block.Target, "ignore-")
	target, err := t.feedItem(t.datastore.Blocks().Get(targetId), feedItemOpts{})
	if err != nil {
		return nil, err
	}

	return &pb.Ignore{
		Block:  block.Id,
		Date:   block.Date,
		User:   t.PeerUser(block.Author),
		Target: target,
	}, nil
}

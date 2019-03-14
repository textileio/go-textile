package core

import (
	"strings"

	"github.com/textileio/go-textile/pb"
)

func (t *Textile) flag(block *pb.Block, opts feedItemOpts) (*pb.Flag, error) {
	if block.Type != pb.Block_FLAG {
		return nil, ErrBlockWrongType
	}

	targetId := strings.TrimPrefix(block.Target, "flag-")
	target, err := t.feedItem(t.datastore.Blocks().Get(targetId), feedItemOpts{})
	if err != nil {
		return nil, err
	}

	return &pb.Flag{
		Block:  block.Id,
		Date:   block.Date,
		User:   t.PeerUser(block.Author),
		Target: target,
	}, nil
}

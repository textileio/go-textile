package core

import (
	"github.com/textileio/textile-go/pb"
)

func (t *Textile) flag(block *pb.Block, opts feedItemOpts) (*pb.Flag, error) {
	if block.Type != pb.Block_FLAG {
		return nil, ErrBlockWrongType
	}

	target, err := t.feedItem(t.datastore.Blocks().Get(block.Target), feedItemOpts{})
	if err != nil {
		return nil, err
	}

	return &pb.Flag{
		Block:  block.Id,
		Date:   block.Date,
		User:   t.User(block.Author),
		Target: target,
	}, nil
}

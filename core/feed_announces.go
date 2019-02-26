package core

import (
	"github.com/textileio/textile-go/pb"
)

func (t *Textile) announce(block *pb.Block, opts feedItemOpts) (*pb.Announce, error) {
	if block.Type != pb.Block_ANNOUNCE {
		return nil, ErrBlockWrongType
	}

	return &pb.Announce{
		Block: block.Id,
		Date:  block.Date,
		User:  t.User(block.Author),
	}, nil
}

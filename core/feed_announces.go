package core

import (
	"github.com/textileio/go-textile/pb"
)

func (t *Textile) announce(block *pb.Block, opts feedItemOpts) (*pb.Announce, error) {
	if block.Type != pb.Block_ANNOUNCE {
		return nil, ErrBlockWrongType
	}

	return &pb.Announce{
		Block: block.Id,
		Date:  block.Date,
		User:  t.PeerUser(block.Author),
	}, nil
}

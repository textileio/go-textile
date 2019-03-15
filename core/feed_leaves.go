package core

import (
	"github.com/textileio/go-textile/pb"
)

func (t *Textile) leave(block *pb.Block, opts feedItemOpts) (*pb.Leave, error) {
	if block.Type != pb.Block_LEAVE {
		return nil, ErrBlockWrongType
	}

	item := &pb.Leave{
		Block: block.Id,
		Date:  block.Date,
		User:  t.PeerUser(block.Author),
	}

	if opts.annotations {
		likes, err := t.Likes(block.Id)
		if err != nil {
			return nil, err
		}
		item.Likes = likes.Items
	} else {
		item.Likes = opts.likes
	}

	return item, nil
}

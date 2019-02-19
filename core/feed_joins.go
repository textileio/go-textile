package core

import (
	"github.com/textileio/textile-go/pb"
)

func (t *Textile) join(block *pb.Block, opts feedItemOpts) (*pb.Join, error) {
	if block.Type != pb.Block_JOIN {
		return nil, ErrBlockWrongType
	}

	item := &pb.Join{
		Block: block.Id,
		Date:  block.Date,
		User:  t.User(block.Author),
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

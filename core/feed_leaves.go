package core

import (
	"github.com/textileio/textile-go/pb"
)

func (t *Textile) leave(block *pb.Block, opts feedItemOpts) (*pb.Leave, error) {
	if block.Type != pb.Block_LEAVE {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.Author)
	item := &pb.Leave{
		Block:    block.Id,
		Date:     block.Date,
		Author:   block.Author,
		Username: username,
		Avatar:   avatar,
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

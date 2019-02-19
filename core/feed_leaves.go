package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) leave(block *repo.Block, opts feedItemOpts) (*pb.Leave, error) {
	if block.Type != repo.LeaveBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)
	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.Leave{
		Block:    block.Id,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
	}

	if opts.annotations {
		likes, err := t.Likes(block.Id)
		if err != nil {
			return nil, err
		}
		info.Likes = likes.Items
	} else {
		info.Likes = opts.likes
	}

	return info, nil
}

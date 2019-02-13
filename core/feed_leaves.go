package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) FeedLeave(block *repo.Block, annotated bool) (*pb.FeedLeave, error) {
	if block.Type != repo.LeaveBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.FeedLeave{
		Block:    block.Id,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
	}

	if annotated {
		likes, err := t.Likes(block.Id)
		if err != nil {
			return nil, err
		}
		info.Likes = likes.Items
	}

	return info, nil
}

package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) Likes(target string) (*pb.FeedLikeList, error) {
	likes := make([]*pb.FeedLike, 0)

	query := fmt.Sprintf("type=%d and target='%s'", repo.LikeBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info, err := t.FeedLike(&block, true)
		if err != nil {
			continue
		}
		likes = append(likes, info)
	}

	return &pb.FeedLikeList{Items: likes}, nil
}

func (t *Textile) FeedLike(block *repo.Block, annotation bool) (*pb.FeedLike, error) {
	if block.Type != repo.LikeBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.FeedLike{
		Id:       block.Id,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
	}

	if !annotation {
		target, err := t.feedItem(t.datastore.Blocks().Get(block.Target), false)
		if err != nil {
			return nil, err
		}
		info.Target = target
	}

	return info, nil
}

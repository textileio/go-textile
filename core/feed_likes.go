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
		info, err := t.feedLike(&block, feedItemOpts{annotations: true})
		if err != nil {
			continue
		}
		likes = append(likes, info)
	}

	return &pb.FeedLikeList{Items: likes}, nil
}

func (t *Textile) FeedLike(blockId string) (*pb.FeedLike, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.feedLike(block, feedItemOpts{annotations: true})
}

func (t *Textile) feedLike(block *repo.Block, opts feedItemOpts) (*pb.FeedLike, error) {
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

	if !opts.annotations {
		target, err := t.feedItem(t.datastore.Blocks().Get(block.Target), feedItemOpts{})
		if err != nil {
			return nil, err
		}
		info.Target = target
	} else {
		info.Target = opts.target
	}

	return info, nil
}

package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) Likes(target string) (*pb.LikeList, error) {
	likes := make([]*pb.Like, 0)

	query := fmt.Sprintf("type=%d and target='%s'", repo.LikeBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info, err := t.like(&block, feedItemOpts{annotations: true})
		if err != nil {
			continue
		}
		likes = append(likes, info)
	}

	return &pb.LikeList{Items: likes}, nil
}

func (t *Textile) Like(blockId string) (*pb.Like, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.like(block, feedItemOpts{annotations: true})
}

func (t *Textile) like(block *repo.Block, opts feedItemOpts) (*pb.Like, error) {
	if block.Type != repo.LikeBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)
	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.Like{
		Id:       block.Id,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
	}

	if opts.target != nil {
		info.Target = opts.target
	} else if !opts.annotations {
		target, err := t.feedItem(t.datastore.Blocks().Get(block.Target), feedItemOpts{})
		if err != nil {
			return nil, err
		}
		info.Target = target
	}

	return info, nil
}

package core

import (
	"fmt"

	"github.com/textileio/go-textile/pb"
)

func (t *Textile) Likes(target string) (*pb.LikeList, error) {
	likes := make([]*pb.Like, 0)

	query := fmt.Sprintf("type=%d and target='%s'", pb.Block_LIKE, target)
	for _, block := range t.Blocks("", -1, query).Items {
		info, err := t.like(block, feedItemOpts{annotations: true})
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

func (t *Textile) like(block *pb.Block, opts feedItemOpts) (*pb.Like, error) {
	if block.Type != pb.Block_LIKE {
		return nil, ErrBlockWrongType
	}

	item := &pb.Like{
		Id:   block.Id,
		Date: block.Date,
		User: t.PeerUser(block.Author),
	}

	if opts.target != nil {
		item.Target = opts.target
	} else if !opts.annotations {
		target, err := t.feedItem(t.datastore.Blocks().Get(block.Target), feedItemOpts{})
		if err != nil {
			return nil, err
		}
		item.Target = target
	}

	return item, nil
}

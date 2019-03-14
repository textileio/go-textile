package core

import (
	"fmt"

	"github.com/textileio/go-textile/pb"
)

func (t *Textile) Comments(target string) (*pb.CommentList, error) {
	comments := make([]*pb.Comment, 0)

	query := fmt.Sprintf("type=%d and target='%s'", pb.Block_COMMENT, target)
	for _, block := range t.Blocks("", -1, query).Items {
		info, err := t.comment(block, feedItemOpts{annotations: true})
		if err != nil {
			continue
		}
		comments = append(comments, info)
	}

	return &pb.CommentList{Items: comments}, nil
}

func (t *Textile) Comment(blockId string) (*pb.Comment, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.comment(block, feedItemOpts{annotations: true})
}

func (t *Textile) comment(block *pb.Block, opts feedItemOpts) (*pb.Comment, error) {
	if block.Type != pb.Block_COMMENT {
		return nil, ErrBlockWrongType
	}

	item := &pb.Comment{
		Id:   block.Id,
		Date: block.Date,
		User: t.PeerUser(block.Author),
		Body: block.Body,
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

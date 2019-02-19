package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) Messages(offset string, limit int, threadId string) (*pb.TextList, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("threadId='%s' and type=%d", threadId, repo.MessageBlock)
	} else {
		query = fmt.Sprintf("type=%d", repo.MessageBlock)
	}

	list := make([]*pb.Text, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		msg, err := t.message(&block, feedItemOpts{annotations: true})
		if err != nil {
			return nil, err
		}
		list = append(list, msg)
	}

	return &pb.TextList{Items: list}, nil
}

func (t *Textile) Message(blockId string) (*pb.Text, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.message(block, feedItemOpts{annotations: true})
}

func (t *Textile) message(block *repo.Block, opts feedItemOpts) (*pb.Text, error) {
	if block.Type != repo.MessageBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)
	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.Text{
		Block:    block.Id,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Body:     block.Body,
	}

	if opts.annotations {
		comments, err := t.Comments(block.Id)
		if err != nil {
			return nil, err
		}
		info.Comments = comments.Items

		likes, err := t.Likes(block.Id)
		if err != nil {
			return nil, err
		}
		info.Likes = likes.Items
	} else {
		info.Comments = opts.comments
		info.Likes = opts.likes
	}

	return info, nil
}

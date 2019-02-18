package core

import (
	"fmt"

	"github.com/textileio/textile-go/pb"
)

func (t *Textile) Messages(offset string, limit int, threadId string) (*pb.TextList, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("threadId='%s' and type=%d", threadId, pb.Block_MESSAGE)
	} else {
		query = fmt.Sprintf("type=%d", pb.Block_MESSAGE)
	}

	list := make([]*pb.Text, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks.Items {
		msg, err := t.message(block, feedItemOpts{annotations: true})
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

func (t *Textile) message(block *pb.Block, opts feedItemOpts) (*pb.Text, error) {
	if block.Type != pb.Block_MESSAGE {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.Author)
	item := &pb.Text{
		Block:    block.Id,
		Date:     block.Date,
		Author:   block.Author,
		Username: username,
		Avatar:   avatar,
		Body:     block.Body,
	}

	if opts.annotations {
		comments, err := t.Comments(block.Id)
		if err != nil {
			return nil, err
		}
		item.Comments = comments.Items

		likes, err := t.Likes(block.Id)
		if err != nil {
			return nil, err
		}
		item.Likes = likes.Items
	} else {
		item.Comments = opts.comments
		item.Likes = opts.likes
	}

	return item, nil
}

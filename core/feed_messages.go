package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) Messages(offset string, limit int, threadId string) (*pb.FeedMessageList, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("threadId='%s' and type=%d", threadId, repo.MessageBlock)
	} else {
		query = fmt.Sprintf("type=%d", repo.MessageBlock)
	}

	list := make([]*pb.FeedMessage, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		msg, err := t.feedMessage(&block, true)
		if err != nil {
			return nil, err
		}
		list = append(list, msg)
	}

	return &pb.FeedMessageList{Items: list}, nil
}

func (t *Textile) FeedMessage(blockId string) (*pb.FeedMessage, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.feedMessage(block, true)
}

func (t *Textile) feedMessage(block *repo.Block, annotated bool) (*pb.FeedMessage, error) {
	if block.Type != repo.MessageBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.FeedMessage{
		Block:    block.Id,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Body:     block.Body,
	}

	if annotated {
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
	}

	return info, nil
}

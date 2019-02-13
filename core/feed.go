package core

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

var feedTypes = []repo.BlockType{
	repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock, repo.CommentBlock, repo.LikeBlock,
}

var annotatedFeedTypes = []repo.BlockType{
	repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock,
}

func (t *Textile) Feed(offset string, limit int, threadId string, annotated bool) ([]*pb.FeedItem, error) {
	var types []repo.BlockType
	if annotated {
		types = annotatedFeedTypes
	} else {
		types = feedTypes
	}

	var query string
	for i, t := range types {
		query += fmt.Sprintf("type=%d", t)
		if i != len(types)-1 {
			query += " or "
		}
	}
	query = "(" + query + ")"
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("(threadId='%s') and %s", threadId, query)
	}

	list := make([]*pb.FeedItem, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		item, err := t.feedItem(&block, annotated)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (t *Textile) feedItem(block *repo.Block, annotated bool) (*pb.FeedItem, error) {
	item := &pb.FeedItem{
		Block: block.Id,
		Body:  &any.Any{},
	}

	var body proto.Message
	var err error
	switch block.Type {
	case repo.JoinBlock:
		item.Body.TypeUrl = "/FeedJoin"
		body, err = t.FeedJoin(block, annotated)
	case repo.LeaveBlock:
		item.Body.TypeUrl = "/FeedLeave"
		body, err = t.FeedLeave(block, annotated)
	case repo.FilesBlock:
		item.Body.TypeUrl = "/FeedFiles"
		body, err = t.feedFile(block, annotated)
	case repo.MessageBlock:
		item.Body.TypeUrl = "/FeedMessage"
		body, err = t.feedMessage(block, annotated)
	case repo.CommentBlock:
		item.Body.TypeUrl = "/FeedComment"
		body, err = t.FeedComment(block, annotated)
	case repo.LikeBlock:
		item.Body.TypeUrl = "/FeedLike"
		body, err = t.FeedLike(block, annotated)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	value, err := proto.Marshal(body)
	if err != nil {
		return nil, err
	}
	item.Body.Value = value

	return item, err
}

package core

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) Comments(target string) (*pb.FeedCommentList, error) {
	comments := make([]*pb.FeedComment, 0)

	query := fmt.Sprintf("type=%d and target='%s'", repo.CommentBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info, err := t.FeedComment(&block, true)
		if err != nil {
			continue
		}
		comments = append(comments, info)
	}

	return &pb.FeedCommentList{Items: comments}, nil
}

func (t *Textile) FeedComment(block *repo.Block, annotation bool) (*pb.FeedComment, error) {
	if block.Type != repo.CommentBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.FeedComment{
		Id:       block.Id,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Body:     block.Body,
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

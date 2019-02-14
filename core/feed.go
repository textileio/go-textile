package core

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

var flatFeedTypes = []repo.BlockType{
	repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock, repo.CommentBlock, repo.LikeBlock,
}

var annotatedFeedTypes = []repo.BlockType{
	repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock,
}

type hybridStack struct {
	top      repo.Block
	children []repo.Block
}

type feedItemOpts struct {
	annotations bool
	comments    []*pb.FeedComment
	likes       []*pb.FeedLike
	target      *pb.FeedItem
}

func (t *Textile) Feed(offset string, limit int, threadId string, feedType pb.FeedType) (*pb.FeedItemList, error) {
	var types []repo.BlockType
	switch feedType {
	case pb.FeedType_FLAT, pb.FeedType_HYBRID:
		types = flatFeedTypes
	case pb.FeedType_ANNOTATED:
		types = annotatedFeedTypes
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

	blocks := t.Blocks(offset, limit, query)
	list := make([]*pb.FeedItem, 0)

	switch feedType {
	case pb.FeedType_FLAT, pb.FeedType_ANNOTATED:
		for _, block := range blocks {
			item, err := t.feedItem(&block, feedItemOpts{
				annotations: feedType == pb.FeedType_ANNOTATED,
			})
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}

	case pb.FeedType_HYBRID:
		stacks := make([]hybridStack, 0)
		for _, block := range blocks {
			targetId := getTargetId(block)
			if len(stacks) == 0 || targetId != getTargetId(stacks[len(stacks)-1].top) {
				// start a new stack
				stacks = append(stacks, hybridStack{top: block})
			} else {
				// append to last
				stacks[len(stacks)-1].children = append(stacks[len(stacks)-1].children, block)
			}
		}

		for _, stack := range stacks {
			item, err := t.feedStackItem(stack)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
	}

	return &pb.FeedItemList{Items: list}, nil
}

func (t *Textile) feedItem(block *repo.Block, opts feedItemOpts) (*pb.FeedItem, error) {
	item := &pb.FeedItem{
		Block:   block.Id,
		Payload: &any.Any{},
	}

	var payload proto.Message
	var err error
	switch block.Type {
	case repo.JoinBlock:
		item.Payload.TypeUrl = "/FeedJoin"
		payload, err = t.feedJoin(block, opts)
	case repo.LeaveBlock:
		item.Payload.TypeUrl = "/FeedLeave"
		payload, err = t.feedLeave(block, opts)
	case repo.FilesBlock:
		item.Payload.TypeUrl = "/FeedFiles"
		payload, err = t.feedFile(block, opts)
	case repo.MessageBlock:
		item.Payload.TypeUrl = "/FeedMessage"
		payload, err = t.feedMessage(block, opts)
	case repo.CommentBlock:
		item.Payload.TypeUrl = "/FeedComment"
		payload, err = t.feedComment(block, opts)
	case repo.LikeBlock:
		item.Payload.TypeUrl = "/FeedLike"
		payload, err = t.feedLike(block, opts)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	value, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}
	item.Payload.Value = value

	return item, nil
}

func (t *Textile) feedStackItem(stack hybridStack) (*pb.FeedItem, error) {
	var comments []*pb.FeedComment
	var likes []*pb.FeedLike

	// does the stack contain the initial target,
	// or is it a continuation stack of annotations?
	// we'll need to load the target in the latter case.
	var target *repo.Block
	var targetId string
	for _, child := range stack.children {
		switch child.Type {
		case repo.CommentBlock:
			targetId = child.Target
			comment, err := t.feedComment(&child, feedItemOpts{annotations: true})
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)

		case repo.LikeBlock:
			targetId = child.Target
			like, err := t.feedLike(&child, feedItemOpts{annotations: true})
			if err != nil {
				return nil, err
			}
			likes = append(likes, like)

		default:
			target = &child
		}
	}

	var initial bool
	if target != nil { // target was in children
		initial = true
	} else if !isAnnotation(stack.top) { // top is target
		initial = true
		target = &stack.top
	} else { // target is not in the stack, load it
		target = t.datastore.Blocks().Get(targetId)
		if target == nil {
			return nil, nil
		}
	}

	targetItem, err := t.feedItem(target, feedItemOpts{
		comments: comments,
		likes:    likes,
	})
	if err != nil {
		return nil, err
	}

	if initial {
		return targetItem, nil
	} else {
		// target gets wrapped with the top block
		return t.feedItem(&stack.top, feedItemOpts{
			target: targetItem,
		})
	}
}

func getTargetId(block repo.Block) string {
	switch block.Type {
	case repo.CommentBlock, repo.LikeBlock:
		return block.Target
	default:
		return block.Id
	}
}

func isAnnotation(block repo.Block) bool {
	switch block.Type {
	case repo.CommentBlock, repo.LikeBlock:
		return true
	default:
		return false
	}
}

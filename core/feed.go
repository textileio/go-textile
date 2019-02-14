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
	id       string
	top      repo.Block
	children []repo.Block
}

type feedItemOpts struct {
	annotations bool
	comments    []*pb.Comment
	likes       []*pb.Like
	target      *pb.FeedItem
}

func (t *Textile) Feed(offset string, limit int, threadId string, mode pb.FeedMode) (*pb.FeedItemList, error) {
	var types []repo.BlockType
	switch mode {
	case pb.FeedMode_FLAT, pb.FeedMode_HYBRID:
		types = flatFeedTypes
	case pb.FeedMode_ANNOTATED:
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
	var count int

	switch mode {
	case pb.FeedMode_FLAT, pb.FeedMode_ANNOTATED:
		for _, block := range blocks {
			item, err := t.feedItem(&block, feedItemOpts{
				annotations: mode == pb.FeedMode_ANNOTATED,
			})
			if err != nil {
				return nil, err
			}
			list = append(list, item)
			count++
		}

	case pb.FeedMode_HYBRID:
		stacks := make([]hybridStack, 0)
		var last *hybridStack
		for _, block := range blocks {
			if len(stacks) > 0 {
				last = &stacks[len(stacks)-1]
			}
			targetId := getTargetId(block)

			if len(stacks) == 0 || targetId != getTargetId(last.top) {
				// start a new stack
				stacks = append(stacks, hybridStack{id: targetId, top: block})
			} else {
				// append to last
				last.children = append(last.children, block)
			}
		}

		for _, stack := range stacks {
			item, err := t.feedStackItem(stack)
			if err != nil {
				return nil, err
			}
			if item == nil {
				continue
			}
			list = append(list, item)
			count += len(stack.children) + 1
		}
	}

	var nextOffset string
	if len(blocks) > 0 {
		nextOffset = blocks[len(blocks)-1].Id

		// see if there's actually more
		if len(t.datastore.Blocks().List(nextOffset, 1, query)) == 0 {
			nextOffset = ""
		}
	}

	return &pb.FeedItemList{
		Items: list,
		Count: int32(count),
		Next:  nextOffset,
	}, nil
}

func (t *Textile) feedItem(block *repo.Block, opts feedItemOpts) (*pb.FeedItem, error) {
	item := &pb.FeedItem{
		Block:   block.Id,
		Thread:  block.ThreadId,
		Payload: &any.Any{},
	}

	var payload proto.Message
	var err error
	switch block.Type {
	case repo.JoinBlock:
		item.Payload.TypeUrl = "/Join"
		payload, err = t.join(block, opts)
	case repo.LeaveBlock:
		item.Payload.TypeUrl = "/Leave"
		payload, err = t.leave(block, opts)
	case repo.FilesBlock:
		item.Payload.TypeUrl = "/Files"
		payload, err = t.file(block, opts)
	case repo.MessageBlock:
		item.Payload.TypeUrl = "/Text"
		payload, err = t.message(block, opts)
	case repo.CommentBlock:
		item.Payload.TypeUrl = "/Comment"
		payload, err = t.comment(block, opts)
	case repo.LikeBlock:
		item.Payload.TypeUrl = "/Like"
		payload, err = t.like(block, opts)
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
	var comments []*pb.Comment
	var likes []*pb.Like

	// Does the stack contain the initial target,
	// or is it a continuation stack of just annotations?
	// We'll need to load the target in the latter case.
	var target *repo.Block
	handleChild := func(child *repo.Block) error {
		switch child.Type {
		case repo.CommentBlock:
			comment, err := t.comment(child, feedItemOpts{annotations: true})
			if err != nil {
				return err
			}
			comments = append(comments, comment)
		case repo.LikeBlock:
			like, err := t.like(child, feedItemOpts{annotations: true})
			if err != nil {
				return err
			}
			likes = append(likes, like)
		default:
			target = child
		}
		return nil
	}
	for _, child := range stack.children {
		if err := handleChild(&child); err != nil {
			return nil, err
		}
	}

	var initial bool
	if target != nil { // target was in children, newer annotations may exist, make target top
		initial = true
		if err := handleChild(&stack.top); err != nil {
			return nil, err
		}
	} else if !isAnnotation(stack.top) { // target is top, newer annotations may exist
		initial = true
		target = &stack.top
	} else { // target needs to be loaded, older annotations may exist
		if t.blockIgnored(stack.id) {
			return nil, nil
		}
		target = t.datastore.Blocks().Get(stack.id)
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

	if !initial {
		// target gets wrapped with the top block
		return t.feedItem(&stack.top, feedItemOpts{
			target: targetItem,
		})
	}

	return targetItem, nil
}

func (t *Textile) blockIgnored(blockId string) bool {
	return len(t.datastore.Blocks().List("", -1, "target='ignore-"+blockId+"'")) > 0
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

package core

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/textileio/go-textile/pb"
)

var flatFeedTypes = []pb.Block_BlockType{
	pb.Block_JOIN,
	pb.Block_LEAVE,
	pb.Block_FILES,
	pb.Block_TEXT,
	pb.Block_COMMENT,
	pb.Block_LIKE,
}

var annotatedFeedTypes = []pb.Block_BlockType{
	pb.Block_JOIN,
	pb.Block_LEAVE,
	pb.Block_FILES,
	pb.Block_TEXT,
}

type feedStack struct {
	id       string
	top      *pb.Block
	children []*pb.Block
}

type feedItemOpts struct {
	annotations bool
	comments    []*pb.Comment
	likes       []*pb.Like
	target      *pb.FeedItem
}

func (t *Textile) Feed(req *pb.FeedRequest) (*pb.FeedItemList, error) {
	var types []pb.Block_BlockType
	switch req.Mode {
	case pb.FeedRequest_CHRONO, pb.FeedRequest_STACKS:
		types = flatFeedTypes
	case pb.FeedRequest_ANNOTATED:
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
	if req.Thread != "" {
		if t.Thread(req.Thread) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("(threadId='%s') and %s", req.Thread, query)
	}

	blocks := t.Blocks(req.Offset, int(req.Limit), query)
	list := make([]*pb.FeedItem, 0)
	var count int

	switch req.Mode {
	case pb.FeedRequest_CHRONO, pb.FeedRequest_ANNOTATED:
		for _, block := range blocks.Items {
			item, err := t.feedItem(block, feedItemOpts{
				annotations: req.Mode == pb.FeedRequest_ANNOTATED,
			})
			if err != nil {
				return nil, err
			}
			list = append(list, item)
			count++
		}

	case pb.FeedRequest_STACKS:
		stacks := make([]feedStack, 0)
		var last *feedStack
		for _, block := range blocks.Items {
			if len(stacks) > 0 {
				last = &stacks[len(stacks)-1]
			} else {
				last = &feedStack{}
			}
			targetId := getTargetId(block)

			if len(stacks) == 0 || targetId != getTargetId(last.top) {
				// start a new stack
				stacks = append(stacks, feedStack{id: targetId, top: block})
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
	if len(blocks.Items) > 0 {
		nextOffset = blocks.Items[len(blocks.Items)-1].Id

		// see if there's actually more
		if len(t.datastore.Blocks().List(nextOffset, 1, query).Items) == 0 {
			nextOffset = ""
		}
	}

	return &pb.FeedItemList{
		Items: list,
		Count: int32(count),
		Next:  nextOffset,
	}, nil
}

func (t *Textile) feedItem(block *pb.Block, opts feedItemOpts) (*pb.FeedItem, error) {
	if block == nil {
		return nil, nil
	}

	item := &pb.FeedItem{
		Block:  block.Id,
		Thread: block.Thread,
		Payload: &any.Any{
			TypeUrl: "/" + strings.Title(strings.ToLower(block.Type.String())),
		},
	}

	var payload proto.Message
	var err error
	switch block.Type {
	case pb.Block_MERGE:
		payload, err = t.merge(block, opts)
	case pb.Block_IGNORE:
		payload, err = t.ignore(block, opts)
	case pb.Block_FLAG:
		payload, err = t.flag(block, opts)
	case pb.Block_JOIN:
		payload, err = t.join(block, opts)
	case pb.Block_ANNOUNCE:
		payload, err = t.announce(block, opts)
	case pb.Block_LEAVE:
		payload, err = t.leave(block, opts)
	case pb.Block_TEXT:
		payload, err = t.message(block, opts)
	case pb.Block_FILES:
		payload, err = t.file(block, opts)
	case pb.Block_COMMENT:
		payload, err = t.comment(block, opts)
	case pb.Block_LIKE:
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

func (t *Textile) feedStackItem(stack feedStack) (*pb.FeedItem, error) {
	var comments []*pb.Comment
	var likes []*pb.Like

	// Does the stack contain the initial target,
	// or is it a continuation stack of just annotations?
	// We'll need to load the target in the latter case.
	var target *pb.Block
	handleChild := func(child *pb.Block) error {
		switch child.Type {
		case pb.Block_COMMENT:
			comment, err := t.comment(child, feedItemOpts{annotations: true})
			if err != nil {
				return err
			}
			comments = append(comments, comment)
		case pb.Block_LIKE:
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
		if err := handleChild(child); err != nil {
			return nil, err
		}
	}

	var initial bool
	if target != nil { // target was in children, newer annotations may exist, make target top
		initial = true
		if err := handleChild(stack.top); err != nil {
			return nil, err
		}
	} else if !isAnnotation(stack.top) { // target is top, newer annotations may exist
		initial = true
		target = stack.top
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
		return t.feedItem(stack.top, feedItemOpts{
			target: targetItem,
		})
	}

	return targetItem, nil
}

func (t *Textile) blockIgnored(blockId string) bool {
	query := fmt.Sprintf("target='%s' and type=%d", blockId, pb.Block_IGNORE)
	return len(t.datastore.Blocks().List("", -1, query).Items) > 0
}

func FeedItemType(item *pb.FeedItem) (pb.Block_BlockType, error) {
	if i, ok := pb.Block_BlockType_value[strings.ToUpper(item.Payload.TypeUrl[1:])]; ok {
		return pb.Block_BlockType(i), nil
	} else {
		return 0, fmt.Errorf("unable to determine block type")
	}
}

type FeedItemPayload interface {
	GetUser() *pb.User
	GetDate() *timestamp.Timestamp
	Reset()
	String() string
	ProtoMessage()
}

func GetFeedItemPayload(item *pb.FeedItem) (FeedItemPayload, error) {
	blockType, err := FeedItemType(item)
	if err != nil {
		return nil, err
	}

	var payload FeedItemPayload
	switch blockType {
	case pb.Block_MERGE:
		payload = new(pb.Merge)
	case pb.Block_IGNORE:
		payload = new(pb.Ignore)
	case pb.Block_FLAG:
		payload = new(pb.Flag)
	case pb.Block_JOIN:
		payload = new(pb.Join)
	case pb.Block_ANNOUNCE:
		payload = new(pb.Announce)
	case pb.Block_LEAVE:
		payload = new(pb.Leave)
	case pb.Block_TEXT:
		payload = new(pb.Text)
	case pb.Block_FILES:
		payload = new(pb.Files)
	case pb.Block_COMMENT:
		payload = new(pb.Comment)
	case pb.Block_LIKE:
		payload = new(pb.Like)
	default:
		return nil, fmt.Errorf("unable to parse payload")
	}

	if err := ptypes.UnmarshalAny(item.Payload, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func getTargetId(block *pb.Block) string {
	switch block.Type {
	case pb.Block_COMMENT, pb.Block_LIKE:
		return block.Target
	default:
		return block.Id
	}
}

func isAnnotation(block *pb.Block) bool {
	switch block.Type {
	case pb.Block_COMMENT, pb.Block_LIKE:
		return true
	default:
		return false
	}
}

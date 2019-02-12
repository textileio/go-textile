package core

import (
	"fmt"
	"time"

	"github.com/textileio/textile-go/repo"
)

type ThreadFeedItemType string

const (
	JoinThreadFeedItem    ThreadFeedItemType = "join"
	LeaveThreadFeedItem   ThreadFeedItemType = "leave"
	FilesThreadFeedItem   ThreadFeedItemType = "files"
	MessageThreadFeedItem ThreadFeedItemType = "message"
	CommentThreadFeedItem ThreadFeedItemType = "comment"
	LikeThreadFeedItem    ThreadFeedItemType = "like"
)

type ThreadFeedItem struct {
	Block   string             `json:"block"`
	Type    ThreadFeedItemType `json:"type"`
	Join    *ThreadJoinInfo    `json:"join,omitempty"`
	Leave   *ThreadLeaveInfo   `json:"leave,omitempty"`
	Files   *ThreadFilesInfo   `json:"files,omitempty"`
	Message *ThreadMessageInfo `json:"message,omitempty"`
	Comment *ThreadCommentInfo `json:"comment,omitempty"`
	Like    *ThreadLikeInfo    `json:"like,omitempty"`
}

var feedTypes = []repo.BlockType{
	repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock, repo.CommentBlock, repo.LikeBlock,
}

var annotatedFeedTypes = []repo.BlockType{
	repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock,
}

func (t *Textile) ThreadFeed(offset string, limit int, threadId string, annotated bool) ([]ThreadFeedItem, error) {
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

	list := make([]ThreadFeedItem, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		item, err := t.threadItem(&block, annotated)
		if err != nil {
			return nil, err
		}
		list = append(list, *item)
	}

	return list, nil
}

func (t *Textile) threadItem(block *repo.Block, annotated bool) (*ThreadFeedItem, error) {
	item := &ThreadFeedItem{
		Block: block.Id,
	}

	var err error
	switch block.Type {
	case repo.JoinBlock:
		item.Type = JoinThreadFeedItem
		item.Join, err = t.ThreadJoin(block, annotated)
	case repo.LeaveBlock:
		item.Type = LeaveThreadFeedItem
		item.Leave, err = t.ThreadLeave(block, annotated)
	case repo.FilesBlock:
		item.Type = FilesThreadFeedItem
		item.Files, err = t.threadFile(block, annotated)
	case repo.MessageBlock:
		item.Type = MessageThreadFeedItem
		item.Message, err = t.threadMessage(block, annotated)
	case repo.CommentBlock:
		item.Type = CommentThreadFeedItem
		item.Comment, err = t.ThreadComment(block, !annotated)
	case repo.LikeBlock:
		item.Type = LikeThreadFeedItem
		item.Like, err = t.ThreadLike(block, !annotated)
	default:
		return nil, nil
	}

	return item, err
}

type ThreadJoinInfo struct {
	Block    string           `json:"block"`
	Date     time.Time        `json:"date"`
	AuthorId string           `json:"author_id"`
	Username string           `json:"username,omitempty"`
	Avatar   string           `json:"avatar,omitempty"`
	Likes    []ThreadLikeInfo `json:"likes"`
}

type ThreadLeaveInfo struct {
	Block    string           `json:"block"`
	Date     time.Time        `json:"date"`
	AuthorId string           `json:"author_id"`
	Username string           `json:"username,omitempty"`
	Avatar   string           `json:"avatar,omitempty"`
	Likes    []ThreadLikeInfo `json:"likes"`
}

func (t *Textile) ThreadJoin(block *repo.Block, annotated bool) (*ThreadJoinInfo, error) {
	if block.Type != repo.JoinBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	info := &ThreadJoinInfo{
		Block:    block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
	}

	if annotated {
		likes, err := t.ThreadLikes(block.Id)
		if err != nil {
			return nil, err
		}
		info.Likes = likes
	}

	return info, nil
}

func (t *Textile) ThreadLeave(block *repo.Block, annotated bool) (*ThreadLeaveInfo, error) {
	if block.Type != repo.LeaveBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	info := &ThreadLeaveInfo{
		Block:    block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
	}

	if annotated {
		likes, err := t.ThreadLikes(block.Id)
		if err != nil {
			return nil, err
		}
		info.Likes = likes
	}

	return info, nil
}

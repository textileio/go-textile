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
)

type ThreadFeedItem struct {
	Block   string             `json:"block"`
	Type    ThreadFeedItemType `json:"type"`
	Join    *ThreadJoinInfo    `json:"join,omitempty"`
	Leave   *ThreadLeaveInfo   `json:"leave,omitempty"`
	Files   *ThreadFilesInfo   `json:"files,omitempty"`
	Message *ThreadMessageInfo `json:"message,omitempty"`
}

func (t *Textile) ThreadFeed(offset string, limit int, threadId string) ([]ThreadFeedItem, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		stm := "(threadId='%s') and (type=%d or type=%d or type=%d or type=%d)"
		query = fmt.Sprintf(stm, threadId, repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock)
	} else {
		stm := "(type=%d or type=%d or type=%d or type=%d)"
		query = fmt.Sprintf(stm, repo.JoinBlock, repo.LeaveBlock, repo.FilesBlock, repo.MessageBlock)
	}

	list := make([]ThreadFeedItem, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		item := ThreadFeedItem{
			Block: block.Id,
		}
		var err error
		switch block.Type {
		case repo.JoinBlock:
			item.Type = JoinThreadFeedItem
			item.Join, err = t.ThreadJoin(block)
		case repo.LeaveBlock:
			item.Type = LeaveThreadFeedItem
			item.Leave, err = t.ThreadLeave(block)
		case repo.FilesBlock:
			item.Type = FilesThreadFeedItem
			item.Files, err = t.threadFile(block)
		case repo.MessageBlock:
			item.Type = MessageThreadFeedItem
			item.Message, err = t.ThreadMessage(block)
		default:
			continue
		}
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
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

func (t *Textile) ThreadJoin(block repo.Block) (*ThreadJoinInfo, error) {
	if block.Type != repo.JoinBlock {
		return nil, ErrBlockWrongType
	}

	likes, err := t.ThreadLikes(block.Id)
	if err != nil {
		return nil, err
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	return &ThreadJoinInfo{
		Block:    block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Likes:    likes,
	}, nil
}

func (t *Textile) ThreadLeave(block repo.Block) (*ThreadLeaveInfo, error) {
	if block.Type != repo.LeaveBlock {
		return nil, ErrBlockWrongType
	}

	likes, err := t.ThreadLikes(block.Id)
	if err != nil {
		return nil, err
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	return &ThreadLeaveInfo{
		Block:    block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Likes:    likes,
	}, nil
}

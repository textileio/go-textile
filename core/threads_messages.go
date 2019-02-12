package core

import (
	"fmt"
	"time"

	"github.com/textileio/textile-go/repo"
)

type ThreadMessageInfo struct {
	Block    string              `json:"block"`
	Date     time.Time           `json:"date"`
	AuthorId string              `json:"author_id"`
	Username string              `json:"username,omitempty"`
	Avatar   string              `json:"avatar,omitempty"`
	Body     string              `json:"body"`
	Comments []ThreadCommentInfo `json:"comments"`
	Likes    []ThreadLikeInfo    `json:"likes"`
}

func (t *Textile) ThreadMessages(offset string, limit int, threadId string) ([]ThreadMessageInfo, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("threadId='%s' and type=%d", threadId, repo.MessageBlock)
	} else {
		query = fmt.Sprintf("type=%d", repo.MessageBlock)
	}

	list := make([]ThreadMessageInfo, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		msg, err := t.threadMessage(&block, true)
		if err != nil {
			return nil, err
		}
		list = append(list, *msg)
	}

	return list, nil
}
func (t *Textile) ThreadMessage(blockId string) (*ThreadMessageInfo, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.threadMessage(block, true)
}

func (t *Textile) threadMessage(block *repo.Block, annotated bool) (*ThreadMessageInfo, error) {
	if block.Type != repo.MessageBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	info := &ThreadMessageInfo{
		Block:    block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Body:     block.Body,
	}

	if annotated {
		comments, err := t.ThreadComments(block.Id)
		if err != nil {
			return nil, err
		}
		info.Comments = comments

		likes, err := t.ThreadLikes(block.Id)
		if err != nil {
			return nil, err
		}
		info.Likes = likes
	}

	return info, nil
}

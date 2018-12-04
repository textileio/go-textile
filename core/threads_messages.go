package core

import (
	"fmt"
	"time"

	"github.com/textileio/textile-go/repo"
)

type ThreadMessageInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
	Body     string    `json:"body"`
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
		msg, err := t.ThreadMessage(block)
		if err != nil {
			return nil, err
		}
		list = append(list, *msg)
	}

	return list, nil
}

func (t *Textile) ThreadMessage(block repo.Block) (*ThreadMessageInfo, error) {
	if block.Type != repo.MessageBlock {
		return nil, ErrBlockWrongType
	}

	return &ThreadMessageInfo{
		Id:       block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: t.ContactUsername(block.AuthorId),
		Body:     block.Body,
	}, nil
}

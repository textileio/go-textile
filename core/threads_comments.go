package core

import (
	"fmt"
	"time"

	"github.com/textileio/textile-go/repo"
)

type ThreadCommentInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
	Avatar   string    `json:"avatar,omitempty"`
	Body     string    `json:"body"`
}

func (t *Textile) ThreadComments(target string) ([]ThreadCommentInfo, error) {
	comments := make([]ThreadCommentInfo, 0)

	query := fmt.Sprintf("type=%d and target='%s'", repo.CommentBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info, err := t.ThreadComment(block)
		if err != nil {
			continue
		}
		comments = append(comments, *info)
	}

	return comments, nil
}

func (t *Textile) ThreadComment(block repo.Block) (*ThreadCommentInfo, error) {
	if block.Type != repo.CommentBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	return &ThreadCommentInfo{
		Id:       block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Body:     block.Body,
	}, nil
}

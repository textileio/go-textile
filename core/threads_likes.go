package core

import (
	"fmt"
	"time"

	"github.com/textileio/textile-go/repo"
)

type ThreadLikeInfo struct {
	Id       string          `json:"id"`
	Date     time.Time       `json:"date"`
	AuthorId string          `json:"author_id"`
	Username string          `json:"username,omitempty"`
	Avatar   string          `json:"avatar,omitempty"`
	Target   *ThreadFeedItem `json:"target"`
}

func (t *Textile) ThreadLikes(target string) ([]ThreadLikeInfo, error) {
	likes := make([]ThreadLikeInfo, 0)

	query := fmt.Sprintf("type=%d and target='%s'", repo.LikeBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info, err := t.ThreadLike(&block, true)
		if err != nil {
			continue
		}
		likes = append(likes, *info)
	}

	return likes, nil
}

func (t *Textile) ThreadLike(block *repo.Block, annotation bool) (*ThreadLikeInfo, error) {
	if block.Type != repo.LikeBlock {
		return nil, ErrBlockWrongType
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	info := &ThreadLikeInfo{
		Id:       block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
	}

	if !annotation {
		target, err := t.threadItem(t.datastore.Blocks().Get(block.Target), false)
		if err != nil {
			return nil, err
		}
		info.Target = target
	}

	return info, nil
}

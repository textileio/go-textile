package core

import (
	"errors"

	"github.com/textileio/textile-go/repo"
)

// ErrBlockNotFound indicates a block was not found in the index
var ErrBlockNotFound = errors.New("block not found")

// GetBlocks paginates blocks
func (t *Textile) Blocks(offset string, limit int, query string) []repo.Block {
	var filtered []repo.Block

	for _, block := range t.datastore.Blocks().List(offset, limit, query) {
		ignored := t.datastore.Blocks().List("", -1, "target='ignore-"+block.Id+"'")
		if len(ignored) == 0 {
			filtered = append(filtered, block)
		}
	}

	return filtered
}

// Block returns block with id
func (t *Textile) Block(id string) (*repo.Block, error) {
	block := t.datastore.Blocks().Get(id)
	if block == nil {
		return nil, ErrBlockNotFound
	}
	return block, nil
}

// BlocksByTarget returns block with parent
func (t *Textile) BlocksByTarget(target string) []repo.Block {
	return t.datastore.Blocks().List("", -1, "target='"+target+"'")
}

// BlockInfo returns block info with id
func (t *Textile) BlockInfo(id string) (*BlockInfo, error) {
	block, err := t.Block(id)
	if err != nil {
		return nil, err
	}

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	return &BlockInfo{
		Id:       block.Id,
		ThreadId: block.ThreadId,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Type:     block.Type.Description(),
		Date:     block.Date,
		Parents:  block.Parents,
		Target:   block.Target,
		Body:     block.Body,
	}, nil
}

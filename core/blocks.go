package core

import (
	"errors"
	"fmt"
	"github.com/textileio/textile-go/repo"
)

// ErrBlockNotFound indicates a block was not found in the index
var ErrBlockNotFound = errors.New("block not found")

// GetBlocks paginates blocks
func (t *Textile) Blocks(offset string, limit int, query string) []repo.Block {
	var filtered []repo.Block
	for _, block := range t.datastore.Blocks().List(offset, limit, query) {
		ignored := t.datastore.Blocks().GetByTarget(fmt.Sprintf("ignore-%s", block.Id))
		if ignored == nil {
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

// BlockByParent returns block with parent
func (t *Textile) BlockByParent(target string) (*repo.Block, error) {
	block := t.datastore.Blocks().GetByTarget(target)
	if block == nil {
		return nil, ErrBlockNotFound
	}
	return block, nil
}

// BlockInfo returns block info with id
func (t *Textile) BlockInfo(id string) (*BlockInfo, error) {
	block, err := t.Block(id)
	if err != nil {
		return nil, err
	}

	return &BlockInfo{
		Id:       block.Id,
		ThreadId: block.ThreadId,
		AuthorId: block.AuthorId,
		Type:     block.Type.Description(),
		Date:     block.Date,
		Parents:  block.Parents,
		Target:   block.Target,
		Body:     block.Body,
	}, nil
}

// BlockData cats file data from ipfs and tries to decrypt it with the provided block
//func (t *Textile) BlockData(path string, block *repo.Block) ([]byte, error) {
//	ciphertext, err := ipfs.DataAtPath(t.node, path)
//	if err != nil {
//		// size migrations
//		parts := strings.Split(path, "/")
//		if len(parts) > 1 && strings.Contains(err.Error(), "no link named") {
//			switch parts[1] {
//			case "small":
//				parts[1] = "thumb"
//			case "medium":
//				parts[1] = "photo"
//			default:
//				return nil, err
//			}
//			ciphertext, err = ipfs.DataAtPath(t.node, strings.Join(parts, "/"))
//			if err != nil {
//				return nil, err
//			}
//		} else {
//			return nil, err
//		}
//	}
//	return crypto.DecryptAES(ciphertext, block.DataKey)
//}

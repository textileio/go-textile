package core

import (
	"errors"

	"github.com/textileio/textile-go/pb"
)

// ErrBlockNotFound indicates a block was not found in the index
var ErrBlockNotFound = errors.New("block not found")

// GetBlocks paginates blocks
func (t *Textile) Blocks(offset string, limit int, query string) *pb.BlockList {
	filtered := &pb.BlockList{Items: make([]*pb.Block, 0)}

	for _, block := range t.datastore.Blocks().List(offset, limit, query).Items {
		ignored := t.datastore.Blocks().List("", -1, "target='ignore-"+block.Id+"'")
		if len(ignored.Items) == 0 {
			filtered.Items = append(filtered.Items, block)
		}
	}

	return filtered
}

// Block returns block with id
func (t *Textile) Block(id string) (*pb.Block, error) {
	block := t.datastore.Blocks().Get(id)
	if block == nil {
		return nil, ErrBlockNotFound
	}
	return block, nil
}

// BlocksByTarget returns block with parent
func (t *Textile) BlocksByTarget(target string) *pb.BlockList {
	return t.datastore.Blocks().List("", -1, "target='"+target+"'")
}

// BlockView returns block with expanded view properties
func (t *Textile) BlockView(id string) (*pb.Block, error) {
	block, err := t.Block(id)
	if err != nil {
		return nil, err
	}

	block.User = t.User(block.Author)
	return block, nil
}

package core

import (
	"errors"
	"github.com/textileio/textile-go/repo"
)

// GetBlock searches for a local block associated with the given target
func (t *Textile) GetBlock(id string) (*repo.Block, error) {
	block := t.datastore.Blocks().Get(id)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
}

// GetBlockByDataId searches for a local block associated with the given data id
func (t *Textile) GetBlockByDataId(dataId string) (*repo.Block, error) {
	if dataId == "" {
		return nil, nil
	}
	block := t.datastore.Blocks().GetByData(dataId)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
}

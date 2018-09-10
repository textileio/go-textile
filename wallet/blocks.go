package wallet

import (
	"errors"
	"github.com/textileio/textile-go/repo"
)

// GetBlock searches for a local block associated with the given target
func (w *Wallet) GetBlock(id string) (*repo.Block, error) {
	block := w.datastore.Blocks().Get(id)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
}

// GetBlockByDataId searches for a local block associated with the given data id
func (w *Wallet) GetBlockByDataId(dataId string) (*repo.Block, error) {
	if dataId == "" {
		return nil, nil
	}
	block := w.datastore.Blocks().GetByDataId(dataId)
	if block == nil {
		return nil, errors.New("block not found locally")
	}
	return block, nil
}

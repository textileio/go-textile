package core

import (
	"errors"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"strings"
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

// GetBlockData cats file data from ipfs and tries to decrypt it with the provided block
func (t *Textile) GetBlockData(path string, block *repo.Block) ([]byte, error) {
	ciphertext, err := ipfs.GetDataAtPath(t.ipfs, path)
	if err != nil {
		// size migrations
		parts := strings.Split(path, "/")
		if len(parts) > 1 && strings.Contains(err.Error(), "no link named") {
			switch parts[1] {
			case "small":
				parts[1] = "thumb"
			case "medium":
				parts[1] = "photo"
			default:
				return nil, err
			}
			ciphertext, err = ipfs.GetDataAtPath(t.ipfs, strings.Join(parts, "/"))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return crypto.DecryptAES(ciphertext, []byte(block.DataKey))
}

package thread

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/wallet/model"
)

// GetBlockDataKey returns the decrypted AES key for a block
func (t *Thread) GetBlockDataKey(block *repo.Block) ([]byte, error) {
	if block.Type != repo.PhotoBlock {
		return nil, errors.New("incorrect block type")
	}
	key, err := t.Decrypt(block.DataKeyCipher)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GetBlockData cats file data from ipfs and tries to decrypt it with the provided block
func (t *Thread) GetBlockData(path string, block *repo.Block) ([]byte, error) {
	key, err := t.GetBlockDataKey(block)
	if err != nil {
		return nil, err
	}
	cipher, err := util.GetDataAtPath(t.ipfs(), path)
	if err != nil {
		return nil, err
	}
	return crypto.DecryptAES(cipher, key)
}

// GetPhotoMetaData returns photo metadata under an id
func (t *Thread) GetPhotoMetaData(id string, block *repo.Block) (*model.PhotoMetadata, error) {
	file, err := t.GetBlockData(fmt.Sprintf("%s/meta", id), block)
	if err != nil {
		return nil, err
	}
	var data *model.PhotoMetadata
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

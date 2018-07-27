package thread

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/wallet/model"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// GetBlockDataKey returns the decrypted AES key for a block
func (t *Thread) GetBlockDataKey(block *repo.Block) ([]byte, error) {
	if block.Type != repo.PhotoBlock {
		return nil, errors.New("incorrect block type")
	}
	key, err := t.Decrypt(block.DataKeyCipher)
	if err != nil {
		log.Errorf("error decrypting key cipher: %s", err)
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
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}
	return crypto.DecryptAES(cipher, key)
}

// GetFileDataBase64 returns file data encoded as base64 under an ipfs path
func (t *Thread) GetBlockDataBase64(path string, block *repo.Block) (string, error) {
	data, err := t.GetBlockData(path, block)
	if err != nil {
		return "error", err
	}
	return libp2pc.ConfigEncodeKey(data), nil
}

// GetPhotoMetaData returns photo metadata under an id
func (t *Thread) GetPhotoMetaData(id string, block *repo.Block) (*model.PhotoMetadata, error) {
	file, err := t.GetBlockData(fmt.Sprintf("%s/meta", id), block)
	if err != nil {
		log.Errorf("error getting meta file %s: %s", id, err)
		return nil, err
	}
	var data *model.PhotoMetadata
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Errorf("error unmarshaling meta file: %s: %s", id, err)
		return nil, err
	}
	return data, nil
}

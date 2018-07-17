package thread

import (
	"encoding/json"
	"fmt"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/util"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// GetBlockData cats file data from ipfs and tries to decrypt it with the provided block
func (t *Thread) GetBlockData(path string, block *repo.Block) ([]byte, error) {
	cipher, err := util.GetDataAtPath(t.ipfs(), path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// decrypt with thread key
	return t.Decrypt(cipher)
}

// GetBlockDataBase64 returns block data encoded as base64 under an ipfs path
func (t *Thread) GetBlockDataBase64(path string, block *repo.Block) (string, error) {
	file, err := t.GetBlockData(path, block)
	if err != nil {
		return "error", err
	}
	return libp2pc.ConfigEncodeKey(file), nil
}

// GetFileKey returns the decrypted AES key for a block
func (t *Thread) GetFileKey(block *repo.Block) (string, error) {
	key, err := t.Decrypt(block.TargetKey)
	if err != nil {
		log.Errorf("error decrypting key: %s", err)
		return "", err
	}
	return string(key), nil
}

// GetFileData cats file data from ipfs and tries to decrypt it with the provided block
func (t *Thread) GetFileData(path string, block *repo.Block) ([]byte, error) {
	// get bytes
	cipher, err := util.GetDataAtPath(t.ipfs(), path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// decrypt the file key
	key, err := t.Decrypt(block.TargetKey)
	if err != nil {
		log.Errorf("error decrypting key: %s", err)
		return nil, err
	}

	// finally, decrypt the file
	return crypto.DecryptAES(cipher, key)
}

// GetFileDataBase64 returns file data encoded as base64 under an ipfs path
func (t *Thread) GetFileDataBase64(path string, block *repo.Block) (string, error) {
	file, err := t.GetFileData(path, block)
	if err != nil {
		return "error", err
	}
	return libp2pc.ConfigEncodeKey(file), nil
}

// GetMetaData returns photo metadata under an id
func (t *Thread) GetPhotoMetaData(id string, block *repo.Block) (*model.PhotoMetadata, error) {
	file, err := t.GetFileData(fmt.Sprintf("%s/meta", id), block)
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

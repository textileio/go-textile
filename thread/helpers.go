package thread

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/photo"
	"github.com/textileio/textile-go/repo"
	"strings"
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
	cipher, err := ipfs.GetDataAtPath(t.ipfs(), path)
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
			cipher, err = ipfs.GetDataAtPath(t.ipfs(), strings.Join(parts, "/"))
		}
		if err != nil {
			return nil, err
		}
	}
	return crypto.DecryptAES(cipher, key)
}

// GetPhotoMetaData returns photo metadata indexed under a block
func (t *Thread) GetPhotoMetaData(id string, block *repo.Block) (*photo.Metadata, error) {
	key, err := t.GetBlockDataKey(block)
	if err != nil {
		return nil, err
	}
	if block.DataMetadataCipher == nil {
		return nil, errors.New("metadata was not indexed")
	}
	metadatab, err := crypto.DecryptAES(block.DataMetadataCipher, key)
	if err != nil {
		return nil, err
	}
	var metadata *photo.Metadata
	if err := json.Unmarshal(metadatab, &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

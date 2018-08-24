package mobile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"image"
	"time"
)

// Photo is a simple meta data wrapper around a photo block
type Photo struct {
	Id       string              `json:"id"`
	BlockId  string              `json:"block_id"`
	Date     time.Time           `json:"date"`
	AuthorId string              `json:"author_id"`
	Caption  string              `json:"caption,omitempty"`
	Username string              `json:"username,omitempty"`
	Metadata *util.PhotoMetadata `json:"metadata,omitempty"`
}

// Photos is a wrapper around a list of photos
type Photos struct {
	Items []Photo `json:"items"`
}

// ImageData is a wrapper around an image data url
type ImageData struct {
	Url string `json:"url"`
}

// AddPhoto adds a photo by path
func (m *Mobile) AddPhoto(path string) (string, error) {
	added, err := core.Node.Wallet.AddPhoto(path)
	if err != nil {
		return "", err
	}
	return toJSON(added)
}

// SharePhoto adds an existing photo to a new thread
func (m *Mobile) AddPhotoToThread(dataId string, key string, threadId string, caption string) (string, error) {
	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", threadId))
	}

	addr, err := thrd.AddPhoto(dataId, caption, []byte(key))
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// SharePhoto adds an existing photo to a new thread
func (m *Mobile) SharePhotoToThread(dataId string, threadId string, caption string) (string, error) {
	block, err := core.Node.Wallet.GetBlockByDataId(dataId)
	if err != nil {
		return "", err
	}
	_, fromThread := core.Node.Wallet.GetThread(block.ThreadId)
	if fromThread == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}
	_, toThread := core.Node.Wallet.GetThread(threadId)
	if toThread == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", threadId))
	}
	key, err := fromThread.Decrypt(block.DataKeyCipher)
	if err != nil {
		return "", err
	}

	// TODO: owner challenge

	addr, err := toThread.AddPhoto(dataId, caption, key)
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// GetPhotos returns thread photo blocks with json encoding
func (m *Mobile) GetPhotos(offsetId string, limit int, threadId string) (string, error) {
	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("thread not found: %s", threadId))
	}

	// build json
	photos := &Photos{Items: make([]Photo, 0)}
	btype := repo.PhotoBlock
	for _, b := range thrd.Blocks(offsetId, limit, &btype) {
		var username, caption string
		var metadata *util.PhotoMetadata
		if b.AuthorUnCipher != nil {
			usernameb, err := thrd.Decrypt(b.AuthorUnCipher)
			if err != nil {
				return "", err
			}
			username = string(usernameb)
		}
		if b.DataCaptionCipher != nil {
			captionb, err := thrd.Decrypt(b.DataCaptionCipher)
			if err != nil {
				return "", err
			}
			caption = string(captionb)
		}
		if b.DataMetadataCipher != nil {
			key, err := thrd.Decrypt(b.DataKeyCipher)
			if err != nil {
				return "", err
			}
			metadatab, err := crypto.DecryptAES(b.DataMetadataCipher, key)
			if err != nil {
				return "", err
			}
			if err := json.Unmarshal(metadatab, &metadata); err != nil {
				log.Warningf("unmarshal photo metadata failed: %s", err)
			}
		}
		authorId, err := util.IdFromEncodedPublicKey(b.AuthorPk)
		if err != nil {
			return "", err
		}
		item := Photo{
			Id:       b.DataId,
			BlockId:  b.Id,
			Date:     b.Date,
			AuthorId: authorId.Pretty(),
			Metadata: metadata,
		}
		if username != "" {
			item.Username = username
		}
		if caption != "" {
			item.Caption = caption
		}
		photos.Items = append(photos.Items, item)
	}

	return toJSON(photos)
}

// IgnorePhoto adds an ignore block targeted at the given block and unpins the photo set locally
func (m *Mobile) IgnorePhoto(blockId string) (string, error) {
	block, err := core.Node.Wallet.GetBlock(blockId)
	if err != nil {
		return "", err
	}
	_, thrd := core.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}

	addr, err := thrd.Ignore(block.Id)
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// AddPhotoComment adds an comment block targeted at the given block
func (m *Mobile) AddPhotoComment(blockId string, body string) (string, error) {
	block, err := core.Node.Wallet.GetBlock(blockId)
	if err != nil {
		return "", err
	}
	_, thrd := core.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}

	addr, err := thrd.AddComment(block.Id, body)
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// AddPhotoLike adds a like block targeted at the given block
func (m *Mobile) AddPhotoLike(blockId string) (string, error) {
	block, err := core.Node.Wallet.GetBlock(blockId)
	if err != nil {
		return "", err
	}
	_, thrd := core.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}

	addr, err := thrd.AddLike(block.Id)
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// GetPhotoData returns a data url of an image under a path
func (m *Mobile) GetPhotoData(id string, path string) (string, error) {
	block, err := core.Node.Wallet.GetBlockByDataId(id)
	if err != nil {
		log.Errorf("could not find block for data id %s: %s", id, err)
		return "", err
	}
	_, thrd := core.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		err := errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadId))
		log.Error(err.Error())
		return "", err
	}
	data, err := thrd.GetBlockData(fmt.Sprintf("%s/%s", id, path), block)
	if err != nil {
		log.Errorf("get block data failed %s: %s", id, err)
		return "", err
	}
	var format string
	meta, err := thrd.GetPhotoMetaData(id, block)
	if err != nil {
		log.Warningf("get indexed photo meta data failed, decoding...")
		_, format, err = image.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			log.Errorf("could not determine image format: %s", err)
			return "", err
		}
	} else {
		format = meta.EncodingFormat
	}
	prefix := getImageDataURLPrefix(util.Format(format))
	encoded := libp2pc.ConfigEncodeKey(data)
	img := &ImageData{Url: prefix + encoded}
	return toJSON(img)
}

// GetPhotoDataForSize returns a data url of an image at or above requested size, or the next best option
func (m *Mobile) GetPhotoDataForMinWidth(id string, minWidth int) (string, error) {
	path := util.ImagePathForSize(util.ImageSizeForMinWidth(minWidth))
	return m.GetPhotoData(id, string(path))
}

// GetPhotoMetadata returns a meta data object for a photo
func (m *Mobile) GetPhotoMetadata(id string) (string, error) {
	block, err := core.Node.Wallet.GetBlockByDataId(id)
	if err != nil {
		log.Errorf("could not find block for data id %s: %s", id, err)
		return "", err
	}
	_, thrd := core.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		err := errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadId))
		log.Error(err.Error())
		return "", err
	}
	meta, err := thrd.GetPhotoMetaData(id, block)
	if err != nil {
		log.Warningf("get photo meta data failed %s: %s", id, err)
		meta = &util.PhotoMetadata{}
	}
	return toJSON(meta)
}

// GetPhotoKey calls core GetPhotoKey
func (m *Mobile) GetPhotoKey(id string) (string, error) {
	return core.Node.Wallet.GetPhotoKey(id)
}

// PhotoThreads call core PhotoThreads
func (m *Mobile) PhotoThreads(id string) (string, error) {
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range core.Node.Wallet.PhotoThreads(id) {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

// getImageDataURLPrefix adds the correct data url prefix to a data url
func getImageDataURLPrefix(format util.Format) string {
	switch format {
	case util.PNG:
		return "data:image/png;base64,"
	case util.GIF:
		return "data:image/gif;base64,"
	default:
		return "data:image/jpeg;base64,"
	}
}

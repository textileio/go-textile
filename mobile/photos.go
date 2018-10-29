package mobile

import (
	"errors"
	"fmt"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/images"
	"github.com/textileio/textile-go/repo"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"time"
)

// Photo is a simple meta data wrapper around a photo block
type Photo struct {
	Id       string           `json:"id"`
	BlockId  string           `json:"block_id"`
	Date     time.Time        `json:"date"`
	AuthorId string           `json:"author_id"`
	Caption  string           `json:"caption,omitempty"`
	Username string           `json:"username,omitempty"`
	Metadata *images.Metadata `json:"metadata,omitempty"`
	Comments []Comment        `json:"comments"`
	Likes    []Like           `json:"likes"`
}

// Photos is a wrapper around a list of photos
type Photos struct {
	Items []Photo `json:"items"`
}

// Annotation represents common annotation fields
type Annotation struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
}

// Comment is a simple wrapper around a comment block
type Comment struct {
	Annotation
	Body string `json:"body"`
}

// Like is a simple wrapper around a like block
type Like struct {
	Annotation
}

// ImageData is a wrapper around an image data url
type ImageData struct {
	Url string `json:"url"`
}

// AddPhoto adds a photo by path
func (m *Mobile) AddPhoto(path string) (string, error) {
	added, err := core.Node.AddImageByPath(path)
	if err != nil {
		return "", err
	}
	return toJSON(added)
}

// SharePhoto adds an existing photo to a new thread
func (m *Mobile) AddPhotoToThread(dataId string, key string, threadId string, caption string) (string, error) {
	_, thrd := core.Node.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", threadId))
	}
	keyb, err := base58.Decode(key)
	if err != nil {
		return "", err
	}
	hash, err := thrd.AddPhoto(dataId, caption, keyb)
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// SharePhoto adds an existing photo to a new thread
func (m *Mobile) SharePhotoToThread(dataId string, threadId string, caption string) (string, error) {
	block, err := core.Node.GetBlockByDataId(dataId)
	if err != nil {
		return "", err
	}
	if block == nil {
		return "", errors.New(fmt.Sprintf("could not find block with data id: %s", dataId))
	}
	_, toThread := core.Node.GetThread(threadId)
	if toThread == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", threadId))
	}
	// TODO: owner challenge
	hash, err := toThread.AddPhoto(dataId, caption, block.DataKey)
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// GetPhotos returns thread photo blocks with json encoding
func (m *Mobile) GetPhotos(offsetId string, limit int, threadId string) (string, error) {
	_, thrd := core.Node.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("thread not found: %s", threadId))
	}

	// build json
	photos := &Photos{Items: make([]Photo, 0)}
	btype := repo.PhotoBlock
	for _, b := range thrd.Blocks(offsetId, limit, &btype, nil) {
		item := Photo{
			Id:       b.DataId,
			BlockId:  b.Id,
			Date:     b.Date,
			AuthorId: b.AuthorId,
			Caption:  b.DataCaption,
			Username: getUsername(b.AuthorId[:8]),
			Metadata: b.DataMetadata,
		}

		// add comments
		item.Comments = make([]Comment, 0)
		ctype := repo.CommentBlock
		for _, c := range thrd.Blocks("", -1, &ctype, &b.Id) {
			comment := Comment{
				Annotation: Annotation{
					Id:       c.Id,
					Date:     c.Date,
					AuthorId: c.AuthorId,
					Username: getUsername(c.AuthorId[:8]),
				},
				Body: c.DataCaption,
			}
			item.Comments = append(item.Comments, comment)
		}

		// add likes
		item.Likes = make([]Like, 0)
		ltype := repo.LikeBlock
		for _, l := range thrd.Blocks("", -1, &ltype, &b.Id) {
			like := Like{
				Annotation: Annotation{
					Id:       l.Id,
					Date:     l.Date,
					AuthorId: l.AuthorId,
					Username: getUsername(l.AuthorId[:8]),
				},
			}
			item.Likes = append(item.Likes, like)
		}

		// collect
		photos.Items = append(photos.Items, item)
	}
	return toJSON(photos)
}

// IgnorePhoto is a semantic helper for mobile, just call IgnoreBlock
func (m *Mobile) IgnorePhoto(blockId string) (string, error) {
	return m.ignoreBlock(blockId)
}

// AddPhotoComment adds an comment block targeted at the given block
func (m *Mobile) AddPhotoComment(blockId string, body string) (string, error) {
	block, err := core.Node.GetBlock(blockId)
	if err != nil {
		return "", err
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}
	hash, err := thrd.AddComment(block.Id, body)
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// IgnorePhotoComment is a semantic helper for mobile, just call IgnoreBlock
func (m *Mobile) IgnorePhotoComment(blockId string) (string, error) {
	return m.ignoreBlock(blockId)
}

// AddPhotoLike adds a like block targeted at the given block
func (m *Mobile) AddPhotoLike(blockId string) (string, error) {
	block, err := core.Node.GetBlock(blockId)
	if err != nil {
		return "", err
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}
	hash, err := thrd.AddLike(block.Id)
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// IgnorePhotoLike is a semantic helper for mobile, just call IgnoreBlock
func (m *Mobile) IgnorePhotoLike(blockId string) (string, error) {
	return m.ignoreBlock(blockId)
}

// GetPhotoData returns a data url of an image under a path
func (m *Mobile) GetPhotoData(id string, path string) (string, error) {
	block, err := core.Node.GetBlockByDataId(id)
	if err != nil {
		log.Errorf("could not find block for data id %s: %s", id, err)
		return "", err
	}
	if block == nil {
		err := errors.New(fmt.Sprintf("could not find block with data id: %s", id))
		log.Error(err.Error())
		return "", err
	}
	data, err := core.Node.GetBlockData(fmt.Sprintf("%s/%s", id, path), block)
	if err != nil {
		log.Errorf("get block data failed %s: %s", id, err)
		return "", err
	}
	format := block.DataMetadata.EncodingFormat
	prefix := getImageDataURLPrefix(images.Format(format))
	encoded := libp2pc.ConfigEncodeKey(data)
	img := &ImageData{Url: prefix + encoded}
	return toJSON(img)
}

// GetPhotoDataForSize returns a data url of an image at or above requested size, or the next best option
func (m *Mobile) GetPhotoDataForMinWidth(id string, minWidth int) (string, error) {
	path := images.ImagePathForSize(images.ImageSizeForMinWidth(minWidth))
	return m.GetPhotoData(id, string(path))
}

// GetPhotoMetadata returns a meta data object for a photo
func (m *Mobile) GetPhotoMetadata(id string) (string, error) {
	block, err := core.Node.GetBlockByDataId(id)
	if err != nil {
		log.Errorf("could not find block for data id %s: %s", id, err)
		return "", err
	}
	if block == nil {
		err := errors.New(fmt.Sprintf("could not find block with data id: %s", id))
		log.Error(err.Error())
		return "", err
	}
	return toJSON(block.DataMetadata)
}

// GetPhotoKey calls core GetPhotoKey
func (m *Mobile) GetPhotoKey(id string) (string, error) {
	key, err := core.Node.GetPhotoKey(id)
	if err != nil {
		return "", err
	}
	return base58.FastBase58Encoding(key), nil
}

// PhotoThreads call core PhotoThreads
func (m *Mobile) PhotoThreads(id string) (string, error) {
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range core.Node.PhotoThreads(id) {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

// ignoreBlock adds an ignore block targeted at the given block and unpins any associated block data
func (m *Mobile) ignoreBlock(blockId string) (string, error) {
	block, err := core.Node.GetBlock(blockId)
	if err != nil {
		return "", err
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}
	hash, err := thrd.Ignore(block.Id)
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// getImageDataURLPrefix adds the correct data url prefix to a data url
func getImageDataURLPrefix(format images.Format) string {
	switch format {
	case images.PNG:
		return "data:image/png;base64,"
	case images.GIF:
		return "data:image/gif;base64,"
	default:
		return "data:image/jpeg;base64,"
	}
}

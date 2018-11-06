package mobile

import (
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
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	thrd := core.Node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}
	keyb, err := base58.Decode(key)
	if err != nil {
		return "", err
	}
	hash, err := thrd.AddFile(dataId, caption, keyb)
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// SharePhoto adds an existing photo to a new thread
func (m *Mobile) SharePhotoToThread(dataId string, threadId string, caption string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	block, err := core.Node.BlockByDataId(dataId)
	if err != nil {
		return "", err
	}
	toThread := core.Node.Thread(threadId)
	if toThread == nil {
		return "", core.ErrThreadNotFound
	}
	// TODO: owner challenge
	hash, err := toThread.AddFile(dataId, caption, block.DataKey)
	if err != nil {
		return "", err
	}
	return hash.B58String(), nil
}

// Photos returns thread photo blocks with json encoding
func (m *Mobile) Photos(offset string, limit int, threadId string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	var pre, query string
	if threadId != "" {
		thrd := core.Node.Thread(threadId)
		if thrd == nil {
			return "", core.ErrThreadNotFound
		}
		pre = fmt.Sprintf("threadId='%s' and ", threadId)
	}
	query = fmt.Sprintf("%stype=%d", pre, repo.FileBlock)

	// build json
	photos := &Photos{Items: make([]Photo, 0)}
	for _, b := range core.Node.Blocks(offset, limit, query) {
		item := Photo{
			Id:       b.DataId,
			BlockId:  b.Id,
			Date:     b.Date,
			AuthorId: b.AuthorId,
			Caption:  b.DataCaption,
			Username: core.Node.ContactUsername(b.AuthorId),
			Metadata: b.DataMetadata,
		}

		// add comments
		cquery := fmt.Sprintf("%stype=%d and dataId='%s'", pre, repo.CommentBlock, b.Id)
		item.Comments = make([]Comment, 0)
		for _, c := range core.Node.Blocks("", -1, cquery) {
			comment := Comment{
				Annotation: Annotation{
					Id:       c.Id,
					Date:     c.Date,
					AuthorId: c.AuthorId,
					Username: core.Node.ContactUsername(c.AuthorId),
				},
				Body: c.DataCaption,
			}
			item.Comments = append(item.Comments, comment)
		}

		// add likes
		lquery := fmt.Sprintf("%stype=%d and dataId='%s'", pre, repo.LikeBlock, b.Id)
		item.Likes = make([]Like, 0)
		for _, l := range core.Node.Blocks("", -1, lquery) {
			like := Like{
				Annotation: Annotation{
					Id:       l.Id,
					Date:     l.Date,
					AuthorId: l.AuthorId,
					Username: core.Node.ContactUsername(l.AuthorId),
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
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	block, err := core.Node.Block(blockId)
	if err != nil {
		return "", err
	}
	thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
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
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	block, err := core.Node.Block(blockId)
	if err != nil {
		return "", err
	}
	thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
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

// PhotoData returns a data url of an image under a path
func (m *Mobile) PhotoData(id string, path string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	block, err := core.Node.BlockByDataId(id)
	if err != nil {
		return "", err
	}
	data, err := core.Node.BlockData(fmt.Sprintf("%s/%s", id, path), block)
	if err != nil {
		return "", err
	}
	format := block.DataMetadata.EncodingFormat
	prefix := getImageDataURLPrefix(images.Format(format))
	encoded := libp2pc.ConfigEncodeKey(data)
	img := &ImageData{Url: prefix + encoded}
	return toJSON(img)
}

// PhotoDataForSize returns a data url of an image at or above requested size, or the next best option
func (m *Mobile) PhotoDataForMinWidth(id string, minWidth int) (string, error) {
	path := images.ImagePathForSize(images.ImageSizeForMinWidth(minWidth))
	return m.PhotoData(id, string(path))
}

// PhotoMetadata returns a meta data object for a photo
func (m *Mobile) PhotoMetadata(id string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	block, err := core.Node.BlockByDataId(id)
	if err != nil {
		return "", err
	}
	return toJSON(block.DataMetadata)
}

// PhotoKey calls core PhotoKey
func (m *Mobile) PhotoKey(id string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	key, err := core.Node.PhotoKey(id)
	if err != nil {
		return "", err
	}
	return base58.FastBase58Encoding(key), nil
}

// PhotoThreads call core PhotoThreads
func (m *Mobile) PhotoThreads(id string) (string, error) {
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
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
	if !core.Node.Started() {
		return "", core.ErrStopped
	}
	block, err := core.Node.Block(blockId)
	if err != nil {
		return "", err
	}
	thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
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

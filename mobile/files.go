package mobile

import (
	"encoding/json"
	"errors"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"io"
	"io/ioutil"
	"os"

	"github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"

	"github.com/textileio/textile-go/core"
)

// FileData is a wrapper around a file data url
type FileData struct {
	Url string `json:"url"`
}

// AddSchema adds a new schema via schema mill
func (m *Mobile) AddSchema(jsonstr string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}
	added, err := m.addSchema(jsonstr)
	if err != nil {
		return "", err
	}

	return toJSON(added)
}

// PrepareFile processes a file by path for a thread, but does NOT share it
func (m *Mobile) PrepareFile(path string, threadId string) (string, error) {
	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	if thrd.Schema == nil {
		return "", core.ErrThreadSchemaRequired
	}

	var result interface{}

	mil, err := getMill(thrd.Schema.Mill, thrd.Schema.Opts)
	if err != nil {
		return "", err
	}
	if mil != nil {
		conf, err := m.getFileConfig(mil, path, "")
		if err != nil {
			return "", err
		}

		added, err := m.node.AddFile(mil, *conf)
		if err != nil {
			return "", err
		}
		result = &added

	} else if len(thrd.Schema.Links) > 0 {
		dir := make(map[string]*repo.File)

		// determine order
		steps, err := schema.Steps(thrd.Schema.Links)
		if err != nil {
			return "", err
		}

		// send each link
		for _, step := range steps {
			mil, err := getMill(step.Link.Mill, step.Link.Opts)
			if err != nil {
				return "", err
			}
			var conf *core.AddFileConfig

			if step.Link.Use == schema.FileTag {
				conf, err = m.getFileConfig(mil, path, "")
				if err != nil {
					return "", err
				}

			} else {
				if dir[step.Link.Use] == nil {
					return "", errors.New(step.Link.Use + " not found")
				}
				conf, err = m.getFileConfig(mil, path, dir[step.Link.Use].Hash)
				if err != nil {
					return "", err
				}
			}
			added, err := m.node.AddFile(mil, *conf)
			if err != nil {
				return "", err
			}
			dir[step.Name] = added
		}
		result = &dir

	} else {
		return "", schema.ErrEmptySchema
	}

	return toJSON(result)
}

// AddFile adds a prepared file to a thread
func (m *Mobile) AddFile(jsonstr string, threadId string, caption string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	var node ipld.Node
	var keys core.Keys

	// parse file or directory
	var dir core.Directory
	if err := json.Unmarshal([]byte(jsonstr), &dir); err != nil {
		return "", err
	}
	var err error
	if len(dir) > 0 {
		node, keys, err = m.node.AddNodeFromDirs([]core.Directory{dir})
		if err != nil {
			return "", err
		}
	} else {
		var file repo.File
		if err := json.Unmarshal([]byte(jsonstr), &file); err != nil {
			return "", err
		}
		node, keys, err = m.node.AddNodeFromFiles([]repo.File{file})
		if err != nil {
			return "", err
		}
	}

	if node == nil {
		return "", errors.New("no files found")
	}

	hash, err := thrd.AddFiles(node, caption, keys)
	if err != nil {
		return "", err
	}

	return m.blockInfo(hash)
}

// AddFileByTarget adds a prepared file to a thread by referencing its top level hash,
// which is the target of an existing files block.
func (m *Mobile) AddFileByTarget(target string, threadId string, caption string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}

	block, err := m.node.BlockByTarget(target)
	if err != nil {
		return "", err
	}

	fsinfo, err := m.node.File(threadId, block.Id)
	if err != nil {
		return "", err
	}

	var dirs []core.Directory
	var files []repo.File

	for _, info := range fsinfo.Files {
		if len(info.Links) > 0 {
			dirs = append(dirs, info.Links)
		} else if info.File != nil {
			files = append(files, *info.File)
		}
	}

	var node ipld.Node
	var keys core.Keys

	if len(dirs) > 0 {
		node, keys, err = m.node.AddNodeFromDirs(dirs)
		if err != nil {
			return "", err
		}
	} else if len(files) > 0 {
		node, keys, err = m.node.AddNodeFromFiles(files)
		if err != nil {
			return "", err
		}
	}

	if node == nil {
		return "", errors.New("no files found")
	}

	hash, err := thrd.AddFiles(node, caption, keys)
	if err != nil {
		return "", err
	}

	return m.blockInfo(hash)
}

// Files calls core Files
func (m *Mobile) Files(threadId string, offset string, limit int) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	files, err := m.node.Files(threadId, offset, limit)
	if err != nil {
		return "", err
	}

	return toJSON(files)
}

func (m *Mobile) addSchema(jsonstr string) (*repo.File, error) {
	var node schema.Node
	if err := json.Unmarshal([]byte(jsonstr), &node); err != nil {
		return nil, err
	}
	data, err := json.Marshal(&node)
	if err != nil {
		return nil, err
	}

	conf := core.AddFileConfig{
		Input: data,
		Media: "application/json",
	}

	return m.node.AddFile(&mill.Schema{}, conf)
}

func (m *Mobile) getFileConfig(mil mill.Mill, path string, use string) (*core.AddFileConfig, error) {
	var reader io.ReadSeeker
	conf := &core.AddFileConfig{}

	if use == "" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
		conf.Name = f.Name()
	} else {
		var file *repo.File
		var err error
		reader, file, err = m.node.FilePlaintext(use)
		if err != nil {
			return nil, err
		}
		conf.Name = file.Name
		conf.Use = file.Checksum
	}

	media, err := m.node.GetMedia(reader, mil)
	if err != nil {
		return nil, err
	}
	conf.Media = media
	reader.Seek(0, 0)

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf.Input = data

	return conf, nil
}

func getMill(id string, opts map[string]string) (mill.Mill, error) {
	switch id {
	case "/blob":
		return &mill.Blob{}, nil
	case "/image/resize":
		width := opts["width"]
		if width == "" {
			return nil, errors.New("missing width")
		}
		quality := opts["quality"]
		if quality == "" {
			quality = "75"
		}
		return &mill.ImageResize{
			Opts: mill.ImageResizeOpts{
				Width:   width,
				Quality: quality,
			},
		}, nil
	case "/image/exif":
		return &mill.ImageExif{}, nil
	default:
		return nil, nil
	}
}

//// FileData returns a data url of a file under a path
//func (m *Mobile) PhotoData(id string, path string) (string, error) {
//	if !m.node.Started() {
//		return "", core.ErrStopped
//	}
//	block, err := m.node.BlockByDataId(id)
//	if err != nil {
//		return "", err
//	}
//	data, err := m.node.BlockData(fmt.Sprintf("%s/%s", id, path), block)
//	if err != nil {
//		return "", err
//	}
//	format := block.DataMetadata.EncodingFormat
//	prefix := getImageDataURLPrefix(images.Format(format))
//	encoded := libp2pc.ConfigEncodeKey(data)
//	img := &ImageData{Url: prefix + encoded}
//	return toJSON(img)
//}

//// FileDataForMinWidth returns a data url of an image at or above requested size, or the next best option
//func (m *Mobile) FileDataForMinWidth(id string, minWidth int) (string, error) {
//	path := images.ImagePathForSize(images.ImageSizeForMinWidth(minWidth))
//	return m.PhotoData(id, string(path))
//}

//// FileMetadata returns meta data object for a photo
//func (m *Mobile) PhotoMetadata(id string) (string, error) {
//	if !m.node.Started() {
//		return "", core.ErrStopped
//	}
//	block, err := m.node.BlockByDataId(id)
//	if err != nil {
//		return "", err
//	}
//	return toJSON(block.DataMetadata)
//}

//// FileThreads call core PhotoThreads
//func (m *Mobile) PhotoThreads(id string) (string, error) {
//	if !m.node.Started() {
//		return "", core.ErrStopped
//	}
//	threads := Threads{Items: make([]Thread, 0)}
//	for _, thrd := range m.node.PhotoThreads(id) {
//		peers := thrd.Peers()
//		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
//		threads.Items = append(threads.Items, item)
//	}
//	return toJSON(threads)
//}

// getImageDataURLPrefix adds the correct data url prefix to a data url
//func getImageDataURLPrefix(format images.Format) string {
//	switch format {
//	case images.PNG:
//		return "data:image/png;base64,"
//	case images.GIF:
//		return "data:image/gif;base64,"
//	default:
//		return "data:image/jpeg;base64,"
//	}
//}

//// FileThreads lists threads which contain a photo (known to the local peer)
//func (t *Textile) FileThreads(id string) []*Thread {
//	blocks := t.datastore.Blocks().List("", -1, "dataId='"+id+"'")
//	if len(blocks) == 0 {
//		return nil
//	}
//	var threads []*Thread
//	for _, block := range blocks {
//		if thrd := t.Thread(block.ThreadId); thrd != nil {
//			threads = append(threads, thrd)
//		}
//	}
//	return threads
//}

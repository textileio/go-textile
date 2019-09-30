package mobile

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/textileio/go-textile/util"

	"github.com/golang/protobuf/proto"
	ipld "github.com/ipfs/go-ipld-format"
	ipfspath "github.com/ipfs/go-path"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema"
)

var fileConfigOpt fileConfigOption

type fileConfigSettings struct {
	Data      []byte
	Path      string
	Plaintext bool
}

type fileConfigOption func(*fileConfigSettings)

func (fileConfigOption) Data(val []byte) fileConfigOption {
	return func(settings *fileConfigSettings) {
		settings.Data = val
	}
}

func (fileConfigOption) Path(val string) fileConfigOption {
	return func(settings *fileConfigSettings) {
		settings.Path = val
	}
}

func (fileConfigOption) Plaintext(val bool) fileConfigOption {
	return func(settings *fileConfigSettings) {
		settings.Plaintext = val
	}
}

func fileConfigOptions(opts ...fileConfigOption) *fileConfigSettings {
	options := &fileConfigSettings{}

	for _, opt := range opts {
		opt(options)
	}
	return options
}

// AddData adds raw data to a thread
func (m *Mobile) AddData(data string, threadId string, caption string, cb ProtoCallback) {
	m.node.WaitAdd(1, "Mobile.AddData")
	go func() {
		defer m.node.WaitDone("Mobile.AddData")

		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			cb.Call(nil, err)
			return
		}
		hash, err := m.addData(decoded, threadId, caption)
		if err != nil {
			cb.Call(nil, err)
			return
		}

		cb.Call(m.blockView(hash))
	}()
}

// AddFiles builds a directory from paths (comma separated) and adds it to the thread
// Note: paths can be file system paths, IPFS hashes, or an existing file hash that may need decryption.
func (m *Mobile) AddFiles(paths string, threadId string, caption string, cb ProtoCallback) {
	m.node.WaitAdd(1, "Mobile.AddFiles")
	go func() {
		defer m.node.WaitDone("Mobile.AddFiles")

		hash, err := m.addFiles(util.SplitString(paths, ","), threadId, caption)
		if err != nil {
			cb.Call(nil, err)
			return
		}

		cb.Call(m.blockView(hash))
	}()
}

// ShareFiles adds an existing file DAG to a thread via its top level hash (data)
func (m *Mobile) ShareFiles(data string, threadId string, caption string, cb ProtoCallback) {
	m.node.WaitAdd(1, "Mobile.ShareFiles")
	go func() {
		defer m.node.WaitDone("Mobile.ShareFiles")

		hash, err := m.shareFiles(data, threadId, caption)
		if err != nil {
			cb.Call(nil, err)
			return
		}

		cb.Call(m.blockView(hash))
	}()
}

// Files calls core Files
func (m *Mobile) Files(threadId string, offset string, limit int) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	files, err := m.node.Files(offset, limit, threadId)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(files)
}

// File calls core File
func (m *Mobile) File(blockId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	file, err := m.node.File(blockId)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(file)
}

// FileContent is the async version of fileContent
func (m *Mobile) FileContent(hash string, cb DataCallback) {
	m.node.WaitAdd(1, "Mobile.FileContent")
	go func() {
		defer m.node.WaitDone("Mobile.FileContent")
		cb.Call(m.fileContent(hash))
	}()
}

// fileContent returns the data and media type of a raw file under a path
func (m *Mobile) fileContent(hash string) ([]byte, string, error) {
	if !m.node.Started() {
		return nil, "", core.ErrStopped
	}

	reader, file, err := m.node.FileContent(hash)
	if err != nil {
		if err == core.ErrFileNotFound || err == ipld.ErrNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	return data, file.Media, nil
}

type img struct {
	hash  string
	width int
}

// ImageFileContentForMinWidth is the async version of imageFileContentForMinWidth
func (m *Mobile) ImageFileContentForMinWidth(pth string, minWidth int, cb DataCallback) {
	m.node.WaitAdd(1, "Mobile.ImageFileContentForMinWidth")
	go func() {
		defer m.node.WaitDone("Mobile.ImageFileContentForMinWidth")
		cb.Call(m.imageFileContentForMinWidth(pth, minWidth))
	}()
}

// imageFileContentForMinWidth returns a data url of an image at or above requested size,
// or the next best option.
// Note: Now that consumers are in control of image sizes via schemas,
// handling this here doesn't feel right. We can eventually push this up to RN, Obj-C, Java.
// Note: pth is <data>/<index>, e.g., "Qm.../0"
func (m *Mobile) imageFileContentForMinWidth(pth string, minWidth int) ([]byte, string, error) {
	if !m.node.Started() {
		return nil, "", core.ErrStopped
	}

	node, err := ipfs.NodeAtPath(m.node.Ipfs(), pth, ipfs.CatTimeout)
	if err != nil {
		if err == ipld.ErrNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	var imgs []img
	for _, link := range node.Links() {
		nd, err := ipfs.NodeAtLink(m.node.Ipfs(), link)
		if err != nil {
			if err == ipld.ErrNotFound {
				return nil, "", nil
			}
			return nil, "", err
		}

		dlink := schema.LinkByName(nd.Links(), core.ValidContentLinkNames)
		if dlink == nil {
			continue
		}

		file, err := m.node.FileMeta(dlink.Cid.Hash().B58String())
		if err != nil {
			if err == core.ErrFileNotFound {
				return nil, "", nil
			}
			return nil, "", err
		}

		if file.Mill == "/image/resize" {
			width := file.Meta.Fields["width"]
			if width != nil {
				imgs = append(imgs, img{
					hash:  file.Hash,
					width: int(width.GetNumberValue()),
				})
			}
		}
	}

	if len(imgs) == 0 {
		return nil, "", nil
	}

	sort.SliceStable(imgs, func(i, j int) bool {
		return imgs[i].width < imgs[j].width
	})

	var hash string
	for _, img := range imgs {
		if img.width >= minWidth {
			hash = img.hash
			break
		}
	}
	if hash == "" {
		hash = imgs[len(imgs)-1].hash
	}

	return m.fileContent(hash)
}

func (m *Mobile) addData(data []byte, threadId string, caption string) (mh.Multihash, error) {
	dir, err := m.buildDirectory(data, "", threadId)
	if err != nil {
		return nil, err
	}

	return m.writeFiles(&pb.DirectoryList{Items: []*pb.Directory{dir}}, threadId, caption)
}

func (m *Mobile) addFiles(paths []string, threadId string, caption string) (mh.Multihash, error) {
	dirs := &pb.DirectoryList{Items: make([]*pb.Directory, 0)}
	for _, pth := range paths {
		dir, err := m.buildDirectory(nil, pth, threadId)
		if err != nil {
			return nil, err
		}
		dirs.Items = append(dirs.Items, dir)
	}

	return m.writeFiles(dirs, threadId, caption)
}

func (m *Mobile) shareFiles(data string, threadId string, caption string) (mh.Multihash, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return nil, core.ErrThreadNotFound
	}

	node, err := ipfs.NodeAtPath(m.node.Ipfs(), data, ipfs.CatTimeout)
	if err != nil {
		return nil, err
	}

	keys, err := m.node.TargetNodeKeys(node)
	if err != nil {
		return nil, err
	}

	hash, err := thrd.AddFiles(node, "", caption, keys.Files)
	if err != nil {
		return nil, err
	}

	m.node.FlushCafes()

	return hash, nil
}

func (m *Mobile) buildDirectory(data []byte, path string, threadId string) (*pb.Directory, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return nil, core.ErrThreadNotFound
	}

	if thrd.Schema == nil {
		return nil, core.ErrThreadSchemaRequired
	}

	dir := &pb.Directory{
		Files: make(map[string]*pb.FileIndex),
	}

	mil, err := getMill(thrd.Schema.Mill, thrd.Schema.Opts)
	if err != nil {
		return nil, err
	}
	if mil != nil {
		conf, err := m.getFileConfig(mil,
			fileConfigOpt.Data(data),
			fileConfigOpt.Path(path),
			fileConfigOpt.Plaintext(thrd.Schema.Plaintext),
		)
		if err != nil {
			return nil, err
		}

		added, err := m.node.AddFileIndex(mil, *conf)
		if err != nil {
			return nil, err
		}
		dir.Files[schema.SingleFileTag] = added

	} else if len(thrd.Schema.Links) > 0 {

		// determine order
		steps, err := schema.Steps(thrd.Schema.Links)
		if err != nil {
			return nil, err
		}

		// send each link
		for _, step := range steps {
			mil, err := getMill(step.Link.Mill, step.Link.Opts)
			if err != nil {
				return nil, err
			}
			var conf *core.AddFileConfig

			if step.Link.Use == schema.FileTag {
				conf, err = m.getFileConfig(mil,
					fileConfigOpt.Data(data),
					fileConfigOpt.Path(path),
					fileConfigOpt.Plaintext(step.Link.Plaintext),
				)
				if err != nil {
					return nil, err
				}

			} else {
				if dir.Files[step.Link.Use] == nil {
					return nil, fmt.Errorf(step.Link.Use + " not found")
				}

				conf, err = m.getFileConfig(mil,
					fileConfigOpt.Data(data),
					fileConfigOpt.Path(dir.Files[step.Link.Use].Hash),
					fileConfigOpt.Plaintext(step.Link.Plaintext),
				)
				if err != nil {
					return nil, err
				}
			}

			added, err := m.node.AddFileIndex(mil, *conf)
			if err != nil {
				return nil, err
			}
			dir.Files[step.Name] = added
		}
	} else {
		return nil, schema.ErrEmptySchema
	}

	return dir, nil
}

func (m *Mobile) getFileConfig(mil mill.Mill, opts ...fileConfigOption) (*core.AddFileConfig, error) {
	var reader io.ReadSeeker
	conf := &core.AddFileConfig{}
	settings := fileConfigOptions(opts...)

	if settings.Data != nil {
		reader = bytes.NewReader(settings.Data)
	} else {
		ref, err := ipfspath.ParsePath(settings.Path)
		if err == nil {
			parts := strings.Split(ref.String(), "/")
			hash := parts[len(parts)-1]
			var file *pb.FileIndex
			reader, file, err = m.node.FileContent(hash)
			if err != nil {
				if err == core.ErrFileNotFound {
					// just cat the data from ipfs
					b, err := ipfs.DataAtPath(m.node.Ipfs(), ref.String())
					if err != nil {
						return nil, err
					}
					reader = bytes.NewReader(b)
					conf.Use = ref.String()
				} else {
					return nil, err
				}
			} else {
				conf.Use = file.Checksum
			}
		} else { // lastly, try and open as an os file
			f, err := os.Open(settings.Path)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			reader = f
			_, conf.Name = filepath.Split(f.Name())
		}
	}

	var err error
	if mil.ID() == "/json" {
		conf.Media = "application/json"
	} else {
		conf.Media, err = m.node.GetMillMedia(reader, mil)
		if err != nil {
			return nil, err
		}
	}
	_, _ = reader.Seek(0, 0)

	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf.Input = input
	conf.Plaintext = settings.Plaintext

	return conf, nil
}

func getMill(id string, opts map[string]string) (mill.Mill, error) {
	switch id {
	case "/blob":
		return &mill.Blob{}, nil
	case "/image/resize":
		width := opts["width"]
		if width == "" {
			return nil, fmt.Errorf("missing width")
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
	case "/json":
		return &mill.Json{}, nil
	default:
		return nil, nil
	}
}

func (m *Mobile) writeFiles(dirs *pb.DirectoryList, threadId string, caption string) (mh.Multihash, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	if len(dirs.Items) == 0 || len(dirs.Items[0].Files) == 0 {
		return nil, fmt.Errorf("no files found")
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return nil, core.ErrThreadNotFound
	}

	var node ipld.Node
	var keys *pb.Keys

	var err error
	file := dirs.Items[0].Files[schema.SingleFileTag]
	if file != nil {
		node, keys, err = m.node.AddNodeFromFiles([]*pb.FileIndex{file})
	} else {
		node, keys, err = m.node.AddNodeFromDirs(dirs)
	}
	if err != nil {
		return nil, err
	}

	if node == nil {
		return nil, fmt.Errorf("no files found")
	}

	hash, err := thrd.AddFiles(node, "", caption, keys.Files)
	if err != nil {
		return nil, err
	}

	m.node.FlushCafes()

	return hash, nil
}

package mobile

import (
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core/coreapi/interface"

	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/schema"
)

// AddSchema adds a new schema via schema mill
func (m *Mobile) AddSchema(jsonstr string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	added, err := m.node.AddFileIndex(&mill.Schema{}, core.AddFileConfig{
		Input: []byte(jsonstr),
		Media: "application/json",
	})
	if err != nil {
		return nil, err
	}

	return proto.Marshal(added)
}

// PrepareFiles processes a file by path for a thread, but does NOT share it
func (m *Mobile) PrepareFiles(path string, threadId string) ([]byte, error) {
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

	var use string
	if ref, err := iface.ParsePath(path); err == nil {
		parts := strings.Split(ref.String(), "/")
		use = parts[len(parts)-1]
	}

	mdir := &pb.MobilePreparedFiles{
		Dir: &pb.Directory{
			Files: make(map[string]*pb.FileIndex),
		},
		Pin: make(map[string]string),
	}

	writeDir := m.RepoPath + "/tmp/"

	mil, err := getMill(thrd.Schema.Mill, thrd.Schema.Opts)
	if err != nil {
		return nil, err
	}
	if mil != nil {
		conf, err := m.getFileConfig(mil, path, use, thrd.Schema.Plaintext)
		if err != nil {
			return nil, err
		}

		added, err := m.node.AddFileIndex(mil, *conf)
		if err != nil {
			return nil, err
		}
		mdir.Dir.Files[schema.SingleFileTag] = added

		if added.Size >= int64(m.node.Config().Cafe.Client.Mobile.P2PWireLimit) {
			mdir.Pin[added.Hash] = writeDir + added.Hash
		}

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
				conf, err = m.getFileConfig(mil, path, use, step.Link.Plaintext)
				if err != nil {
					return nil, err
				}

			} else {
				if mdir.Dir.Files[step.Link.Use] == nil {
					return nil, errors.New(step.Link.Use + " not found")
				}

				conf, err = m.getFileConfig(mil, path, mdir.Dir.Files[step.Link.Use].Hash, step.Link.Plaintext)
				if err != nil {
					return nil, err
				}
			}

			added, err := m.node.AddFileIndex(mil, *conf)
			if err != nil {
				return nil, err
			}
			mdir.Dir.Files[step.Name] = added

			if added.Size >= int64(m.node.Config().Cafe.Client.Mobile.P2PWireLimit) {
				mdir.Pin[added.Hash] = writeDir + added.Hash
			}
		}
	} else {
		return nil, schema.ErrEmptySchema
	}

	for hash, pth := range mdir.Pin {
		if err := m.writeFileData(hash, pth); err != nil {
			return nil, err
		}
	}

	return proto.Marshal(mdir)
}

// PrepareFilesAsync is the async flavor of PrepareFiles
func (m *Mobile) PrepareFilesAsync(path string, threadId string, cb Callback) {
	go func() {
		cb.Call(m.PrepareFiles(path, threadId))
	}()
}

// AddFiles adds a prepared file to a thread
func (m *Mobile) AddFiles(dir []byte, threadId string, caption string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return nil, core.ErrThreadNotFound
	}

	var node ipld.Node
	var keys *pb.Keys

	mdir := new(pb.Directory)
	if err := proto.Unmarshal(dir, mdir); err != nil {
		return nil, err
	}
	if len(mdir.Files) == 0 {
		return nil, errors.New("no files found")
	}

	var err error
	file := mdir.Files[schema.SingleFileTag]
	if file != nil {

		node, keys, err = m.node.AddNodeFromFiles([]*pb.FileIndex{file})
		if err != nil {
			return nil, err
		}

	} else {

		rdir := &pb.Directory{Files: make(map[string]*pb.FileIndex)}
		for k, file := range mdir.Files {
			rdir.Files[k] = file
		}

		node, keys, err = m.node.AddNodeFromDirs(&pb.DirectoryList{Items: []*pb.Directory{rdir}})
		if err != nil {
			return nil, err
		}
	}

	if node == nil {
		return nil, errors.New("no files found")
	}

	hash, err := thrd.AddFiles(node, caption, keys.Files)
	if err != nil {
		return nil, err
	}

	return m.blockView(hash)
}

// AddFilesByTarget adds a prepared file to a thread by referencing its top level hash
func (m *Mobile) AddFilesByTarget(target string, threadId string, caption string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return nil, core.ErrThreadNotFound
	}

	node, err := ipfs.NodeAtPath(m.node.Ipfs(), target)
	if err != nil {
		return nil, err
	}

	keys, err := m.node.TargetNodeKeys(node)
	if err != nil {
		return nil, err
	}

	hash, err := thrd.AddFiles(node, caption, keys.Files)
	if err != nil {
		return nil, err
	}

	return m.blockView(hash)
}

// Files calls core Files
func (m *Mobile) Files(offset string, limit int, threadId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	files, err := m.node.Files(offset, limit, threadId)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(files)
}

// FileData returns a data url of a raw file under a path
func (m *Mobile) FileData(hash string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	reader, file, err := m.node.FileData(hash)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	prefix := "data:" + file.Media + ";base64,"
	img := &pb.MobileFileData{
		Url: prefix + base64.StdEncoding.EncodeToString(data),
	}

	return proto.Marshal(img)
}

type img struct {
	hash  string
	width int
}

// ImageFileDataForMinWidth returns a data url of an image at or above requested size,
// or the next best option.
// Note: Now that consumers are in control of image sizes via schemas,
// handling this here doesn't feel right. We can eventually push this up to RN, Obj-C, Java.
// Note: pth is <target>/<index>, e.g., "Qm.../0"
func (m *Mobile) ImageFileDataForMinWidth(pth string, minWidth int) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	node, err := ipfs.NodeAtPath(m.node.Ipfs(), pth)
	if err != nil {
		return nil, err
	}

	var imgs []img
	for _, link := range node.Links() {
		nd, err := ipfs.NodeAtLink(m.node.Ipfs(), link)
		if err != nil {
			return nil, err
		}

		dlink := schema.LinkByName(nd.Links(), core.DataLinkName)
		if dlink == nil {
			continue
		}

		file, err := m.node.FileIndex(dlink.Cid.Hash().B58String())
		if err != nil {
			return nil, err
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
		return nil, errors.New("no image files found")
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

	return m.FileData(hash)
}

func (m *Mobile) getFileConfig(mil mill.Mill, path string, use string, plaintext bool) (*core.AddFileConfig, error) {
	var reader io.ReadSeeker
	conf := &core.AddFileConfig{}

	if use == "" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f

		_, file := filepath.Split(f.Name())
		conf.Name = file

	} else {
		var file *pb.FileIndex
		var err error
		reader, file, err = m.node.FileData(use)
		if err != nil {
			return nil, err
		}

		conf.Name = file.Name
		conf.Use = file.Checksum
	}

	var err error
	if mil.ID() == "/json" {
		conf.Media = "application/json"
	} else {
		conf.Media, err = m.node.GetMedia(reader, mil)
		if err != nil {
			return nil, err
		}
	}
	reader.Seek(0, 0)

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf.Input = data
	conf.Plaintext = plaintext

	return conf, nil
}

func (m *Mobile) writeFileData(hash string, pth string) error {
	if err := os.MkdirAll(filepath.Dir(pth), os.ModePerm); err != nil {
		return err
	}

	data, err := m.node.DataAtPath(hash)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pth, data, 0644)
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
	case "/json":
		return &mill.Json{}, nil
	default:
		return nil, nil
	}
}

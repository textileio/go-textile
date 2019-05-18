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

	"github.com/golang/protobuf/proto"
	ipld "github.com/ipfs/go-ipld-format"
	ipfspath "github.com/ipfs/go-path"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema"
)

// PrepareFiles processes base64 encoded data for a thread, but does NOT share it
func (m *Mobile) PrepareFiles(data string, threadId string, cb Callback) {
	go func() {
		cb.Call(m.PrepareFilesSync(data, threadId))
	}()
}

// PrepareFilesByPath processes a file by path for a thread, but does NOT share it
func (m *Mobile) PrepareFilesByPath(path string, threadId string, cb Callback) {
	go func() {
		cb.Call(m.PrepareFilesByPathSync(path, threadId))
	}()
}

// PrepareFiles processes base64 encoded data for a thread, but does NOT share it
func (m *Mobile) PrepareFilesSync(data string, threadId string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	dec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	thrd := m.node.Thread(threadId)
	if thrd == nil {
		return nil, core.ErrThreadNotFound
	}

	if thrd.Schema == nil {
		return nil, core.ErrThreadSchemaRequired
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
		conf, err := m.getFileConfig(mil, dec, "", thrd.Schema.Plaintext)
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
				conf, err = m.getFileConfig(mil, dec, "", step.Link.Plaintext)
				if err != nil {
					return nil, err
				}

			} else {
				if mdir.Dir.Files[step.Link.Use] == nil {
					return nil, fmt.Errorf(step.Link.Use + " not found")
				}

				conf, err = m.getFileConfig(mil, dec, mdir.Dir.Files[step.Link.Use].Hash, step.Link.Plaintext)
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
		if err := m.writeFileContent(hash, pth); err != nil {
			return nil, err
		}
	}

	return proto.Marshal(mdir)
}

// PrepareFilesByPath processes a file by path for a thread, but does NOT share it
func (m *Mobile) PrepareFilesByPathSync(path string, threadId string) ([]byte, error) {
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
	if ref, err := ipfspath.ParsePath(path); err == nil {
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
		conf, err := m.getFileConfigByPath(mil, path, use, thrd.Schema.Plaintext)
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
				conf, err = m.getFileConfigByPath(mil, path, use, step.Link.Plaintext)
				if err != nil {
					return nil, err
				}

			} else {
				if mdir.Dir.Files[step.Link.Use] == nil {
					return nil, fmt.Errorf(step.Link.Use + " not found")
				}

				conf, err = m.getFileConfigByPath(mil, path, mdir.Dir.Files[step.Link.Use].Hash, step.Link.Plaintext)
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
		if err := m.writeFileContent(hash, pth); err != nil {
			return nil, err
		}
	}

	return proto.Marshal(mdir)
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
		return nil, fmt.Errorf("no files found")
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
		return nil, fmt.Errorf("no files found")
	}

	hash, err := thrd.AddFiles(node, caption, keys.Files)
	if err != nil {
		return nil, err
	}

	if thrd.Key == "account" {
		if err := m.node.SetAvatar(); err != nil {
			return nil, err
		}
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

	if thrd.Key == "account" {
		if err := m.node.SetAvatar(); err != nil {
			return nil, err
		}
	}

	return m.blockView(hash)
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

// FileContent returns a data url of a raw file under a path
func (m *Mobile) FileContent(hash string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	reader, file, err := m.node.FileContent(hash)
	if err != nil {
		if err == core.ErrFileNotFound || err == ipld.ErrNotFound {
			return "", nil
		}
		return "", err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	prefix := "data:" + file.Media + ";base64,"
	return prefix + base64.StdEncoding.EncodeToString(data), nil
}

type img struct {
	hash  string
	width int
}

// ImageFileContentForMinWidth returns a data url of an image at or above requested size,
// or the next best option.
// Note: Now that consumers are in control of image sizes via schemas,
// handling this here doesn't feel right. We can eventually push this up to RN, Obj-C, Java.
// Note: pth is <target>/<index>, e.g., "Qm.../0"
func (m *Mobile) ImageFileContentForMinWidth(pth string, minWidth int) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	node, err := ipfs.NodeAtPath(m.node.Ipfs(), pth)
	if err != nil {
		if err == ipld.ErrNotFound {
			return "", nil
		}
		return "", err
	}

	var imgs []img
	for _, link := range node.Links() {
		nd, err := ipfs.NodeAtLink(m.node.Ipfs(), link)
		if err != nil {
			if err == ipld.ErrNotFound {
				return "", nil
			}
			return "", err
		}

		dlink := schema.LinkByName(nd.Links(), core.ValidContentLinkNames)
		if dlink == nil {
			continue
		}

		file, err := m.node.FileIndex(dlink.Cid.Hash().B58String())
		if err != nil {
			if err == core.ErrFileNotFound {
				return "", nil
			}
			return "", err
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
		return "", nil
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

	return m.FileContent(hash)
}

func (m *Mobile) getFileConfig(mil mill.Mill, data []byte, use string, plaintext bool) (*core.AddFileConfig, error) {
	var reader io.ReadSeeker
	conf := &core.AddFileConfig{}

	if use == "" {
		reader = bytes.NewReader(data)
	} else {
		var file *pb.FileIndex
		var err error
		reader, file, err = m.node.FileContent(use)
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

	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf.Input = input
	conf.Plaintext = plaintext

	return conf, nil
}

func (m *Mobile) getFileConfigByPath(mil mill.Mill, path string, use string, plaintext bool) (*core.AddFileConfig, error) {
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
		reader, file, err = m.node.FileContent(use)
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

	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	conf.Input = input
	conf.Plaintext = plaintext

	return conf, nil
}

func (m *Mobile) writeFileContent(hash string, pth string) error {
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

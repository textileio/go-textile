package mobile

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	ipfspath "github.com/ipfs/go-path"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
)

// Profile calls core Profile
func (m *Mobile) Profile() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	self := m.node.Profile()
	if self == nil {
		return nil, fmt.Errorf("profile not found")
	}

	return proto.Marshal(self)
}

// Name calls core Name
func (m *Mobile) Name() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	return m.node.Name(), nil
}

// SetName calls core SetName
func (m *Mobile) SetName(username string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	return m.node.SetName(username)
}

// Avatar calls core Avatar
func (m *Mobile) Avatar() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	return m.node.Avatar(), nil
}

// SetAvatar calls core SetAvatar
func (m *Mobile) SetAvatar(path string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	thrd := m.node.AccountThread()
	if thrd == nil {
		return fmt.Errorf("account thread not found")
	}

	var reader io.ReadSeeker
	var file *pb.FileIndex
	var name, media, use string

	// first, see if this is an existing file
	ref, err := ipfspath.ParsePath(path)
	if err == nil {
		parts := strings.Split(ref.String(), "/")
		hash := parts[len(parts)-1]
		reader, file, err = m.node.FileData(hash)
		if err != nil {
			if err == core.ErrFileNotFound {
				// just cat the data from ipfs
				b, err := ipfs.DataAtPath(m.node.Ipfs(), hash)
				if err != nil {
					return err
				}
				reader = bytes.NewReader(b)
			}
			return err
		} else {
			name = file.Name
			media = file.Media
			use = file.Checksum
		}
	} else { // lastly, try and open as an os file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		reader = f
		_, name = filepath.Split(f.Name())
	}

	if media == "" {
		media, err = tmpGetMedia(reader)
		if err != nil {
			return err
		}
		reader.Seek(0, 0)
	}

	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	large, err := m.node.AddFileIndex(&mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   thrd.Schema.Links["large"].Opts["width"],
			Quality: thrd.Schema.Links["large"].Opts["quality"],
		},
	}, core.AddFileConfig{
		Input:     input,
		Use:       use,
		Media:     media,
		Name:      name,
		Plaintext: thrd.Schema.Links["large"].Plaintext,
	})
	if err != nil {
		return err
	}

	small, err := m.node.AddFileIndex(&mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   thrd.Schema.Links["small"].Opts["width"],
			Quality: thrd.Schema.Links["small"].Opts["quality"],
		},
	}, core.AddFileConfig{
		Input:     input,
		Use:       use,
		Media:     media,
		Name:      name,
		Plaintext: thrd.Schema.Links["small"].Plaintext,
	})
	if err != nil {
		return err
	}

	dir := map[string]*pb.FileIndex{"large": large, "small": small}
	dirs := &pb.DirectoryList{Items: []*pb.Directory{{Files: dir}}}
	node, keys, err := m.node.AddNodeFromDirs(dirs)
	if err != nil {
		return err
	}

	if _, err := thrd.AddFiles(node, "", keys.Files); err != nil {
		return err
	}

	return m.node.SetAvatar()
}

func tmpGetMedia(reader io.Reader) (string, error) {
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	return http.DetectContentType(buffer[:n]), nil
}

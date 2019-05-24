package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	ipld "github.com/ipfs/go-ipld-format"
	uio "github.com/ipfs/go-unixfs/io"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	m "github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/schema"
)

var ErrFileNotFound = fmt.Errorf("file not found")
var ErrMissingMetaLink = fmt.Errorf("meta link not in node")
var ErrMissingContentLink = fmt.Errorf("content link not in node")

const MetaLinkName = "meta"
const ContentLinkName = "content"

var ValidMetaLinkNames = []string{"f", "meta"}
var ValidContentLinkNames = []string{"d", "content"}

type AddFileConfig struct {
	Input     []byte `json:"input"`
	Use       string `json:"use"`
	Media     string `json:"media"`
	Name      string `json:"name"`
	Plaintext bool   `json:"plaintext"`
}

func (t *Textile) AddFileIndex(mill m.Mill, conf AddFileConfig) (*pb.FileIndex, error) {
	var source string
	if conf.Use != "" {
		source = conf.Use
	} else {
		source = t.checksum(conf.Input, conf.Plaintext)
	}

	opts, err := mill.Options(map[string]interface{}{
		"plaintext": conf.Plaintext,
	})
	if err != nil {
		return nil, err
	}

	if efile := t.datastore.Files().GetBySource(mill.ID(), source, opts); efile != nil {
		return efile, nil
	}

	res, err := mill.Mill(conf.Input, conf.Name)
	if err != nil {
		return nil, err
	}

	check := t.checksum(res.File, conf.Plaintext)
	if efile := t.datastore.Files().GetByPrimary(mill.ID(), check); efile != nil {
		return efile, nil
	}

	model := &pb.FileIndex{
		Mill:     mill.ID(),
		Checksum: check,
		Source:   source,
		Opts:     opts,
		Media:    conf.Media,
		Name:     conf.Name,
		Size:     int64(len(res.File)),
		Added:    ptypes.TimestampNow(),
		Meta:     pb.ToStruct(res.Meta),
	}

	var reader *bytes.Reader
	if mill.Encrypt() && !conf.Plaintext {
		key, err := crypto.GenerateAESKey()
		if err != nil {
			return nil, err
		}
		ciphertext, err := crypto.EncryptAES(res.File, key)
		if err != nil {
			return nil, err
		}
		model.Key = base58.FastBase58Encoding(key)
		reader = bytes.NewReader(ciphertext)
	} else {
		reader = bytes.NewReader(res.File)
	}

	hash, err := ipfs.AddData(t.node, reader, mill.Pin())
	if err != nil {
		return nil, err
	}
	model.Hash = hash.Hash().B58String()

	if err := t.datastore.Files().Add(model); err != nil {
		if db.ConflictError(err) {
			// we may have lost the race
			return t.datastore.Files().Get(model.Hash), nil
		}
		return nil, err
	}

	// Return the model fetched from the datastore to ensure
	// consistent date formatting and therefore consistent
	// directory hashes.
	return t.datastore.Files().Get(model.Hash), nil
}

func (t *Textile) GetMedia(reader io.Reader, mill m.Mill) (string, error) {
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	media := http.DetectContentType(buffer[:n])

	return media, mill.AcceptMedia(media)
}

func (t *Textile) AddSchema(jsonstr string, name string) (*pb.FileIndex, error) {
	var node pb.Node
	if err := jsonpb.UnmarshalString(jsonstr, &node); err != nil {
		return nil, err
	}

	data, err := pbMarshaler.MarshalToString(&node)
	if err != nil {
		return nil, err
	}

	return t.AddFileIndex(&m.Schema{}, AddFileConfig{
		Input: []byte(data),
		Media: "application/json",
		Name:  name,
	})
}

func (t *Textile) AddNodeFromFiles(files []*pb.FileIndex) (ipld.Node, *pb.Keys, error) {
	keys := &pb.Keys{Files: make(map[string]string)}
	outer := uio.NewDirectory(t.node.DAG)

	for i, file := range files {
		link := strconv.Itoa(i)
		if err := t.fileNode(file, outer, link); err != nil {
			return nil, nil, err
		}
		keys.Files["/"+link+"/"] = file.Key
	}

	node, err := outer.GetNode()
	if err != nil {
		return nil, nil, err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return nil, nil, err
	}
	return node, keys, nil
}

func (t *Textile) AddNodeFromDirs(dirs *pb.DirectoryList) (ipld.Node, *pb.Keys, error) {
	keys := &pb.Keys{Files: make(map[string]string)}
	outer := uio.NewDirectory(t.node.DAG)

	for i, dir := range dirs.Items {
		inner := uio.NewDirectory(t.node.DAG)
		olink := strconv.Itoa(i)

		for link, file := range dir.Files {
			if err := t.fileNode(file, inner, link); err != nil {
				return nil, nil, err
			}
			keys.Files["/"+olink+"/"+link+"/"] = file.Key
		}

		node, err := inner.GetNode()
		if err != nil {
			return nil, nil, err
		}
		if err := ipfs.PinNode(t.node, node, false); err != nil {
			return nil, nil, err
		}

		id := node.Cid().Hash().B58String()
		if err := ipfs.AddLinkToDirectory(t.node, outer, olink, id); err != nil {
			return nil, nil, err
		}
	}

	node, err := outer.GetNode()
	if err != nil {
		return nil, nil, err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return nil, nil, err
	}
	return node, keys, nil
}

func (t *Textile) FileIndex(hash string) (*pb.FileIndex, error) {
	file := t.datastore.Files().Get(hash)
	if file == nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

func (t *Textile) FileMeta(hash string) (*pb.FileIndex, error) {
	file := t.datastore.Files().Get(hash)
	if file == nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

func (t *Textile) FileContent(hash string) (io.ReadSeeker, *pb.FileIndex, error) {
	var err error
	var file *pb.FileIndex
	var reader io.ReadSeeker
	file, err = t.FileMeta(hash)
	if err != nil {
		return nil, nil, err
	}
	reader, err = t.FileIndexContent(file)
	return reader, file, err
}

func (t *Textile)  FileIndexContent(file  *pb.FileIndex) (io.ReadSeeker, error) {
	fd, err := ipfs.DataAtPath(t.node, file.Hash)
	if err != nil {
		return nil, err
	}

	var plaintext []byte
	if file.Key != "" {
		key, err := base58.Decode(file.Key)
		if err != nil {
			return nil, err
		}
		plaintext, err = crypto.DecryptAES(fd, key)
		if err != nil {
			return nil, err
		}
	} else {
		plaintext = fd
	}

	return bytes.NewReader(plaintext), nil
}

func (t *Textile) TargetNodeKeys(node ipld.Node) (*pb.Keys, error) {
	keys := &pb.Keys{Files: make(map[string]string)}

	for i, link := range node.Links() {
		fn, err := ipfs.NodeAtLink(t.node, link)
		if err != nil {
			return nil, err
		}
		if err := t.fileNodeKeys(fn, i, &keys.Files); err != nil {
			return nil, err
		}
	}

	return keys, nil
}

func (t *Textile) fileNode(file *pb.FileIndex, dir uio.Directory, link string) error {
	if t.datastore.Files().Get(file.Hash) == nil {
		return ErrFileNotFound
	}

	// remove locally indexed targets
	file.Targets = nil

	plaintext, err := pbMarshaler.MarshalToString(file)
	if err != nil {
		return err
	}

	var reader io.Reader
	if file.Key != "" {
		key, err := base58.Decode(file.Key)
		if err != nil {
			return err
		}

		ciphertext, err := crypto.EncryptAES([]byte(plaintext), key)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(ciphertext)
	} else {
		reader = strings.NewReader(plaintext)
	}

	pair := uio.NewDirectory(t.node.DAG)
	if _, err := ipfs.AddDataToDirectory(t.node, pair, MetaLinkName, reader); err != nil {
		return err
	}

	if err := ipfs.AddLinkToDirectory(t.node, pair, ContentLinkName, file.Hash); err != nil {
		return err
	}

	node, err := pair.GetNode()
	if err != nil {
		return err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return err
	}

	return ipfs.AddLinkToDirectory(t.node, dir, link, node.Cid().Hash().B58String())
}

func (t *Textile) fileIndexForPair(pair ipld.Node) (*pb.FileIndex, error) {
	c, err := ipfs.ResolveLinkByNames(pair, ValidContentLinkNames)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return t.datastore.Files().Get(c.Cid.Hash().B58String()), nil
}

func (t *Textile) checksum(plaintext []byte, willEncrypt bool) string {
	var add int
	if willEncrypt {
		add = 1
	}
	plaintext = append(plaintext, byte(add))
	sum := sha256.Sum256(plaintext)
	return base58.FastBase58Encoding(sum[:])
}

func (t *Textile) fileNodeKeys(node ipld.Node, index int, keys *map[string]string) error {
	vkeys := *keys

	if looksLikeFileNode(node) {
		key, err := t.fileLinkKey(node)
		if err != nil {
			return err
		}

		vkeys["/"+strconv.Itoa(index)+"/"] = key
	} else {
		for _, link := range node.Links() {
			n, err := ipfs.NodeAtLink(t.node, link)
			if err != nil {
				return err
			}

			key, err := t.fileLinkKey(n)
			if err != nil {
				return err
			}

			vkeys["/"+strconv.Itoa(index)+"/"+link.Name+"/"] = key
		}
	}
	keys = &vkeys

	return nil
}

func (t *Textile) fileLinkKey(inode ipld.Node) (string, error) {
	dlink := schema.LinkByName(inode.Links(), ValidContentLinkNames)
	if dlink == nil {
		return "", ErrMissingContentLink
	}

	file := t.datastore.Files().Get(dlink.Cid.Hash().B58String())
	if file == nil {
		return "", ErrFileNotFound
	}
	return file.Key, nil
}

// looksLikeFileNode returns whether or not a node appears to
// be a textile node. It doesn't inspect the actual data.
func looksLikeFileNode(node ipld.Node) bool {
	links := node.Links()
	if len(links) != 2 {
		return false
	}
	if schema.LinkByName(links, ValidMetaLinkNames) == nil ||
		schema.LinkByName(links, ValidContentLinkNames) == nil {
		return false
	}
	return true
}

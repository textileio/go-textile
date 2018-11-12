package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	m "github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"io"
	"net/http"
	"strconv"
	"time"
)

var ErrFileNotFound = errors.New("file not found")

type Directory map[string]repo.File

type Keys map[string]string

func (t *Textile) MediaType(reader io.Reader, mill m.Mill) (string, error) {
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	media := http.DetectContentType(buffer[:n])

	return media, mill.AcceptMedia(media)
}

func (t *Textile) AddFile(input []byte, name string, media string, mill m.Mill) (*repo.File, error) {
	res, err := mill.Mill(input, name)
	if err != nil {
		return nil, err
	}

	check := t.checksum(res.File)
	if efile := t.datastore.Files().GetByPrimary(mill.ID(), check); efile != nil {
		return efile, nil
	}

	model := &repo.File{
		Mill:     mill.ID(),
		Checksum: check,
		Media:    media,
		Size:     len(res.File),
		Added:    time.Now(),
		Meta:     res.Meta,
	}

	var reader *bytes.Reader
	if mill.Encrypt() {
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
		return nil, err
	}

	return model, nil
}

func (t *Textile) AddNodeFromFiles(files []repo.File) (ipld.Node, Keys, error) {
	keys := make(Keys)
	outer := uio.NewDirectory(t.node.DAG)

	for i, file := range files {
		link := strconv.Itoa(i)
		if err := t.FileNode(file, outer, link); err != nil {
			return nil, nil, err
		}
		keys["/"+link] = file.Key
	}

	node, err := outer.GetNode()
	if err != nil {
		return nil, nil, err
	}
	if err := ipfs.PinNode(t.node, node); err != nil {
		return nil, nil, err
	}
	return node, keys, nil
}

func (t *Textile) AddNodeFromDirs(dirs []Directory) (ipld.Node, Keys, error) {
	keys := make(Keys)
	outer := uio.NewDirectory(t.node.DAG)

	for i, dir := range dirs {
		inner := uio.NewDirectory(t.node.DAG)
		olink := strconv.Itoa(i)

		for link, file := range dir {
			if err := t.FileNode(file, inner, link); err != nil {
				return nil, nil, err
			}
			keys["/"+olink+"/"+link] = file.Key
		}

		node, err := inner.GetNode()
		if err != nil {
			return nil, nil, err
		}
		if err := ipfs.PinNode(t.node, node); err != nil {
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
	if err := ipfs.PinNode(t.node, node); err != nil {
		return nil, nil, err
	}
	return node, keys, nil
}

func (t *Textile) FileNode(file repo.File, dir uio.Directory, link string) error {
	if t.datastore.Files().Get(file.Hash) == nil {
		return ErrFileNotFound
	}

	// include encrypted file info as well
	plaintext, err := json.Marshal(&file)
	if err != nil {
		return err
	}
	key, err := base58.Decode(file.Key)
	if err != nil {
		return err
	}
	ciphertext, err := crypto.EncryptAES(plaintext, key)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(ciphertext)

	pair := uio.NewDirectory(t.node.DAG)
	if _, err := ipfs.AddDataToDirectory(t.node, pair, "info", reader); err != nil {
		return err
	}
	if err := ipfs.AddLinkToDirectory(t.node, pair, "data", file.Hash); err != nil {
		return err
	}

	node, err := pair.GetNode()
	if err != nil {
		return err
	}
	if err := ipfs.PinNode(t.node, node); err != nil {
		return err
	}
	return ipfs.AddLinkToDirectory(t.node, dir, link, node.Cid().Hash().B58String())
}

func (t *Textile) checksum(plaintext []byte) string {
	sum := sha256.Sum256(plaintext)
	return base58.FastBase58Encoding(sum[:])
}

//// PhotoThreads lists threads which contain a photo (known to the local peer)
//func (t *Textile) PhotoThreads(id string) []*Thread {
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
//
//// PhotoKey returns the AES key for a photo set
//func (t *Textile) PhotoKey(id string) ([]byte, error) {
//	block, err := t.BlockByDataId(id)
//	if err != nil {
//		return nil, err
//	}
//	return block.DataKey, nil
//}

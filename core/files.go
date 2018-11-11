package core

import (
	"bytes"
	"crypto/sha256"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	m "github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type Directory map[string]repo.File

func (t *Textile) AddFile(file multipart.File, name string, mill m.Mill) (*repo.File, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	media := http.DetectContentType(buffer[:n])

	if err := mill.AcceptMedia(media); err != nil {
		return nil, err
	}

	file.Seek(0, 0)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	check := t.checksum(data)
	if efile := t.datastore.Files().GetByPrimary(mill.ID(), check); efile != nil {
		return efile, nil
	}

	file.Seek(0, 0)
	res, err := mill.Mill(file, name)
	if err != nil {
		return nil, err
	}

	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, err
	}
	ciphertext, err := crypto.EncryptAES(res.File, key)
	if err != nil {
		return nil, err
	}

	hash, err := ipfs.AddData(t.node, bytes.NewReader(ciphertext), false)
	if err != nil {
		return nil, err
	}

	model := &repo.File{
		Mill:     mill.ID(),
		Checksum: check,
		Hash:     hash.Hash().B58String(),
		Key:      base58.FastBase58Encoding(key),
		Media:    media,
		Size:     len(res.File),
		Added:    time.Now(),
		Meta:     res.Meta,
	}
	if err := t.datastore.Files().Add(model); err != nil {
		return nil, err
	}
	return model, nil
}

//func (t *Textile) NewDir() uio.Directory {
//	return uio.NewDirectory(t.node.DAG)
//}
//
//func (t *Textile) LoadDir(hash string) (uio.Directory, error) {
//	id, err := cid.Parse(hash)
//	if err != nil {
//		return nil, err
//	}
//	node, err := ipfs.CidNode(t.node, id)
//	if err != nil {
//		return nil, err
//	}
//	return uio.NewDirectoryFromNode(t.node.DAG, node)
//}
//
//func (t *Textile) AddFileToDir(dir uio.Directory, fileId string, link string) error {
//	file := t.datastore.Files().Get(fileId)
//	if file == nil {
//		return ErrFileNotFound
//	}
//	return ipfs.AddLinkToDirectory(t.node, dir, link, file.Hash)
//}

// AddFile adds a file to ipfs, it is NOT saved to a thread
//func (t *Textile) AddFile(file []byte, mill string) (*repo.File, error) {
//	check := t.checksum(file)
//
//	// check if exists
//	if efile := t.datastore.Files().Get(check); efile != nil {
//		return efile, nil
//	}
//
//	key, err := crypto.GenerateAESKey()
//	if err != nil {
//		return nil, err
//	}
//	ciphertext, err := crypto.EncryptAES(file, key)
//	if err != nil {
//		return nil, err
//	}
//
//	// persist
//	id, err := ipfs.AddData(t.node, bytes.NewReader(ciphertext), false)
//	if err != nil {
//		return nil, err
//	}
//
//	model := &repo.File{
//		Mill:     mill.ID,
//		Checksum: check,
//		Hash:     id.Hash().B58String(),
//		Key:      base58.FastBase58Encoding(key),
//		Media:    media,
//		Size:     len(file),
//		Added:    time.Now(),
//		Meta:     mill.Meta,
//	}
//	if err := t.datastore.Files().Add(model); err != nil {
//		return nil, err
//	}
//	return model, nil
//}

//func (t *Textile) SaveFile(fileId string, caption string) error {
//	file := t.datastore.Files().Get(fileId)
//	if file == nil {
//		return ErrFileNotFound
//	}
//
//	// save to account thread
//	thrd := t.AccountThread()
//	if thrd == nil {
//		return ErrThreadNotFound
//	}
//	thrd.AddFiles(file.Hash, caption, file.Key)
//}

func (t *Textile) checksum(plaintext []byte) string {
	sum := sha256.Sum256(plaintext)
	return base58.FastBase58Encoding(sum[:])
}

//func (d *Directory) Pin() (mh.Multihash, error) {
//	node, err := d.dir.LinkNode()
//	if err != nil {
//		return nil, err
//	}
//
//	// local pin
//	if err := ipfs.PinDirectory(d.node.node, node); err != nil {
//		return nil, err
//	}
//
//	// cafe pins
//	hash := node.Cid().Hash().B58String()
//	if err := d.node.cafeOutbox.Add(hash, repo.CafeStoreRequest); err != nil {
//		return nil, err
//	}
//
//	return node.Cid().Hash(), nil
//}

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

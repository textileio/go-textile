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
	"net/http"
	"time"
)

type Directory map[string]repo.File

func (t *Textile) FileMedia(reader io.Reader, mill m.Mill) (string, error) {
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

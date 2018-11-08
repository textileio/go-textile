package core

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"io/ioutil"
	"mime/multipart"
	"time"
)

// ErrFileExists indicates that a file is already locally indexed
var ErrFileExists = errors.New("file exists")

// ErrFileNotFound indicates that a file is not locally indexed
var ErrFileNotFound = errors.New("file not found")

func (t *Textile) NewDir() uio.Directory {
	return uio.NewDirectory(t.node.DAG)
}

func (t *Textile) AddDirFile(dir uio.Directory, fileId string, link string) error {
	file := t.datastore.Files().Get(fileId)
	if file == nil {
		return ErrFileNotFound
	}
	return ipfs.AddDirectoryLink(t.node, dir, link, file.Hash)
}

// AddFile adds a file
func (t *Textile) AddFile(file multipart.File, schema string, key []byte) (*repo.File, error) {
	// check schema
	schemaHash := t.fileSchemas[schema]
	if schemaHash == nil {
		return nil, ErrFileSchemaNotFound
	}

	// checksum
	plaintext, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	check := t.checksum(plaintext)

	// check if exists
	if t.datastore.Files().Get(check) != nil {
		return nil, ErrFileExists
	}

	// get a key if needed
	if key == nil {
		var err error
		key, err = crypto.GenerateAESKey()
		if err != nil {
			return nil, err
		}
	}

	// encrypt and add to local ipfs
	ciphertext, err := crypto.EncryptAES(plaintext, key)
	if err != nil {
		return nil, err
	}
	id, err := ipfs.AddData(t.node, bytes.NewReader(ciphertext), false)
	if err != nil {
		return nil, err
	}
	hash := id.Hash().B58String()

	// add to local file index
	model := &repo.File{
		Id:     check,
		Hash:   hash,
		Schema: schemaHash.B58String(),
		Key:    base58.FastBase58Encoding(key),
		Added:  time.Now(),
	}
	if err := t.datastore.Files().Add(model); err != nil {
		return nil, err
	}
	return model, nil
}

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
//	node, err := d.dir.GetNode()
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

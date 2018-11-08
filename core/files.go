package core

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"io/ioutil"
	"mime/multipart"
	"time"
)

// ErrFileExists indicates that a file is already locally indexed
var ErrFileExists = errors.New("file exists")

// ErrFileNotFound indicates that a file is not locally indexed
var ErrFileNotFound = errors.New("file not found")

// AddFile adds a file
func (t *Textile) AddFile(file multipart.File, schema string, key []byte) (*repo.File, error) {
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

	// add store requests for cafe pins
	if err := t.cafeOutbox.Add(hash, repo.CafeStoreRequest); err != nil {
		return nil, err
	}

	// add to local file index
	model := &repo.File{
		Id:     check,
		Hash:   hash,
		Schema: schema,
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

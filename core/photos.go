package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/archive"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/photo"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"os"
	"path/filepath"
	"strings"
)

// AddPhoto add a photo to the local ipfs node
func (t *Textile) AddPhoto(path string) (*AddDataResult, error) {
	// get a key to encrypt with
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, err
	}

	// read file from disk
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode image
	reader, format, size, err := photo.DecodeImage(file)
	if err != nil {
		return nil, err
	}
	var encodingFormat photo.Format
	if *format == photo.GIF {
		encodingFormat = photo.GIF
	} else {
		encodingFormat = photo.JPEG
	}

	// make all image sizes
	reader.Seek(0, 0)
	thumb, err := photo.EncodeImage(reader, encodingFormat, photo.ThumbnailSize)
	if err != nil {
		return nil, err
	}
	reader.Seek(0, 0)
	small, err := photo.EncodeImage(reader, encodingFormat, photo.SmallSize)
	if err != nil {
		return nil, err
	}
	reader.Seek(0, 0)
	medium, err := photo.EncodeImage(reader, encodingFormat, photo.MediumSize)
	if err != nil {
		return nil, err
	}
	reader.Seek(0, 0)
	large, err := photo.EncodeImage(reader, encodingFormat, photo.LargeSize)
	if err != nil {
		return nil, err
	}

	// make meta data
	fpath := file.Name()
	ext := strings.ToLower(filepath.Ext(fpath))
	reader.Seek(0, 0)
	meta, err := photo.MakeMetadata(reader, fpath, ext, *format, encodingFormat, size.X, size.Y, t.version)
	if err != nil {
		return nil, err
	}
	metab, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	// get public key
	accnt, err := t.Account()
	if err != nil {
		return nil, err
	}
	addrb := []byte(accnt.Address())

	// encrypt files
	thumbcipher, err := crypto.EncryptAES(thumb, key)
	if err != nil {
		return nil, err
	}
	smallcipher, err := crypto.EncryptAES(small, key)
	if err != nil {
		return nil, err
	}
	mediumcipher, err := crypto.EncryptAES(medium, key)
	if err != nil {
		return nil, err
	}
	largecipher, err := crypto.EncryptAES(large, key)
	if err != nil {
		return nil, err
	}
	metacipher, err := crypto.EncryptAES(metab, key)
	if err != nil {
		return nil, err
	}
	addrcipher, err := crypto.EncryptAES(addrb, key)
	if err != nil {
		return nil, err
	}

	// create a virtual directory for the photo
	dirb := uio.NewDirectory(t.ipfs.DAG)
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader(thumbcipher), "thumb"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader(smallcipher), "small"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader(mediumcipher), "medium"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader(largecipher), "photo"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader(metacipher), "meta"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader(addrcipher), "pk"); err != nil {
		return nil, err
	}

	// pin the directory
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := ipfs.PinDirectory(t.ipfs, dir, []string{"medium", "photo"}); err != nil {
		return nil, err
	}
	result := &AddDataResult{Id: dir.Cid().Hash().B58String(), Key: string(key)}

	// if not mobile, create a pin request
	// on mobile, we let the OS handle the archive directly
	if !t.Mobile() {
		//t.cafeRequestQueue.Add(result.Id)
		return result, nil
	}

	// make an archive for remote pinning by the OS
	apath := filepath.Join(t.repoPath, "tmp", result.Id)
	result.Archive, err = archive.NewArchive(&apath)
	if err != nil {
		return nil, err
	}
	defer result.Archive.Close()

	// add files
	if err := result.Archive.AddFile(thumbcipher, "thumb"); err != nil {
		return nil, err
	}
	if err := result.Archive.AddFile(smallcipher, "small"); err != nil {
		return nil, err
	}
	if err := result.Archive.AddFile(mediumcipher, "medium"); err != nil {
		return nil, err
	}
	if err := result.Archive.AddFile(largecipher, "photo"); err != nil {
		return nil, err
	}
	if err := result.Archive.AddFile(metacipher, "meta"); err != nil {
		return nil, err
	}
	if err := result.Archive.AddFile(addrcipher, "pk"); err != nil {
		return nil, err
	}

	// all done
	return result, nil
}

// PhotoThreads lists threads which contain a photo (known to the local peer)
func (t *Textile) PhotoThreads(id string) []*Thread {
	blocks := t.datastore.Blocks().List("", -1, "dataId='"+id+"'")
	if len(blocks) == 0 {
		return nil
	}
	var threads []*Thread
	for _, block := range blocks {
		if _, thrd := t.GetThread(block.ThreadId); thrd != nil {
			threads = append(threads, thrd)
		}
	}
	return threads
}

// GetPhotoKey returns a the decrypted AES key for a photo set
func (t *Textile) GetPhotoKey(id string) (string, error) {
	block, err := t.GetBlockByDataId(id)
	if err != nil {
		return "", err
	}
	_, thrd := t.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadId))
	}
	key, err := thrd.GetBlockDataKey(block)
	if err != nil {
		return "", err
	}
	return string(key), nil
}

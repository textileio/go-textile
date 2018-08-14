package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cafe "github.com/textileio/textile-go/core/cafe"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/wallet/thread"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"os"
	"path/filepath"
	"strings"
)

// AddPhoto add a photo to the local ipfs node
func (w *Wallet) AddPhoto(path string) (*AddDataResult, error) {
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
	reader, format, size, err := util.DecodeImage(file)
	if err != nil {
		return nil, err
	}
	var encodingFormat util.Format
	if *format == util.GIF {
		encodingFormat = util.GIF
	} else {
		encodingFormat = util.JPEG
	}

	// make all image sizes
	reader.Seek(0, 0)
	thumb, err := util.EncodeImage(reader, encodingFormat, util.ThumbnailSize)
	if err != nil {
		return nil, err
	}
	reader.Seek(0, 0)
	small, err := util.EncodeImage(reader, encodingFormat, util.SmallSize)
	if err != nil {
		return nil, err
	}
	reader.Seek(0, 0)
	medium, err := util.EncodeImage(reader, encodingFormat, util.MediumSize)
	if err != nil {
		return nil, err
	}
	reader.Seek(0, 0)
	large, err := util.EncodeImage(reader, encodingFormat, util.LargeSize)
	if err != nil {
		return nil, err
	}

	// make meta data
	fpath := file.Name()
	ext := strings.ToLower(filepath.Ext(fpath))
	reader.Seek(0, 0)
	meta, err := util.MakeMetadata(reader, fpath, ext, *format, encodingFormat, size.X, size.Y, w.version)
	if err != nil {
		return nil, err
	}
	metab, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	// get public key
	mpk, err := w.GetPubKey()
	if err != nil {
		return nil, err
	}
	mpkb, err := mpk.Bytes()
	if err != nil {
		return nil, err
	}

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
	mpkcipher, err := crypto.EncryptAES(mpkb, key)
	if err != nil {
		return nil, err
	}

	// create a virtual directory for the photo
	dirb := uio.NewDirectory(w.ipfs.DAG)
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader(thumbcipher), "thumb"); err != nil {
		return nil, err
	}
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader(smallcipher), "small"); err != nil {
		return nil, err
	}
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader(mediumcipher), "medium"); err != nil {
		return nil, err
	}
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader(largecipher), "photo"); err != nil {
		return nil, err
	}
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader(metacipher), "meta"); err != nil {
		return nil, err
	}
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader(mpkcipher), "pk"); err != nil {
		return nil, err
	}

	// pin the directory
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := util.PinDirectory(w.ipfs, dir, []string{"medium", "photo"}); err != nil {
		return nil, err
	}
	result := &AddDataResult{Id: dir.Cid().Hash().B58String(), Key: string(key)}

	// if not mobile, create a pin request
	// on mobile, we let the OS handle the archive directly
	if !w.isMobile {
		if err := w.putPinRequest(result.Id); err != nil {
			return nil, err
		}
		return result, nil
	}

	// make an archive for remote pinning by the OS
	apath := filepath.Join(w.repoPath, "tmp", result.Id)
	result.Archive, err = cafe.NewArchive(&apath)
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
	if err := result.Archive.AddFile(mpkcipher, "pk"); err != nil {
		return nil, err
	}

	// all done
	return result, nil
}

// PhotoThreads lists threads which contain a photo (known to the local peer)
func (w *Wallet) PhotoThreads(id string) []*thread.Thread {
	blocks := w.datastore.Blocks().List("", -1, "dataId='"+id+"'")
	if len(blocks) == 0 {
		return nil
	}
	var threads []*thread.Thread
	for _, block := range blocks {
		if _, thrd := w.GetThread(block.ThreadId); thrd != nil {
			threads = append(threads, thrd)
		}
	}
	return threads
}

// GetPhotoKey returns a the decrypted AES key for a photo set
func (w *Wallet) GetPhotoKey(id string) (string, error) {
	block, err := w.GetBlockByDataId(id)
	if err != nil {
		return "", err
	}
	_, thrd := w.GetThread(block.ThreadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadId))
	}
	key, err := thrd.GetBlockDataKey(block)
	if err != nil {
		return "", err
	}
	return string(key), nil
}

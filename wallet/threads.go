package wallet

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"image"
	"bytes"
	"image/jpeg"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/disintegration/imaging"
	"github.com/textileio/textile-go/crypto"
	uio "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/unixfs/io"
	"github.com/textileio/textile-go/net"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	"github.com/textileio/textile-go/repo"
	"fmt"
	"gx/ipfs/QmSFihvoND3eDaAYRCeLgLPt62yCPgMZs1NSZmKFEtJQQw/go-libp2p-floodsub"
	"context"
	"io"
	"encoding/json"
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"sync"
)

type Thread struct {
	Id        string
	Name      string
	PrivKey   libp2pc.PrivKey
	Head      string
	LeftCh    chan struct{}
	leaveCh   chan struct{}
	repoPath  string
	ipfs      *core.IpfsNode
	blocks    repo.BlockStore
	update    func (head string) error
	mux       sync.Mutex
	listening bool
}

type Block struct {
	Id           string
	Target       string
	Parents      []string
	TargetKey    []byte
	ThreadPubKey []byte
	Type         BlockType
	Date         time.Time
}

type BlockType int
const (
	InviteBlock BlockType = iota
	PhotoBlock
	CommentBlock
	LikeBlock
)

func (bt BlockType) Bytes() []byte {
	return []byte(strconv.Itoa(int(bt)))
}

const thumbnailWidth = 300

type PhotoMetadata struct {
	FileMetadata
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lon,omitempty"`
}

// ThreadUpdate is used to notify listeners about updates in a thread
type ThreadUpdate struct {
	Id       string `json:"id"`
	Thread   string `json:"thread"`
	ThreadID string `json:"thread_id"`
}

// GetFile cats data from ipfs and tries to decrypt it with the provided block
// e.g., Qm../thumb, Qm../photo, Qm../meta, Qm../caption
func (t *Thread) GetFile(path string, block *Block) ([]byte, error) {
	// get bytes
	cypher, err := GetDataAtPath(t.ipfs, path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// decrypt the file key
	key, err := t.Decrypt(block.TargetKey)
	if err != nil {
		log.Errorf("error decrypting key: %s", err)
		return nil, err
	}

	// finally, decrypt the file
	return crypto.DecryptAES(cypher, key)
}

// GetFileBase64 returns data encoded as base64 under an ipfs path
func (t *Thread) GetFileBase64(path string, block *Block) (string, error) {
	file, err := t.GetFile(path, block)
	if err != nil {
		return "error", err
	}
	return base64.StdEncoding.EncodeToString(file), nil
}

// GetMetaData returns photo metadata under an id
func (t *Thread) GetPhotoMetaData(id string, block *Block) (*PhotoMetadata, error) {
	file, err := t.GetFile(fmt.Sprintf("%s/meta", id), block)
	if err != nil {
		log.Errorf("error getting meta file %s: %s", id, err)
		return nil, err
	}
	var data *PhotoMetadata
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Errorf("error unmarshaling meta file: %s: %s", id, err)
		return nil, err
	}
	return data, nil
}

// GetLastHash return the caption under an id
func (t *Thread) GetCaption(id string, block *Block) (string, error) {
	file, err := t.GetFile(fmt.Sprintf("%s/caption", id), block)
	if err != nil {
		log.Errorf("error getting caption file %s: %s", id, err)
		return "", err
	}
	return string(file), nil
}

// JoinRoom with a given id
func (t *Thread) Subscribe(datac chan ThreadUpdate) {
	if t.listening {
		return
	}
	t.listening = true
	sub, err := t.ipfs.Floodsub.Subscribe(t.Id)
	if err != nil {
		log.Errorf("error creating subscription: %s", err)
		return
	}
	log.Infof("joined thread: %s\n", t.Id)

	t.leaveCh = make(chan struct{})
	t.LeftCh = make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	leave := func() {
		cancel()
		close(t.LeftCh)
		t.listening = false
		log.Infof("left thread: %s\n", sub.Topic())
	}

	defer func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("thread update channel already closed")
			}
		}()
		close(datac)
	}()
	go func() {
		for {
			// unload new message
			msg, err := sub.Next(ctx)
			if err == io.EOF || err == context.Canceled {
				log.Debugf("thread subscription ended: %s", err)
				return
			} else if err != nil {
				log.Debugf(err.Error())
				return
			}

			// handle the update
			go func(msg *floodsub.Message) {
				if err = t.preHandleBlock(msg, datac); err != nil {
					log.Errorf("error handling room update: %s", err)
				}
			}(msg)
		}
	}()

	// block so we can shutdown with the leave room signal
	for {
		select {
		case <-t.leaveCh:
			leave()
			return
		case <-t.ipfs.Context().Done():
			leave()
			return
		}
	}
}

// Unsubscribe leaves the thread
func (t *Thread) Unsubscribe() {
	if t.leaveCh == nil {
		return
	}
	close(t.leaveCh)
}

// Listening indicates whether or not we are listening in the thread
func (t *Thread) Listening() bool {
	return t.listening
}

// AddPhoto adds a block for a photo to this thread
func (t *Thread) AddPhoto(id string, caption string, key []byte) (*AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	// encrypt AES key with thread pk
	keycypher, err := t.Encrypt(key)
	if err != nil {
		return nil, err
	}
	threadkey, err := t.PrivKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	typeb := PhotoBlock.Bytes() // silly?
	dateb := getNowBytes()

	// encrypt caption with thread pk
	captioncypher, err := t.Encrypt([]byte(caption))
	if err != nil {
		return nil, err
	}

	// create a virtual directory for the new block
	dirb := uio.NewDirectory(t.ipfs.DAG)
	err = addFileToDirectory(t.ipfs, dirb, []byte(id), "target")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(t.ipfs, dirb, []byte(t.Head), "parents")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(t.ipfs, dirb, keycypher, "key")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(t.ipfs, dirb, threadkey, "pk")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(t.ipfs, dirb, typeb, "type")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(t.ipfs, dirb, dateb, "date")
	if err != nil {
		return nil, err
	}
	err = addFileToDirectory(t.ipfs, dirb, captioncypher, "caption")
	if err != nil {
		return nil, err
	}

	// pin it
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := pinDirectory(t.ipfs, dir, []string{}); err != nil {
		return nil, err
	}
	bid := dir.Cid().Hash().B58String()

	// index it
	block, err := t.indexBlock(bid)
	if err != nil {
		return nil, err
	}

	// update head
	if err := t.update(bid); err != nil {
		return nil, err
	}

	// post it
	go func(bid string) {
		err = t.post([]byte(id))
		if err != nil {
			log.Errorf("error posting block %s: %s", bid, err)
		}
		log.Debugf("posted block %s to %s", bid, t.Id)
	}(block.Id)

	// create and init a new multipart request
	request := &net.MultipartRequest{}
	request.Init(filepath.Join(t.repoPath, "tmp"), bid)

	// add files to request
	if err := request.AddFile([]byte(id), "target"); err != nil {
		return nil, err
	}
	if err := request.AddFile([]byte(t.Head), "parents"); err != nil {
		return nil, err
	}
	if err := request.AddFile(keycypher, "key"); err != nil {
		return nil, err
	}
	if err := request.AddFile(threadkey, "pk"); err != nil {
		return nil, err
	}
	if err := request.AddFile(typeb, "type"); err != nil {
		return nil, err
	}
	if err := request.AddFile(dateb, "date"); err != nil {
		return nil, err
	}
	if err := request.AddFile(captioncypher, "caption"); err != nil {
		return nil, err
	}

	// finish request
	if err := request.Finish(); err != nil {
		return nil, err
	}

	// all done
	return &AddResult{Id: block.Id, RemoteRequest: request}, nil
}

// Blocks paginates photos from the datastore
func (t *Thread) Blocks(offsetId string, limit int) []Block {
	log.Debugf("listing blocks: offsetId: %s, limit: %d, thread: %s", offsetId, limit, t.Name)
	query := fmt.Sprintf("pk='%s' and type=%d", t.Id, PhotoBlock)
	list := t.blocks.List(offsetId, limit, query)
	log.Debugf("found %d photos in thread %s", len(list), t.Name)
	return list
}

func (t *Thread) Encrypt(data []byte) ([]byte, error) {
	return crypto.Encrypt(t.PrivKey.GetPublic(), data)
}

func (t *Thread) Decrypt(data []byte) ([]byte, error) {
	return crypto.Decrypt(t.PrivKey, data)
}

// Publish publishes HEAD
func (t *Thread) Publish() {
	if t.Head == "" {
		return
	}
	log.Debugf("publishing thread %s...", t.Head, t.Name)
	if err := t.post([]byte(t.Head)); err != nil {
		log.Errorf("error publishing %s: %s", t.Head, err)
		return
	}
	log.Debugf("published %s to %s thread", t.Head, t.Id)
}

func (t *Thread) post(payload []byte) error {
	return t.ipfs.Floodsub.Publish(t.Id, payload)
}

// preHandleBlock tries to recursively process an update sent to a thread
func (t *Thread) preHandleBlock(msg *floodsub.Message, datac chan ThreadUpdate) error {
	// unpack from
	from := msg.GetFrom().Pretty()
	if from == t.ipfs.Identity.Pretty() {
		return nil
	}

	// unpack message data
	data := string(msg.GetData())

	// determine if this is from a relay node
	tmp := strings.Split(data, ":")
	var id string
	if len(tmp) > 1 && tmp[0] == "relay" {
		id = tmp[1]
		from = fmt.Sprintf("relay:%s", from)
	} else {
		id = tmp[0]
	}
	log.Debugf("got block from %s in thread %s", from, t.Id)

	// recurse back in time starting at this hash
	err := t.handleBlock(id, datac)
	if err != nil {
		return err
	}

	return nil
}

// handleBlock tries to process a block
func (t *Thread) handleBlock(id string, datac chan ThreadUpdate) error {
	// first update?
	if id == "" {
		log.Debugf("found genesis block, aborting")
		return nil
	}
	log.Debugf("handling block: %s...", id)

	// check if we aleady have this block
	block := t.blocks.Get(id)
	if block != nil {
		log.Debugf("block %s exists, aborting", id)
		return nil
	}

	log.Debugf("pinning block %s...", id)
	if err := pinPath(t.ipfs, id, true); err != nil {
		return err
	}

	log.Debugf("indexing %s...", id)
	block, err := t.indexBlock(id)
	if err != nil {
		return err
	}

	// don't block on the send since nobody might be listening
	select {
	case datac <- ThreadUpdate{Id: id, Thread: t.Name, ThreadID: t.Id}:
	default:
	}
	defer func() {
		if recover() != nil {
			log.Error("update channel already closed")
		}
	}()

	// check last block
	// TODO: handle multi parents from 3-way merge
	return t.handleBlock(block.Parents[0], datac)
}

func (t *Thread) indexBlock(id string) (*Block, error) {
	target, err := GetDataAtPath(t.ipfs, fmt.Sprintf("%s/target", id))
	if err != nil {
		return nil, err
	}
	parents, err := GetDataAtPath(t.ipfs, fmt.Sprintf("%s/parents", id))
	if err != nil {
		return nil, err
	}
	key, err := GetDataAtPath(t.ipfs, fmt.Sprintf("%s/key", id))
	if err != nil {
		return nil, err
	}
	pk, err := GetDataAtPath(t.ipfs, fmt.Sprintf("%s/pk", id))
	if err != nil {
		return nil, err
	}
	typeb, err := GetDataAtPath(t.ipfs, fmt.Sprintf("%s/type", id))
	if err != nil {
		return nil, err
	}
	dateb, err := GetDataAtPath(t.ipfs, fmt.Sprintf("%s/date", id))
	if err != nil {
		return nil, err
	}
	typei := binary.BigEndian.Uint64(typeb)
	datei := binary.BigEndian.Uint64(dateb)
	block := &Block{
		Id:           id,
		Target:       string(target),
		Parents:      strings.Split(string(parents), ","),
		TargetKey:    key,
		ThreadPubKey: pk,
		Type:         BlockType(int(typei)),
		Date:         time.Unix(int64(datei), 0),
	}
	if err := t.blocks.Add(block); err != nil {
		return nil, err
	}
	return block, nil
}

func makeThumbnail(photo *os.File, width int) ([]byte, error) {
	img, _, err := image.Decode(photo)
	if err != nil {
		return nil, err
	}
	thumb := imaging.Resize(img, width, 0, imaging.Lanczos)
	buff := new(bytes.Buffer)
	if err = jpeg.Encode(buff, thumb, nil); err != nil {
		return nil, err
	}
	photo.Seek(0, 0) // be kind, rewind
	return buff.Bytes(), nil
}

// TODO: get image size info
func getMetadata(photo *os.File, path string, ext string, username string) (PhotoMetadata, error) {
	var created time.Time
	var lat, lon float64
	x, err := exif.Decode(photo)
	if err == nil {
		// time taken
		createdTmp, err := x.DateTime()
		if err == nil {
			created = createdTmp
		}
		// coords taken
		latTmp, lonTmp, err := x.LatLong()
		if err == nil {
			lat, lon = latTmp, lonTmp
		}
	}
	meta := PhotoMetadata{
		FileMetadata: FileMetadata{
			Metadata: Metadata{
				Username: username,
				Created: created,
				Added: time.Now(),
			},
			Name: strings.TrimSuffix(filepath.Base(path), ext),
			Ext: ext,
		},
		Latitude: lat,
		Longitude: lon,
	}
	photo.Seek(0, 0) // be kind, rewind
	return meta, nil
}

// WaitForRoom to join
//func (t *Thread) WaitForInvite() {
//	// we're in a lonesome state here, we can just sub to our own
//	// peer id and hope somebody sends us a priv key to join a thread with
//	rid := t.ipfs.Identity.Pretty()
//	sub, err := t.ipfs.Floodsub.Subscribe(rid)
//	if err != nil {
//		log.Errorf("error creating subscription: %s", err)
//		return
//	}
//	log.Infof("waiting for invite at own peer id: %s\n", rid)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	cancelCh := make(chan struct{})
//	go func() {
//		for {
//			msg, err := sub.Next(ctx)
//			if err == io.EOF || err == context.Canceled {
//				log.Debugf("wait subscription ended: %s", err)
//				return
//			} else if err != nil {
//				log.Debugf(err.Error())
//				return
//			}
//			from := msg.GetFrom().Pretty()
//			log.Infof("got pairing request from: %s\n", from)
//
//			// get private peer key and decrypt the phrase
//			sk, err := t.UnmarshalPrivatePeerKey()
//			if err != nil {
//				log.Errorf("error unmarshaling priv peer key: %s", err)
//				return
//			}
//			p, err := crypto.Decrypt(sk, msg.GetData())
//			if err != nil {
//				log.Errorf("error decrypting msg data: %s", err)
//				return
//			}
//			ps := string(p)
//			log.Debugf("decrypted mnemonic phrase as: %s\n", ps)
//
//			// create a new album for the room
//			// TODO: let user name this or take phone's name, e.g., bob's iphone
//			// TODO: or auto name it, cause this means only one pairing can happen
//			t.Wallet.AddThread("mobile", ps)
//
//			// we're done
//			close(cancelCh)
//		}
//	}()
//
//	for {
//		select {
//		case <-cancelCh:
//			cancel()
//			return
//		case <-t.IpfsNode.Context().Done():
//			cancel()
//			return
//		}
//	}
//}
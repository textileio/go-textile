package thread

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/op/go-logging"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net"
	nm "github.com/textileio/textile-go/net/model"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/util"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var log = logging.MustGetLogger("thread")

const EmptyFileString = "__nil__"

// Config is used to construct a Thread
type Config struct {
	RepoPath   string
	Ipfs       func() *core.IpfsNode
	Blocks     func() repo.BlockStore
	Peers      func() repo.PeerStore
	GetHead    func() (string, error)
	UpdateHead func(head string) error
	Publish    func(payload []byte) error
	Send       func(message *pb.Message, peerId string) error
}

// ThreadUpdate is used to notify listeners about updates in a thread
type Update struct {
	Id         string         `json:"id"`
	Type       repo.BlockType `json:"type"`
	TargetId   string         `json:"target_id"`
	ThreadId   string         `json:"thread_id"`
	ThreadName string         `json:"thread_name"`
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id         string
	Name       string
	PrivKey    libp2pc.PrivKey
	updates    chan Update
	repoPath   string
	ipfs       func() *core.IpfsNode
	blocks     func() repo.BlockStore
	peers      func() repo.PeerStore
	GetHead    func() (string, error)
	updateHead func(head string) error
	publish    func(payload []byte) error
	send       func(message *pb.Message, peerId string) error
	mux        sync.Mutex
}

// NewThread create a new Thread from a repo model and config
func NewThread(model *repo.Thread, config *Config) (*Thread, error) {
	sk, err := libp2pc.UnmarshalPrivateKey(model.PrivKey)
	if err != nil {
		return nil, err
	}
	return &Thread{
		Id:         model.Id,
		Name:       model.Name,
		PrivKey:    sk,
		updates:    make(chan Update),
		repoPath:   config.RepoPath,
		ipfs:       config.Ipfs,
		blocks:     config.Blocks,
		peers:      config.Peers,
		GetHead:    config.GetHead,
		updateHead: config.UpdateHead,
		publish:    config.Publish,
		send:       config.Send,
	}, nil
}

func (t *Thread) Close() {
	close(t.updates)
}

// AddInvite creates an invite block for the given recipient
func (t *Thread) AddInvite(target libp2pc.PubKey) (*nm.AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// get the peer id from the pub key
	targetId, err := peer.IDFromPublicKey(target)
	if err != nil {
		return nil, err
	}
	targetpk, err := target.Bytes()
	if err != nil {
		return nil, err
	}

	// encypt thread secret with the recipient's public key
	threadsk, err := t.PrivKey.Bytes()
	if err != nil {
		return nil, err
	}
	threadskcipher, err := crypto.Encrypt(target, threadsk)
	if err != nil {
		return nil, err
	}

	// add a block
	res, err := t.addBlock(repo.InviteBlock, []byte(targetId.Pretty()), threadskcipher, nil)
	if err != nil {
		return nil, err
	}

	// create new peer for posting, but don't add yet. will get added if they accept.
	newPeer := repo.Peer{
		Row:      ksuid.New().String(),
		Id:       targetId.Pretty(),
		ThreadId: t.Id,
		PubKey:   targetpk,
	}

	// post it
	t.PostHead([]repo.Peer{newPeer})

	// all done
	return res, nil
}

// AddExternalInvite creates an invite block for the given recipient
func (t *Thread) AddExternalInvite() (*nm.AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// generate an aes key
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, err
	}

	// encypt thread secret with the key
	threadsk, err := t.PrivKey.Bytes()
	if err != nil {
		return nil, err
	}
	threadskcipher, err := crypto.EncryptAES(threadsk, key)
	if err != nil {
		return nil, err
	}

	// add a block
	res, err := t.addBlock(repo.ExternalInviteBlock, []byte("unlimited"), threadskcipher, nil)
	if err != nil {
		return nil, err
	}

	// post it
	t.PostHead(t.Peers())

	// all done
	return &nm.AddResult{Id: res.Id, Key: key, RemoteRequest: res.RemoteRequest}, nil
}

// Join creates a join block
func (t *Thread) Join(from libp2pc.PubKey, id string) (*nm.AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// add a block
	res, err := t.addBlock(repo.JoinBlock, []byte(id), []byte(EmptyFileString), nil)
	if err != nil {
		return nil, err
	}

	// add new peer
	fromb, err := from.Bytes()
	if err != nil {
		return nil, err
	}
	pid, err := peer.IDFromPublicKey(from)
	if err != nil {
		return nil, err
	}
	newPeer := &repo.Peer{
		Row:      ksuid.New().String(),
		Id:       pid.Pretty(),
		ThreadId: t.Id,
		PubKey:   fromb,
	}
	if err := t.peers().Add(newPeer); err != nil {
		return nil, err
	}

	// post it
	t.PostHead(t.Peers())

	// all done
	return res, nil
}

// Leave creates a leave block
func (t *Thread) Leave() (*nm.AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// add a block
	res, err := t.addBlock(repo.LeaveBlock, []byte(t.ipfs().Identity.Pretty()), []byte(EmptyFileString), nil)
	if err != nil {
		return nil, err
	}

	// post it
	t.PostHead(t.Peers())

	// delete blocks
	if err := t.blocks().DeleteByThread(t.Id); err != nil {
		return nil, err
	}
	// delete peers
	if err := t.peers().DeleteByThread(t.Id); err != nil {
		return nil, err
	}

	// all done
	return res, nil
}

// AddPhoto adds a block for a photo to this thread
func (t *Thread) AddPhoto(id string, caption string, key []byte) (*nm.AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// encrypt AES key with thread pk
	keycipher, err := t.Encrypt(key)
	if err != nil {
		return nil, err
	}

	// encrypt caption with thread pk
	captioncipher, err := t.Encrypt([]byte(caption))
	if err != nil {
		return nil, err
	}

	// add a block
	res, err := t.addBlock(repo.DataBlock, []byte(id), keycipher, captioncipher)
	if err != nil {
		return nil, err
	}

	// post it
	t.PostHead(t.Peers())

	// all done
	return res, nil
}

// GetBlockData cats file data from ipfs and tries to decrypt it with the provided block
func (t *Thread) GetBlockData(path string, block *repo.Block) ([]byte, error) {
	// get bytes
	cipher, err := util.GetDataAtPath(t.ipfs(), path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// decrypt with thread key
	return t.Decrypt(cipher)
}

// GetBlockDataBase64 returns block data encoded as base64 under an ipfs path
func (t *Thread) GetBlockDataBase64(path string, block *repo.Block) (string, error) {
	file, err := t.GetBlockData(path, block)
	if err != nil {
		return "error", err
	}
	return libp2pc.ConfigEncodeKey(file), nil
}

// GetFileKey returns the decrypted AES key for a block
func (t *Thread) GetFileKey(block *repo.Block) (string, error) {
	key, err := t.Decrypt(block.TargetKey)
	if err != nil {
		log.Errorf("error decrypting key: %s", err)
		return "", err
	}
	return string(key), nil
}

// GetFileData cats file data from ipfs and tries to decrypt it with the provided block
func (t *Thread) GetFileData(path string, block *repo.Block) ([]byte, error) {
	// get bytes
	cipher, err := util.GetDataAtPath(t.ipfs(), path)
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
	return crypto.DecryptAES(cipher, key)
}

// GetFileDataBase64 returns file data encoded as base64 under an ipfs path
func (t *Thread) GetFileDataBase64(path string, block *repo.Block) (string, error) {
	file, err := t.GetFileData(path, block)
	if err != nil {
		return "error", err
	}
	return libp2pc.ConfigEncodeKey(file), nil
}

// GetMetaData returns photo metadata under an id
func (t *Thread) GetPhotoMetaData(id string, block *repo.Block) (*model.PhotoMetadata, error) {
	file, err := t.GetFileData(fmt.Sprintf("%s/meta", id), block)
	if err != nil {
		log.Errorf("error getting meta file %s: %s", id, err)
		return nil, err
	}
	var data *model.PhotoMetadata
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Errorf("error unmarshaling meta file: %s: %s", id, err)
		return nil, err
	}
	return data, nil
}

// Blocks paginates blocks from the datastore
func (t *Thread) Blocks(offsetId string, limit int, bType repo.BlockType) []repo.Block {
	query := fmt.Sprintf("pk='%s' and type=%d", t.Id, bType)
	return t.blocks().List(offsetId, limit, query)
}

// Encrypt data with thread public key
func (t *Thread) Encrypt(data []byte) ([]byte, error) {
	return crypto.Encrypt(t.PrivKey.GetPublic(), data)
}

// Decrypt data with thread secret key
func (t *Thread) Decrypt(data []byte) ([]byte, error) {
	return crypto.Decrypt(t.PrivKey, data)
}

// Sign data with thread secret key
func (t *Thread) Sign(data []byte) ([]byte, error) {
	return t.PrivKey.Sign(data)
}

// Verify data with thread public key
func (t *Thread) Verify(data []byte, sig []byte) error {
	good, err := t.PrivKey.GetPublic().Verify(data, sig)
	if err != nil || !good {
		return errors.New("bad signature")
	}
	return nil
}

// PostHead publishes HEAD to peers
func (t *Thread) PostHead(peers []repo.Peer) error {
	if len(peers) == 0 {
		return nil
	}
	head, err := t.GetHead()
	if err != nil {
		log.Errorf("failed to get HEAD for %s: %s", t.Id, err)
		return err
	}
	block := t.blocks().Get(head)
	if block == nil {
		return nil
	}
	message, err := t.getMessageForBlock(block)
	if err != nil {
		log.Errorf("failed to get message for block %s: %s", block.Id, err)
		return err
	}

	log.Debugf("posting %s in thread %s...", block.Id, t.Name)
	wg := sync.WaitGroup{}
	for _, p := range peers {
		wg.Add(1)
		go func(peerId string) {
			if err := t.send(message, peerId); err != nil {
				log.Errorf("error sending block %s to peer %s: %s", block.Id, peerId, err)
			}
			wg.Done()
		}(p.Id)
	}
	wg.Wait()
	log.Debugf("posted to %d peers", len(peers))
	return nil
}

// Peers returns locally known peers in this thread
func (t *Thread) Peers() []repo.Peer {
	query := fmt.Sprintf("thread='%s'", t.Id)
	return t.peers().List("", -1, query)
}

// Updates returns a read-only channel of updates
func (t *Thread) Updates() <-chan Update {
	return t.updates
}

// handleBlock tries to process a block
func (t *Thread) HandleBlock(id string) error {
	// first update?
	if id == "" {
		log.Debugf("found genesis block, aborting")
		return nil
	}

	// check if we aleady have this block
	block := t.blocks().Get(id)
	if block != nil {
		log.Debugf("block %s exists, aborting", id)
		return nil
	}

	// pin it
	if err := util.PinPath(t.ipfs(), id, true); err != nil {
		return err
	}

	// index it
	block, err := t.indexBlock(id)
	if err != nil {
		return err
	}

	// update current head
	if err := t.updateHead(id); err != nil {
		return err
	}

	// don't block on the send since nobody could be listening
	defer func() {
		if recover() != nil {
			log.Error("update channel already closed")
		}
	}()
	select {
	case t.updates <- Update{
		Id:         id,
		Type:       block.Type,
		TargetId:   block.Target,
		ThreadId:   t.Id,
		ThreadName: t.Name,
	}:
	default:
	}

	log.Debugf("handled block: %s", id)

	// check last block
	// TODO: handle multi parents from 3-way merge
	return t.HandleBlock(block.Parents[0])
}

// indexBlock attempts to download the block and index it in the local db
func (t *Thread) indexBlock(id string) (*repo.Block, error) {
	target, err := util.GetDataAtPath(t.ipfs(), fmt.Sprintf("%s/target", id))
	if err != nil {
		return nil, err
	}
	parents, err := util.GetDataAtPath(t.ipfs(), fmt.Sprintf("%s/parents", id))
	if err != nil {
		return nil, err
	}
	key, err := util.GetDataAtPath(t.ipfs(), fmt.Sprintf("%s/key", id))
	if err != nil {
		return nil, err
	}
	pk, err := util.GetDataAtPath(t.ipfs(), fmt.Sprintf("%s/pk", id))
	if err != nil {
		return nil, err
	}
	ppk, err := util.GetDataAtPath(t.ipfs(), fmt.Sprintf("%s/ppk", id))
	if err != nil {
		return nil, err
	}
	typeb, err := util.GetDataAtPath(t.ipfs(), fmt.Sprintf("%s/type", id))
	if err != nil {
		return nil, err
	}
	typei, err := strconv.ParseInt(string(typeb), 10, 0)
	if err != nil {
		return nil, err
	}
	dateb, err := util.GetDataAtPath(t.ipfs(), fmt.Sprintf("%s/date", id))
	if err != nil {
		return nil, err
	}
	datei, err := strconv.ParseInt(string(dateb), 10, 0)
	if err != nil {
		return nil, err
	}
	block := &repo.Block{
		Id:           id,
		Target:       string(target),
		Parents:      strings.Split(string(parents), ","),
		TargetKey:    key,
		ThreadPubKey: string(pk),
		PeerPubKey:   string(ppk),
		Type:         repo.BlockType(int(typei)),
		Date:         time.Unix(int64(datei), 0),
	}
	if err := t.blocks().Add(block); err != nil {
		return nil, err
	}
	return block, nil
}

func (t *Thread) getMessageForBlock(block *repo.Block) (*pb.Message, error) {
	// translate to pb
	// TODO: throw out repo.Block
	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}
	pblock := &pb.ThreadBlock{
		Id:           block.Id,
		Target:       block.Target,
		Parents:      block.Parents,
		TargetKey:    block.TargetKey,
		ThreadPubKey: block.ThreadPubKey,
		Date:         date,
	}

	var wrap proto.Message
	var wrapType pb.Message_Type
	switch block.Type {
	case repo.InviteBlock:
		wrap = &pb.ThreadInvite{Block: pblock, Type: pb.ThreadInvite_INTERNAL}
		wrapType = pb.Message_THREAD_INVITE
	case repo.ExternalInviteBlock:
		wrap = &pb.ThreadInvite{Block: pblock, Type: pb.ThreadInvite_EXTERNAL}
		wrapType = pb.Message_THREAD_INVITE
	case repo.JoinBlock:
		wrap = &pb.ThreadJoin{Block: pblock}
		wrapType = pb.Message_THREAD_JOIN
	case repo.LeaveBlock:
		wrap = &pb.ThreadJoin{Block: pblock}
		wrapType = pb.Message_THREAD_LEAVE
	case repo.DataBlock:
		wrap = &pb.ThreadData{Block: pblock}
		wrapType = pb.Message_THREAD_DATA
	case repo.AnnotationBlock:
		wrap = &pb.ThreadAnnotation{Block: pblock}
		wrapType = pb.Message_THREAD_ANNOTATION
	}

	// sign it
	serialized, err := proto.Marshal(wrap)
	if err != nil {
		return nil, err
	}
	threadSignature, err := t.PrivKey.Sign(serialized)
	if err != nil {
		return nil, err
	}
	issuerSignature, err := t.ipfs().PrivateKey.Sign(serialized)
	if err != nil {
		return nil, err
	}
	pkb, err := t.ipfs().PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	signed := &pb.SignedThreadBlock{
		Id:              block.Id,
		Data:            serialized,
		IssuerSignature: issuerSignature,
		ThreadSignature: threadSignature,
		ThreadId:        t.Id,
		ThreadName:      t.Name,
		IssuerPubKey:    pkb,
	}

	// create the message
	payload, err := ptypes.MarshalAny(signed)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Type: wrapType, Payload: payload}, nil
}

func (t *Thread) addBlock(bt repo.BlockType, target []byte, key []byte, caption []byte) (*nm.AddResult, error) {
	// get current HEAD
	head, err := t.GetHead()
	if err != nil {
		return nil, err
	}
	now := util.GetNowBytes()
	if caption == nil {
		caption = []byte(EmptyFileString)
	}

	// get our own public key
	peerpk, err := t.ipfs().PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// create a virtual directory for the new block
	dirb := uio.NewDirectory(t.ipfs().DAG)
	err = util.AddFileToDirectory(t.ipfs(), dirb, target, "target")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(head), "parents")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, key, "key")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(t.Id), "pk")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, peerpk, "ppk")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, bt.Bytes(), "type")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, now, "date")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, caption, "caption")
	if err != nil {
		return nil, err
	}

	// pin it
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := util.PinDirectory(t.ipfs(), dir, []string{}); err != nil {
		return nil, err
	}
	bid := dir.Cid().Hash().B58String()

	// index it
	block, err := t.indexBlock(bid)
	if err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(bid); err != nil {
		return nil, err
	}

	// create and init a new multipart request
	request := &net.PinRequest{}
	request.Init(filepath.Join(t.repoPath, "tmp"), block.Id)

	// add files to request
	if err := request.AddFile(target, "target"); err != nil {
		return nil, err
	}
	if err := request.AddFile([]byte(head), "parents"); err != nil {
		return nil, err
	}
	if err := request.AddFile(key, "key"); err != nil {
		return nil, err
	}
	if err := request.AddFile([]byte(t.Id), "pk"); err != nil {
		return nil, err
	}
	if err := request.AddFile(peerpk, "ppk"); err != nil {
		return nil, err
	}
	if err := request.AddFile(bt.Bytes(), "type"); err != nil {
		return nil, err
	}
	if err := request.AddFile(now, "date"); err != nil {
		return nil, err
	}
	if err := request.AddFile(caption, "caption"); err != nil {
		return nil, err
	}

	// finish request
	if err := request.Finish(); err != nil {
		return nil, err
	}

	// all done
	return &nm.AddResult{Id: block.Id, RemoteRequest: request}, nil
}

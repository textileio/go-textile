package thread

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/op/go-logging"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/util"
	"gx/ipfs/QmSFihvoND3eDaAYRCeLgLPt62yCPgMZs1NSZmKFEtJQQw/go-libp2p-floodsub"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	uio "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/unixfs/io"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var log = logging.MustGetLogger("thread")

// ErrInvalidBlock is used to reject invalid block updates
var ErrInvalidBlock = errors.New("block is not a valid token")

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
	Id       string `json:"id"`
	Thread   string `json:"thread"`
	ThreadID string `json:"thread_id"`
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id         string
	Name       string
	PrivKey    libp2pc.PrivKey
	LeftCh     chan struct{}
	leaveCh    chan struct{}
	repoPath   string
	ipfs       func() *core.IpfsNode
	blocks     func() repo.BlockStore
	peers      func() repo.PeerStore
	GetHead    func() (string, error)
	updateHead func(head string) error
	publish    func(payload []byte) error
	send       func(message *pb.Message, peerId string) error
	mux        sync.Mutex
	listening  bool
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

// AddInvite creates an invite block for the given recipient
func (t *Thread) AddInvite(target libp2pc.PubKey) (*model.AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// get current HEAD
	head, err := t.GetHead()
	if err != nil {
		return nil, err
	}

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
	threadskcypher, err := crypto.Encrypt(target, threadsk)
	if err != nil {
		return nil, err
	}

	// type and date
	typeb := repo.InviteBlock.Bytes() // silly?
	dateb := util.GetNowBytes()

	// create a virtual directory for the new block
	dirb := uio.NewDirectory(t.ipfs().DAG)
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(targetId.Pretty()), "target")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(head), "parents")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, threadskcypher, "key")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(t.Id), "pk")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, typeb, "type")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, dateb, "date")
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

	// add new peer
	newPeer := &repo.Peer{
		Row:      ksuid.New().String(),
		Id:       targetId.Pretty(),
		ThreadId: t.Id,
		PubKey:   targetpk,
	}
	if err := t.peers().Add(newPeer); err != nil {
		return nil, err
	}

	// post it
	go t.PostHead()

	// all done
	return &model.AddResult{Id: block.Id}, nil
}

// AddPhoto adds a block for a photo to this thread
func (t *Thread) AddPhoto(id string, caption string, key []byte) (*model.AddResult, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// get current HEAD
	head, err := t.GetHead()
	if err != nil {
		return nil, err
	}

	// encrypt AES key with thread pk
	keycypher, err := t.Encrypt(key)
	if err != nil {
		return nil, err
	}

	// type and date
	typeb := repo.PhotoBlock.Bytes() // silly?
	dateb := util.GetNowBytes()

	// encrypt caption with thread pk
	captioncypher, err := t.Encrypt([]byte(caption))
	if err != nil {
		return nil, err
	}

	// create a virtual directory for the new block
	dirb := uio.NewDirectory(t.ipfs().DAG)
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(id), "target")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(head), "parents")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, keycypher, "key")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, []byte(t.Id), "pk")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, typeb, "type")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, dateb, "date")
	if err != nil {
		return nil, err
	}
	err = util.AddFileToDirectory(t.ipfs(), dirb, captioncypher, "caption")
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

	// post it
	go t.PostHead()

	// create and init a new multipart request
	request := &net.MultipartRequest{}
	request.Init(filepath.Join(t.repoPath, "tmp"), bid)

	// add files to request
	if err := request.AddFile([]byte(id), "target"); err != nil {
		return nil, err
	}
	if err := request.AddFile([]byte(head), "parents"); err != nil {
		return nil, err
	}
	if err := request.AddFile(keycypher, "key"); err != nil {
		return nil, err
	}
	if err := request.AddFile([]byte(t.Id), "pk"); err != nil {
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
	return &model.AddResult{Id: block.Id, RemoteRequest: request}, nil
}

// GetBlockData cats file data from ipfs and tries to decrypt it with the provided block
func (t *Thread) GetBlockData(path string, block *repo.Block) ([]byte, error) {
	// get bytes
	cypher, err := util.GetDataAtPath(t.ipfs(), path)
	if err != nil {
		log.Errorf("error getting file data: %s", err)
		return nil, err
	}

	// decrypt with thread key
	return t.Decrypt(cypher)
}

// GetBlockDataBase64 returns block data encoded as base64 under an ipfs path
func (t *Thread) GetBlockDataBase64(path string, block *repo.Block) (string, error) {
	file, err := t.GetBlockData(path, block)
	if err != nil {
		return "error", err
	}
	return base64.StdEncoding.EncodeToString(file), nil
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
	cypher, err := util.GetDataAtPath(t.ipfs(), path)
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

// GetFileDataBase64 returns file data encoded as base64 under an ipfs path
func (t *Thread) GetFileDataBase64(path string, block *repo.Block) (string, error) {
	file, err := t.GetFileData(path, block)
	if err != nil {
		return "error", err
	}
	return base64.StdEncoding.EncodeToString(file), nil
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

// Subscribe joins the thread
func (t *Thread) Subscribe(datac chan Update) {
	if t.listening {
		return
	}
	t.listening = true
	sub, err := t.ipfs().Floodsub.Subscribe(t.Id)
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
		case <-t.ipfs().Context().Done():
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
func (t *Thread) Verify(data []byte, sig []byte) (bool, error) {
	return t.PrivKey.GetPublic().Verify(data, sig)
}

// PostHead publishes HEAD to peers
func (t *Thread) PostHead() error {
	peers := t.GetPeers("", -1)
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
		return err
	}

	log.Debugf("posting %s in thread %s...", block.Id, t.Name)
	for _, p := range peers {
		if err := t.send(message, p.Id); err != nil {
			log.Errorf("error sending block %s to peer %s: %s", block.Id, p.Id, err)
			return err
		}
	}
	log.Debugf("posted to %d peers", len(peers))
	return nil
}

// GetPeers returns locally known peers in this thread
func (t *Thread) GetPeers(offset string, limit int) []repo.Peer {
	query := fmt.Sprintf("thread='%s'", t.Id)
	return t.peers().List(offset, limit, query)
}

// preHandleBlock tries to recursively process an update sent to a thread
func (t *Thread) preHandleBlock(msg *floodsub.Message, datac chan Update) error {
	// unpack from
	from := msg.GetFrom().Pretty()
	if from == t.ipfs().Identity.Pretty() {
		return nil
	}

	// unpack message data
	data := string(msg.GetData())

	// determine if this is from a relay node
	tmp := strings.Split(data, ":")
	var tokenStr string
	if len(tmp) > 1 && tmp[0] == "relay" {
		tokenStr = tmp[1]
		from = fmt.Sprintf("relay:%s", from)
	} else {
		tokenStr = tmp[0]
	}
	var id string

	// parse token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*crypto.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.PrivKey.GetPublic(), nil
	})
	if err != nil {
		return err
	}
	// validate
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok = claims["jti"].(string); !ok || id == "" {
			return ErrInvalidBlock
		}
		// TODO validate pk, iss, iat
	} else {
		return ErrInvalidBlock
	}

	log.Debugf("got block %s from %s in thread %s", id, from, t.Id)

	// exit if block is just a ping
	if id == "ping" {
		return nil
	}

	// recurse back in time starting at this hash
	if err := t.HandleBlock(id, datac); err != nil {
		return err
	}

	return nil
}

// handleBlock tries to process a block
func (t *Thread) HandleBlock(id string, datac chan Update) error {
	// first update?
	if id == "" {
		log.Debugf("found genesis block, aborting")
		return nil
	}
	log.Debugf("handling block: %s...", id)

	// check if we aleady have this block
	block := t.blocks().Get(id)
	if block != nil {
		log.Debugf("block %s exists, aborting", id)
		return nil
	}

	log.Debugf("pinning block %s...", id)
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

	// don't block on the send since nobody might be listening
	select {
	case datac <- Update{Id: id, Thread: t.Name, ThreadID: t.Id}:
	default:
	}
	defer func() {
		if recover() != nil {
			log.Error("update channel already closed")
		}
	}()

	log.Debugf("handled block: %s", id)

	// check last block
	// TODO: handle multi parents from 3-way merge
	return t.HandleBlock(block.Parents[0], datac)
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
		Type:         repo.BlockType(int(typei)),
		Date:         time.Unix(int64(datei), 0),
	}
	if err := t.blocks().Add(block); err != nil {
		return nil, err
	}
	return block, nil
}

// signBlock generated a valid JWT based on a thread block
func (t *Thread) signBlock(block *repo.Block) (string, error) {
	var blockId string
	var date time.Time
	if block != nil {
		blockId = block.Id
		date = block.Date
	} else {
		blockId = "ping"
		date = time.Now()
	}
	claims := jwt.StandardClaims{
		Id:       blockId, // block cid
		Subject:  t.Id,    // thread id (pk, base64)
		Issuer:   t.ipfs().Identity.Pretty(),
		IssuedAt: date.Unix(),
	}
	token, err := jwt.NewWithClaims(crypto.SigningMethodEd25519i, claims).SignedString(t.PrivKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (t *Thread) getMessageForBlock(block *repo.Block) (*pb.Message, error) {
	// translate to pb
	// TODO: throw out repo.Block
	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}
	pblock := &pb.Block{
		Id:           block.Id,
		Target:       block.Target,
		Parents:      block.Parents,
		TargetKey:    block.TargetKey,
		ThreadPubKey: block.ThreadPubKey,
		Type:         pb.Block_Type(int32(block.Type)),
		Date:         date,
	}

	// sign it
	serialized, err := proto.Marshal(pblock)
	if err != nil {
		return nil, err
	}
	signature, err := t.PrivKey.Sign(serialized)
	if err != nil {
		return nil, err
	}
	pkb, err := t.ipfs().PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	signed := &pb.SignedBlock{
		Id:           block.Id,
		Data:         serialized,
		Signature:    signature,
		ThreadId:     t.Id,
		ThreadName:   t.Name,
		IssuerPubKey: pkb,
	}

	// create the message
	payload, err := ptypes.MarshalAny(signed)
	if err != nil {
		return nil, err
	}
	return &pb.Message{
		MessageType: pb.Message_THREAD_BLOCK,
		Payload:     payload,
	}, nil
}

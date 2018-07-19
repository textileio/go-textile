package thread

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/util"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"strings"
	"sync"
	"time"
)

var log = logging.MustGetLogger("thread")

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
	Index      repo.Block `json:"index"`
	ThreadId   string     `json:"thread_id"`
	ThreadName string     `json:"thread_name"`
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

// Updates returns a read-only channel of updates
func (t *Thread) Updates() <-chan Update {
	return t.updates
}

// Close shutsdown the update channel
func (t *Thread) Close() {
	close(t.updates)
}

// Blocks paginates blocks from the datastore
func (t *Thread) Blocks(offsetId string, limit int, bType repo.BlockType) []repo.Block {
	query := fmt.Sprintf("threadId='%s' and type=%d", t.Id, bType)
	return t.blocks().List(offsetId, limit, query)
}

// Peers returns locally known peers in this thread
func (t *Thread) Peers() []repo.Peer {
	query := fmt.Sprintf("threadId='%s'", t.Id)
	return t.peers().List("", -1, query)
}

// Encrypt data with thread public key
func (t *Thread) Encrypt(data []byte) ([]byte, error) {
	return crypto.Encrypt(t.PrivKey.GetPublic(), data)
}

// Decrypt data with thread secret key
func (t *Thread) Decrypt(data []byte) ([]byte, error) {
	return crypto.Decrypt(t.PrivKey, data)
}

// Verify verifies a signed block
func (t *Thread) Verify(signed *pb.SignedThreadBlock) error {
	return crypto.Verify(t.PrivKey.GetPublic(), signed.Block, signed.ThreadSig)
}

// verifyAuthor checks the signature in the header
func (t *Thread) verifyAuthor(signed *pb.SignedThreadBlock, header *pb.ThreadBlockHeader) error {
	authorPk, err := libp2pc.UnmarshalPublicKey(header.AuthorPk)
	if err != nil {
		return err
	}
	if err := crypto.Verify(authorPk, signed.Block, signed.AuthorSig); err != nil {
		return err
	}
	return nil
}

// newBlockHeader creates a new header
func (t *Thread) newBlockHeader(date time.Time) (*pb.ThreadBlockHeader, error) {
	// get current HEAD
	head, err := t.GetHead()
	if err != nil {
		return nil, err
	}

	// get our own public key
	threadPk, err := t.PrivKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// get our own public key
	authorPk, err := t.ipfs().PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// get now date
	pdate, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	return &pb.ThreadBlockHeader{
		Date:     pdate,
		Parents:  strings.Split(string(head), ","),
		ThreadPk: threadPk,
		AuthorPk: authorPk,
	}, nil
}

// addBlock adds to ipfs
func (t *Thread) addBlock(message *pb.Message) (mh.Multihash, error) {
	// marshal to bytes
	messageb, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	// pin it
	return util.PinData(t.ipfs(), bytes.NewReader(messageb))
}

// commitBlock seals and signs the content of a block and adds it to ipfs
func (t *Thread) commitBlock(content proto.Message, mt pb.Message_Type) (*pb.Message, mh.Multihash, error) {
	// sign it
	serialized, err := proto.Marshal(content)
	if err != nil {
		return nil, nil, err
	}
	threadSig, err := t.PrivKey.Sign(serialized)
	if err != nil {
		return nil, nil, err
	}
	authorSig, err := t.ipfs().PrivateKey.Sign(serialized)
	if err != nil {
		return nil, nil, err
	}
	signed := &pb.SignedThreadBlock{
		Block:     serialized,
		ThreadSig: threadSig,
		AuthorSig: authorSig,
	}

	// create the message
	payload, err := ptypes.MarshalAny(signed)
	if err != nil {
		return nil, nil, err
	}
	message := &pb.Message{Type: mt, Payload: payload}

	// add to ipfs
	addr, err := t.addBlock(message)
	if err != nil {
		return nil, nil, err
	}

	return message, addr, nil
}

// indexBlock stores off index info for this block type
func (t *Thread) indexBlock(id string, header *pb.ThreadBlockHeader, blockType repo.BlockType, dataConf *repo.DataBlockConfig) error {
	// add a new one
	date, err := ptypes.Timestamp(header.Date)
	if err != nil {
		return err
	}
	if dataConf == nil {
		dataConf = new(repo.DataBlockConfig)
	}
	index := &repo.Block{
		Id:       id,
		Date:     date,
		Parents:  header.Parents,
		ThreadId: libp2pc.ConfigEncodeKey(header.ThreadPk),
		AuthorPk: libp2pc.ConfigEncodeKey(header.AuthorPk),
		Type:     blockType,

		// off-chain data links
		DataId:            dataConf.DataId,
		DataKeyCipher:     dataConf.DataKeyCipher,
		DataCaptionCipher: dataConf.DataCaptionCipher,
	}
	if err := t.blocks().Add(index); err != nil {
		return err
	}

	// update current head
	if err := t.updateHead(index.Id); err != nil {
		return err
	}

	// notify listeners
	t.pushUpdate(*index)

	return nil
}

// followParents tries to follow a chain of block ids, processing along the way
func (t *Thread) followParents(parents []string) error {
	// TODO: follow all parent paths
	parent := parents[0]

	// first update?
	if parent == "" {
		log.Debugf("found genesis block, aborting")
		return nil
	}

	// check if we aleady have this block indexed
	index := t.blocks().Get(parent)
	if index != nil {
		log.Debugf("block %s exists, aborting", parent)
		return nil
	}

	// download it
	serialized, err := util.GetDataAtPath(t.ipfs(), parent)
	if err != nil {
		return err
	}
	message := new(pb.Message)
	if err := proto.Unmarshal(serialized, message); err != nil {
		return err
	}
	signed := new(pb.SignedThreadBlock)
	err = ptypes.UnmarshalAny(message.Payload, signed)
	if err != nil {
		return err
	}

	// verify thread sig
	if err := t.Verify(signed); err != nil {
		return err
	}

	// handle each type
	switch message.Type {
	case pb.Message_THREAD_INVITE:
		if _, err = t.HandleInviteBlock(message, signed, nil); err != nil {
			return err
		}
	case pb.Message_THREAD_EXTERNAL_INVITE:
		if _, err = t.HandleInviteBlock(message, signed, nil); err != nil {
			return err
		}
	case pb.Message_THREAD_JOIN:
		if _, err = t.HandleInviteBlock(message, signed, nil); err != nil {
			return err
		}
	case pb.Message_THREAD_LEAVE:
		if _, err = t.HandleInviteBlock(message, signed, nil); err != nil {
			return err
		}
	case pb.Message_THREAD_DATA:
		if _, err = t.HandleInviteBlock(message, signed, nil); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("invalid message type: %s", message.Type))
	}
	return nil
}

// post publishes a message with content id to peers
func (t *Thread) post(message *pb.Message, id string, peers []repo.Peer) error {
	if len(peers) == 0 {
		return nil
	}
	log.Debugf("posting %s in thread %s...", id, t.Name)
	wg := sync.WaitGroup{}
	for _, p := range peers {
		wg.Add(1)
		go func(peerId string) {
			if err := t.send(message, peerId); err != nil {
				log.Errorf("error sending block %s to peer %s: %s", id, peerId, err)
			}
			wg.Done()
		}(p.Id)
	}
	wg.Wait()
	log.Debugf("posted to %d peers", len(peers))
	return nil
}

// pushUpdate pushes thread updates to UI listeners
func (t *Thread) pushUpdate(index repo.Block) {
	defer func() {
		if recover() != nil {
			log.Error("update channel closed")
		}
	}()
	select {
	case t.updates <- Update{
		Index:      index,
		ThreadId:   t.Id,
		ThreadName: t.Name,
	}:
	default:
	}
}

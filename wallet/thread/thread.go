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
	"github.com/textileio/textile-go/util"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"strings"
	"sync"
	"time"
)

var log = logging.MustGetLogger("thread")

// Update is used to notify listeners about updates in a thread
type Update struct {
	Block      repo.Block `json:"block"`
	ThreadId   string     `json:"thread_id"`
	ThreadName string     `json:"thread_name"`
}

// Info reports info about a thread
type Info struct {
	Head        *repo.Block `json:"head,omitempty"`
	BlockCount  int         `json:"block_count"`
	LatestPhoto *repo.Block `json:"latest_photo,omitempty"`
	PhotoCount  int         `json:"photo_count"`
}

// Config is used to construct a Thread
type Config struct {
	RepoPath      string
	Ipfs          func() *core.IpfsNode
	Blocks        func() repo.BlockStore
	Peers         func() repo.PeerStore
	Notifications func() repo.NotificationStore
	GetHead       func() (string, error)
	UpdateHead    func(head string) error
	Publish       func(payload []byte) error
	Send          func(message *pb.Envelope, peerId string, hash *string) error
	NewEnvelope   func(message *pb.Message) (*pb.Envelope, error)
	PutPinRequest func(id string) error
	GetUsername   func() (*string, error)
	SendUpdate    func(update Update)
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id            string
	Name          string
	PrivKey       libp2pc.PrivKey
	repoPath      string
	ipfs          func() *core.IpfsNode
	blocks        func() repo.BlockStore
	peers         func() repo.PeerStore
	notifications func() repo.NotificationStore
	GetHead       func() (string, error)
	updateHead    func(head string) error
	publish       func(payload []byte) error
	send          func(message *pb.Envelope, peerId string, hash *string) error
	newEnvelope   func(message *pb.Message) (*pb.Envelope, error)
	putPinRequest func(id string) error
	getUsername   func() (*string, error)
	sendUpdate    func(update Update)
	mux           sync.Mutex
}

// NewThread create a new Thread from a repo model and config
func NewThread(model *repo.Thread, config *Config) (*Thread, error) {
	sk, err := libp2pc.UnmarshalPrivateKey(model.PrivKey)
	if err != nil {
		return nil, err
	}
	return &Thread{
		Id:            model.Id,
		Name:          model.Name,
		PrivKey:       sk,
		repoPath:      config.RepoPath,
		ipfs:          config.Ipfs,
		blocks:        config.Blocks,
		peers:         config.Peers,
		notifications: config.Notifications,
		GetHead:       config.GetHead,
		updateHead:    config.UpdateHead,
		publish:       config.Publish,
		send:          config.Send,
		newEnvelope:   config.NewEnvelope,
		putPinRequest: config.PutPinRequest,
		getUsername:   config.GetUsername,
		sendUpdate:    config.SendUpdate,
	}, nil
}

// Info returns thread info
func (t *Thread) Info() (*Info, error) {
	// block info
	var head, latestPhoto *repo.Block
	headId, err := t.GetHead()
	if err != nil {
		return nil, err
	}
	if headId != "" {
		head = t.blocks().Get(headId)
	}
	blocks := t.blocks().Count(fmt.Sprintf("threadId='%s'", t.Id))

	// photo specific info
	query := fmt.Sprintf("threadId='%s' and type=%d", t.Id, repo.PhotoBlock)
	latestPhotos := t.blocks().List("", 1, query)
	if len(latestPhotos) > 0 {
		latestPhoto = &latestPhotos[0]
	}
	photos := t.blocks().Count(fmt.Sprintf("threadId='%s' and type=%d", t.Id, repo.PhotoBlock))

	return &Info{
		Head:        head,
		BlockCount:  blocks,
		LatestPhoto: latestPhoto,
		PhotoCount:  photos,
	}, nil
}

// Blocks paginates blocks from the datastore
func (t *Thread) Blocks(offsetId string, limit int, btype *repo.BlockType, dataId *string) []repo.Block {
	var query string
	if btype != nil {
		query = fmt.Sprintf("threadId='%s' and type=%d", t.Id, *btype)
	} else {
		query = fmt.Sprintf("threadId='%s'", t.Id)
	}
	if dataId != nil {
		query += fmt.Sprintf(" and dataId='%s'", *dataId)
	}
	all := t.blocks().List(offsetId, limit, query)
	if btype == nil {
		return all
	}
	var filtered []repo.Block
	for _, block := range all {
		ignored := t.blocks().GetByDataId(fmt.Sprintf("ignore-%s", block.Id))
		if ignored == nil {
			filtered = append(filtered, block)
		}
	}
	return filtered
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

// FollowParents tries to follow a list of chains of block ids, processing along the way
func (t *Thread) FollowParents(parents []string, from *peer.ID) ([]repo.Peer, error) {
	var joins []repo.Peer
	for _, parent := range parents {
		joined, err := t.followParent(parent, from)
		if err != nil {
			return nil, err
		}
		if joined != nil {
			joins = append(joins, *joined)
		}
	}
	return joins, nil
}

// followParent tries to follow a chain of block ids, processing along the way
func (t *Thread) followParent(parent string, from *peer.ID) (*repo.Peer, error) {
	// first update?
	if parent == "" {
		log.Debugf("found genesis block, aborting")
		return nil, nil
	}

	// check if we aleady have this block indexed
	index := t.blocks().Get(parent)
	if index != nil {
		log.Debugf("block %s exists, aborting", parent)
		return nil, nil
	}

	// download it
	serialized, err := util.GetDataAtPath(t.ipfs(), parent)
	if err != nil {
		return nil, err
	}
	env := new(pb.Envelope)
	message := new(pb.Message)
	if err := proto.Unmarshal(serialized, env); err != nil {
		return nil, err
	}
	if env.Message != nil {
		// verify author sig
		messageb, err := proto.Marshal(env.Message)
		if err != nil {
			return nil, err
		}
		authorPk, err := libp2pc.UnmarshalPublicKey(env.Pk)
		if err != nil {
			return nil, err
		}
		if err := crypto.Verify(authorPk, messageb, env.Sig); err != nil {
			return nil, err
		}
		message = env.Message
	} else {
		// might be a merge block
		if err := proto.Unmarshal(serialized, message); err != nil {
			return nil, err
		}
	}
	if message.Payload == nil {
		return nil, errors.New("nil message payload")
	}

	// verify thread sig
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(message.Payload, signed); err != nil {
		return nil, err
	}
	if err := t.Verify(signed); err != nil {
		return nil, err
	}

	// handle each type
	var joined *repo.Peer
	switch message.Type {
	case pb.Message_THREAD_JOIN:
		var err error
		_, joined, err = t.HandleJoinBlock(from, env, signed, nil, true)
		if err != nil {
			return nil, err
		}
	case pb.Message_THREAD_LEAVE:
		if _, err := t.HandleLeaveBlock(from, env, signed, nil, true); err != nil {
			return nil, err
		}
	case pb.Message_THREAD_DATA:
		if _, err := t.HandleDataBlock(from, env, signed, nil, true); err != nil {
			return nil, err
		}
	case pb.Message_THREAD_ANNOTATION:
		if _, err := t.HandleAnnotationBlock(from, env, signed, nil, true); err != nil {
			return nil, err
		}
	case pb.Message_THREAD_IGNORE:
		if _, err := t.HandleIgnoreBlock(from, env, signed, nil, true); err != nil {
			return nil, err
		}
	case pb.Message_THREAD_MERGE:
		if _, err := t.HandleMergeBlock(from, message, signed, nil, true); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("invalid message type: %s", message.Type))
	}
	return joined, nil
}

// newBlockHeader creates a new header
func (t *Thread) newBlockHeader() (*pb.ThreadBlockHeader, error) {
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

	// encrypt our own username with thread pk
	var authorUnCipher []byte
	authorUn, err := t.getUsername()
	if err != nil {
		return nil, err
	}
	if authorUn != nil {
		authorUnCipher, err = t.Encrypt([]byte(*authorUn))
		if err != nil {
			return nil, err
		}
	}

	// get now date
	pdate, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	return &pb.ThreadBlockHeader{
		Date:           pdate,
		Parents:        strings.Split(string(head), ","),
		ThreadPk:       threadPk,
		AuthorPk:       authorPk,
		AuthorUnCipher: authorUnCipher,
	}, nil
}

// addBlock adds to ipfs
func (t *Thread) addBlock(envelope *pb.Envelope) (mh.Multihash, error) {
	// marshal to bytes
	messageb, err := proto.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	// pin it
	id, err := util.PinData(t.ipfs(), bytes.NewReader(messageb))
	if err != nil {
		return nil, err
	}

	// add a pin request
	if err := t.putPinRequest(id.Hash().B58String()); err != nil {
		log.Warningf("pin request exists: %s", id.Hash().B58String())
	}

	return id.Hash(), nil
}

// commitBlock seals and signs the content of a block and adds it to ipfs
func (t *Thread) commitBlock(content proto.Message, mt pb.Message_Type) (*pb.Envelope, mh.Multihash, error) {
	// sign it
	serializedContent, err := proto.Marshal(content)
	if err != nil {
		return nil, nil, err
	}
	threadSig, err := t.PrivKey.Sign(serializedContent)
	if err != nil {
		return nil, nil, err
	}
	signed := &pb.SignedThreadBlock{
		Block:     serializedContent,
		ThreadSig: threadSig,
	}

	// create the message
	payload, err := ptypes.MarshalAny(signed)
	if err != nil {
		return nil, nil, err
	}
	message := &pb.Message{Type: mt, Payload: payload}
	envelope, err := t.newEnvelope(message)
	if err != nil {
		return nil, nil, err
	}

	// add to ipfs
	addr, err := t.addBlock(envelope)
	if err != nil {
		return nil, nil, err
	}

	return envelope, addr, nil
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
		Id:                   id,
		Date:                 date,
		Parents:              header.Parents,
		ThreadId:             libp2pc.ConfigEncodeKey(header.ThreadPk),
		AuthorPk:             libp2pc.ConfigEncodeKey(header.AuthorPk),
		AuthorUsernameCipher: header.AuthorUnCipher,
		Type:                 blockType,

		// off-chain data links
		DataId:             dataConf.DataId,
		DataKeyCipher:      dataConf.DataKeyCipher,
		DataCaptionCipher:  dataConf.DataCaptionCipher,
		DataMetadataCipher: dataConf.DataMetadataCipher,
	}
	if err := t.blocks().Add(index); err != nil {
		return err
	}

	// notify listeners
	t.pushUpdate(*index)

	return nil
}

// handleHead determines whether or not a thread can be fast-forwarded or if a merge block is needed
// - parents are the parents of the incoming chain
func (t *Thread) handleHead(inboundId string, parents []string) (mh.Multihash, error) {
	// get current HEAD
	head, err := t.GetHead()
	if err != nil {
		return nil, err
	}

	// fast-forward is possible if current HEAD is equal to one of the incoming parents
	var fastForwardable bool
	if head == "" {
		fastForwardable = true
	} else {
		for _, parent := range parents {
			if head == parent {
				fastForwardable = true
			}
		}
	}
	if fastForwardable {
		// no need for a merge
		log.Debugf("fast-forwarded to %s", inboundId)
		if err := t.updateHead(inboundId); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// needs merge
	return t.Merge(inboundId)
}

// post publishes a message with content id to peers
func (t *Thread) post(env *pb.Envelope, id string, peers []repo.Peer) {
	if len(peers) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	for _, p := range peers {
		wg.Add(1)
		go func(peerId string) {
			if err := t.send(env, peerId, &id); err != nil {
				log.Errorf("error sending block %s to peer %s: %s", id, peerId, err)
			}
			wg.Done()
		}(p.Id)
	}
	wg.Wait()
}

// pushUpdate pushes thread updates to UI listeners
func (t *Thread) pushUpdate(index repo.Block) {
	t.sendUpdate(Update{
		Block:      index,
		ThreadId:   t.Id,
		ThreadName: t.Name,
	})
}

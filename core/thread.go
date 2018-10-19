package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"strings"
	"sync"
	"time"
)

// ThreadUpdate is used to notify listeners about updates in a thread
type ThreadUpdate struct {
	Block      repo.Block `json:"block"`
	ThreadId   string     `json:"thread_id"`
	ThreadName string     `json:"thread_name"`
}

// ThreadInfo reports info about a thread
type ThreadInfo struct {
	Head        *repo.Block `json:"head,omitempty"`
	BlockCount  int         `json:"block_count"`
	LatestPhoto *repo.Block `json:"latest_photo,omitempty"`
	PhotoCount  int         `json:"photo_count"`
}

// ThreadConfig is used to construct a Thread
type ThreadConfig struct {
	RepoPath   string
	Node       func() *core.IpfsNode
	Datastore  repo.Datastore
	Service    *ThreadsService
	CafeQueue  *CafeRequestQueue
	SendUpdate func(update ThreadUpdate)
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id         string
	Name       string
	privKey    libp2pc.PrivKey
	repoPath   string
	node       func() *core.IpfsNode
	datastore  repo.Datastore
	service    *ThreadsService
	cafeQueue  *CafeRequestQueue
	sendUpdate func(update ThreadUpdate)
	mux        sync.Mutex
}

// NewThread create a new Thread from a repo model and config
func NewThread(model *repo.Thread, config *ThreadConfig) (*Thread, error) {
	sk, err := libp2pc.UnmarshalPrivateKey(model.PrivKey)
	if err != nil {
		return nil, err
	}
	return &Thread{
		Id:         model.Id,
		Name:       model.Name,
		privKey:    sk,
		repoPath:   config.RepoPath,
		node:       config.Node,
		datastore:  config.Datastore,
		service:    config.Service,
		cafeQueue:  config.CafeQueue,
		sendUpdate: config.SendUpdate,
	}, nil
}

// Info returns thread info
func (t *Thread) Info() (*ThreadInfo, error) {
	// block info
	var head, latestPhoto *repo.Block
	headId, err := t.Head()
	if err != nil {
		return nil, err
	}
	if headId != "" {
		head = t.datastore.Blocks().Get(headId)
	}
	blocks := t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s'", t.Id))

	// photo specific info
	query := fmt.Sprintf("threadId='%s' and type=%d", t.Id, repo.PhotoBlock)
	latestPhotos := t.datastore.Blocks().List("", 1, query)
	if len(latestPhotos) > 0 {
		latestPhoto = &latestPhotos[0]
	}
	photos := t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s' and type=%d", t.Id, repo.PhotoBlock))

	// send back summary
	return &ThreadInfo{
		Head:        head,
		BlockCount:  blocks,
		LatestPhoto: latestPhoto,
		PhotoCount:  photos,
	}, nil
}

// Head returns content id of the latest update
func (t *Thread) Head() (string, error) {
	mod := t.datastore.Threads().Get(t.Id)
	if mod == nil {
		return "", errors.New(fmt.Sprintf("could not re-load thread: %s", t.Id))
	}
	return mod.Head, nil
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
	all := t.datastore.Blocks().List(offsetId, limit, query)
	if btype == nil {
		return all
	}
	var filtered []repo.Block
	for _, block := range all {
		ignored := t.datastore.Blocks().GetByData(fmt.Sprintf("ignore-%s", block.Id))
		if ignored == nil {
			filtered = append(filtered, block)
		}
	}
	return filtered
}

// Peers returns locally known peers in this thread
func (t *Thread) Peers() []repo.ThreadPeer {
	return t.datastore.ThreadPeers().ListByThread(t.Id)
}

// Encrypt data with thread public key
func (t *Thread) Encrypt(data []byte) ([]byte, error) {
	return crypto.Encrypt(t.privKey.GetPublic(), data)
}

// Decrypt data with thread secret key
func (t *Thread) Decrypt(data []byte) ([]byte, error) {
	return crypto.Decrypt(t.privKey, data)
}

// Verify verifies a signed block
func (t *Thread) Verify(signed *pb.SignedThreadBlock) error {
	return crypto.Verify(t.privKey.GetPublic(), signed.Block, signed.ThreadSig)
}

// FollowParents tries to follow a list of chains of block ids, processing along the way
func (t *Thread) FollowParents(parents []string, from *peer.ID) ([]repo.ThreadPeer, error) {
	var joins []repo.ThreadPeer
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
func (t *Thread) followParent(parent string, from *peer.ID) (*repo.ThreadPeer, error) {
	// first update?
	if parent == "" {
		log.Debugf("found genesis block, aborting")
		return nil, nil
	}

	// check if we aleady have this block indexed
	index := t.datastore.Blocks().Get(parent)
	if index != nil {
		log.Debugf("block %s exists, aborting", parent)
		return nil, nil
	}

	// download it
	serialized, err := ipfs.GetDataAtPath(t.node(), parent)
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
	var joined *repo.ThreadPeer
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
	head, err := t.Head()
	if err != nil {
		return nil, err
	}

	// get our own public key
	threadPk, err := t.privKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// get our own public key
	authorPk, err := t.node().PrivateKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// encrypt our own username with thread pk
	var authorUnCipher []byte
	authorUn, _ := t.datastore.Profile().GetUsername()
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
	id, err := ipfs.PinData(t.node(), bytes.NewReader(messageb))
	if err != nil {
		return nil, err
	}

	// add a store request
	t.cafeQueue.Add(id.Hash().B58String(), repo.CafeStoreRequest)

	return id.Hash(), nil
}

// commitBlock seals and signs the content of a block and adds it to ipfs
func (t *Thread) commitBlock(content proto.Message, mtype pb.Message_Type) (*pb.Envelope, mh.Multihash, error) {
	// create the block
	env, err := t.service.NewBlock(t.privKey, mtype, content)
	if err != nil {
		return nil, nil, err
	}

	// add to ipfs
	addr, err := t.addBlock(env)
	if err != nil {
		return nil, nil, err
	}

	return env, addr, nil
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
	threadId, err := ipfs.IDFromPublicKeyBytes(header.ThreadPk)
	if err != nil {
		return err
	}
	authorId, err := ipfs.IDFromPublicKeyBytes(header.AuthorPk)
	if err != nil {
		return err
	}
	index := &repo.Block{
		Id:                   id,
		Date:                 date,
		Parents:              header.Parents,
		ThreadId:             threadId.Pretty(),
		AuthorId:             authorId.Pretty(),
		AuthorUsernameCipher: header.AuthorUnCipher,
		Type:                 blockType,

		// off-chain data links
		DataId:             dataConf.DataId,
		DataKeyCipher:      dataConf.DataKeyCipher,
		DataCaptionCipher:  dataConf.DataCaptionCipher,
		DataMetadataCipher: dataConf.DataMetadataCipher,
	}
	if err := t.datastore.Blocks().Add(index); err != nil {
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
	head, err := t.Head()
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

// updateHead updates the ref to the content id of the latest update
func (t *Thread) updateHead(head string) error {
	if err := t.datastore.Threads().UpdateHead(t.Id, head); err != nil {
		return err
	}

	// update head on cafe backups
	t.cafeQueue.Add(t.Id, repo.CafeStoreThreadRequest)
	return nil
}

// post publishes a message with content id to peers
func (t *Thread) post(env *pb.Envelope, id string, peers []repo.ThreadPeer) {
	if len(peers) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	for _, tp := range peers {
		wg.Add(1)
		go func(tp repo.ThreadPeer) {
			if err := t.sendMessage(&tp, env); err != nil {
				log.Errorf("error sending block %s to peer %s: %s", id, tp.Id, err)
			}
			wg.Done()
		}(tp)
	}
	wg.Wait()
}

// sendMessage sends a message directly to a peer, falling back to cafe inbox based delivery
func (t *Thread) sendMessage(tpeer *repo.ThreadPeer, env *pb.Envelope) error {
	//t.threadsService.SendMessage
	//t.cafeService.DeliverMessage
	return nil
}

// pushUpdate pushes thread updates to UI listeners
func (t *Thread) pushUpdate(index repo.Block) {
	t.sendUpdate(ThreadUpdate{
		Block:      index,
		ThreadId:   t.Id,
		ThreadName: t.Name,
	})
}

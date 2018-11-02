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
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Head       *repo.Block `json:"head,omitempty"`
	PeerCount  int         `json:"peer_cnt"`
	BlockCount int         `json:"block_cnt"`
	FileCount  int         `json:"file_cnt"`
}

// ThreadConfig is used to construct a Thread
type ThreadConfig struct {
	RepoPath      string
	Node          func() *core.IpfsNode
	Datastore     repo.Datastore
	Service       func() *ThreadsService
	ThreadsOutbox *ThreadsOutbox
	CafeOutbox    *CafeOutbox
	SendUpdate    func(update ThreadUpdate)
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id            string
	Name          string
	privKey       libp2pc.PrivKey
	repoPath      string
	node          func() *core.IpfsNode
	datastore     repo.Datastore
	service       func() *ThreadsService
	threadsOutbox *ThreadsOutbox
	cafeOutbox    *CafeOutbox
	sendUpdate    func(update ThreadUpdate)
	mux           sync.Mutex
}

// NewThread create a new Thread from a repo model and config
func NewThread(model *repo.Thread, config *ThreadConfig) (*Thread, error) {
	sk, err := libp2pc.UnmarshalPrivateKey(model.PrivKey)
	if err != nil {
		return nil, err
	}
	return &Thread{
		Id:            model.Id,
		Name:          model.Name,
		privKey:       sk,
		repoPath:      config.RepoPath,
		node:          config.Node,
		datastore:     config.Datastore,
		service:       config.Service,
		threadsOutbox: config.ThreadsOutbox,
		cafeOutbox:    config.CafeOutbox,
		sendUpdate:    config.SendUpdate,
	}, nil
}

// Info returns thread info
func (t *Thread) Info() (*ThreadInfo, error) {
	// block info
	var head *repo.Block
	headId, err := t.Head()
	if err != nil {
		return nil, err
	}
	if headId != "" {
		head = t.datastore.Blocks().Get(headId)
	}
	blocks := t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s'", t.Id))
	files := t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s' and type=%d", t.Id, repo.PhotoBlock))

	// send back summary
	return &ThreadInfo{
		Id:         t.Id,
		Name:       t.Name,
		Head:       head,
		PeerCount:  len(t.Peers()) + 1,
		BlockCount: blocks,
		FileCount:  files,
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

// followParents tries to follow a list of chains of block ids, processing along the way
func (t *Thread) followParents(parents []string) error {
	for _, parent := range parents {
		if parent == "" {
			log.Debugf("found genesis block, aborting")
			continue
		}
		hash, err := mh.FromB58String(parent)
		if err != nil {
			return err
		}
		if err := t.followParent(hash); err != nil {
			log.Errorf("failed to follow parent %s: %s", parent, err)
			continue
		}
	}
	return nil
}

// followParent tries to follow a chain of block ids, processing along the way
func (t *Thread) followParent(parent mh.Multihash) error {
	// download and decrypt
	ciphertext, err := ipfs.DataAtPath(t.node(), parent.B58String())
	if err != nil {
		return err
	}
	block, err := t.handleBlock(parent, ciphertext)
	if err != nil {
		return err
	}
	if block == nil {
		// exists, abort
		return nil
	}

	// handle each type
	switch block.Type {
	case pb.ThreadBlock_JOIN:
		_, err = t.handleJoinBlock(parent, block)
	case pb.ThreadBlock_LEAVE:
		err = t.handleLeaveBlock(parent, block)
	case pb.ThreadBlock_DATA:
		_, err = t.handleDataBlock(parent, block)
	case pb.ThreadBlock_ANNOTATION:
		_, err = t.handleAnnotationBlock(parent, block)
	case pb.ThreadBlock_IGNORE:
		_, err = t.handleIgnoreBlock(parent, block)
	case pb.ThreadBlock_MERGE:
		err = t.handleMergeBlock(parent, block)
	default:
		return errors.New(fmt.Sprintf("invalid message type: %s", block.Type))
	}
	if err != nil {
		return err
	}

	// back prop
	return t.followParents(block.Header.Parents)
}

// addOrUpdatePeer collects thread peers, saving them as contacts and
// saving their cafe inboxes for offline message delivery
func (t *Thread) addOrUpdatePeer(pid peer.ID, username string, inboxes []string) {
	// add to this thread's peer list
	t.datastore.ThreadPeers().Add(&repo.ThreadPeer{
		Id:       pid.Pretty(),
		ThreadId: t.Id,
		Welcomed: false,
	})

	// add to contacts
	t.datastore.Contacts().AddOrUpdate(&repo.Contact{
		Id:       pid.Pretty(),
		Username: username,
		Inboxes:  inboxes,
		Added:    time.Now(),
	})
}

// newBlockHeader creates a new header
func (t *Thread) newBlockHeader() (*pb.ThreadBlockHeader, error) {
	// get current HEAD
	head, err := t.Head()
	if err != nil {
		return nil, err
	}

	// build the header
	pdate, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}
	var parents []string
	if head != "" {
		parents = strings.Split(head, ",")
	}
	return &pb.ThreadBlockHeader{
		Date:    pdate,
		Parents: parents,
		Author:  t.node().Identity.Pretty(),
	}, nil
}

// commitResult wraps the results of a block commit
type commitResult struct {
	hash       mh.Multihash
	ciphertext []byte
	header     *pb.ThreadBlockHeader
}

// commitBlock encrypts a block with thread key (or custom method if provided) and adds it to ipfs
func (t *Thread) commitBlock(msg proto.Message, mtype pb.ThreadBlock_Type, encrypt func(plaintext []byte) ([]byte, error)) (*commitResult, error) {
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	block := &pb.ThreadBlock{
		Header: header,
		Type:   mtype,
	}
	if msg != nil {
		payload, err := ptypes.MarshalAny(msg)
		if err != nil {
			return nil, err
		}
		block.Payload = payload
	}
	plaintext, err := proto.Marshal(block)
	if err != nil {
		return nil, err
	}

	// encrypt, falling back to thread key
	if encrypt == nil {
		encrypt = t.Encrypt
	}
	ciphertext, err := encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	// add to ipfs
	hash, err := t.addBlock(ciphertext)
	if err != nil {
		return nil, err
	}

	return &commitResult{hash, ciphertext, header}, nil
}

// addBlock adds to ipfs
func (t *Thread) addBlock(ciphertext []byte) (mh.Multihash, error) {
	// pin it
	id, err := ipfs.PinData(t.node(), bytes.NewReader(ciphertext))
	if err != nil {
		return nil, err
	}

	// add a store request
	t.cafeOutbox.Add(id.Hash().B58String(), repo.CafeStoreRequest)

	return id.Hash(), nil
}

// handleBlock receives an incoming encrypted block
func (t *Thread) handleBlock(hash mh.Multihash, ciphertext []byte) (*pb.ThreadBlock, error) {
	// check if we aleady have this block indexed
	index := t.datastore.Blocks().Get(hash.B58String())
	if index != nil {
		return nil, nil
	}

	// decrypt
	block := new(pb.ThreadBlock)
	plaintext, err := t.Decrypt(ciphertext)
	if err != nil {
		// might be a merge block
		err2 := proto.Unmarshal(ciphertext, block)
		if err2 != nil || block.Type != pb.ThreadBlock_MERGE {
			return nil, err
		}
	} else {
		if err := proto.Unmarshal(plaintext, block); err != nil {
			return nil, err
		}
	}

	// nil payload only allowed for some types
	if block.Payload == nil && block.Type != pb.ThreadBlock_MERGE && block.Type != pb.ThreadBlock_LEAVE {
		return nil, errors.New("nil message payload")
	}

	// add to ipfs
	if _, err := t.addBlock(ciphertext); err != nil {
		return nil, err
	}
	return block, nil
}

// indexBlock stores off index info for this block type
func (t *Thread) indexBlock(commit *commitResult, blockType repo.BlockType, dataConf *repo.DataBlockConfig) error {
	// add a new one
	date, err := ptypes.Timestamp(commit.header.Date)
	if err != nil {
		return err
	}
	if dataConf == nil {
		dataConf = new(repo.DataBlockConfig)
	}
	index := &repo.Block{
		Id:       commit.hash.B58String(),
		Date:     date,
		Parents:  commit.header.Parents,
		ThreadId: t.Id,
		AuthorId: commit.header.Author,
		Type:     blockType,

		// off-chain data links
		DataId:       dataConf.DataId,
		DataKey:      dataConf.DataKey,
		DataCaption:  dataConf.DataCaption,
		DataMetadata: dataConf.DataMetadata,
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
func (t *Thread) handleHead(inbound mh.Multihash, parents []string) (mh.Multihash, error) {
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
		log.Debugf("fast-forwarded to %s", inbound.B58String())
		if err := t.updateHead(inbound); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// needs merge
	return t.Merge(inbound)
}

// updateHead updates the ref to the content id of the latest update
func (t *Thread) updateHead(head mh.Multihash) error {
	if err := t.datastore.Threads().UpdateHead(t.Id, head.B58String()); err != nil {
		return err
	}

	// update head on cafe backups
	t.cafeOutbox.Add(t.Id, repo.CafeStoreThreadRequest)

	return nil
}

// sendWelcome sends the latest HEAD block to a set of peers
func (t *Thread) sendWelcome() error {
	// get unwelcomed peers
	peers := t.datastore.ThreadPeers().ListUnwelcomedByThread(t.Id)
	if len(peers) == 0 {
		return nil
	}

	// get current HEAD
	head, err := t.Head()
	if err != nil {
		return err
	}
	if head == "" {
		return nil
	}

	// download it
	ciphertext, err := ipfs.DataAtPath(t.node(), head)
	if err != nil {
		return err
	}

	// post it
	hash, err := mh.FromB58String(head)
	if err != nil {
		return err
	}
	res := &commitResult{hash: hash, ciphertext: ciphertext}
	if err := t.post(res, peers); err != nil {
		return err
	}

	// mark each as welcomed
	if err := t.datastore.ThreadPeers().WelcomeByThread(t.Id); err != nil {
		return err
	}
	for _, p := range peers {
		log.Debugf("WELCOME sent to %s at %s", p.Id, head)
	}
	return nil
}

// post publishes an encrypted message to thread peers
func (t *Thread) post(commit *commitResult, peers []repo.ThreadPeer) error {
	if len(peers) == 0 {
		// flush the storage queue—this is normally done in a thread
		// via thread message queue handling, but that won't run if there's
		// no peers to send the message to.
		t.cafeOutbox.Flush()
		return nil
	}
	env, err := t.service().NewEnvelope(t.Id, commit.hash, commit.ciphertext)
	if err != nil {
		return err
	}
	for _, tp := range peers {
		pid, err := peer.IDB58Decode(tp.Id)
		if err != nil {
			return err
		}
		if err := t.threadsOutbox.Add(pid, env); err != nil {
			return err
		}
	}

	// flush the message queue
	go t.threadsOutbox.Flush()

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

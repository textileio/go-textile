package core

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	libp2pc "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
	mh "gx/ipfs/QmerPMzPk1mJVowm8KgmoknWa4yCYvvugMPsgWmDNUvDLW/go-multihash"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/repo/config"
)

// ErrContactNotFound indicates a local contact was not found
var ErrContactNotFound = errors.New("contact not found")

// ErrNotShareable indicates the thread does not allow invites, at least for _you_
var ErrNotShareable = errors.New("thread is not shareable")

// ErrNotReadable indicates the thread is not readable
var ErrNotReadable = errors.New("thread is not readable")

// ErrNotAnnotatable indicates the thread is not annotatable (comments/likes)
var ErrNotAnnotatable = errors.New("thread is not annotatable")

// ErrNotWritable indicates the thread is not writable (files/messages)
var ErrNotWritable = errors.New("thread is not writable")

// ErrThreadSchemaRequired indicates files where added without a thread schema
var ErrThreadSchemaRequired = errors.New("thread schema required to add files")

// ErrJsonSchemaRequired indicates json files where added without a json schema
var ErrJsonSchemaRequired = errors.New("thread schema does not allow json files")

// ErrInvalidFileNode indicates files where added via a nil ipld node
var ErrInvalidFileNode = errors.New("invalid files node")

// ErrBlockWrongType indicates a block was requested as a type other than its own
var ErrBlockWrongType = errors.New("block type is not the type requested")

// errReloadFailed indicates an error occurred during thread reload
var errThreadReload = errors.New("could not re-load thread")

// ThreadConfig is used to construct a Thread
type ThreadConfig struct {
	RepoPath    string
	Config      *config.Config
	Account     *keypair.Full
	Node        func() *core.IpfsNode
	Datastore   repo.Datastore
	Service     func() *ThreadsService
	BlockOutbox *BlockOutbox
	CafeOutbox  *CafeOutbox
	AddContact  func(*pb.Contact) error
	PushUpdate  func(*pb.Block, string)
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id          string
	Key         string // app key, usually UUID
	Name        string
	Schema      *pb.Node
	schemaId    string
	initiator   string
	ttype       pb.Thread_Type
	sharing     pb.Thread_Sharing
	members     []string
	privKey     libp2pc.PrivKey
	repoPath    string
	config      *config.Config
	account     *keypair.Full
	node        func() *core.IpfsNode
	datastore   repo.Datastore
	service     func() *ThreadsService
	blockOutbox *BlockOutbox
	cafeOutbox  *CafeOutbox
	addContact  func(*pb.Contact) error
	pushUpdate  func(*pb.Block, string)
	mux         sync.Mutex
}

// NewThread create a new Thread from a repo model and config
func NewThread(model *pb.Thread, conf *ThreadConfig) (*Thread, error) {
	sk, err := ipfs.UnmarshalPrivateKey(model.Sk)
	if err != nil {
		return nil, err
	}

	var sch *pb.Node
	if model.Schema != "" {
		sch, err = loadSchema(conf.Node(), model.Schema)
		if err != nil {
			return nil, err
		}
	}

	return &Thread{
		Id:          model.Id,
		Key:         model.Key,
		Name:        model.Name,
		Schema:      sch,
		schemaId:    model.Schema,
		initiator:   model.Initiator,
		ttype:       model.Type,
		sharing:     model.Sharing,
		members:     model.Members,
		privKey:     sk,
		repoPath:    conf.RepoPath,
		config:      conf.Config,
		account:     conf.Account,
		node:        conf.Node,
		datastore:   conf.Datastore,
		service:     conf.Service,
		blockOutbox: conf.BlockOutbox,
		cafeOutbox:  conf.CafeOutbox,
		addContact:  conf.AddContact,
		pushUpdate:  conf.PushUpdate,
	}, nil
}

// Head returns content id of the latest update
func (t *Thread) Head() (string, error) {
	mod := t.datastore.Threads().Get(t.Id)
	if mod == nil {
		return "", errThreadReload
	}
	return mod.Head, nil
}

// Peers returns locally known peers in this thread
func (t *Thread) Peers() []pb.ThreadPeer {
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
			log.Warningf("failed to follow parent %s: %s", parent, err)
			continue
		}
	}

	return nil
}

// followParent tries to follow a chain of block ids, processing along the way
func (t *Thread) followParent(parent mh.Multihash) error {
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
		log.Debugf("%s exists, aborting", parent)
		return nil
	}

	log.Debugf("handling %s from %s", block.Type.String(), block.Header.Author)

	switch block.Type {
	case pb.Block_MERGE:
		err = t.handleMergeBlock(parent, block)
	case pb.Block_IGNORE:
		_, err = t.handleIgnoreBlock(parent, block)
	case pb.Block_FLAG:
		_, err = t.handleFlagBlock(parent, block)
	case pb.Block_JOIN:
		_, err = t.handleJoinBlock(parent, block)
	case pb.Block_ANNOUNCE:
		_, err = t.handleAnnounceBlock(parent, block)
	case pb.Block_LEAVE:
		err = t.handleLeaveBlock(parent, block)
	case pb.Block_TEXT:
		_, err = t.handleMessageBlock(parent, block)
	case pb.Block_FILES:
		_, err = t.handleFilesBlock(parent, block)
	case pb.Block_COMMENT:
		_, err = t.handleCommentBlock(parent, block)
	case pb.Block_LIKE:
		_, err = t.handleLikeBlock(parent, block)
	default:
		return errors.New(fmt.Sprintf("invalid message type: %s", block.Type))
	}
	if err != nil {
		return err
	}

	return t.followParents(block.Header.Parents)
}

// addOrUpdateContact collects thread peers and saves them as contacts
func (t *Thread) addOrUpdateContact(contact *pb.Contact) error {
	if contact.Id == t.node().Identity.Pretty() {
		return nil
	}

	if err := t.datastore.ThreadPeers().Add(&pb.ThreadPeer{
		Id:       contact.Id,
		Thread:   t.Id,
		Welcomed: false,
	}); err != nil {
		if !repo.ConflictError(err) {
			return err
		}
	}

	return t.addContact(contact)
}

// newBlockHeader creates a new header
func (t *Thread) newBlockHeader() (*pb.ThreadBlockHeader, error) {
	head, err := t.Head()
	if err != nil {
		return nil, err
	}

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
		Address: t.account.Address(),
	}, nil
}

// commitResult wraps the results of a block commit
type commitResult struct {
	hash       mh.Multihash
	ciphertext []byte
	header     *pb.ThreadBlockHeader
}

// commitBlock encrypts a block with thread key (or custom method if provided) and adds it to ipfs
func (t *Thread) commitBlock(
	msg proto.Message,
	mtype pb.Block_BlockType,
	encrypt func(plaintext []byte) ([]byte, error)) (*commitResult, error) {

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

	hash, err := t.addBlock(ciphertext)
	if err != nil {
		return nil, err
	}

	return &commitResult{hash, ciphertext, header}, nil
}

// addBlock adds to ipfs
func (t *Thread) addBlock(ciphertext []byte) (mh.Multihash, error) {
	id, err := ipfs.AddData(t.node(), bytes.NewReader(ciphertext), true)
	if err != nil {
		return nil, err
	}

	if err := t.cafeOutbox.Add(id.Hash().B58String(), pb.CafeRequest_STORE); err != nil {
		return nil, err
	}

	return id.Hash(), nil
}

// handleBlock receives an incoming encrypted block
func (t *Thread) handleBlock(hash mh.Multihash, ciphertext []byte) (*pb.ThreadBlock, error) {
	index := t.datastore.Blocks().Get(hash.B58String())
	if index != nil {
		return nil, nil
	}

	block := new(pb.ThreadBlock)
	plaintext, err := t.Decrypt(ciphertext)
	if err != nil {
		// might be a merge block
		err2 := proto.Unmarshal(ciphertext, block)
		if err2 != nil || block.Type != pb.Block_MERGE {
			return nil, err
		}
	} else {
		if err := proto.Unmarshal(plaintext, block); err != nil {
			return nil, err
		}
	}

	// nil payload only allowed for some types
	if block.Payload == nil && block.Type != pb.Block_MERGE && block.Type != pb.Block_LEAVE {
		return nil, errors.New("nil message payload")
	}

	if _, err := t.addBlock(ciphertext); err != nil {
		return nil, err
	}
	return block, nil
}

// indexBlock stores off index info for this block type
func (t *Thread) indexBlock(commit *commitResult, blockType pb.Block_BlockType, target string, body string) error {
	block := &pb.Block{
		Id:      commit.hash.B58String(),
		Type:    blockType,
		Date:    commit.header.Date,
		Parents: commit.header.Parents,
		Thread:  t.Id,
		Author:  commit.header.Author,
		Target:  target,
		Body:    body,
	}
	if err := t.datastore.Blocks().Add(block); err != nil {
		return err
	}

	t.pushUpdate(block, t.Key)
	return nil
}

// handleHead determines whether or not a thread can be fast-forwarded or if a merge block is needed
// - parents are the parents of the incoming chain
func (t *Thread) handleHead(inbound mh.Multihash, parents []string) (mh.Multihash, error) {
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
	return t.merge(inbound)
}

// updateHead updates the ref to the content id of the latest update
func (t *Thread) updateHead(head mh.Multihash) error {
	if err := t.datastore.Threads().UpdateHead(t.Id, head.B58String()); err != nil {
		return err
	}

	return t.cafeOutbox.Add(t.Id, pb.CafeRequest_STORE_THREAD)
}

// sendWelcome sends the latest HEAD block to a set of peers
func (t *Thread) sendWelcome() error {
	peers := t.datastore.ThreadPeers().ListUnwelcomedByThread(t.Id)
	if len(peers) == 0 {
		return nil
	}

	head, err := t.Head()
	if err != nil {
		return err
	}
	if head == "" {
		return nil
	}

	ciphertext, err := ipfs.DataAtPath(t.node(), head)
	if err != nil {
		return err
	}

	hash, err := mh.FromB58String(head)
	if err != nil {
		return err
	}
	res := &commitResult{hash: hash, ciphertext: ciphertext}
	if err := t.post(res, peers); err != nil {
		return err
	}

	if err := t.datastore.ThreadPeers().WelcomeByThread(t.Id); err != nil {
		return err
	}
	for _, p := range peers {
		log.Debugf("WELCOME sent to %s at %s", p.Id, head)
	}
	return nil
}

// post publishes an encrypted message to thread peers
func (t *Thread) post(commit *commitResult, peers []pb.ThreadPeer) error {
	if len(peers) == 0 {
		// flush the storage queue—this is normally done in a thread
		// via thread message queue handling, but that won't run if there's
		// no peers to send the message to.
		go t.cafeOutbox.Flush()
		return nil
	}

	// add account signature
	sig, err := t.account.Sign(commit.ciphertext)
	if err != nil {
		return err
	}
	env, err := t.service().NewEnvelope(t.Id, commit.hash, commit.ciphertext, sig)
	if err != nil {
		return err
	}

	for _, tp := range peers {
		pid, err := peer.IDB58Decode(tp.Id)
		if err != nil {
			return err
		}

		if err := t.blockOutbox.Add(pid, env); err != nil {
			return err
		}
	}

	go t.blockOutbox.Flush()

	return nil
}

// readable returns whether or not this thread is readable from the
// perspective of the given address
func (t *Thread) readable(addr string) bool {
	if addr == t.initiator {
		return true
	}
	switch t.ttype {
	case pb.Thread_PRIVATE:
		return false // should not happen
	case pb.Thread_READ_ONLY:
		return t.member(addr)
	case pb.Thread_PUBLIC:
		return t.member(addr)
	case pb.Thread_OPEN:
		return t.member(addr)
	default:
		return false
	}
}

// annotatable returns whether or not this thread is annotatable from the
// perspective of the given address
func (t *Thread) annotatable(addr string) bool {
	if addr == t.initiator {
		return true
	}
	switch t.ttype {
	case pb.Thread_PRIVATE:
		return false // should not happen
	case pb.Thread_READ_ONLY:
		return false
	case pb.Thread_PUBLIC:
		return t.member(addr)
	case pb.Thread_OPEN:
		return t.member(addr)
	default:
		return false
	}
}

// writable returns whether or not this thread can accept files from the
// perspective of the given address
func (t *Thread) writable(addr string) bool {
	if addr == t.initiator {
		return true
	}
	switch t.ttype {
	case pb.Thread_PRIVATE:
		return false // should not happen
	case pb.Thread_READ_ONLY:
		return false
	case pb.Thread_PUBLIC:
		return false
	case pb.Thread_OPEN:
		return t.member(addr)
	default:
		return false
	}
}

// shareable returns whether or not this thread is shareable from one address to another
func (t *Thread) shareable(from string, to string) bool {
	if from == to {
		return true
	}
	switch t.sharing {
	case pb.Thread_NOT_SHARED:
		return false
	case pb.Thread_INVITE_ONLY:
		return from == t.initiator && t.member(to)
	case pb.Thread_SHARED:
		return t.member(from) && t.member(to)
	default:
		return false
	}
}

// member returns whether or not the given address is a thread member
// NOTE: Thread members are a fixed set of textile addresses specified
// when a thread is created. If empty, _everyone_ is a member.
func (t *Thread) member(addr string) bool {
	if len(t.members) == 0 || addr == t.initiator {
		return true
	}
	for _, m := range t.members {
		if m == addr {
			return true
		}
	}
	return false
}

// loadSchema loads a schema from a local file
func loadSchema(node *core.IpfsNode, id string) (*pb.Node, error) {
	data, err := ipfs.DataAtPath(node, id)
	if err != nil {
		return nil, err
	}

	var sch pb.Node
	if err := jsonpb.UnmarshalString(string(data), &sch); err != nil {
		log.Errorf("failed to unmarshal thread schema %s: %s", id, err)
		return nil, err
	}
	return &sch, nil
}

package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/schema"
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

// ThreadUpdate is used to notify listeners about updates in a thread
type ThreadUpdate struct {
	Block      BlockInfo   `json:"block"`
	ThreadId   string      `json:"thread_id"`
	ThreadKey  string      `json:"thread_key"`
	ThreadName string      `json:"thread_name"`
	Info       interface{} `json:"info,omitempty"`
}

// ThreadInfo reports info about a thread
type ThreadInfo struct {
	Id         string       `json:"id"`
	Key        string       `json:"key"`
	Name       string       `json:"name"`
	Schema     *schema.Node `json:"schema,omitempty"`
	SchemaId   string       `json:"schema_id,omitempty"`
	Initiator  string       `json:"initiator"`
	Type       string       `json:"type"`
	Sharing    string       `json:"sharing"`
	Members    []string     `json:"members,omitempty"`
	State      string       `json:"state"`
	Head       *BlockInfo   `json:"head,omitempty"`
	PeerCount  int          `json:"peer_cnt"`
	BlockCount int          `json:"block_cnt"`
	FileCount  int          `json:"file_cnt"`
}

// ThreadInviteInfo reports info about a thread invite
type ThreadInviteInfo struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username,omitempty"`
	Avatar   string    `json:"avatar,omitempty"`
	Date     time.Time `json:"date"`
}

// BlockInfo is a more readable version of repo.Block
type BlockInfo struct {
	Id       string    `json:"id"`
	ThreadId string    `json:"thread_id"`
	AuthorId string    `json:"author_id,omitempty"`
	Username string    `json:"username,omitempty"`
	Avatar   string    `json:"avatar,omitempty"`
	Type     string    `json:"type"`
	Date     time.Time `json:"date"`
	Parents  []string  `json:"parents"`
	Target   string    `json:"target,omitempty"`
	Body     string    `json:"body,omitempty"`
}

// ThreadConfig is used to construct a Thread
type ThreadConfig struct {
	RepoPath           string
	Config             *config.Config
	Node               func() *core.IpfsNode
	Datastore          repo.Datastore
	Service            func() *ThreadsService
	ThreadsOutbox      *ThreadsOutbox
	CafeOutbox         *CafeOutbox
	SendUpdate         func(update ThreadUpdate)
	ContactDisplayInfo func(id string) (string, string)
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id                 string
	Key                string // app key, usually UUID
	Name               string
	Schema             *schema.Node
	schemaId           string
	initiator          string
	ttype              repo.ThreadType
	sharing            repo.ThreadSharing
	members            []string
	privKey            libp2pc.PrivKey
	repoPath           string
	config             *config.Config
	node               func() *core.IpfsNode
	datastore          repo.Datastore
	service            func() *ThreadsService
	threadsOutbox      *ThreadsOutbox
	cafeOutbox         *CafeOutbox
	sendUpdate         func(update ThreadUpdate)
	contactDisplayInfo func(id string) (string, string)
	mux                sync.Mutex
}

// NewThread create a new Thread from a repo model and config
func NewThread(model *repo.Thread, conf *ThreadConfig) (*Thread, error) {
	sk, err := libp2pc.UnmarshalPrivateKey(model.PrivKey)
	if err != nil {
		return nil, err
	}

	var sch *schema.Node
	if model.Schema != "" {
		sch, err = loadSchema(conf.Node(), model.Schema)
		if err != nil {
			return nil, err
		}
	}

	return &Thread{
		Id:                 model.Id,
		Key:                model.Key,
		Name:               model.Name,
		Schema:             sch,
		schemaId:           model.Schema,
		initiator:          model.Initiator,
		ttype:              model.Type,
		sharing:            model.Sharing,
		members:            model.Members,
		privKey:            sk,
		repoPath:           conf.RepoPath,
		config:             conf.Config,
		node:               conf.Node,
		datastore:          conf.Datastore,
		service:            conf.Service,
		threadsOutbox:      conf.ThreadsOutbox,
		cafeOutbox:         conf.CafeOutbox,
		sendUpdate:         conf.SendUpdate,
		contactDisplayInfo: conf.ContactDisplayInfo,
	}, nil
}

// Info returns thread info
func (t *Thread) Info() (*ThreadInfo, error) {
	mod := t.datastore.Threads().Get(t.Id)
	if mod == nil {
		return nil, errThreadReload
	}

	var head *BlockInfo
	if mod.Head != "" {
		h := t.datastore.Blocks().Get(mod.Head)
		if h != nil {

			username, avatar := t.contactDisplayInfo(h.AuthorId)

			head = &BlockInfo{
				Id:       h.Id,
				ThreadId: h.ThreadId,
				AuthorId: h.AuthorId,
				Username: username,
				Avatar:   avatar,
				Type:     h.Type.Description(),
				Date:     h.Date,
				Parents:  h.Parents,
				Target:   h.Target,
				Body:     h.Body,
			}
		}
	}

	state, err := t.State()
	if err != nil {
		return nil, err
	}

	blocks := t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s'", t.Id))
	files := t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s' and type=%d", t.Id, repo.FilesBlock))

	return &ThreadInfo{
		Id:         t.Id,
		Key:        t.Key,
		Name:       t.Name,
		Schema:     t.Schema,
		SchemaId:   t.schemaId,
		Initiator:  t.initiator,
		Type:       mod.Type.Description(),
		Sharing:    mod.Sharing.Description(),
		Members:    mod.Members,
		State:      state.Description(),
		Head:       head,
		PeerCount:  len(t.Peers()) + 1,
		BlockCount: blocks,
		FileCount:  files,
	}, nil
}

// State returns the current thread state
func (t *Thread) State() (repo.ThreadState, error) {
	mod := t.datastore.Threads().Get(t.Id)
	if mod == nil {
		return -1, errThreadReload
	}
	return mod.State, nil
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
		return nil
	}

	switch block.Type {
	case pb.ThreadBlock_MERGE:
		err = t.handleMergeBlock(parent, block)
	case pb.ThreadBlock_IGNORE:
		_, err = t.handleIgnoreBlock(parent, block)
	case pb.ThreadBlock_FLAG:
		_, err = t.handleFlagBlock(parent, block)
	case pb.ThreadBlock_JOIN:
		_, err = t.handleJoinBlock(parent, block)
	case pb.ThreadBlock_ANNOUNCE:
		_, err = t.handleAnnounceBlock(parent, block)
	case pb.ThreadBlock_LEAVE:
		err = t.handleLeaveBlock(parent, block)
	case pb.ThreadBlock_MESSAGE:
		_, err = t.handleMessageBlock(parent, block)
	case pb.ThreadBlock_FILES:
		_, err = t.handleFilesBlock(parent, block)
	case pb.ThreadBlock_COMMENT:
		_, err = t.handleCommentBlock(parent, block)
	case pb.ThreadBlock_LIKE:
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
func (t *Thread) addOrUpdateContact(contact *repo.Contact) error {
	if err := t.datastore.ThreadPeers().Add(&repo.ThreadPeer{
		Id:       contact.Id,
		ThreadId: t.Id,
		Welcomed: false,
	}); err != nil {
		if !repo.ConflictError(err) {
			return err
		}
	}

	ex := t.datastore.Contacts().Get(contact.Id)
	if ex == nil || ex.Updated.UnixNano() < contact.Updated.UnixNano() {
		return t.datastore.Contacts().AddOrUpdate(contact)
	}
	return nil
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
		Address: t.config.Account.Address,
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

	if err := t.cafeOutbox.Add(id.Hash().B58String(), repo.CafeStoreRequest); err != nil {
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

	if _, err := t.addBlock(ciphertext); err != nil {
		return nil, err
	}
	return block, nil
}

// indexBlock stores off index info for this block type
func (t *Thread) indexBlock(commit *commitResult, blockType repo.BlockType, target string, body string) error {
	date, err := ptypes.Timestamp(commit.header.Date)
	if err != nil {
		return err
	}
	index := &repo.Block{
		Id:       commit.hash.B58String(),
		Type:     blockType,
		Date:     date,
		Parents:  commit.header.Parents,
		ThreadId: t.Id,
		AuthorId: commit.header.Author,
		Target:   target,
		Body:     body,
	}
	if err := t.datastore.Blocks().Add(index); err != nil {
		return err
	}

	username, avatar := t.contactDisplayInfo(index.AuthorId)

	t.pushUpdate(BlockInfo{
		Id:       index.Id,
		ThreadId: index.ThreadId,
		AuthorId: index.AuthorId,
		Username: username,
		Avatar:   avatar,
		Type:     index.Type.Description(),
		Date:     index.Date,
		Parents:  index.Parents,
		Target:   index.Target,
		Body:     index.Body,
	})

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

	return t.cafeOutbox.Add(t.Id, repo.CafeStoreThreadRequest)
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
func (t *Thread) post(commit *commitResult, peers []repo.ThreadPeer) error {
	if len(peers) == 0 {
		// flush the storage queue—this is normally done in a thread
		// via thread message queue handling, but that won't run if there's
		// no peers to send the message to.
		go t.cafeOutbox.Flush()
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

	go t.threadsOutbox.Flush()

	return nil
}

// pushUpdate pushes thread updates to UI listeners
func (t *Thread) pushUpdate(index BlockInfo) {
	t.sendUpdate(ThreadUpdate{
		Block:      index,
		ThreadId:   t.Id,
		ThreadKey:  t.Key,
		ThreadName: t.Name,
	})
}

// readable returns whether or not this thread is readable from the
// perspective of the given address
func (t *Thread) readable(addr string) bool {
	if addr == t.initiator {
		return true
	}
	switch t.ttype {
	case repo.PrivateThread:
		return false // should not happen
	case repo.ReadOnlyThread:
		return t.member(addr)
	case repo.PublicThread:
		return t.member(addr)
	case repo.OpenThread:
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
	case repo.PrivateThread:
		return false // should not happen
	case repo.ReadOnlyThread:
		return false
	case repo.PublicThread:
		return t.member(addr)
	case repo.OpenThread:
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
	case repo.PrivateThread:
		return false // should not happen
	case repo.ReadOnlyThread:
		return false
	case repo.PublicThread:
		return false
	case repo.OpenThread:
		return t.member(addr)
	default:
		return false
	}
}

// shareable returns whether or not this thread is shareable from one address to another
func (t *Thread) shareable(from string, to string) bool {
	switch t.sharing {
	case repo.NotSharedThread:
		return false
	case repo.InviteOnlyThread:
		return from == t.initiator && t.member(to)
	case repo.SharedThread:
		return t.member(from) && t.member(to)
	default:
		return false
	}
}

// member returns whether or not the given address is a thread member
// NOTE: Thread members are a fixed set of textile addresses specified
// when a thread is created. If empty, _everyone_ is a member.
func (t *Thread) member(address string) bool {
	if len(t.members) == 0 {
		return true
	}
	for _, m := range t.members {
		if m == address {
			return true
		}
	}
	return false
}

// loadSchema loads a schema from a local file
func loadSchema(node *core.IpfsNode, id string) (*schema.Node, error) {
	data, err := ipfs.DataAtPath(node, id)
	if err != nil {
		return nil, err
	}

	var sch schema.Node
	if err := json.Unmarshal(data, &sch); err != nil {
		log.Errorf("failed to unmarshal thread schema %s: %s", id, err)
		return nil, err
	}
	return &sch, nil
}

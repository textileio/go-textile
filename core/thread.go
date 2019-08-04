package core

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ipfs/go-ipfs/core"
	ipld "github.com/ipfs/go-ipld-format"
	uio "github.com/ipfs/go-unixfs/io"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/repo/config"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/schema"
	"github.com/textileio/go-textile/util"
)

// blockLinkName is the name of the link containing the encrypted block
const blockLinkName = "block"

// parentsLinkName is the name of the link used to reference parent nodes
const parentsLinkName = "parents"

// targetLinkName is the name of the link used to reference another node
const targetLinkName = "target"

// dataLinkName is the name of the link used to reference a data node
const dataLinkName = "data"

// ErrInvalidNode indicates the thread node is not valid
var ErrInvalidNode = fmt.Errorf("thread node is not valid")

// ErrNotShareable indicates the thread does not allow invites, at least for _you_
var ErrNotShareable = fmt.Errorf("thread is not shareable")

// ErrNotReadable indicates the thread is not readable
var ErrNotReadable = fmt.Errorf("thread is not readable")

// ErrNotAnnotatable indicates the thread is not annotatable (comments/likes)
var ErrNotAnnotatable = fmt.Errorf("thread is not annotatable")

// ErrNotWritable indicates the thread is not writable (files/messages)
var ErrNotWritable = fmt.Errorf("thread is not writable")

// ErrThreadSchemaRequired indicates files where added without a thread schema
var ErrThreadSchemaRequired = fmt.Errorf("thread schema required to add files")

// ErrJsonSchemaRequired indicates json files where added without a json schema
var ErrJsonSchemaRequired = fmt.Errorf("thread schema does not allow json files")

// ErrInvalidFileNode indicates files where added via a nil ipld node
var ErrInvalidFileNode = fmt.Errorf("invalid files node")

// ErrBlockWrongType indicates a block was requested as a type other than its own
var ErrBlockWrongType = fmt.Errorf("block type is not the type requested")

// errReloadFailed indicates an error occurred during thread reload
var errThreadReload = fmt.Errorf("could not re-load thread")

// ThreadConfig is used to construct a Thread
type ThreadConfig struct {
	RepoPath       string
	Config         *config.Config
	Account        *keypair.Full
	Node           func() *core.IpfsNode
	Datastore      repo.Datastore
	Service        func() *ThreadsService
	BlockOutbox    *BlockOutbox
	BlockDownloads *BlockDownloads
	CafeOutbox     *CafeOutbox
	AddPeer        func(*pb.Peer) error
	PushUpdate     func(*pb.Block, string)
}

// Thread is the primary mechanism representing a collecion of data / files / photos
type Thread struct {
	Id             string
	Key            string // app key, usually UUID
	Name           string
	PrivKey        libp2pc.PrivKey
	Schema         *pb.Node
	schemaId       string
	initiator      string
	ttype          pb.Thread_Type
	sharing        pb.Thread_Sharing
	whitelist      []string
	repoPath       string
	config         *config.Config
	account        *keypair.Full
	node           func() *core.IpfsNode
	datastore      repo.Datastore
	service        func() *ThreadsService
	blockOutbox    *BlockOutbox
	cafeOutbox     *CafeOutbox
	blockDownloads *BlockDownloads
	addPeer        func(*pb.Peer) error
	pushUpdate     func(*pb.Block, string)
	lock           sync.Mutex
}

// NewThread create a new Thread from a repo model and config
func NewThread(model *pb.Thread, conf *ThreadConfig) (*Thread, error) {
	sk, err := ipfs.UnmarshalPrivateKey(model.Sk)
	if err != nil {
		return nil, err
	}

	thrd := &Thread{
		Id:             model.Id,
		Key:            model.Key,
		Name:           model.Name,
		schemaId:       model.Schema,
		initiator:      model.Initiator,
		ttype:          model.Type,
		sharing:        model.Sharing,
		whitelist:      model.Whitelist,
		PrivKey:        sk,
		repoPath:       conf.RepoPath,
		config:         conf.Config,
		account:        conf.Account,
		node:           conf.Node,
		datastore:      conf.Datastore,
		service:        conf.Service,
		blockOutbox:    conf.BlockOutbox,
		blockDownloads: conf.BlockDownloads,
		cafeOutbox:     conf.CafeOutbox,
		addPeer:        conf.AddPeer,
		pushUpdate:     conf.PushUpdate,
	}

	err = thrd.loadSchema()
	if err != nil {
		return nil, err
	}
	return thrd, nil
}

// Heads returns the node ids of the current HEADs
func (t *Thread) Heads() ([]string, error) {
	mod := t.datastore.Threads().Get(t.Id)
	if mod == nil {
		return nil, errThreadReload
	}
	return util.SplitString(mod.Head, ","), nil
}

// LatestFiles returns the most recent files block
func (t *Thread) LatestFiles() *pb.Block {
	query := fmt.Sprintf("threadId='%s' and type=%d", t.Id, pb.Block_FILES)
	list := t.datastore.Blocks().List("", 1, query)
	if len(list.Items) == 0 {
		return nil
	}
	return list.Items[0]
}

// Peers returns locally known peers in this thread
func (t *Thread) Peers() []pb.ThreadPeer {
	return t.datastore.ThreadPeers().ListByThread(t.Id)
}

// Encrypt data with thread public key
func (t *Thread) Encrypt(data []byte) ([]byte, error) {
	return crypto.Encrypt(t.PrivKey.GetPublic(), data)
}

// Decrypt data with thread secret key
func (t *Thread) Decrypt(data []byte) ([]byte, error) {
	return crypto.Decrypt(t.PrivKey, data)
}

// UpdateSchema sets a new schema hash on the model and loads its node
func (t *Thread) UpdateSchema(hash string) error {
	err := t.datastore.Threads().UpdateSchema(t.Id, hash)
	if err != nil {
		return err
	}
	t.Schema = nil
	return t.loadSchema()
}

// followParents follows a list of node links, queueing block downloads along the way
// Note: Returns a final list of existing parent hashes that were reached during the tree traversal
func (t *Thread) followParents(parents []string) []string {
	if len(parents) == 0 {
		log.Debugf("found genesis block, aborting")
		return nil
	}
	final := make(map[string]struct{})

	wg := sync.WaitGroup{}
	for _, parent := range parents {
		if parent == "" {
			continue // some old blocks may contain empty string parents
		}

		wg.Add(1)
		go func(p string) {
			ends, err := t.followParent(p)
			if err != nil {
				log.Warningf("failed to follow parent %s: %s", p, err)
			}
			for _, p := range ends {
				final[p] = struct{}{}
			}
			wg.Done()
		}(parent)
	}
	wg.Wait()

	var list []string
	for p := range final {
		list = append(list, p)
	}
	return list
}

// followParent tries to follow a tree of blocks, processing along the way
func (t *Thread) followParent(parent string) ([]string, error) {
	node, err := ipfs.NodeAtPath(t.node(), parent, ipfs.CatTimeout)
	if err != nil {
		return nil, err
	}

	bnode := &blockNode{}
	if len(node.Links()) == 0 {
		// older block, has to be downloaded now because the links are not yet known
		bnode.hash = parent

		// avoid an unneeded download, pre-check existence
		index := t.datastore.Blocks().Get(bnode.hash)
		if index != nil {
			log.Debugf("%s exists, aborting", bnode.hash)
			return []string{bnode.hash}, nil
		}

		bnode.ciphertext, err = ipfs.DataAtPath(t.node(), bnode.hash)
		if err != nil {
			return nil, err
		}
	} else {
		bnode, err = extractNode(t.node(), node, false)
		if err != nil {
			return nil, err
		}
	}

	if bnode.ciphertext == nil {
		// content is not yet known, download it later
		err = t.blockDownloads.Add(&pb.Block{
			Id:      bnode.hash,
			Thread:  t.Id,
			Parents: bnode.parents,
			Target:  bnode.target,
			Data:    bnode.data,
			Status:  pb.Block_PENDING,
		})
		if err != nil {
			if db.ConflictError(err) {
				log.Debugf("%s exists, aborting", bnode.hash)
				return []string{parent}, nil
			}
			return nil, err
		}
	} else {
		// old block, handle now
		_, err = t.handle(bnode, false)
		if err != nil {
			return nil, err
		}
	}

	// old block, handle it now
	return t.followParents(bnode.parents), nil
}

// blockNode represents the components of a block wrapped by an ipld node
type blockNode struct {
	hash       string
	ciphertext []byte
	parents    []string
	target     string
	data       string
}

// handleResult returns info extracted from an encrypted block
type handleResult struct {
	body      string
	oldTarget string
	oldData   string
}

// handle receives a downloaded block allowing w/ it node links
func (t *Thread) handle(bnode *blockNode, replace bool) (*pb.Block, error) {
	block, err := t.unmarshalBlock(bnode.ciphertext)
	if err != nil {
		return nil, err
	}
	_, err = t.addBlock(bnode.ciphertext, false)
	if err != nil {
		return nil, err
	}

	msg := "handling " + block.Type.String()
	if block.Header.Author != "" {
		msg += " from " + block.Header.Author
	}
	log.Debug(msg)

	var res handleResult
	switch block.Type {
	case pb.Block_MERGE:
		res, err = t.handleMergeBlock(block)
	case pb.Block_IGNORE:
		res, err = t.handleIgnoreBlock(bnode, block)
	case pb.Block_FLAG:
		res, err = t.handleFlagBlock(block)
	case pb.Block_JOIN:
		res, err = t.handleJoinBlock(block)
	case pb.Block_ANNOUNCE:
		res, err = t.handleAnnounceBlock(block)
	case pb.Block_LEAVE:
		res, err = t.handleLeaveBlock(block)
	case pb.Block_TEXT:
		res, err = t.handleMessageBlock(block)
	case pb.Block_FILES:
		res, err = t.handleFilesBlock(bnode, block)
	case pb.Block_COMMENT:
		res, err = t.handleCommentBlock(block)
	case pb.Block_LIKE:
		res, err = t.handleLikeBlock(block)
	default:
		err = fmt.Errorf("invalid type: %s", block.Type)
	}
	if err != nil {
		return nil, err
	}

	// handle old block fields
	if len(block.Header.Parents) > 0 {
		bnode.parents = block.Header.Parents
	}
	if res.oldTarget != "" {
		bnode.target = res.oldTarget
	}
	if res.oldData != "" {
		bnode.data = res.oldData
	}

	index := &pb.Block{
		Id:      bnode.hash,
		Thread:  t.Id,
		Author:  block.Header.Author,
		Type:    block.Type,
		Date:    block.Header.Date,
		Parents: bnode.parents,
		Target:  bnode.target,
		Data:    bnode.data,
		Body:    res.body,
		Status:  pb.Block_READY,
	}
	err = t.indexBlock(index, replace)
	if err != nil {
		return nil, err
	}

	return index, nil
}

// addOrUpdatePeer collects and saves thread peers
func (t *Thread) addOrUpdatePeer(peer *pb.Peer, welcomed bool) error {
	if peer.Id == t.node().Identity.Pretty() {
		return nil
	}

	err := t.datastore.ThreadPeers().Add(&pb.ThreadPeer{
		Id:       peer.Id,
		Thread:   t.Id,
		Welcomed: welcomed,
	})
	if err != nil {
		if !db.ConflictError(err) {
			return err
		}
	}

	return t.addPeer(peer)
}

// newBlockHeader creates a new header
func (t *Thread) newBlockHeader() (*pb.ThreadBlockHeader, error) {
	pdate, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	return &pb.ThreadBlockHeader{
		Date:    pdate,
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
func (t *Thread) commitBlock(msg proto.Message, mtype pb.Block_BlockType, add bool, encrypt func(plaintext []byte) ([]byte, error)) (*commitResult, error) {
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

	hash, err := t.addBlock(ciphertext, !add)
	if err != nil {
		return nil, err
	}

	return &commitResult{
		hash:       hash,
		ciphertext: ciphertext,
		header:     header,
	}, nil
}

// addBlock adds to ipfs
func (t *Thread) addBlock(ciphertext []byte, hashOnly bool) (mh.Multihash, error) {
	id, err := ipfs.AddData(t.node(), bytes.NewReader(ciphertext), true, hashOnly)
	if err != nil {
		return nil, err
	}
	hash := id.Hash().B58String()

	if !hashOnly {
		err = t.cafeOutbox.Add(hash, pb.CafeRequest_STORE, cafeReqOpt.SyncGroup(hash))
		if err != nil {
			return nil, err
		}
	}

	return id.Hash(), nil
}

// unmarshalBlock decrypts and unmarshals an encrypted block
func (t *Thread) unmarshalBlock(ciphertext []byte) (*pb.ThreadBlock, error) {
	block := new(pb.ThreadBlock)
	plaintext, err := t.Decrypt(ciphertext)
	if err != nil {
		// might be a merge block
		err2 := proto.Unmarshal(ciphertext, block)
		if err2 != nil || block.Type != pb.Block_MERGE {
			return nil, err
		}
	} else {
		err = proto.Unmarshal(plaintext, block)
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

// commitNode writes the block to an IPLD node
func (t *Thread) commitNode(index *pb.Block, additionalParents []string, addIndex bool) (mh.Multihash, error) {
	dir := uio.NewDirectory(t.node().DAG)

	// add block
	err := ipfs.AddLinkToDirectory(t.node(), dir, blockLinkName, index.Id)
	if err != nil {
		return nil, err
	}

	// add parents
	heads, err := t.Heads()
	if err != nil {
		return nil, err
	}
	for _, p := range additionalParents {
		heads = append(heads, p)
	}
	pdir := uio.NewDirectory(t.node().DAG)
	for i, p := range heads {
		err = ipfs.AddLinkToDirectory(t.node(), pdir, strconv.Itoa(i), p)
		if err != nil {
			return nil, err
		}
	}
	pnode, err := pdir.GetNode()
	if err != nil {
		return nil, err
	}
	err = ipfs.PinNode(t.node(), pnode, false)
	if err != nil {
		return nil, err
	}
	pnodeId := pnode.Cid().Hash().B58String()
	err = ipfs.AddLinkToDirectory(t.node(), dir, parentsLinkName, pnodeId)
	if err != nil {
		return nil, err
	}

	// add target
	if index.Target != "" {
		err = ipfs.AddLinkToDirectory(t.node(), dir, targetLinkName, index.Target)
		if err != nil {
			return nil, err
		}
	}

	// add data
	if index.Data != "" {
		err = ipfs.AddLinkToDirectory(t.node(), dir, dataLinkName, index.Data)
		if err != nil {
			return nil, err
		}
	}

	// pin whole thing
	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}
	err = ipfs.PinNode(t.node(), node, false)
	if err != nil {
		return nil, err
	}
	nhash := node.Cid().Hash()

	if addIndex {
		// replace index with final values
		err = t.datastore.Blocks().Replace(&pb.Block{
			Id:       index.Id,
			Thread:   index.Thread,
			Author:   index.Author,
			Type:     index.Type,
			Date:     index.Date,
			Parents:  heads,
			Target:   index.Target,
			Data:     index.Data,
			Body:     index.Body,
			Status:   pb.Block_READY,
			Attempts: index.Attempts,
		})
		if err != nil {
			return nil, err
		}

		err = t.updateHead([]string{nhash.B58String()}, true)
		if err != nil {
			return nil, err
		}
	}

	// store nodes
	nodeId := nhash.B58String()
	group := cafeReqOpt.Group(nodeId)
	syncGroup := cafeReqOpt.SyncGroup(nodeId)
	err = t.cafeOutbox.Add(nodeId, pb.CafeRequest_STORE, group, syncGroup)
	if err != nil {
		return nil, err
	}
	err = t.cafeOutbox.Add(pnodeId, pb.CafeRequest_STORE, group, syncGroup)
	if err != nil {
		return nil, err
	}

	return nhash, nil
}

// indexBlock stores off index info for this block type
func (t *Thread) indexBlock(index *pb.Block, replace bool) error {
	var err error
	if replace {
		err = t.datastore.Blocks().Replace(index)
	} else {
		err = t.datastore.Blocks().Add(index)
	}
	if err != nil {
		return err
	}

	t.pushUpdate(index, t.Key)
	return nil
}

// handleHead determines what the next set of HEADs will be
// One of three situations will occur:
// 1) fast-forward: the inbound leaves are identical to current heads (heads -> inbound)
// 2) partial-split: at least one inbound leaf matches a current head (heads -> [inbound..., unmatched...])
// 3) full-split: zero inbound leaves match a current head (heads -> [inbound..., current...])
func (t *Thread) handleHead(inbound []string, leaves []string) error {
	heads, err := t.Heads()
	if err != nil {
		return err
	}

	var next []string
	unique := make(map[string]struct{})
	add := func(v string) {
		if _, ok := unique[v]; !ok {
			unique[v] = struct{}{}
			next = append(next, v)
		}
	}

	for _, i := range inbound {
		add(i)
	}
outer:
	for _, h := range heads {
		for _, l := range leaves {
			if h == l {
				continue outer
			}
		}
		// this head was not found, keep it
		add(h)
	}

	return t.updateHead(next, true)
}

// addHead adds an additional (usually temporary) head
func (t *Thread) addHead(head string) error {
	heads, err := t.Heads()
	if err != nil {
		return err
	}
	return t.updateHead(append(heads, head), false)
}

// updateHead updates the ref to the content id of the latest update
func (t *Thread) updateHead(heads []string, store bool) error {
	err := t.datastore.Threads().UpdateHead(t.Id, heads)
	if err != nil {
		return err
	}
	if !store {
		return nil
	}
	return t.store()
}

// sendWelcome sends the latest HEAD block to a set of peers
func (t *Thread) sendWelcome() error {
	peers := t.datastore.ThreadPeers().ListUnwelcomedByThread(t.Id)
	if len(peers) == 0 {
		return nil
	}

	heads, err := t.Heads()
	if err != nil {
		return err
	}
	if len(heads) == 0 {
		return nil
	}
	for _, head := range heads {
		ndata, err := ipfs.ObjectAtPath(t.node(), head)
		if err != nil {
			return err
		}
		block, err := blockCIDFromNode(t.node(), head)
		if err != nil {
			return err
		}

		ciphertext, err := ipfs.DataAtPath(t.node(), block)
		if err != nil {
			return err
		}
		sig, err := t.account.Sign(ciphertext)
		if err != nil {
			return err
		}
		env, err := t.service().NewEnvelope(t.Id, ndata, ciphertext, sig)
		if err != nil {
			return err
		}

		for _, tp := range peers {
			err = t.blockOutbox.Add(tp.Id, env)
			if err != nil {
				return err
			}
			log.Debugf("WELCOME sent to %s at %s", tp.Id, head)
		}
	}

	return t.datastore.ThreadPeers().WelcomeByThread(t.Id)
}

// post publishes an encrypted message to thread peers
func (t *Thread) post(index *pb.Block) error {
	nhash, err := t.commitNode(index, nil, index.Type != pb.Block_ADD)
	if err != nil {
		return err
	}
	ndata, err := ipfs.ObjectAtPath(t.node(), nhash.B58String())
	if err != nil {
		return err
	}
	ciphertext, err := ipfs.DataAtPath(t.node(), index.Id)
	if err != nil {
		return err
	}

	sig, err := t.account.Sign(ciphertext)
	if err != nil {
		return err
	}
	env, err := t.service().NewEnvelope(t.Id, ndata, ciphertext, sig)
	if err != nil {
		return err
	}

	var peers []pb.ThreadPeer
	if index.Type == pb.Block_ADD {
		if index.Body != "" {
			peers = []pb.ThreadPeer{{Id: index.Body}}
		}
	} else {
		peers = t.Peers()
	}

	for _, tp := range peers {
		err = t.blockOutbox.Add(tp.Id, env)
		if err != nil {
			return err
		}
	}

	// delete add blocks as they are no longer needed
	if index.Type == pb.Block_ADD {
		return t.datastore.Blocks().Delete(index.Id)
	}

	return nil
}

// store adds a store thread request
func (t *Thread) store() error {
	return t.cafeOutbox.Add(t.Id, pb.CafeRequest_STORE_THREAD)
}

// readable returns whether or not this thread is readable from the
// perspective of the given address
func (t *Thread) readable(addr string) bool {
	if addr == "" || addr == t.initiator {
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
	if addr == "" || addr == t.initiator {
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
	if addr == "" || addr == t.initiator {
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
// NOTE: Thread whitelist are a fixed set of textile addresses specified
// when a thread is created. If empty, _everyone_ is a member.
func (t *Thread) member(addr string) bool {
	if len(t.whitelist) == 0 || addr == t.initiator {
		return true
	}
	for _, m := range t.whitelist {
		if m == addr {
			return true
		}
	}
	return false
}

// loadSchema loads and attaches a schema from the network
func (t *Thread) loadSchema() error {
	if t.schemaId == "" || t.Schema != nil {
		return nil
	}

	data, err := ipfs.DataAtPath(t.node(), t.schemaId)
	if err != nil {
		if err == ipld.ErrNotFound {
			return nil
		}
		return err
	}

	var sch pb.Node
	err = jsonpb.UnmarshalString(string(data), &sch)
	if err != nil {
		return err
	}
	t.Schema = &sch

	// pin/repin to ensure remotely added schemas are readily accessible
	_, err = ipfs.AddData(t.node(), bytes.NewReader(data), true, false)
	if err != nil {
		return err
	}

	return nil
}

// validateNode ensures that the node contains the correct links
func validateNode(node ipld.Node) error {
	links := node.Links()
	if schema.LinkByName(links, []string{blockLinkName}) == nil {
		return ErrInvalidNode
	}
	return nil
}

// extractNode pulls out block components from an ipld node
func extractNode(ipfsNode *core.IpfsNode, node ipld.Node, downloadBlock bool) (*blockNode, error) {
	err := validateNode(node)
	if err != nil {
		return nil, err
	}
	err = ipfs.PinNode(ipfsNode, node, false)
	if err != nil {
		return nil, err
	}

	bnode := &blockNode{}
	links := node.Links()

	// get parents
	plink := schema.LinkByName(links, []string{parentsLinkName})
	pnode, err := ipfs.NodeAtLink(ipfsNode, plink)
	if err != nil {
		return nil, err
	}
	err = ipfs.PinNode(ipfsNode, pnode, false)
	if err != nil {
		return nil, err
	}
	for _, l := range pnode.Links() {
		bnode.parents = append(bnode.parents, l.Cid.Hash().B58String())
	}

	// get target
	tlink := schema.LinkByName(links, []string{targetLinkName})
	if tlink != nil {
		bnode.target = tlink.Cid.Hash().B58String()
	}

	// get data
	dlink := schema.LinkByName(links, []string{dataLinkName})
	if dlink != nil {
		bnode.data = dlink.Cid.Hash().B58String()
	}

	// get block
	blink := schema.LinkByName(links, []string{blockLinkName})
	bnode.hash = blink.Cid.Hash().B58String()
	if downloadBlock {
		bnode.ciphertext, err = ipfs.DataAtPath(ipfsNode, bnode.hash)
		if err != nil {
			return nil, err
		}
	}

	return bnode, nil
}

// blockCIDFromNode returns the inner block id from its ipld wrapper
func blockCIDFromNode(ipfsNode *core.IpfsNode, nhash string) (string, error) {
	node, err := ipfs.NodeAtPath(ipfsNode, nhash, ipfs.DefaultTimeout)
	if err != nil {
		return "", err
	}
	if len(node.Links()) == 0 {
		// old block
		return nhash, nil
	}
	link := schema.LinkByName(node.Links(), []string{blockLinkName})
	if link == nil {
		return "", ErrInvalidNode
	}
	return link.Cid.Hash().B58String(), nil
}

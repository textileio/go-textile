package core

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ipfs/go-ipfs/core"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	mh "github.com/multiformats/go-multihash"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/service"
)

// ErrInvalidThreadBlock is a catch all error for malformed / invalid blocks
var ErrInvalidThreadBlock = fmt.Errorf("invalid thread block")

// ThreadService is a libp2p service for orchestrating a collection of files
// with annotations amongst a group of peers
type ThreadsService struct {
	service          *service.Service
	datastore        repo.Datastore
	getThread        func(string) *Thread
	addThread        func([]byte) (mh.Multihash, error)
	removeThread     func(string) (mh.Multihash, error)
	sendNotification func(*pb.Notification) error
	online           bool
}

// NewThreadsService returns a new threads service
func NewThreadsService(
	account *keypair.Full,
	node func() *core.IpfsNode,
	datastore repo.Datastore,
	getThread func(string) *Thread,
	addThread func([]byte) (mh.Multihash, error),
	removeThread func(string) (mh.Multihash, error),
	sendNotification func(*pb.Notification) error,
) *ThreadsService {
	handler := &ThreadsService{
		datastore:        datastore,
		getThread:        getThread,
		addThread:        addThread,
		removeThread:     removeThread,
		sendNotification: sendNotification,
	}
	handler.service = service.NewService(account, handler, node)
	return handler
}

// Protocol returns the handler protocol
func (h *ThreadsService) Protocol() protocol.ID {
	return protocol.ID("/textile/threads/2.0.0")
}

// Start begins online services
func (h *ThreadsService) Start() {
	h.service.Start()
}

// Ping pings another peer
func (h *ThreadsService) Ping(pid peer.ID) (service.PeerStatus, error) {
	return h.service.Ping(pid)
}

// Handle is called by the underlying service handler method
func (h *ThreadsService) Handle(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	if env.Message.Type != pb.Message_THREAD_ENVELOPE {
		return nil, nil
	}
	tenv := new(pb.ThreadEnvelope)
	err := ptypes.UnmarshalAny(env.Message.Payload, tenv)
	if err != nil {
		return nil, err
	}

	// check for an account signature
	var accountPeer bool
	if tenv.Sig != nil {
		err = h.service.Account.Verify(tenv.Ciphertext, tenv.Sig)
		if err == nil {
			accountPeer = true
		}
	}

	hash, err := mh.FromB58String(tenv.Hash)
	if err != nil {
		return nil, err
	}

	thrd := h.getThread(tenv.Thread)
	if thrd == nil {
		// this might be a direct invite
		err = h.handleAdd(hash, tenv, accountPeer)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	block, err := thrd.handleBlock(hash, tenv.Ciphertext)
	if err != nil {
		if err == ErrBlockExists {
			// exists, abort
			log.Debugf("%s exists, aborting", hash.B58String())
			return nil, nil
		}
		return nil, err
	}

	if accountPeer {
		log.Debugf("handling %s from account peer %s", block.Type.String(), block.Header.Author)
	} else {
		if block.Header.Author != "" {
			log.Debugf("handling %s from %s", block.Type.String(), block.Header.Author)
		} else {
			log.Debugf("handling %s", block.Type.String())
		}
	}

	var leave bool
	switch block.Type {
	case pb.Block_MERGE:
		err = h.handleMerge(thrd, hash, block)
	case pb.Block_IGNORE:
		err = h.handleIgnore(thrd, hash, block)
	case pb.Block_FLAG:
		err = h.handleFlag(thrd, hash, block)
	case pb.Block_JOIN:
		err = h.handleJoin(thrd, hash, block, accountPeer)
	case pb.Block_ANNOUNCE:
		err = h.handleAnnounce(thrd, hash, block)
	case pb.Block_LEAVE:
		err = h.handleLeave(thrd, hash, block, accountPeer)
		if accountPeer {
			leave = true // we will leave as well
		}
	case pb.Block_TEXT:
		err = h.handleMessage(thrd, hash, block)
	case pb.Block_FILES:
		err = h.handleFiles(thrd, hash, block)
	case pb.Block_COMMENT:
		err = h.handleComment(thrd, hash, block)
	case pb.Block_LIKE:
		err = h.handleLike(thrd, hash, block)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	parents, err := thrd.followParents(block.Header.Parents)
	if err != nil {
		return nil, err
	}

	_, err = thrd.handleHead(hash, parents)
	if err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop
	err = thrd.sendWelcome()
	if err != nil {
		return nil, err
	}

	// we may be auto-leaving
	if leave {
		_, err = h.removeThread(thrd.Id)
		if err != nil {
			return nil, err
		}
	}

	// flush cafe queue _at the very end_
	go thrd.cafeOutbox.Flush()

	return nil, nil
}

// HandleStream is called by the underlying service handler method
func (h *ThreadsService) HandleStream(pid peer.ID, env *pb.Envelope) (chan *pb.Envelope, chan error, chan interface{}) {
	// no-op
	return make(chan *pb.Envelope), make(chan error), make(chan interface{})
}

// SendMessage sends a message to a peer
func (h *ThreadsService) SendMessage(ctx context.Context, pid peer.ID, env *pb.Envelope) error {
	return h.service.SendMessage(ctx, pid, env)
}

// NewEnvelope signs and wraps an encypted block for transport
func (h *ThreadsService) NewEnvelope(threadId string, hash mh.Multihash, ciphertext []byte, sig []byte) (*pb.Envelope, error) {
	tenv := &pb.ThreadEnvelope{
		Thread:     threadId,
		Hash:       hash.B58String(),
		Ciphertext: ciphertext,
		Sig:        sig,
	}
	return h.service.NewEnvelope(pb.Message_THREAD_ENVELOPE, tenv, nil, false)
}

// handleAdd receives an invite message
func (h *ThreadsService) handleAdd(hash mh.Multihash, tenv *pb.ThreadEnvelope, accountPeer bool) error {
	plaintext, err := crypto.Decrypt(h.service.Node().PrivateKey, tenv.Ciphertext)
	if err != nil {
		// wasn't an invite, abort
		return nil
	}
	block := new(pb.ThreadBlock)
	err = proto.Unmarshal(plaintext, block)
	if err != nil {
		return err
	}
	if block.Type != pb.Block_ADD {
		return ErrInvalidThreadBlock
	}
	msg := new(pb.ThreadAdd)
	err = ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return err
	}

	if accountPeer {
		log.Debugf("handling %s from account peer %s", block.Type.String(), block.Header.Author)

		// can auto-join
		_, err = h.addThread(plaintext)
		if err != nil {
			return err
		}
		return nil
	} else {
		if block.Header.Author != "" {
			log.Debugf("handling %s from %s", block.Type.String(), block.Header.Author)
		} else {
			log.Debugf("handling %s", block.Type.String())
		}
	}

	err = h.datastore.Invites().Add(&pb.Invite{
		Id:      hash.B58String(),
		Block:   plaintext,
		Name:    msg.Thread.Name,
		Inviter: msg.Inviter,
		Date:    block.Header.Date,
	})
	if err != nil {
		if !db.ConflictError(err) {
			return err
		}
		// exists, abort
		return nil
	}

	note := h.newNotification(block.Header, pb.Notification_INVITE_RECEIVED)
	note.SubjectDesc = msg.Thread.Name
	note.Subject = tenv.Thread
	note.Block = hash.B58String() // invite block
	note.Body = "invited you to join"

	return h.sendNotification(note)
}

// handleMerge receives a merge message
func (h *ThreadsService) handleMerge(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	return thrd.handleMergeBlock(hash, block)
}

// handleIgnore receives an ignore message
func (h *ThreadsService) handleIgnore(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	_, err := thrd.handleIgnoreBlock(hash, block)
	return err
}

// handleFlag receives a flag message
func (h *ThreadsService) handleFlag(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	_, err := thrd.handleFlagBlock(hash, block)
	return err
}

// handleJoin receives a join message
func (h *ThreadsService) handleJoin(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock, accountPeer bool) error {
	_, err := thrd.handleJoinBlock(hash, block)
	if err != nil {
		return err
	}

	var ntype pb.Notification_Type
	if accountPeer {
		ntype = pb.Notification_ACCOUNT_PEER_JOINED
	} else {
		ntype = pb.Notification_PEER_JOINED
	}
	note := h.newNotification(block.Header, ntype)
	note.SubjectDesc = thrd.Name
	note.Subject = thrd.Id
	note.Block = hash.B58String()
	note.Body = "joined"

	return h.sendNotification(note)
}

// handleAnnounce receives an announce message
func (h *ThreadsService) handleAnnounce(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	_, err := thrd.handleAnnounceBlock(hash, block)
	return err
}

// handleLeave receives a leave message
func (h *ThreadsService) handleLeave(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock, accountPeer bool) error {
	if err := thrd.handleLeaveBlock(hash, block); err != nil {
		return err
	}

	var ntype pb.Notification_Type
	if accountPeer {
		ntype = pb.Notification_ACCOUNT_PEER_LEFT
	} else {
		ntype = pb.Notification_PEER_LEFT
	}
	note := h.newNotification(block.Header, ntype)
	note.SubjectDesc = thrd.Name
	note.Subject = thrd.Id
	note.Block = hash.B58String()
	note.Body = "left"

	return h.sendNotification(note)
}

// handleMessage receives a message message
func (h *ThreadsService) handleMessage(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleMessageBlock(hash, block)
	if err != nil {
		return err
	}

	note := h.newNotification(block.Header, pb.Notification_MESSAGE_ADDED)
	note.Body = msg.Body
	note.Block = hash.B58String()
	note.SubjectDesc = thrd.Name
	note.Subject = thrd.Id

	return h.sendNotification(note)
}

// handleData receives a files message
func (h *ThreadsService) handleFiles(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleFilesBlock(hash, block)
	if err != nil {
		return err
	}

	note := h.newNotification(block.Header, pb.Notification_FILES_ADDED)
	note.Target = msg.Target
	note.Body = "added " + threadSubject(thrd.Schema.Name)
	note.Block = hash.B58String()
	note.SubjectDesc = thrd.Name
	note.Subject = thrd.Id

	return h.sendNotification(note)
}

// handleComment receives a comment message
func (h *ThreadsService) handleComment(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleCommentBlock(hash, block)
	if err != nil {
		return err
	}

	target := h.datastore.Blocks().Get(msg.Target)
	if target == nil {
		return nil
	}
	var desc string
	if target.Author == h.service.Node().Identity.Pretty() {
		desc = "your " + threadSubject(thrd.Schema.Name)
	} else {
		desc = "a " + threadSubject(thrd.Schema.Name)
	}

	note := h.newNotification(block.Header, pb.Notification_COMMENT_ADDED)
	note.Body = fmt.Sprintf("commented on %s: \"%s\"", desc, msg.Body)
	note.Block = hash.B58String()
	note.Target = target.Target
	note.SubjectDesc = thrd.Name
	note.Subject = thrd.Id

	return h.sendNotification(note)
}

// handleLike receives a like message
func (h *ThreadsService) handleLike(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleLikeBlock(hash, block)
	if err != nil {
		return err
	}

	target := h.datastore.Blocks().Get(msg.Target)
	if target == nil {
		return nil
	}
	var desc string
	if target.Author == h.service.Node().Identity.Pretty() {
		desc = "your " + threadSubject(thrd.Schema.Name)
	} else {
		desc = "a " + threadSubject(thrd.Schema.Name)
	}

	note := h.newNotification(block.Header, pb.Notification_LIKE_ADDED)
	note.Body = "liked " + desc
	note.Block = hash.B58String()
	note.Target = target.Target
	note.SubjectDesc = thrd.Name
	note.Subject = thrd.Id

	return h.sendNotification(note)
}

// newNotification returns new thread notification
func (h *ThreadsService) newNotification(header *pb.ThreadBlockHeader, ntype pb.Notification_Type) *pb.Notification {
	return &pb.Notification{
		Id:    ksuid.New().String(),
		Date:  header.Date,
		Actor: header.Author,
		Type:  ntype,
	}
}

// threadSubject returns the thread subject
func threadSubject(name string) string {
	var sub string
	if name != "" {
		sub = name + " "
	}
	return sub + "files"
}

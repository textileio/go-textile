package core

import (
	"context"
	"errors"
	"fmt"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/service"
)

// ErrInvalidThreadBlock is a catch all error for malformed / invalid blocks
var ErrInvalidThreadBlock = errors.New("invalid thread block")

// ThreadService is a libp2p service for orchestrating a collection of files
// with annotations amongst a group of peers
type ThreadsService struct {
	service          *service.Service
	datastore        repo.Datastore
	getThread        func(id string) *Thread
	sendNotification func(note *repo.Notification) error
	online           bool
}

// NewThreadsService returns a new threads service
func NewThreadsService(
	account *keypair.Full,
	node func() *core.IpfsNode,
	datastore repo.Datastore,
	getThread func(id string) *Thread,
	sendNotification func(note *repo.Notification) error,
) *ThreadsService {
	handler := &ThreadsService{
		datastore:        datastore,
		getThread:        getThread,
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
	if err := ptypes.UnmarshalAny(env.Message.Payload, tenv); err != nil {
		return nil, err
	}
	hash, err := mh.FromB58String(tenv.Hash)
	if err != nil {
		return nil, err
	}

	thrd := h.getThread(tenv.Thread)
	if thrd == nil {
		// this might be a direct invite
		if err := h.handleInvite(hash, tenv); err != nil {
			return nil, err
		}
		return nil, nil
	}

	block, err := thrd.handleBlock(hash, tenv.Ciphertext)
	if err != nil {
		return nil, err
	}
	if block == nil {
		// exists, abort
		return nil, nil
	}

	switch block.Type {
	case pb.ThreadBlock_MERGE:
		log.Debugf("handling MERGE from %s", block.Header.Author)
		err = h.handleMerge(thrd, hash, block)
	case pb.ThreadBlock_IGNORE:
		log.Debugf("handling IGNORE from %s", block.Header.Author)
		err = h.handleIgnore(thrd, hash, block)
	case pb.ThreadBlock_FLAG:
		log.Debugf("handling FLAG from %s", block.Header.Author)
		err = h.handleFlag(thrd, hash, block)
	case pb.ThreadBlock_JOIN:
		log.Debugf("handling JOIN from %s", block.Header.Author)
		err = h.handleJoin(thrd, hash, block)
	case pb.ThreadBlock_ANNOUNCE:
		log.Debugf("handling ANNOUNCE from %s", block.Header.Author)
		err = h.handleAnnounce(thrd, hash, block)
	case pb.ThreadBlock_LEAVE:
		log.Debugf("handling LEAVE from %s", block.Header.Author)
		err = h.handleLeave(thrd, hash, block)
	case pb.ThreadBlock_MESSAGE:
		log.Debugf("handling MESSAGE from %s", block.Header.Author)
		err = h.handleMessage(thrd, hash, block)
	case pb.ThreadBlock_FILES:
		log.Debugf("handling FILES from %s", block.Header.Author)
		err = h.handleFiles(thrd, hash, block)
	case pb.ThreadBlock_COMMENT:
		log.Debugf("handling COMMENT from %s", block.Header.Author)
		err = h.handleComment(thrd, hash, block)
	case pb.ThreadBlock_LIKE:
		log.Debugf("handling LIKE from %s", block.Header.Author)
		err = h.handleLike(thrd, hash, block)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := thrd.followParents(block.Header.Parents); err != nil {
		return nil, err
	}

	if _, err := thrd.handleHead(hash, block.Header.Parents); err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop
	if err := thrd.sendWelcome(); err != nil {
		return nil, err
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
func (h *ThreadsService) NewEnvelope(threadId string, hash mh.Multihash, ciphertext []byte) (*pb.Envelope, error) {
	tenv := &pb.ThreadEnvelope{
		Thread:     threadId,
		Hash:       hash.B58String(),
		Ciphertext: ciphertext,
	}
	return h.service.NewEnvelope(pb.Message_THREAD_ENVELOPE, tenv, nil, false)
}

// handleInvite receives an invite message
func (h *ThreadsService) handleInvite(hash mh.Multihash, tenv *pb.ThreadEnvelope) error {
	plaintext, err := crypto.Decrypt(h.service.Node().PrivateKey, tenv.Ciphertext)
	if err != nil {
		// wasn't an invite, abort
		return ErrInvalidThreadBlock
	}
	block := new(pb.ThreadBlock)
	if err := proto.Unmarshal(plaintext, block); err != nil {
		return err
	}
	if block.Type != pb.ThreadBlock_INVITE {
		return ErrInvalidThreadBlock
	}
	msg := new(pb.ThreadInvite)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return err
	}

	log.Debugf("handling THREAD_INVITE from %s", block.Header.Author)

	date, err := ptypes.Timestamp(block.Header.Date)
	if err != nil {
		return err
	}

	if err := h.datastore.ThreadInvites().Add(&repo.ThreadInvite{
		Id:      hash.B58String(),
		Block:   plaintext,
		Name:    msg.Name,
		Contact: protoContactToRepo(msg.Contact),
		Date:    date,
	}); err != nil {
		if !repo.ConflictError(err) {
			return err
		}
		// exists, abort
		return nil
	}

	notification, err := h.newNotification(block.Header, repo.InviteReceivedNotification)
	if err != nil {
		return err
	}
	notification.Subject = msg.Name
	notification.SubjectId = tenv.Thread
	notification.BlockId = hash.B58String() // invite block
	notification.Body = "invited you to join"
	return h.sendNotification(notification)
}

// handleMerge receives a merge message
func (h *ThreadsService) handleMerge(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	return thrd.handleMergeBlock(hash, block)
}

// handleIgnore receives an ignore message
func (h *ThreadsService) handleIgnore(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if _, err := thrd.handleIgnoreBlock(hash, block); err != nil {
		return err
	}
	return nil
}

// handleFlag receives a flag message
func (h *ThreadsService) handleFlag(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if _, err := thrd.handleFlagBlock(hash, block); err != nil {
		return err
	}
	return nil
}

// handleJoin receives a join message
func (h *ThreadsService) handleJoin(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if _, err := thrd.handleJoinBlock(hash, block); err != nil {
		return err
	}

	notification, err := h.newNotification(block.Header, repo.PeerJoinedNotification)
	if err != nil {
		return err
	}
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	notification.BlockId = hash.B58String()
	notification.Body = "joined"
	return h.sendNotification(notification)
}

// handleAnnounce receives an announce message
func (h *ThreadsService) handleAnnounce(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if _, err := thrd.handleAnnounceBlock(hash, block); err != nil {
		return err
	}
	return nil
}

// handleLeave receives a leave message
func (h *ThreadsService) handleLeave(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if err := thrd.handleLeaveBlock(hash, block); err != nil {
		return err
	}

	notification, err := h.newNotification(block.Header, repo.PeerLeftNotification)
	if err != nil {
		return err
	}
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	notification.BlockId = hash.B58String()
	notification.Body = "left"
	return h.sendNotification(notification)
}

// handleMessage receives a message message
func (h *ThreadsService) handleMessage(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleMessageBlock(hash, block)
	if err != nil {
		return err
	}

	notification, err := h.newNotification(block.Header, repo.MessageAddedNotification)
	if err != nil {
		return err
	}
	notification.Body = msg.Body
	notification.BlockId = hash.B58String()
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	return h.sendNotification(notification)
}

// handleData receives a files message
func (h *ThreadsService) handleFiles(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleFilesBlock(hash, block)
	if err != nil {
		return err
	}

	notification, err := h.newNotification(block.Header, repo.FilesAddedNotification)
	if err != nil {
		return err
	}
	notification.Target = msg.Target
	notification.Body = "added a " + threadSubject(thrd.Schema.Name)
	notification.BlockId = hash.B58String()
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	return h.sendNotification(notification)
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
	if target.AuthorId == h.service.Node().Identity.Pretty() {
		desc = "your " + threadSubject(thrd.Schema.Name)
	} else {
		desc = "a " + threadSubject(thrd.Schema.Name)
	}
	notification, err := h.newNotification(block.Header, repo.CommentAddedNotification)
	if err != nil {
		return err
	}
	notification.Body = fmt.Sprintf("commented on %s: \"%s\"", desc, msg.Body)
	notification.BlockId = hash.B58String()
	notification.Target = target.Target
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	return h.sendNotification(notification)
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
	if target.AuthorId == h.service.Node().Identity.Pretty() {
		desc = "your " + threadSubject(thrd.Schema.Name)
	} else {
		desc = "a " + threadSubject(thrd.Schema.Name)
	}
	notification, err := h.newNotification(block.Header, repo.LikeAddedNotification)
	if err != nil {
		return err
	}
	notification.Body = "liked " + desc
	notification.BlockId = hash.B58String()
	notification.Target = target.Target
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	return h.sendNotification(notification)
}

// newNotification returns new thread notification
func (h *ThreadsService) newNotification(header *pb.ThreadBlockHeader, ntype repo.NotificationType) (*repo.Notification, error) {
	date, err := ptypes.Timestamp(header.Date)
	if err != nil {
		return nil, err
	}
	return &repo.Notification{
		Id:      ksuid.New().String(),
		Date:    date,
		ActorId: header.Author,
		Type:    ntype,
	}, nil
}

// threadSubject returns the thread subject
func threadSubject(name string) string {
	if name == "" {
		return "file"
	}
	return name
}

package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/service"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
)

// ErrInvalidThreadBlock is a catch all error for malformed / invalid blocks
var ErrInvalidThreadBlock = errors.New("invalid thread block")

// ThreadService is a libp2p service for orchestrating a collection of files
// with annotations amongst a group of peers
type ThreadsService struct {
	service          *service.Service
	datastore        repo.Datastore
	getThread        func(id string) (*int, *Thread)
	sendNotification func(note *repo.Notification) error
}

// NewThreadsService returns a new threads service
func NewThreadsService(
	account *keypair.Full,
	node *core.IpfsNode,
	datastore repo.Datastore,
	getThread func(id string) (*int, *Thread),
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
	return protocol.ID("/textile/threads/1.0.0")
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

	// lookup thread
	_, thrd := h.getThread(tenv.Thread)
	if thrd == nil {
		// this might be a direct invite
		if err := h.handleInvite(hash, tenv); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// decrypt and handle
	block, err := thrd.handleBlock(hash, tenv.CipherBlock)
	if err != nil {
		return nil, err
	}
	if block == nil {
		// exists, abort
		return nil, nil
	}

	// select a handler
	switch block.Type {
	case pb.ThreadBlock_JOIN:
		err = h.handleJoin(thrd, hash, block)
	case pb.ThreadBlock_LEAVE:
		err = h.handleLeave(thrd, hash, block)
	case pb.ThreadBlock_DATA:
		err = h.handleData(thrd, hash, block)
	case pb.ThreadBlock_ANNOTATION:
		err = h.handleAnnotation(thrd, hash, block)
	case pb.ThreadBlock_IGNORE:
		err = h.handleIgnore(thrd, hash, block)
	case pb.ThreadBlock_MERGE:
		err = h.handleMerge(thrd, hash, block)
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// back prop
	if err := thrd.followParents(block.Header.Parents); err != nil {
		return nil, err
	}

	// handle HEAD
	if _, err := thrd.handleHead(hash, block.Header.Parents); err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop
	if err := thrd.sendWelcome(); err != nil {
		return nil, err
	}

	// flush cafe queue at the very end
	go thrd.cafeOutbox.Flush()

	return nil, nil
}

// SendMessage sends a message to a peer
func (h *ThreadsService) SendMessage(pid peer.ID, env *pb.Envelope) error {
	return h.service.SendMessage(pid, env)
}

// NewEnvelope signs and wraps an encypted block for transport
func (h *ThreadsService) NewEnvelope(threadId string, hash mh.Multihash, ciphertext []byte) (*pb.Envelope, error) {
	tenv := &pb.ThreadEnvelope{
		Thread:      threadId,
		Hash:        hash.B58String(),
		CipherBlock: ciphertext,
	}
	return h.service.NewEnvelope(pb.Message_THREAD_ENVELOPE, tenv, nil, false)
}

// handleInvite receives an invite message
func (h *ThreadsService) handleInvite(hash mh.Multihash, tenv *pb.ThreadEnvelope) error {
	// attempt decrypt w/ own keys
	plaintext, err := crypto.Decrypt(h.service.Node.PrivateKey, tenv.CipherBlock)
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

	// pin locally for use later
	// NOTE: as an enhancement, we could maintain an sql table for "pending threads",
	// or make "pending" a type of thread, put sk here into those rows so we don't have
	// to re-download and unpack this ciphertext
	// TODO: unpin when invite joined / ignored
	// TODO: delete notification when joined / ignored
	if _, err := ipfs.PinData(h.service.Node, bytes.NewReader(tenv.CipherBlock)); err != nil {
		return err
	}

	// send notification
	notification, err := h.newNotification(block.Header, repo.ReceivedInviteNotification)
	if err != nil {
		return err
	}
	notification.Subject = msg.Name
	notification.SubjectId = tenv.Thread
	notification.BlockId = hash.B58String() // invite block
	notification.Body = "invited you to join"
	return h.sendNotification(notification)
}

// handleJoin receives a join message
func (h *ThreadsService) handleJoin(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if _, err := thrd.handleJoinBlock(hash, block); err != nil {
		return err
	}

	// send notification
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

// handleLeave receives a leave message
func (h *ThreadsService) handleLeave(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if err := thrd.handleLeaveBlock(hash, block); err != nil {
		return err
	}

	// send notification
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

// handleData receives a data message
func (h *ThreadsService) handleData(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleDataBlock(hash, block)
	if err != nil {
		return err
	}

	// send notification
	var notification *repo.Notification
	switch msg.Type {
	case pb.ThreadData_PHOTO:
		notification, err = h.newNotification(block.Header, repo.PhotoAddedNotification)
		if err != nil {
			return err
		}
		notification.DataId = msg.Data
		notification.Body = "added a photo"
	case pb.ThreadData_TEXT:
		notification, err = h.newNotification(block.Header, repo.TextAddedNotification)
		if err != nil {
			return err
		}
		notification.Body = msg.Caption
	}
	notification.BlockId = hash.B58String()
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	return h.sendNotification(notification)
}

// handleAnnotation receives an annotation message
func (h *ThreadsService) handleAnnotation(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	msg, err := thrd.handleAnnotationBlock(hash, block)
	if err != nil {
		return err
	}

	// send notification
	dataBlock := h.datastore.Blocks().Get(msg.Data)
	if dataBlock == nil {
		return nil
	}
	var target string
	// NOTE: not restricted to photo annotations here, just currently only thing possible
	if dataBlock.AuthorId == h.service.Node.Identity.Pretty() {
		target = "your photo"
	} else {
		target = "a photo"
	}
	var notification *repo.Notification
	switch msg.Type {
	case pb.ThreadAnnotation_COMMENT:
		notification, err = h.newNotification(block.Header, repo.CommentAddedNotification)
		if err != nil {
			return err
		}
		notification.Body = fmt.Sprintf("commented on %s: \"%s\"", target, msg.Caption)
	case pb.ThreadAnnotation_LIKE:
		notification, err = h.newNotification(block.Header, repo.LikeAddedNotification)
		if err != nil {
			return err
		}
		notification.Body = "liked " + target
	}
	notification.BlockId = hash.B58String()
	notification.DataId = dataBlock.DataId
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	return h.sendNotification(notification)
}

// handleIgnore receives an ignore message
func (h *ThreadsService) handleIgnore(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	if _, err := thrd.handleIgnoreBlock(hash, block); err != nil {
		return err
	}
	return nil
}

// handleMerge receives a merge message
func (h *ThreadsService) handleMerge(thrd *Thread, hash mh.Multihash, block *pb.ThreadBlock) error {
	return thrd.handleMergeBlock(hash, block)
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

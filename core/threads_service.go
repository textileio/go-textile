package core

import (
	"bytes"
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
	"github.com/textileio/go-textile/ipfs"
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
	addThread        func([]byte, []string) (mh.Multihash, error)
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
	addThread func([]byte, []string) (mh.Multihash, error),
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
	return h.service.Ping(pid.Pretty())
}

// Handle is called by the underlying service handler method
func (h *ThreadsService) Handle(env *pb.Envelope, pid peer.ID) (*pb.Envelope, error) {
	if env.Message.Type != pb.Message_THREAD_ENVELOPE {
		return nil, nil
	}

	tenv := new(pb.ThreadEnvelope)
	err := ptypes.UnmarshalAny(env.Message.Payload, tenv)
	if err != nil {
		return nil, err
	}

	bnode := &blockNode{}
	var nhash string
	if tenv.Ciphertext != nil {
		// old block
		bnode.hash = tenv.Hash
		bnode.ciphertext = tenv.Ciphertext
		nhash = bnode.hash
	} else {
		id, err := ipfs.AddObject(h.service.Node(), bytes.NewReader(tenv.Node), false)
		if err != nil {
			return nil, err
		}
		node, err := ipfs.NodeAtCid(h.service.Node(), *id)
		if err != nil {
			return nil, err
		}
		bnode, err = extractNode(h.service.Node(), node, true)
		if err != nil {
			return nil, err
		}
		nhash = node.Cid().Hash().B58String()
	}

	// check for an account signature
	var accountPeer bool
	if tenv.Sig != nil {
		err = h.service.Account.Verify(bnode.ciphertext, tenv.Sig)
		if err == nil {
			accountPeer = true
		}
	}

	thread := h.getThread(tenv.Thread)
	if thread == nil {
		// this might be a direct invite
		err = h.handleAdd(tenv.Thread, bnode, accountPeer)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	index := h.datastore.Blocks().Get(bnode.hash)
	if index != nil {
		log.Debugf("%s exists, aborting")
		return nil, nil
	}
	index, err = thread.handle(bnode, false)
	if err != nil {
		return nil, err
	}

	// naively add the inbound head, it will be cleaned up later after following parents,
	// but we don't want to lose it in the meantime
	err = thread.addHead(nhash)
	if err != nil {
		return nil, err
	}

	// some updates generate a notification
	note := &pb.Notification{
		Id:          ksuid.New().String(),
		Date:        index.Date,
		Actor:       index.Author,
		Subject:     thread.Id,
		SubjectDesc: thread.Name,
		Block:       bnode.hash,
		Target:      index.Target,
		Body:        index.Body,
	}

	send := true
	switch index.Type {
	case pb.Block_JOIN:
		if accountPeer {
			note.Type = pb.Notification_ACCOUNT_PEER_JOINED
		} else {
			note.Type = pb.Notification_PEER_JOINED
		}
		note.Body = "joined"
	case pb.Block_LEAVE:
		if accountPeer {
			note.Type = pb.Notification_ACCOUNT_PEER_LEFT
		} else {
			note.Type = pb.Notification_PEER_LEFT
		}
		note.Body = "left"
	case pb.Block_TEXT:
		note.Type = pb.Notification_MESSAGE_ADDED
	case pb.Block_FILES:
		note.Type = pb.Notification_FILES_ADDED
		if note.Body == "" { // might be caption
			note.Body = "added data"
		}
	case pb.Block_COMMENT:
		note.Type = pb.Notification_COMMENT_ADDED
	case pb.Block_LIKE:
		note.Type = pb.Notification_LIKE_ADDED
		note.Body = "added a like"
	default:
		send = false
	}
	if send {
		err = h.sendNotification(note)
		if err != nil {
			return nil, err
		}
	}

	// we may be auto-leaving
	if index.Type == pb.Block_LEAVE && accountPeer {
		_, err = h.removeThread(thread.Id)
		if err != nil {
			log.Warningf("failed to remove thread %s: %s", thread.Id, err)
			return nil, err
		}

		go thread.cafeOutbox.Flush()
		return nil, nil
	}

	// handle the thread tail in the background
	go func() {
		leaves := thread.followParents(bnode.parents)
		err = thread.handleHead([]string{nhash}, leaves)
		if err != nil {
			log.Warningf("failed to handle head %s: %s", nhash, err)
			return
		}

		// handle newly discovered peers during back prop
		err = thread.sendWelcome()
		if err != nil {
			log.Warningf("error sending welcome: %s", err)
			return
		}

		// flush cafe queue _at the very end_
		thread.cafeOutbox.Flush()
	}()

	return nil, nil
}

// HandleStream is called by the underlying service handler method
func (h *ThreadsService) HandleStream(env *pb.Envelope, pid peer.ID) (chan *pb.Envelope, chan error, chan interface{}) {
	// no-op
	return make(chan *pb.Envelope), make(chan error), make(chan interface{})
}

// SendMessage sends a message to a peer
func (h *ThreadsService) SendMessage(ctx context.Context, peerId string, env *pb.Envelope) error {
	return h.service.SendMessage(ctx, peerId, env)
}

// NewEnvelope signs and wraps an encypted block for transport
func (h *ThreadsService) NewEnvelope(threadId string, node []byte, sig []byte) (*pb.Envelope, error) {
	tenv := &pb.ThreadEnvelope{
		Thread: threadId,
		Node:   node,
		Sig:    sig,
	}
	return h.service.NewEnvelope(pb.Message_THREAD_ENVELOPE, tenv, nil, false)
}

// handleAdd receives an invite message
func (h *ThreadsService) handleAdd(threadId string, bnode *blockNode, accountPeer bool) error {
	plaintext, err := crypto.Decrypt(h.service.Node().PrivateKey, bnode.ciphertext)
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
		_, err = h.addThread(plaintext, bnode.parents)
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
		Id:      bnode.hash,
		Block:   plaintext,
		Name:    msg.Thread.Name,
		Inviter: msg.Inviter,
		Date:    block.Header.Date,
		Parents: bnode.parents,
	})
	if err != nil {
		if !db.ConflictError(err) {
			return err
		}
		// exists, abort
		return nil
	}

	return h.sendNotification(&pb.Notification{
		Id:          ksuid.New().String(),
		Date:        block.Header.Date,
		Actor:       block.Header.Author,
		Subject:     threadId,
		SubjectDesc: msg.Thread.Name,
		Block:       bnode.hash,
		Type:        pb.Notification_INVITE_RECEIVED,
		Body:        "invited you to join",
	})
}

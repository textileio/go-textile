package thread

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"strings"
	"time"
)

// Join creates an outgoing join block
func (t *Thread) Join(inviterPk libp2pc.PubKey, blockId string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	inviterPkb, err := inviterPk.Bytes()
	if err != nil {
		return nil, err
	}

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadJoin{
		Header:    header,
		InviterPk: inviterPkb,
		BlockId:   blockId,
	}

	// commit to ipfs
	message, addr, err := t.commitBlock(content, pb.Message_THREAD_JOIN)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// index it locally
	if err := t.indexBlock(id, header, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// add new peer
	inviterPid, err := peer.IDFromPublicKey(inviterPk)
	if err != nil {
		return nil, err
	}
	newPeer := &repo.Peer{
		Row:      ksuid.New().String(),
		Id:       inviterPid.Pretty(),
		ThreadId: t.Id,
		PubKey:   inviterPkb,
	}
	if err := t.peers().Add(newPeer); err != nil {
		// TODO: #202 (Properly handle database/sql errors)
		log.Warningf("peer with id %s already exists in thread %s", newPeer.Id, t.Id)
	}

	// post it
	t.post(message, id, t.Peers())

	log.Debugf("joined %s via invite %s", t.Id, blockId)

	// all done
	return addr, nil
}

// HandleJoinBlock handles an incoming join block
func (t *Thread) HandleJoinBlock(message *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadJoin, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadJoin)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, err
		}
	}

	// add to ipfs
	addr, err := t.AddBlock(message)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// check if we aleady have this block indexed
	// (should only happen if a misbehaving peer keeps sending the same block)
	index := t.blocks().Get(id)
	if index != nil {
		return nil, nil
	}

	// get the invitee id
	inviteePk, err := libp2pc.UnmarshalPublicKey(content.Header.AuthorPk)
	if err != nil {
		return nil, err
	}
	inviteeId, err := peer.IDFromPublicKey(inviteePk)
	if err != nil {
		return nil, err
	}

	// add invitee as a new local peer.
	// double-check not self in case we're re-discovering the thread
	self := inviteeId.Pretty() == t.ipfs().Identity.Pretty()
	if !self {
		newPeer := &repo.Peer{
			Row:      ksuid.New().String(),
			Id:       inviteeId.Pretty(),
			ThreadId: libp2pc.ConfigEncodeKey(content.Header.ThreadPk),
			PubKey:   content.Header.AuthorPk,
		}
		if err := t.peers().Add(newPeer); err != nil {
			// TODO: #202 (Properly handle database/sql errors)
			log.Warningf("peer with id %s already exists in thread %s", newPeer.Id, t.Id)
		}
	}

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	if err := t.FollowParents(content.Header.Parents); err != nil {
		return nil, err
	}

	// short circuit if we're traversing history
	if following {
		return addr, nil
	}

	// send welcome direct to this peer if the invitee could use a merge, i.e., we have newer updates
	head, err := t.GetHead()
	if err != nil {
		return nil, err
	}
	if head != content.BlockId {
		if _, err := t.Welcome(inviteeId.Pretty()); err != nil {
			return nil, err
		}
	}

	// handle HEAD
	if _, err := t.handleHead(id, content.Header.Parents); err != nil {
		return nil, err
	}

	return addr, nil
}

// Join creates an outgoing join block, which is NOT commited to the hash chain
func (t *Thread) Welcome(peerId string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	target := t.peers().GetById(peerId)
	if target == nil {
		return nil, errors.New(fmt.Sprintf("peer not found: %s", peerId))
	}

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadWelcome{
		Header: header,
	}

	// commit to ipfs
	message, addr, err := t.commitBlock(content, pb.Message_THREAD_WELCOME)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// post it
	t.post(message, id, []repo.Peer{*target})

	log.Debugf("welcomed %s at update %s", peerId, strings.Join(header.Parents, ","))

	// all done
	return addr, nil
}

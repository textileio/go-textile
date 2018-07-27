package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
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
func (t *Thread) HandleJoinBlock(message *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadJoin) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadJoin)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, err
		}
	}

	// add to ipfs
	addr, err := t.addBlock(message)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// check if we aleady have this block indexed
	// (should only happen if a misbehaving peer keeps sending the same block)
	index := t.blocks().Get(id)
	if index != nil {
		return nil, err
	}

	// get the inviter id
	inviteePk, err := libp2pc.UnmarshalPublicKey(content.Header.AuthorPk)
	if err != nil {
		return nil, err
	}
	inviteeId, err := peer.IDFromPublicKey(inviteePk)
	if err != nil {
		return nil, err
	}

	// add issuer as a new local peer
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

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	if err := t.FollowParents(content.Header.Parents); err != nil {
		return nil, err
	}

	return addr, nil
}

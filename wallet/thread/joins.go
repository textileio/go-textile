package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// welcome represents an intention to welcome a peer at a certain block
// (these are needed in order to remember the point at which to welcome a
// given peer while recursively following parents)
type peerWelcome struct {
	peer      repo.Peer
	atBlockId string
}

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
func (t *Thread) HandleJoinBlock(from *peer.ID, message *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadJoin, following bool) (mh.Multihash, *repo.Peer, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadJoin)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, nil, err
		}
	}

	// add to ipfs
	addr, err := t.addBlock(message)
	if err != nil {
		return nil, nil, err
	}
	id := addr.B58String()

	// check if we aleady have this block indexed
	// (should only happen if a misbehaving peer keeps sending the same block)
	index := t.blocks().Get(id)
	if index != nil {
		return nil, nil, err
	}

	// get the invitee id
	authorPk, err := libp2pc.UnmarshalPublicKey(content.Header.AuthorPk)
	if err != nil {
		return nil, nil, err
	}
	authorId, err := peer.IDFromPublicKey(authorPk)
	if err != nil {
		return nil, nil, err
	}

	// add invitee as a new local peer.
	// double-check not self in case we're re-discovering the thread
	var joined *repo.Peer
	var isNew bool
	self := authorId.Pretty() == t.ipfs().Identity.Pretty()
	if !self {
		joined = &repo.Peer{
			Row:      ksuid.New().String(),
			Id:       authorId.Pretty(),
			ThreadId: libp2pc.ConfigEncodeKey(content.Header.ThreadPk),
			PubKey:   content.Header.AuthorPk,
		}
		if err := t.peers().Add(joined); err != nil {
			// TODO: #202 (Properly handle database/sql errors)
			log.Warningf("peer with id %s already exists in thread %s", joined.Id, t.Id)
		} else {
			isNew = true
		}
	}

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.JoinBlock, nil); err != nil {
		return nil, nil, err
	}

	// back prop
	newPeers, err := t.FollowParents(content.Header.Parents, from)
	if err != nil {
		return nil, nil, err
	}

	// short circuit if we're traversing history as a new peer
	if following {
		// if a new peer is discovered during back prop, we'll need to send a welcome
		// but not until _after_ HEAD has been updated at the update entry point, where
		// the new peers will be collected
		// NOTE: if from == nil, we've started with an invite, in which case there is
		// no need to handle new peers in this manner (they're sent OUR join)
		if joined != nil && isNew && from != nil && joined.Id != from.Pretty() {
			return addr, joined, nil
		}
		return addr, nil, nil
	}

	// send latest direct to this peer if they could use a merge, i.e., we have newer updates
	head, err := t.GetHead()
	if err != nil {
		return nil, nil, err
	}
	if joined != nil && isNew && head != content.BlockId {
		if err := t.sendWelcome(*joined); err != nil {
			return nil, nil, err
		}
	}

	// handle HEAD
	if _, err := t.handleHead(id, content.Header.Parents); err != nil {
		return nil, nil, err
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return nil, nil, err
		}
	}

	return addr, nil, nil
}

// welcome sends the latest HEAD block
func (t *Thread) sendWelcome(joined repo.Peer) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	// get current HEAD
	head, err := t.GetHead()
	if err != nil {
		return err
	}

	// download it
	serialized, err := util.GetDataAtPath(t.ipfs(), head)
	if err != nil {
		return err
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(serialized, env); err != nil {
		return err
	}

	log.Debugf("welcoming %s at update %s", joined.Id, head)

	// post it
	t.post(env, head, []repo.Peer{joined})

	// all done
	return nil
}

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
)

// JoinInitial creates an outgoing join block for an emtpy thread
func (t *Thread) JoinInitial() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	inviterPkb, err := t.PrivKey.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadJoin{
		Header:    header,
		InviterPk: inviterPkb,
	}

	// commit to ipfs
	_, addr, err := t.commitBlock(content, pb.Message_THREAD_JOIN)
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

	log.Debugf("added JOIN to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// Join creates an outgoing join block
func (t *Thread) Join(inviterPk libp2pc.PubKey, blockId string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	inviterPkb, err := inviterPk.Bytes()
	if err != nil {
		return nil, err
	}
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadJoin{
		Header:    header,
		InviterPk: inviterPkb,
		BlockId:   blockId,
	}

	// commit to ipfs
	env, addr, err := t.commitBlock(content, pb.Message_THREAD_JOIN)
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
	self := inviterPid.Pretty() == t.ipfs().Identity.Pretty()
	if !self {
		newPeer := &repo.Peer{
			Row:      ksuid.New().String(),
			Id:       inviterPid.Pretty(),
			ThreadId: t.Id,
			PubKey:   inviterPkb,
		}
		if err := t.peers().Add(newPeer); err != nil {
			// TODO: #202 (Properly handle database/sql errors)
		}
	}

	// post it
	t.post(env, id, t.Peers())

	log.Debugf("added JOIN to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// HandleJoinBlock handles an incoming join block
func (t *Thread) HandleJoinBlock(from *peer.ID, env *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadJoin, following bool) (mh.Multihash, *repo.Peer, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadJoin)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, nil, err
		}
	}

	// add to ipfs
	addr, err := t.addBlock(env)
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
		if joined != nil && from != nil && joined.Id != from.Pretty() {
			return addr, joined, nil
		}
		return addr, nil, nil
	}

	// send latest direct to this peer if they could use a merge, i.e., we have newer updates
	head, err := t.GetHead()
	if err != nil {
		return nil, nil, err
	}
	if joined != nil && head != content.BlockId {
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

	// check head
	if head == "" {
		return nil
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
	if env.Message == nil {
		// might be a merge block
		message := new(pb.Message)
		if err := proto.Unmarshal(serialized, message); err != nil {
			return err
		}
		var err error
		env, err = t.newEnvelope(message)
		if err != nil {
			return err
		}
	}

	log.Debugf("WELCOME sent to %s at %s", joined.Id, head)

	// post it
	t.post(env, head, []repo.Peer{joined})

	// all done
	return nil
}

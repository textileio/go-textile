package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// Leave creates an outgoing leave block
func (t *Thread) Leave() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadLeave{
		Header: header,
	}

	// commit to ipfs
	env, addr, err := t.commitBlock(content, pb.Message_THREAD_LEAVE)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// index it locally
	if err := t.indexBlock(id, header, repo.LeaveBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// post it
	t.post(env, id, t.Peers())

	// delete blocks
	if err := t.blocks().DeleteByThreadId(t.Id); err != nil {
		return nil, err
	}

	// delete peers
	if err := t.peers().DeleteByThreadId(t.Id); err != nil {
		return nil, err
	}

	// delete notifications
	if err := t.notifications().DeleteBySubjectId(t.Id); err != nil {
		return nil, err
	}

	log.Debugf("added LEAVE to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// HandleLeaveBlock handles an incoming leave block
func (t *Thread) HandleLeaveBlock(from *peer.ID, env *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadLeave, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadLeave)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, err
		}
	}

	// add to ipfs
	addr, err := t.addBlock(env)
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

	// remove peer
	authorId, err := ipfs.IDFromPublicKeyBytes(content.Header.AuthorPk)
	if err != nil {
		return nil, err
	}
	if err := t.peers().Delete(authorId.Pretty(), t.Id); err != nil {
		return nil, err
	}
	if err := t.notifications().DeleteByActorId(t.Id); err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.LeaveBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	newPeers, err := t.FollowParents(content.Header.Parents, from)
	if err != nil {
		return nil, err
	}

	// handle HEAD
	if following {
		return addr, nil
	}
	if _, err := t.handleHead(id, content.Header.Parents); err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return nil, err
		}
	}

	return addr, nil
}

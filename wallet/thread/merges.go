package thread

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	"sort"
	"time"
)

// Merge adds a merge block, which are kept local until subsequent updates, avoiding possibly endless echoes
func (t *Thread) Merge(head string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}

	// add a second parent
	header.Parents = append(header.Parents, head)
	// sort to ensure a deterministic (the order may be reversed on other peers)
	sort.Strings(header.Parents)
	content := &pb.ThreadMerge{
		Parents:  header.Parents,
		ThreadPk: header.ThreadPk,
	}

	// commit envelope to ipfs
	env, addr, err := t.commitBlock(content, pb.Message_THREAD_MERGE)
	if err != nil {
		return nil, err
	}

	// commit envelope contents which does not include author info (we want to decuplicate merge blocks between peers)
	// the resulting hash will be used to index locally, but the content needs to be available
	// on the network in the event it's encountered during a back prop
	ser, err := proto.Marshal(env.Message)
	if err != nil {
		return nil, err
	}
	cid, err := util.PinData(t.ipfs(), bytes.NewReader(ser))
	if err != nil {
		return nil, err
	}
	id := cid.Hash().B58String()

	// add a pin request
	if err := t.putPinRequest(id); err != nil {
		log.Warningf("pin request exists: %s", id)
	}

	// index it locally
	if err := t.indexBlock(id, header, repo.MergeBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	log.Debugf("adding MERGE to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// HandleMergeBlock handles an incoming merge block
func (t *Thread) HandleMergeBlock(from *peer.ID, message *pb.Message, signed *pb.SignedThreadBlock, content *pb.ThreadMerge, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadMerge)
		if err := proto.Unmarshal(signed.Block, content); err != nil {
			return nil, err
		}
	}

	// this time on the receiver end, determine hash of content, not env
	ser, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	cid, err := util.PinData(t.ipfs(), bytes.NewReader(ser))
	if err != nil {
		return nil, err
	}
	id := cid.Hash().B58String()

	// add a pin request
	if err := t.putPinRequest(id); err != nil {
		log.Warningf("pin request exists: %s", id)
	}

	// check if we aleady have this block indexed
	index := t.blocks().Get(id)
	if index != nil {
		return nil, nil
	}

	// index it locally
	// (create a fake header for the indexing step)
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}
	header.Parents = content.Parents
	if err := t.indexBlock(id, header, repo.MergeBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	newPeers, err := t.FollowParents(content.Parents, from)
	if err != nil {
		return nil, err
	}

	// handle HEAD
	if following {
		return cid.Hash(), nil
	}
	if _, err := t.handleHead(id, content.Parents); err != nil {
		return nil, err
	}

	// handle newly discovered peers during back prop, after updating HEAD
	for _, newPeer := range newPeers {
		if err := t.sendWelcome(newPeer); err != nil {
			return nil, err
		}
	}

	return cid.Hash(), nil
}

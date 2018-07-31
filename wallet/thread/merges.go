package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	"time"
)

// Merge adds an outgoing merge block
func (t *Thread) Merge(head string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	header, err := t.newBlockHeader(time.Now())
	if err != nil {
		return nil, err
	}
	header.Parents = append(header.Parents, head)
	content := &pb.ThreadMerge{
		Header: header,
	}

	// commit to ipfs
	message, addr, err := t.commitBlock(content, pb.Message_THREAD_MERGE)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// index it locally
	if err := t.indexBlock(id, header, repo.MergeBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// post it
	t.post(message, id, t.Peers())

	log.Debugf("adding merge to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// HandleMergeBlock handles an incoming merge block
func (t *Thread) HandleMergeBlock(message *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadMerge, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadMerge)
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
		return nil, nil
	}

	// index it locally
	if err := t.indexBlock(id, content.Header, repo.MergeBlock, nil); err != nil {
		return nil, err
	}

	// back prop
	if err := t.FollowParents(content.Header.Parents); err != nil {
		return nil, err
	}

	// handle HEAD
	if following {
		return addr, nil
	}
	if _, err := t.handleHead(id, content.Header.Parents); err != nil {
		return nil, err
	}

	return addr, nil
}

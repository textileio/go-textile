package thread

import (
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// AddComment adds an outgoing comment block
func (t *Thread) AddComment(dataId string, body string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// encrypt comment body with thread pk
	var bodyCipher []byte
	if body != "" {
		var err error
		bodyCipher, err = t.Encrypt([]byte(body))
		if err != nil {
			return nil, err
		}
	}

	// build block
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadAnnotation{
		Header:        header,
		Type:          pb.ThreadAnnotation_COMMENT,
		DataId:        dataId,
		CaptionCipher: bodyCipher,
	}

	// commit to ipfs
	env, addr, err := t.commitBlock(content, pb.Message_THREAD_ANNOTATION)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId:            dataId,
		DataCaptionCipher: bodyCipher,
	}
	if err := t.indexBlock(id, header, repo.CommentBlock, dconf); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// post it
	t.post(env, id, t.Peers())

	log.Debugf("added ANNOTATION to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// AddLike adds an outgoing like block
func (t *Thread) AddLike(dataId, body string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadAnnotation{
		Header: header,
		Type:   pb.ThreadAnnotation_LIKE,
	}

	// commit to ipfs
	env, addr, err := t.commitBlock(content, pb.Message_THREAD_ANNOTATION)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId: dataId,
	}
	if err := t.indexBlock(id, header, repo.LikeBlock, dconf); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// post it
	t.post(env, id, t.Peers())

	log.Debugf("added ANNOTATION to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// HandleAnnotationBlock handles an incoming data block
func (t *Thread) HandleAnnotationBlock(from *peer.ID, env *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadAnnotation, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadAnnotation)
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

	// get the author id
	authorPk, err := libp2pc.UnmarshalPublicKey(content.Header.AuthorPk)
	if err != nil {
		return nil, err
	}
	authorId, err := peer.IDFromPublicKey(authorPk)
	if err != nil {
		return nil, err
	}

	// add author as a new local peer, just in case we haven't found this peer yet.
	// double-check not self in case we're re-discovering the thread
	self := authorId.Pretty() == t.ipfs().Identity.Pretty()
	if !self {
		newPeer := &repo.Peer{
			Row:      ksuid.New().String(),
			Id:       authorId.Pretty(),
			ThreadId: libp2pc.ConfigEncodeKey(content.Header.ThreadPk),
			PubKey:   content.Header.AuthorPk,
		}
		if err := t.peers().Add(newPeer); err != nil {
			// TODO: #202 (Properly handle database/sql errors)
		}
	}

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId: content.DataId,
	}
	var atype repo.BlockType
	switch content.Type {
	case pb.ThreadAnnotation_COMMENT:
		atype = repo.CommentBlock
		dconf.DataCaptionCipher = content.CaptionCipher
	case pb.ThreadAnnotation_LIKE:
		atype = repo.LikeBlock
	}
	if err := t.indexBlock(id, content.Header, atype, dconf); err != nil {
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

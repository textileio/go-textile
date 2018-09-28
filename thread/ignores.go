package thread

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"strings"
)

// Ignore adds an outgoing ignore block targeted at another block to ignore
func (t *Thread) Ignore(blockId string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// adding an ignore specific prefix here to ensure future flexibility
	dataId := fmt.Sprintf("ignore-%s", blockId)

	// build block
	header, err := t.newBlockHeader()
	if err != nil {
		return nil, err
	}
	content := &pb.ThreadIgnore{
		Header: header,
		DataId: dataId,
	}

	// commit to ipfs
	env, addr, err := t.commitBlock(content, pb.Message_THREAD_IGNORE)
	if err != nil {
		return nil, err
	}
	id := addr.B58String()

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId: content.DataId,
	}
	if err := t.indexBlock(id, header, repo.IgnoreBlock, dconf); err != nil {
		return nil, err
	}

	// unpin dataId if present and not part of another thread
	t.unpinBlockData(blockId)

	// update head
	if err := t.updateHead(id); err != nil {
		return nil, err
	}

	// post it
	t.post(env, id, t.Peers())

	// delete notifications
	if err := t.notifications().DeleteByBlockId(blockId); err != nil {
		return nil, err
	}

	log.Debugf("added IGNORE to %s: %s", t.Id, id)

	// all done
	return addr, nil
}

// HandleIgnoreBlock handles an incoming ignore block
func (t *Thread) HandleIgnoreBlock(from *peer.ID, env *pb.Envelope, signed *pb.SignedThreadBlock, content *pb.ThreadIgnore, following bool) (mh.Multihash, error) {
	// unmarshal if needed
	if content == nil {
		content = new(pb.ThreadIgnore)
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

	// delete notifications
	blockId := strings.Replace(content.DataId, "ignore-", "", 1)
	if err := t.notifications().DeleteByBlockId(blockId); err != nil {
		return nil, err
	}

	// index it locally
	dconf := &repo.DataBlockConfig{
		DataId: content.DataId,
	}
	if err := t.indexBlock(id, content.Header, repo.IgnoreBlock, dconf); err != nil {
		return nil, err
	}

	// unpin dataId if present and not part of another thread
	t.unpinBlockData(blockId)

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

func (t *Thread) unpinBlockData(blockId string) {
	block := t.blocks().Get(blockId)
	if block != nil && block.DataId != "" {
		all := t.blocks().List("", -1, "dataId='"+block.DataId+"'")
		if len(all) == 1 {
			// safe to unpin

			switch block.Type {
			case repo.PhotoBlock:
				// unpin image paths
				path := fmt.Sprintf("%s/thumb", block.DataId)
				if err := ipfs.UnpinPath(t.ipfs(), path); err != nil {
					log.Warningf("failed to unpin %s: %s", path, err)
				}
				path = fmt.Sprintf("%s/small", block.DataId)
				if err := ipfs.UnpinPath(t.ipfs(), path); err != nil {
					log.Warningf("failed to unpin %s: %s", path, err)
				}
				path = fmt.Sprintf("%s/meta", block.DataId)
				if err := ipfs.UnpinPath(t.ipfs(), path); err != nil {
					log.Warningf("failed to unpin %s: %s", path, err)
				}
				path = fmt.Sprintf("%s/pk", block.DataId)
				if err := ipfs.UnpinPath(t.ipfs(), path); err != nil {
					log.Warningf("failed to unpin %s: %s", path, err)
				}
			}
		}
	}
}

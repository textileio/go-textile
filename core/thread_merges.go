package core

import (
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema"
)

// merge updates the head block node with an additional parent
func (t *Thread) merge(head string, inbound string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}

	var blockId string
	if t.datastore.Blocks().Get(head) != nil {
		// old block
		blockId = head
	} else {
		links, err := ipfs.LinksAtPath(t.node(), head)
		if err != nil {
			return nil, err
		}
		blink := schema.LinkByName(links, []string{blockLinkName})
		blockId = blink.Cid.Hash().B58String()
	}
	hash, err := t.commitNode(blockId, []string{inbound}, true)
	if err != nil {
		return nil, err
	}

	log.Debugf("merged %s into %s: %s", inbound, t.Id, hash.B58String())

	return hash, nil
}

// handleMergeBlock handles an incoming merge block
func (t *Thread) handleMergeBlock(hash mh.Multihash, block *pb.ThreadBlock, parents []string) error {
	if !t.readable(t.config.Account.Address) {
		return ErrNotReadable
	}
	if !t.readable(block.Header.Address) {
		return ErrNotReadable
	}

	return t.indexBlock(&commitResult{
		hash:    hash,
		header:  block.Header,
		parents: parents,
	}, pb.Block_MERGE, "", "")
}

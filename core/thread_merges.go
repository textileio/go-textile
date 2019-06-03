package core

import (
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/pb"
)

// handleMergeBlock handles an incoming merge block
// Deprecated
func (t *Thread) handleMergeBlock(hash mh.Multihash, block *pb.ThreadBlock) error {
	if !t.readable(t.config.Account.Address) {
		return ErrNotReadable
	}
	if !t.readable(block.Header.Address) {
		return ErrNotReadable
	}

	return nil
}

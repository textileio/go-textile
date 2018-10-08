package storage

import (
	"bytes"
	"github.com/textileio/textile-go/ipfs"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
)

type CafeStorage struct {
	ipfs     *core.IpfsNode
	repoPath string
	store    func(id *cid.Cid) error
}

func NewCafeStorage(ipfs *core.IpfsNode, repoPath string, store func(id *cid.Cid) error) *CafeStorage {
	return &CafeStorage{
		ipfs:     ipfs,
		repoPath: repoPath,
		store:    store,
	}
}

func (s *CafeStorage) Store(message []byte) (ma.Multiaddr, error) {
	// pin the message
	id, err := ipfs.PinData(s.ipfs, bytes.NewReader(message))
	if err != nil {
		return nil, err
	}

	// ask cafe to store it
	if err := s.store(id); err != nil {
		return nil, err
	}

	// return addr for pointer
	return ipfs.MultiaddrFromId(id.Hash().B58String())
}

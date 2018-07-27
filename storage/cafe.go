package storage

import (
	"bytes"
	"github.com/textileio/textile-go/util"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
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
	id, err := util.PinData(s.ipfs, bytes.NewReader(message))
	if err != nil {
		return nil, err
	}

	// ask cafe to store it
	if err := s.store(id); err != nil {
		return nil, err
	}

	// return addr for pointer
	return util.MultiaddrFromId(id.Hash().B58String())
}

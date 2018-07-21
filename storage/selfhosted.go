package storage

import (
	"bytes"
	"context"
	"github.com/textileio/textile-go/wallet/util"
	routing "gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"sync"
	"time"
)

type SelfHostedStorage struct {
	ipfs     *core.IpfsNode
	repoPath string
	store    func(peerId string, ids []cid.Cid) error
}

func NewSelfHostedStorage(ipfs *core.IpfsNode, repoPath string, store func(peerId string, ids []cid.Cid) error) *SelfHostedStorage {
	return &SelfHostedStorage{
		ipfs:     ipfs,
		repoPath: repoPath,
		store:    store,
	}
}

func (s *SelfHostedStorage) Store(ciphertext []byte) (ma.Multiaddr, error) {
	// pin the message
	id, err := util.PinData(s.ipfs, bytes.NewReader(ciphertext))
	if err != nil {
		return nil, err
	}

	// ask key neighbors to store the message content
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	dht := s.ipfs.Routing.(*routing.IpfsDHT)
	peers, err := dht.GetClosestPeers(ctx, id.KeyString())
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for p := range peers {
		wg.Add(1)
		go func(pid peer.ID) {
			defer wg.Done()
			s.store(pid.Pretty(), []cid.Cid{*id})
		}(p)
	}
	wg.Wait()

	return util.MultiaddrFromId(id.Hash().B58String())
}

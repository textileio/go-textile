package storage

import (
	"github.com/textileio/textile-go/wallet/util"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"sync"
	routing "gx/ipfs/QmVW4cqbibru3hXA1iRmg85Fk7z9qML9k176CYQaMXVCrP/go-libp2p-kad-dht"
	"context"
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

func (s *SelfHostedStorage) Store(peerID peer.ID, ciphertext []byte) (ma.Multiaddr, error) {
	// create a virtual directory for the message
	dirb := uio.NewDirectory(s.ipfs.DAG)
	if err := util.AddFileToDirectory(s.ipfs, dirb, ciphertext, "msg"); err != nil {
		return nil, err
	}

	// pin the directory
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := util.PinDirectory(s.ipfs, dir, []string{}); err != nil {
		return nil, err
	}

	// ask key neighbors to store the message content
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	dht := s.ipfs.Routing.(*routing.IpfsDHT)
	peers, err := dht.GetClosestPeers(ctx, dir.Cid().KeyString())
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for p := range peers {
		wg.Add(1)
		go func(pid peer.ID) {
			defer wg.Done()
			s.store(pid.Pretty(), []cid.Cid{*dir.Cid()})
		}(p)
	}
	wg.Wait()

	return ma.NewMultiaddr("/ipfs/" + dir.Cid().Hash().B58String() + "/")
}

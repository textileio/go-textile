package storage

import (
	"github.com/textileio/textile-go/wallet/util"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
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

	// TODO: figure out where to push (thinking we can use a random subset of nodes to push to)
	//for _, peer := range s.pushNodes {
	//	go s.store(peer.Pretty(), []cid.Cid{*dir.Cid()})
	//}
	return ma.NewMultiaddr("/ipfs/" + dir.Cid().Hash().B58String() + "/")
}

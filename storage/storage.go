package storage

import (
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
)

type OfflineMessagingStorage interface {
	Store(message []byte) (ma.Multiaddr, error)
}

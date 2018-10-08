package storage

import (
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
)

type OfflineMessagingStorage interface {
	Store(message []byte) (ma.Multiaddr, error)
}

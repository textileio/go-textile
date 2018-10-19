package core

import (
	"fmt"
	"github.com/textileio/textile-go/repo"
)

// Overview is a wallet overview object
type Overview struct {
	SwarmSize        int `json:"swarm_size"`
	AccountPeerCount int `json:"account_peer_count"`
	ThreadCount      int `json:"thread_count"`
	PhotoCount       int `json:"photo_count"`
	ContactCount     int `json:"contact_count"`
}

// Overview returns an overview object
func (t *Textile) Overview() (*Overview, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}

	// collect stats
	swarm, err := t.Peers()
	if err != nil {
		return nil, err
	}
	threads := t.datastore.Threads().Count("")
	photos := t.datastore.Blocks().Count(fmt.Sprintf("type=%d", repo.PhotoBlock))
	contacts := t.datastore.ThreadPeers().Count(true)

	return &Overview{
		SwarmSize:        len(swarm),
		AccountPeerCount: 0,
		ThreadCount:      threads,
		PhotoCount:       photos,
		ContactCount:     contacts,
	}, nil
}

package core

import (
	"fmt"
	"github.com/textileio/textile-go/repo"
)

// Overview is a wallet overview object
type Overview struct {
	SwarmSize    int `json:"swarm_size"`
	DeviceCount  int `json:"device_count"`
	ThreadCount  int `json:"thread_count"`
	PhotoCount   int `json:"photo_count"`
	ContactCount int `json:"contact_count"`
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
	devices := t.datastore.AccountPeers().Count("")
	threads := t.datastore.Threads().Count("")
	photos := t.datastore.Blocks().Count(fmt.Sprintf("type=%d", repo.PhotoBlock))
	contacts := t.datastore.ThreadPeers().Count("", true)

	return &Overview{
		SwarmSize:    len(swarm),
		DeviceCount:  devices,
		ThreadCount:  threads,
		PhotoCount:   photos,
		ContactCount: contacts,
	}, nil
}

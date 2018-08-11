package wallet

import (
	"fmt"
	"github.com/textileio/textile-go/repo"
)

type Stats struct {
	SwarmSize   int `json:"swarm_size"`
	DeviceCount int `json:"device_count"`
	ThreadCount int `json:"thread_count"`
	PhotoCount  int `json:"photo_count"`
	PeerCount   int `json:"peer_count"`
}

func (w *Wallet) GetStats() (*Stats, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}

	// collect stats
	swarm, err := w.Peers()
	if err != nil {
		return nil, err
	}
	devices := w.datastore.Devices().Count("")
	threads := w.datastore.Threads().Count("")
	photos := w.datastore.Blocks().Count(fmt.Sprintf("type=%d", repo.PhotoBlock))
	peers := w.datastore.Peers().Count("", true)

	return &Stats{
		SwarmSize:   len(swarm),
		DeviceCount: devices,
		ThreadCount: threads,
		PhotoCount:  photos,
		PeerCount:   peers,
	}, nil
}

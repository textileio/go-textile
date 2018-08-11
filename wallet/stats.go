package wallet

import (
	"fmt"
	"github.com/textileio/textile-go/repo"
)

type Stats struct {
	SwarmSize   int `json:"swarm_size"`
	DeviceCount int `json:"device_count"`
	ThreadCount int `json:"photo_count"`
	PhotoCount  int `json:"photo_count"`
	PeerCount   int `json:"peer_count"`
}

func (w *Wallet) GetStats() (*Stats, error) {
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

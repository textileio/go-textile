package wallet

import (
	"fmt"
	"github.com/textileio/textile-go/repo"
)

// Overview is a wallet overview object
type Overview struct {
	SwarmSize     int `json:"swarm_size"`
	DeviceCount   int `json:"device_count"`
	ThreadCount   int `json:"thread_count"`
	PhotoCount    int `json:"photo_count"`
	ContactsCount int `json:"contacts_count"`
}

// Overview returns an overview object
func (w *Wallet) Overview() (*Overview, error) {
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
	contacts := w.datastore.Peers().Count("", true)

	return &Overview{
		SwarmSize:     len(swarm),
		DeviceCount:   devices,
		ThreadCount:   threads,
		PhotoCount:    photos,
		ContactsCount: contacts,
	}, nil
}

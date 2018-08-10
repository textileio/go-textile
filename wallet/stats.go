package wallet

type Stats struct {
	SwarmSize   int `json:"swarm_size"`
	DeviceCount int `json:"device_count"`
	ThreadCount int `json:"photo_count"`
	PhotoCount  int `json:"photo_count"`
	PeerCount   int `json:"peer_count"`
}

//func (w *Wallet) GetStats() (string, error) {
//
//}

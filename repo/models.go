package repo

import "time"

type SettingsData struct {
	Version *string `json:"version"`
}

type PhotoSet struct {
	Cid       string    `json:"cid"`
	Timestamp time.Time `json:"timestamp"`
}

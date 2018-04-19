package repo

import (
	"github.com/textileio/textile-go/repo/photos"
)

type SettingsData struct {
	Version *string `json:"version"`
}

type PhotoSet struct {
	Cid      string          `json:"cid"`
	LastCid  string          `json:"last_cid"`
	MetaData photos.Metadata `json:"metadata"`
}

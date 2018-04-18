package repo

import (
	"github.com/textileio/textile-go/repo/wallet"
)

type SettingsData struct {
	Version *string `json:"version"`
}

type PhotoSet struct {
	Cid      string           `json:"cid"`
	LastCid  string           `json:"last_cid"`
	MetaData wallet.PhotoData `json:"metadata"`
}

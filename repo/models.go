package repo

import (
	"github.com/textileio/textile-go/repo/photos"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

type SettingsData struct {
	Version *string `json:"version"`
}

type PhotoSet struct {
	Cid      string          `json:"cid"`
	LastCid  string          `json:"last_cid"`
	AlbumID  string          `json:"album_id"`
	MetaData photos.Metadata `json:"metadata"`
	Caption  string          `json:"caption,omitempty"`
	IsLocal  bool            `json:"is_local"`
}

type PhotoAlbum struct {
	Id       string         `json:"id"`
	Key      libp2p.PrivKey `json:"key"`
	Mnemonic string         `json:"mnemonic"`
	Name     string         `json:"name"`
}

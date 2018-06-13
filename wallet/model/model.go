package model

import (
	"github.com/textileio/textile-go/net"
	"time"
)

const ThumbnailWidth = 300

type Metadata struct {
	Username string    `json:"un,omitempty"`
	Created  time.Time `json:"cts,omitempty"`
	Added    time.Time `json:"ats,omitempty"`
}

type FileMetadata struct {
	Metadata
	Name string `json:"name,omitempty"`
	Ext  string `json:"ext,omitempty"`
}

type AddResult struct {
	Id            string
	Key           []byte
	RemoteRequest *net.MultipartRequest
}

type PhotoMetadata struct {
	FileMetadata
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lon,omitempty"`
}

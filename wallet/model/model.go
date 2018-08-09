package model

import "time"

type ImageSize int

const (
	ThumbnailSize ImageSize = 100
	SmallSize               = 320
	MediumSize              = 800
	LargeSize               = 1600
)

type Profile struct {
	Id       string `json:"id"`
	Username string `json:"username,omitempty"`
	AvatarId string `json:"avatar_id,omitempty"`
}

type Metadata struct {
	Version string    `json:"version"`
	PeerId  string    `json:"peer_id"`
	Created time.Time `json:"created,omitempty"`
	Added   time.Time `json:"added"`
}

type FileMetadata struct {
	Metadata
	Name string `json:"name"`
	Ext  string `json:"extension"`
}

type PhotoMetadata struct {
	FileMetadata
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	OriginalFormat string  `json:"original_format"`
	EncodingFormat string  `json:"encoding_format"`
	Latitude       float64 `json:"latitude,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
}

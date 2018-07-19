package model

import "time"

type Profile struct {
	Id       string `json:"id"`
	Username string `json:"username,omitempty"`
	AvatarId string `json:"avatar_id"`
}

const ThumbnailWidth = 300

type Metadata struct {
	Version  string    `json:"version"`
	Username string    `json:"username,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Added    time.Time `json:"added"`
}

type FileMetadata struct {
	Metadata
	Name string `json:"name"`
	Ext  string `json:"extension"`
}

type PhotoMetadata struct {
	FileMetadata
	Format          string  `json:"format"`
	ThumbnailFormat string  `json:"format_thumb"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Latitude        float64 `json:"latitude,omitempty"`
	Longitude       float64 `json:"longitude,omitempty"`
}

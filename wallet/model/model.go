package model

import (
	"time"
)

type Profile struct {
	Id       string `json:"id"`
	Username string `json:"un,omitempty"`
}

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

type PhotoMetadata struct {
	FileMetadata
	Format          string  `json:"fmt"`
	ThumbnailFormat string  `json:"fmt_tn"`
	Latitude        float64 `json:"lat,omitempty"`
	Longitude       float64 `json:"lon,omitempty"`
}

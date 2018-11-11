package images

//
//import (
//	"github.com/rwcarlsen/goexif/exif"
//	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
//	"io"
//	"path/filepath"
//	"strings"
//	"time"
//)
//
//var log = logging.Logger("tex-images")
//
//// Format enumerates the type of images currently supported
//type Format string
//
//const (
//	JPEG Format = "jpeg"
//	PNG  Format = "png"
//	GIF  Format = "gif"
//)
//
//// ImageSize enumerates the sizes that are currently auto encoded
//type ImageSize int
//
//const (
//	ThumbnailSize ImageSize = 100
//	SmallSize               = 320
//	MediumSize              = 800
//	LargeSize               = 1600
//)
//
//// ImageSizeForMinWidth is used to return a supported size for the given min width
//func ImageSizeForMinWidth(width int) ImageSize {
//	if width <= 100 {
//		return ThumbnailSize
//	} else if width <= 320 {
//		return SmallSize
//	} else if width <= 800 {
//		return MediumSize
//	} else {
//		return LargeSize
//	}
//}
//
//// Imagepath enumerates the DAG paths for supported sizes
//type ImagePath string
//
//const (
//	ThumbnailPath ImagePath = "/thumb"
//	SmallPath               = "/small"
//	MediumPath              = "/medium"
//	LargePath               = "/photo"
//)
//
//// ImagePathForSize returns the DAG path for a given supported size
//func ImagePathForSize(size ImageSize) ImagePath {
//	switch size {
//	case ThumbnailSize:
//		return ThumbnailPath
//	case SmallSize:
//		return SmallPath
//	case MediumSize:
//		return MediumPath
//	default:
//		return LargePath
//	}
//}
//
//// Metadata (mostly exif data) stripped from images, pre-encoding.
//// NOTE: Metadata is encrypted and stored alongside encoded images (in the photo set DAG).
//type Metadata struct {
//	Version        string    `json:"version"`
//	Created        time.Time `json:"created,omitempty"`
//	Added          time.Time `json:"added"`
//	Name           string    `json:"name"`
//	Ext            string    `json:"extension"`
//	Width          int       `json:"width"`
//	Height         int       `json:"height"`
//	OriginalFormat string    `json:"original_format"`
//	EncodingFormat string    `json:"encoding_format"`
//	Latitude       float64   `json:"latitude,omitempty"`
//	Longitude      float64   `json:"longitude,omitempty"`
//}
//
//// NewMetadata returns a new image meta data object
//func NewMetadata(
//	reader io.Reader,
//	path string,
//	ext string,
//	format Format,
//	encodingFormat Format,
//	width int,
//	height int,
//	version string,
//) (Metadata, error) {
//	var created time.Time
//	var lat, lon float64
//	x, err := exif.Decode(reader)
//	if err == nil {
//		// time taken
//		createdTmp, err := x.DateTime()
//		if err == nil {
//			created = createdTmp
//		}
//		// coords taken
//		latTmp, lonTmp, err := x.LatLong()
//		if err == nil {
//			lat, lon = latTmp, lonTmp
//		}
//	}
//	return Metadata{
//		Version:        version,
//		Created:        created,
//		Added:          time.Now(),
//		Name:           strings.TrimSuffix(filepath.Base(path), ext),
//		Ext:            ext,
//		OriginalFormat: string(format),
//		EncodingFormat: string(encodingFormat),
//		Width:          width,
//		Height:         height,
//		Latitude:       lat,
//		Longitude:      lon,
//	}, nil
//}
//

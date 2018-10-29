package images

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

var log = logging.Logger("tex-images")

// Format enumerates the type of images currently supported
type Format string

const (
	JPEG Format = "jpeg"
	PNG  Format = "png"
	GIF  Format = "gif"
)

// ImageSize enumerates the sizes that are currently auto encoded
type ImageSize int

const (
	ThumbnailSize ImageSize = 100
	SmallSize               = 320
	MediumSize              = 800
	LargeSize               = 1600
)

// ImageSizeForMinWidth is used to return a supported size for the given min width
func ImageSizeForMinWidth(width int) ImageSize {
	if width <= 100 {
		return ThumbnailSize
	} else if width <= 320 {
		return SmallSize
	} else if width <= 800 {
		return MediumSize
	} else {
		return LargeSize
	}
}

// Imagepath enumerates the DAG paths for supported sizes
type ImagePath string

const (
	ThumbnailPath ImagePath = "/thumb"
	SmallPath               = "/small"
	MediumPath              = "/medium"
	LargePath               = "/photo"
)

// ImagePathForSize returns the DAG path for a given supported size
func ImagePathForSize(size ImageSize) ImagePath {
	switch size {
	case ThumbnailSize:
		return ThumbnailPath
	case SmallSize:
		return SmallPath
	case MediumSize:
		return MediumPath
	default:
		return LargePath
	}
}

// Metadata (mostly exif data) stripped from images, pre-encoding.
// NOTE: Metadata is encrypted and stored alongside encoded images (in the photo set DAG).
type Metadata struct {
	Version        string    `json:"version"`
	Created        time.Time `json:"created,omitempty"`
	Added          time.Time `json:"added"`
	Name           string    `json:"name"`
	Ext            string    `json:"extension"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	OriginalFormat string    `json:"original_format"`
	EncodingFormat string    `json:"encoding_format"`
	Latitude       float64   `json:"latitude,omitempty"`
	Longitude      float64   `json:"longitude,omitempty"`
}

// NewMetadata returns a new image meta data object
func NewMetadata(
	reader io.Reader,
	path string,
	ext string,
	format Format,
	encodingFormat Format,
	width int,
	height int,
	version string,
) (Metadata, error) {
	var created time.Time
	var lat, lon float64
	x, err := exif.Decode(reader)
	if err == nil {
		// time taken
		createdTmp, err := x.DateTime()
		if err == nil {
			created = createdTmp
		}
		// coords taken
		latTmp, lonTmp, err := x.LatLong()
		if err == nil {
			lat, lon = latTmp, lonTmp
		}
	}
	return Metadata{
		Version:        version,
		Created:        created,
		Added:          time.Now(),
		Name:           strings.TrimSuffix(filepath.Base(path), ext),
		Ext:            ext,
		OriginalFormat: string(format),
		EncodingFormat: string(encodingFormat),
		Width:          width,
		Height:         height,
		Latitude:       lat,
		Longitude:      lon,
	}, nil
}

// DecodeImage returns a cleaned reader from an image file
func DecodeImage(file multipart.File) (*bytes.Reader, *Format, *image.Point, error) {
	img, formatStr, err := image.Decode(file)
	if err != nil {
		return nil, nil, nil, err
	}
	format := Format(formatStr)
	size := img.Bounds().Size()

	var reader *bytes.Reader
	file.Seek(0, 0)
	if format != "gif" {
		// decode exif
		exf := DecodeExif(file)
		img, err = correctOrientation(img, exf)
		if err != nil {
			return nil, nil, nil, err
		}

		// re-encoding will remove exif
		reader, err = encodeSingleImage(img, format)
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		fileb, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, nil, nil, err
		}
		reader = bytes.NewReader(fileb)
	}

	return reader, &format, &size, nil
}

// EncodeImage creates a jpeg|gif thumbnail from an image
// - jpeg quality is currently the default (75/100)
func EncodeImage(reader io.Reader, format Format, size ImageSize) ([]byte, error) {
	var result []byte
	width := int(size)
	switch format {
	case JPEG:
		img, _, err := image.Decode(reader)
		if err != nil {
			return nil, err
		}
		if img.Bounds().Size().X < width {
			width = img.Bounds().Size().X
		}
		resized := imaging.Resize(img, width, 0, imaging.Lanczos)
		buff := new(bytes.Buffer)
		if err = jpeg.Encode(buff, resized, nil); err != nil {
			return nil, err
		}
		result = buff.Bytes()
	case GIF:
		img, err := gif.DecodeAll(reader)
		if err != nil {
			return nil, err
		}
		if len(img.Image) == 0 {
			return nil, errors.New("gif does not have any frames")
		}
		firstFrame := img.Image[0].Bounds()
		if firstFrame.Dx() < width {
			width = firstFrame.Dx()
		}
		rect := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
		rgba := image.NewRGBA(rect)
		for index, frame := range img.Image {
			bounds := frame.Bounds()
			draw.Draw(rgba, bounds, frame, bounds.Min, draw.Over)
			img.Image[index] = imageToPaletted(imaging.Resize(rgba, width, 0, imaging.Lanczos))
		}
		img.Config.Width = img.Image[0].Bounds().Dx()
		img.Config.Height = img.Image[0].Bounds().Dy()
		buff := new(bytes.Buffer)
		if err = gif.EncodeAll(buff, img); err != nil {
			return nil, err
		}
		result = buff.Bytes()
	}
	return result, nil
}

// DecodeExif returns exif data from a reader if present
func DecodeExif(reader io.Reader) *exif.Exif {
	exf, err := exif.Decode(reader)
	if err != nil {
		return nil
	}
	return exf
}

// correctOrientation returns a copy of an image (jpg|png|gif) with exif removed
func correctOrientation(img image.Image, exf *exif.Exif) (image.Image, error) {
	if exf == nil {
		return img, nil
	}
	orient, err := exf.Get(exif.Orientation)
	if err != nil && err != exif.TagNotPresentError(exif.Orientation) {
		return nil, err
	}
	if orient != nil {
		log.Debugf("image had orientation %s", orient.String())
		img = reverseOrientation(img, orient.String())
	} else {
		log.Debugf("had no orientation - using 1")
		img = reverseOrientation(img, "1")
	}
	return img, nil
}

// encodeSingleImage creates a reader from an image
func encodeSingleImage(img image.Image, format Format) (*bytes.Reader, error) {
	writer := &bytes.Buffer{}
	var err error
	switch format {
	case JPEG:
		err = jpeg.Encode(writer, img, &jpeg.Options{Quality: 100})
	case PNG:
		// NOTE: while PNGs don't technically have exif data,
		// they can contain meta data with sensitive info
		err = png.Encode(writer, img)
	default:
		err = errors.New("unrecognized image format")
	}
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(writer.Bytes()), nil
}

// reverseOrientation transforms the given orientation to 1
func reverseOrientation(img image.Image, orientation string) *image.NRGBA {
	switch orientation {
	case "1":
		return imaging.Clone(img)
	case "2":
		return imaging.FlipV(img)
	case "3":
		return imaging.Rotate180(img)
	case "4":
		return imaging.Rotate180(imaging.FlipV(img))
	case "5":
		return imaging.Rotate270(imaging.FlipV(img))
	case "6":
		return imaging.Rotate270(img)
	case "7":
		return imaging.Rotate90(imaging.FlipV(img))
	case "8":
		return imaging.Rotate90(img)
	}
	log.Warningf("unknown orientation %s, expected 1-8", orientation)
	return imaging.Clone(img)
}

// imageToPaletted convert Image to Paletted for GIF handling
func imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}

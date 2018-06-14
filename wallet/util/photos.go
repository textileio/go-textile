package util

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/textileio/textile-go/wallet/model"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ThumbnailFormat int

const (
	JPEG ThumbnailFormat = iota
	GIF
)

// DecodeImage returns a cleaned reader from an image file
func DecodeImage(file *os.File) (*bytes.Reader, string, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	var reader *bytes.Reader
	file.Seek(0, 0)
	if format != "gif" {
		// decode exif
		exf := DecodeExif(file)
		img, err = correctOrientation(img, exf)
		if err != nil {
			return nil, "", err
		}

		// re-encoding will remove exif
		reader, err = encodeSingleImage(img, format)
		if err != nil {
			return nil, "", err
		}
	} else {
		fileb, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, "", err
		}
		reader = bytes.NewReader(fileb)
	}
	return reader, format, nil
}

// GetMetaData reads any available meta/exif data from a photo
// TODO: get image size info
func GetMetadata(reader io.Reader, path string, ext string, username string) (model.PhotoMetadata, error) {
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
	meta := model.PhotoMetadata{
		FileMetadata: model.FileMetadata{
			Metadata: model.Metadata{
				Username: username,
				Created:  created,
				Added:    time.Now(),
			},
			Name: strings.TrimSuffix(filepath.Base(path), ext),
			Ext:  ext,
		},
		Latitude:  lat,
		Longitude: lon,
	}
	return meta, nil
}

// MakeThumbnail creates a jpeg|gif thumbnail from an image
func MakeThumbnail(reader io.Reader, format ThumbnailFormat, width int) ([]byte, error) {
	var result []byte
	switch format {
	case JPEG:
		img, _, err := image.Decode(reader)
		if err != nil {
			return nil, err
		}
		thumb := imaging.Resize(img, width, 0, imaging.Lanczos)
		buff := new(bytes.Buffer)
		if err = jpeg.Encode(buff, thumb, nil); err != nil {
			return nil, err
		}
		result = buff.Bytes()
	case GIF:
		img, err := gif.DecodeAll(reader)
		if err != nil {
			return nil, err
		}
		firstFrame := img.Image[0].Bounds()
		rect := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
		rgba := image.NewRGBA(rect)
		for index, frame := range img.Image {
			bounds := frame.Bounds()
			draw.Draw(rgba, bounds, frame, bounds.Min, draw.Over)
			img.Image[index] = imageToPaletted(imaging.Resize(rgba, width, 0, imaging.Box))
		}
		aspect := float64(img.Config.Width) / float64(img.Config.Height)
		img.Config.Width = width
		img.Config.Height = int(float64(width) / aspect)
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
	if err != nil {
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
func encodeSingleImage(img image.Image, format string) (*bytes.Reader, error) {
	writer := &bytes.Buffer{}
	var err error
	switch format {
	case "jpeg":
		err = jpeg.Encode(writer, img, &jpeg.Options{Quality: 100})
	case "png":
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

func imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}

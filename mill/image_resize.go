package mill

import (
	"bytes"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"time"
)

// Format enumerates the type of images currently supported
type Format string

const (
	JPEG Format = "jpeg"
	PNG  Format = "png"
	GIF  Format = "gif"
)

type ImageResizeOpts struct {
	Width   int `json:"width"`
	Quality int `json:"quality"`
}

type ImageResize struct {
	Opts ImageResizeOpts
}

func (m *ImageResize) ID() string {
	return "/image/resize"
}

func (m *ImageResize) AcceptMedia(media string) error {
	return accepts([]string{
		"image/jpeg",
		"image/png",
		"image/gif",
	}, media)
}

func (m *ImageResize) Mill(file multipart.File, name string) (*Result, error) {
	img, formatStr, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	format := Format(formatStr)
	size := img.Bounds().Size()

	clean, err := removeExif(file, img, format)
	if err != nil {
		return nil, err
	}

	buff, err := encodeImage(clean, format, m.Opts.Width, m.Opts.Quality)
	if err != nil {
		return nil, err
	}

	return &Result{
		File: buff.Bytes(),
		Meta: map[string]interface{}{
			"width":   size.X,
			"height":  size.Y,
			"created": time.Now(),
		},
	}, nil
}

// removeExif strips exif data from an image
func removeExif(reader io.Reader, img image.Image, format Format) (io.Reader, error) {
	if format == GIF {
		return reader, nil
	}

	exf, _ := exif.Decode(reader)
	var err error
	img, err = correctOrientation(img, exf)
	if err != nil {
		return nil, err
	}

	// re-encoding will remove any exif
	return encodeSingleImage(img, format)
}

// encodeImage creates a jpeg|gif from reader (quality applies to jpeg only)
// NOTE: format is the reader image format, destination format is chosen accordingly.
func encodeImage(reader io.Reader, format Format, width int, quality int) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)

	if format != GIF {
		// encode to jpeg
		img, _, err := image.Decode(reader)
		if err != nil {
			return nil, err
		}

		if img.Bounds().Size().X < width {
			width = img.Bounds().Size().X
		}

		resized := imaging.Resize(img, width, 0, imaging.Lanczos)
		if err = jpeg.Encode(buff, resized, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}

	} else {
		// encode to gif
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

		if err = gif.EncodeAll(buff, img); err != nil {
			return nil, err
		}
	}

	return buff, nil
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
		img = reverseOrientation(img, orient.String())
	} else {
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

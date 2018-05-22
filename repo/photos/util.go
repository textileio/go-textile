package photos

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// ImagePathWithoutExif makes a copy of image (jpg,png or gif) and applies
// all necessary operation to reverse its orientation to 1
// The result is a image with corrected orientation and without
// exif data.
// ReadImage makes a copy of image (jpg,png or gif) and applies
// all necessary operation to reverse its orientation to 1
// The result is a image with corrected orientation and without
// exif data.
func ImagePathWithoutExif(fpath string) (*bytes.Buffer, error) {
	var img image.Image
	var err error
	writer := &bytes.Buffer{}
	filetype := strings.ToLower(filepath.Ext(fpath))
	// deal with image
	ifile, err := os.Open(fpath)
	if err != nil {
		log.Errorf("could not open file for image transformation: %s", fpath)
		return nil, err
	}
	defer ifile.Close()
	if filetype == ".jpg" || filetype == ".jpeg" {
		img, err = jpeg.Decode(ifile)
		if err != nil {
			return nil, err
		}
	} else if filetype == ".png" {
		img, err = png.Decode(ifile)
		if err != nil {
			return nil, err
		}
	} else if filetype == ".gif" {
		img, err = gif.Decode(ifile)
		if err != nil {
			return nil, err
		}
	}
	// deal with exif
	efile, err := os.Open(fpath)
	if err != nil {
		log.Debugf("could not open file for exif decoder: %s", fpath)
	}
	defer efile.Close()
	x, err := exif.Decode(efile)
	if err != nil {
		if x == nil {
			// ignore - image exif data has been already stripped
		}
		log.Debugf("failed reading exif data in [%s]: %s", fpath, err.Error())
	}
	if x != nil {
		orient, _ := x.Get(exif.Orientation)
		if orient != nil {
			log.Infof("%s had orientation %s", fpath, orient.String())
			img = reverseOrientation(img, orient.String())
		} else {
			log.Errorf("%s had no orientation - implying 1", fpath)
			img = reverseOrientation(img, "1")
		}
	}
	if filetype == ".jpg" || filetype == ".jpeg" {
		err = jpeg.Encode(writer, img, nil)
		if err != nil {
			return nil, err
		}
	} else if filetype == ".png" {
		err = png.Encode(writer, img)
		if err != nil {
			return nil, err
		}
	} else if filetype == ".gif" {
		err = gif.Encode(writer, img, nil)
		if err != nil {
			return nil, err
		}
	}
	return writer, nil
}

// reverseOrientation amply`s what ever operation is necessary to transform given orientation
// to the orientation 1
func reverseOrientation(img image.Image, o string) *image.NRGBA {
	switch o {
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
	log.Errorf("unknown orientation %s, expect 1-8", o)
	return imaging.Clone(img)
}

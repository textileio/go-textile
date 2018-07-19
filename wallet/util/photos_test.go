package util_test

import (
	"bytes"
	"fmt"
	"github.com/textileio/textile-go/wallet/model"
	. "github.com/textileio/textile-go/wallet/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testImage struct {
	path        string
	name        string
	ext         string
	format      Format
	thumbFormat Format
	hasExif     bool
	width       int
	height      int
}

var images = []testImage{
	{
		path:        "../testdata/image.jpg",
		name:        "image",
		ext:         ".jpg",
		format:      JPEG,
		thumbFormat: JPEG,
		hasExif:     true,
		width:       3024,
		height:      4032,
	},
	{
		path:        "../testdata/image.png",
		name:        "image",
		ext:         ".png",
		format:      PNG,
		thumbFormat: JPEG,
		hasExif:     false,
		width:       3024,
		height:      4032,
	},
	{
		path:        "../testdata/image.gif",
		name:        "image",
		ext:         ".gif",
		format:      GIF,
		thumbFormat: GIF,
		hasExif:     false,
		width:       320,
		height:      240,
	},
}

func Test_DecodeImage(t *testing.T) {
	for _, i := range images {
		file, err := os.Open(i.path)
		if err != nil {
			t.Fatal(err)
		}

		reader, format, size, err := DecodeImage(file)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
		if *format != i.format {
			t.Errorf("wrong format")
		}
		if size.X != i.width {
			t.Errorf("wrong width")
		}
		if size.Y != i.height {
			t.Errorf("wrong height")
		}

		// ensure exif was removed
		reader.Seek(0, 0)
		exf2 := DecodeExif(reader)
		if exf2 != nil {
			t.Error("exif data not removed")
		}
	}
}

func Test_MakeMetadata(t *testing.T) {
	for _, i := range images {
		file, err := os.Open(i.path)
		if err != nil {
			t.Fatal(err)
		}

		fpath := file.Name()
		ext := strings.ToLower(filepath.Ext(fpath))

		meta, err := MakeMetadata(file, fpath, ext, i.format, i.thumbFormat, i.width, i.height, "bob", "1.0.0")
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
		if meta.Name != "image" {
			t.Error("bad photo meta name")
		}
		if meta.Ext != i.ext {
			t.Error("bad photo meta extension")
		}
		if meta.Width != i.width {
			t.Error("bad photo meta width")
		}
		if meta.Height != i.height {
			t.Error("bad photo meta height")
		}
		if meta.Added.IsZero() {
			t.Error("bad photo meta added")
		}
		if (i.hasExif && meta.Latitude == 0) || (!i.hasExif && meta.Latitude != 0) {
			t.Error("bad photo meta latitude")
		}
		if (i.hasExif && meta.Longitude == 0) || (!i.hasExif && meta.Longitude != 0) {
			t.Error("bad photo meta longitude")
		}
	}
}

func Test_MakeThumbnail(t *testing.T) {
	for _, i := range images {
		file, err := os.Open(i.path)
		if err != nil {
			t.Fatal(err)
		}

		var thumbFormat Format
		var thumbExt string
		if i.format == "gif" {
			thumbFormat = GIF
			thumbExt = ".gif"
		} else {
			thumbFormat = JPEG
			thumbExt = ".jpeg"
		}
		fileb, err := ioutil.ReadAll(file)
		reader := bytes.NewReader(fileb)

		thumb, err := MakeThumbnail(reader, thumbFormat, model.ThumbnailWidth)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
		err = ioutil.WriteFile(fmt.Sprintf("/tmp/img_%s%s", i.format, thumbExt), thumb, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
}

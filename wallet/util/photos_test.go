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
	path    string
	name    string
	ext     string
	format  string
	hasExif bool
}

var images = []testImage{
	{
		path:    "../testdata/image.jpg",
		name:    "image",
		ext:     ".jpg",
		format:  "jpeg",
		hasExif: true,
	},
	{
		path:    "../testdata/image.png",
		name:    "image",
		ext:     ".png",
		format:  "png",
		hasExif: false,
	},
	{
		path:    "../testdata/image.gif",
		name:    "image",
		ext:     ".gif",
		format:  "gif",
		hasExif: false,
	},
}

func Test_DecodeImage(t *testing.T) {
	for _, i := range images {
		file, err := os.Open(i.path)
		if err != nil {
			t.Fatal(err)
		}

		reader, format, err := DecodeImage(file)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
		if format != i.format {
			t.Errorf("wrong format")
		}

		// ensure exif was removed
		reader.Seek(0, 0)
		exf2 := DecodeExif(reader)
		if exf2 != nil {
			t.Error("exif data not removed")
		}
	}
}

func Test_GetMetadata(t *testing.T) {
	for _, i := range images {
		file, err := os.Open(i.path)
		if err != nil {
			t.Fatal(err)
		}

		fpath := file.Name()
		ext := strings.ToLower(filepath.Ext(fpath))

		meta, err := MakeMetadata(file, fpath, ext, "bob")
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

		var thumbFormat ThumbnailFormat
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

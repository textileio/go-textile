package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

type Archive struct {
	Path string `json:"path"`

	wr io.Writer    `json:"-"`
	gw *gzip.Writer `json:"-"`
	tw *tar.Writer  `json:"-"`
}

func NewArchive(path *string) (*Archive, error) {
	var writer io.Writer
	var file *os.File

	if path == nil {
		writer = &bytes.Buffer{}
	} else {
		name := fmt.Sprintf("%s.tar.gz", *path)
		var err error
		file, err = os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		writer = file
	}

	gw := gzip.NewWriter(writer)
	tw := tar.NewWriter(gw)
	archive := &Archive{wr: writer, gw: gw, tw: tw}
	if file != nil {
		archive.Path = file.Name()
	}

	return archive, nil
}

func (a *Archive) AddFile(blob []byte, fname string) error {
	header := &tar.Header{
		Name:     fname,
		Mode:     0644,
		Size:     int64(len(blob)),
		Typeflag: tar.TypeReg,
	}
	if err := a.tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := a.tw.Write(blob); err != nil {
		return err
	}
	return nil
}

func (a *Archive) VirtualReader() io.Reader {
	if buf, ok := a.wr.(*bytes.Buffer); ok {
		return buf
	}
	return nil
}

func (a *Archive) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	if err := a.gw.Close(); err != nil {
		return err
	}
	if file, ok := a.wr.(*os.File); ok {
		file.Close()
	}
	return nil
}

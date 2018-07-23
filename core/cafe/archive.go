package client

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
)

type Archive struct {
	Id   string `json:"id"`
	Path string `json:"path"`

	file *os.File     `json:"-"`
	gw   *gzip.Writer `json:"-"`
	tw   *tar.Writer  `json:"-"`
}

func NewArchive(id string, dir string) (*Archive, error) {
	path := fmt.Sprintf("%s.tar.gz", filepath.Join(dir, id))
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	gw := gzip.NewWriter(file)
	tw := tar.NewWriter(gw)
	return &Archive{Id: id, Path: path, file: file, gw: gw, tw: tw}, nil
}

func (a *Archive) AddFile(blob []byte, fname string) error {
	header := &tar.Header{
		Name: fname,
		Mode: 0644,
		Size: int64(len(blob)),
	}
	if err := a.tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := a.tw.Write(blob); err != nil {
		return err
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
	return a.file.Close()
}

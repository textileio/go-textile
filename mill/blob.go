package mill

import (
	"io/ioutil"
	"mime/multipart"
)

type Blob struct{}

func (m *Blob) ID() string {
	return "/blob"
}

func (m *Blob) AcceptMedia(media string) error {
	return nil
}

func (m *Blob) Mill(file multipart.File, name string) (*Result, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &Result{File: data}, nil
}

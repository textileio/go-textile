package net

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var nl = "\r\n"

type MultipartRequest struct {
	Created     time.Time `json:"created"`
	Boundary    string    `json:"boundary"`
	PayloadPath string    `json:"payload_path"`
}

func (m *MultipartRequest) Init(dir string, boundary string) {
	m.Boundary = boundary
	m.Created = time.Now()
	m.PayloadPath = filepath.Join(dir, boundary)
}

func (m *MultipartRequest) AddFile(b []byte, fname string) error {
	// create file if not yet exists
	f, err := os.OpenFile(m.PayloadPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// create boundary header
	header := fmt.Sprintf("--%s", m.Boundary)
	header += nl
	header += fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"", fname)
	header += nl
	header += "Content-Type: application/octet-stream"
	header += nl
	header += nl

	// write the header and data part
	if _, err = f.Write([]byte(header)); err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		return err
	}
	if _, err = f.Write([]byte(nl)); err != nil {
		return err
	}

	return nil
}

func (m *MultipartRequest) Finish() error {
	// open file, will error if not first inited
	f, err := os.OpenFile(m.PayloadPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// create boundary footer
	footer := fmt.Sprintf("--%s--", m.Boundary)
	footer += nl

	// write the footer
	if _, err = f.Write([]byte(footer)); err != nil {
		return err
	}

	return nil
}

package net

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var nl = "\r\n"

type MultipartRequest struct {
	Boundary    string
	PayloadPath string
}

func (m *MultipartRequest) Init(dir string, boundary string) {
	m.Boundary = boundary
	m.PayloadPath = filepath.Join(dir, boundary)
}

func (m *MultipartRequest) AddFile(b []byte, fname string) error {
	// create file if not yet exists
	file, err := os.OpenFile(m.PayloadPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// create boundary header
	header := fmt.Sprintf("--%s", m.Boundary)
	header += nl
	header += fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"", fname)
	header += nl
	header += "Content-Type: application/octet-stream"
	header += nl
	header += nl

	// write the header and data part
	if _, err = file.Write([]byte(header)); err != nil {
		return err
	}
	if _, err = file.Write(b); err != nil {
		return err
	}
	if _, err = file.Write([]byte(nl)); err != nil {
		return err
	}

	return nil
}

func (m *MultipartRequest) Finish() error {
	// open file, will error if not first inited
	file, err := os.OpenFile(m.PayloadPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// create boundary footer
	footer := fmt.Sprintf("--%s--", m.Boundary)
	footer += nl

	// write the footer
	if _, err = file.Write([]byte(footer)); err != nil {
		return err
	}

	return nil
}

func (m *MultipartRequest) Send(url string) error {
	// open file, will error if not first inited
	file, err := os.Open(m.PayloadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// include boundary in header
	req, err := http.NewRequest("POST", url, file)
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", m.Boundary))

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

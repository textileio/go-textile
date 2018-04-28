package models

import (
	"io"
	"io/ioutil"
	"encoding/json"
)

type Response struct {
	Status     int `json:"status"`
	ResourceID string `json:"resource_id"`
	Token      string `json:"token"`
	Error      string `json:"error"`
}

func (r *Response) Read(body io.ReadCloser) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, r)
}

package mill

import (
	"encoding/json"
)

type Json struct{}

func (m *Json) ID() string {
	return "/json"
}

func (m *Json) Encrypt() bool {
	return true
}

func (m *Json) Pin() bool {
	return false
}

func (m *Json) AcceptMedia(media string) error {
	return accepts([]string{
		"application/json",
	}, media)
}

func (m *Json) Options() (string, error) {
	return "", nil
}

func (m *Json) Mill(input []byte, name string) (*Result, error) {
	var any map[string]interface{}
	if err := json.Unmarshal(input, &any); err != nil {
		return nil, err
	}

	if len(any) == 0 {
		return nil, ErrEmptyJsonFile
	}

	data, err := json.Marshal(&any)
	if err != nil {
		return nil, err
	}

	return &Result{File: data}, nil
}

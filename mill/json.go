package mill

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type JsonOpts struct {
	Schema string `json:"schema"`
}

type Json struct {
	Opts JsonOpts
}

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
	return hashOpts(m.Opts)
}

func (m *Json) Mill(input []byte, name string) (*Result, error) {
	sch := gojsonschema.NewStringLoader(m.Opts.Schema)
	doc := gojsonschema.NewStringLoader(string(input))

	result, err := gojsonschema.Validate(sch, doc)
	if err != nil {
		return nil, err
	}

	if !result.Valid() {
		var errs string
		for _, err := range result.Errors() {
			errs += fmt.Sprintf("- %s\n", err)
		}
		return nil, errors.New(errs)
	}

	return &Result{File: input}, nil
}

package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alecthomas/jsonschema"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"io"
)

// ErrSchemaValidationFailed indicates dag schema validation failed
var ErrSchemaValidationFailed = errors.New("schema validation failed")

// Node describes a DAG node
type Node struct {
	Pin    bool                   `json:"pin"`
	Use    string                 `json:"use"`
	Mill   string                 `json:"mill"`
	Opts   map[string]interface{} `json:"opts,omitempty"`
	Schema *jsonschema.Schema     `json:"schema,omitempty"`
	Nodes  map[string]*Node       `json:"nodes,omitempty"`
}

// LinkByName find a link w/ the given name in the provided list
func LinkByName(links []*ipld.Link, name string) *ipld.Link {
	for _, l := range links {
		if l.Name == name {
			return l
		}
	}
	return nil
}

// MustReflectAndMarshal panic if the json schema reflection fails
func MustReflectAndMarshal(any interface{}) io.Reader {
	data, err := json.Marshal(jsonschema.Reflect(any))
	if err != nil {
		panic("invalid file schema")
	}
	return bytes.NewReader(data)
}

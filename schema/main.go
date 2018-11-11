package schema

import (
	"errors"
	"github.com/alecthomas/jsonschema"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
)

// ErrSchemaValidationFailed indicates dag schema validation failed
var ErrSchemaValidationFailed = errors.New("schema validation failed")

// Node describes a DAG node
type Node struct {
	Pin    bool                   `json:"pin"`
	Use    string                 `json:"use,omitempty"`
	Mill   string                 `json:"mill,omitempty"`
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

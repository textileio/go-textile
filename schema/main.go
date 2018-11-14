package schema

import (
	"errors"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"

	"github.com/alecthomas/jsonschema"
)

// ErrSchemaValidationFailed indicates dag schema validation failed
var ErrSchemaValidationFailed = errors.New("schema validation failed")

// Node describes a DAG node
type Node struct {
	Pin    bool               `json:"pin"`
	Use    string             `json:"use,omitempty"`
	Mill   string             `json:"mill,omitempty"`
	Opts   map[string]string  `json:"opts,omitempty"`
	Schema *jsonschema.Schema `json:"schema,omitempty"`
	Nodes  map[string]*Node   `json:"nodes,omitempty"`
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

package schema

import (
	"errors"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"

	"github.com/alecthomas/jsonschema"
)

// ErrSchemaValidationFailed indicates dag schema validation failed
var ErrSchemaValidationFailed = errors.New("schema validation failed")

// ErrEmptySchema indicates a schema is empty
var ErrEmptySchema = errors.New("schema is empty")

// ErrLinkOrderNotSolvable
var ErrLinkOrderNotSolvable = errors.New("link order is not solvable")

// FileTag indicates the link should "use" the input file as source
const FileTag = ":file"

// Node describes a DAG node
type Node struct {
	Pin    bool               `json:"pin"`
	Mill   string             `json:"mill,omitempty"`
	Opts   map[string]string  `json:"opts,omitempty"`
	Schema *jsonschema.Schema `json:"schema,omitempty"`
	Links  map[string]*Link   `json:"links,omitempty"`
}

// Link is a sub-node which can "use" input from other sub-nodes
type Link struct {
	Use    string             `json:"use,omitempty"`
	Pin    bool               `json:"pin"`
	Mill   string             `json:"mill,omitempty"`
	Opts   map[string]string  `json:"opts,omitempty"`
	Schema *jsonschema.Schema `json:"schema,omitempty"`
}

// Step is an ordered name-link pair
type Step struct {
	Name string
	Link *Link
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

// Steps returns link steps in the order they should be processed
func Steps(links map[string]*Link) ([]Step, error) {
	var steps []Step
	run := links
	i := 0
	for {
		if i > len(links) {
			return nil, ErrLinkOrderNotSolvable
		}
		next := orderLinks(run, &steps)
		if len(next) == 0 {
			break
		}
		run = next
		i++
	}
	return steps, nil
}

// orderLinks attempts to place all links in steps, returning any unused
// whose source is not yet in steps
func orderLinks(links map[string]*Link, steps *[]Step) map[string]*Link {
	unused := make(map[string]*Link)
	for name, link := range links {
		if link.Use == FileTag {
			*steps = append([]Step{{Name: name, Link: link}}, *steps...)
		} else {
			useAt := -1
			for i, s := range *steps {
				if link.Use == s.Name {
					useAt = i
					break
				}
			}
			if useAt >= 0 {
				*steps = append(*steps, Step{Name: name, Link: link})
			} else {
				unused[name] = link
			}
		}
	}
	return unused
}

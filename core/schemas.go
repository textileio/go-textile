package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alecthomas/jsonschema"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"io"
	"strings"
	"time"
)

// ErrFileSchemaNotFound indicates schema is not found in the loaded list
var ErrFileSchemaNotFound = errors.New("file schema not found")

// ErrFileSchemaValidationFailed indicates file schema validation failed
var ErrFileSchemaValidationFailed = errors.New("file schema validation failed")

// ErrDAGSchemaNotFound indicates schema is not found in the loaded list
var ErrDAGSchemaNotFound = errors.New("DAG schema not found")

// ErrDAGSchemaValidationFailed indicates dag schema validation failed
var ErrDAGSchemaValidationFailed = errors.New("DAG schema validation failed")

// ErrFileNodeMissingSchema indicates an ipld file node is missing a schema link
var ErrFileNodeMissingSchema = errors.New("file node missing schema")

// ErrFileNodeMissingData indicates an ipld file node is missing a data link
var ErrFileNodeMissingData = errors.New("file node missing data")

type Schemas map[string]mh.Multihash

var fileSchemas = []FileSchema{&Blob{}, &Image{}, &ImageExif{}}

func (t *Textile) FileSchema(id string) mh.Multihash {
	return t.fileSchemas[strings.ToTitle(id)]
}

func (t *Textile) loadFileSchemas(schemas []FileSchema) error {
	// basic files
	for _, schema := range schemas {
		cid, err := ipfs.AddData(t.node, schema.MustReflect(), true)
		if err != nil {
			return err
		}
		t.fileSchemas[schema.ID()] = cid.Hash()
	}
	return nil
}

type FileSchema interface {
	ID() string
	MustReflect() io.Reader
}

type Blob struct{}

func (s *Blob) ID() string {
	return "Blob"
}

func (s *Blob) MustReflect() io.Reader {
	return mustReflectAndMarshal(s)
}

type Image struct {
	Width   int `json:"width"`
	Quality int `json:"quality"`
}

func (s *Image) ID() string {
	return "Image"
}

func (s *Image) MustReflect() io.Reader {
	return mustReflectAndMarshal(s)
}

type ImageExif struct {
	Created        time.Time `json:"created,omitempty"`
	Added          time.Time `json:"added"`
	Name           string    `json:"name"`
	Ext            string    `json:"extension"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	OriginalFormat string    `json:"original_format"`
	EncodingFormat string    `json:"encoding_format"`
	Latitude       float64   `json:"latitude,omitempty"`
	Longitude      float64   `json:"longitude,omitempty"`
}

func (s *ImageExif) ID() string {
	return "ImageExif"
}

func (s *ImageExif) MustReflect() io.Reader {
	return mustReflectAndMarshal(s)
}

func mustReflectAndMarshal(any interface{}) io.Reader {
	data, err := json.Marshal(jsonschema.Reflect(any))
	if err != nil {
		panic("invalid file schema")
	}
	return bytes.NewReader(data)
}

// DAG Schemas

var dagSchemas = []*DAGSchema{
	{
		ID: "Account",
	},
	{
		ID:  "Photos",
		Pin: true,
		Nodes: map[DAGLink]*DAGSchema{
			"raw": {
				FileSchemaID: "Blob",
			},
			"exif": {
				FileSchemaID: "ImageExif",
			},
			"large": {
				FileSchemaID: "Image",
				Props: DAGProps{
					"width":   1600,
					"quality": 75,
				},
			},
			"medium": {
				FileSchemaID: "Image",
				Props: DAGProps{
					"width":   800,
					"quality": 75,
				},
			},
			"small": {
				FileSchemaID: "Image",
				Props: DAGProps{
					"width":   320,
					"quality": 75,
				},
			},
			"thumb": {
				FileSchemaID: "Image",
				Props: DAGProps{
					"width":   100,
					"quality": 75,
				},
				Pin: true,
			},
		},
	},
}

func (t *Textile) DAGSchema(id string) mh.Multihash {
	return t.dagSchemas[strings.ToTitle(id)]
}

func (t *Textile) loadDAGSchemas(schemas []*DAGSchema) error {
	// basic dags
	for _, schema := range schemas {
		schema.Load(t.fileSchemas)
		cid, err := ipfs.AddData(t.node, schema.MustMarshal(), true)
		if err != nil {
			return err
		}
		t.dagSchemas[schema.ID] = cid.Hash()
	}
	return nil
}

type DAGLink string

type DAGProps map[string]interface{}

type DAGSchema struct {
	ID           string                 `json:"$id,omitempty"`
	FileSchemaID string                 `json:"file_schema_id,omitempty"`
	FileSchema   mh.Multihash           `json:"file_schema,omitempty"`
	Props        DAGProps               `json:"props,omitempty"`
	Pin          bool                   `json:"pin"`
	Nodes        map[DAGLink]*DAGSchema `json:"nodes,omitempty"`
}

func (d *DAGSchema) Load(files Schemas) {
	loadDAGSchema(d, files)
}

func (d *DAGSchema) MustMarshal() io.Reader {
	data, err := json.Marshal(d)
	if err != nil {
		panic("invalid DAG schema")
	}
	return bytes.NewReader(data)
}

func loadDAGSchema(node *DAGSchema, files Schemas) {
	node.FileSchema = files[node.FileSchemaID]
	for _, n := range node.Nodes {
		loadDAGSchema(n, files)
	}
}

func linkByName(links []*ipld.Link, name string) *ipld.Link {
	for _, l := range links {
		if l.Name == name {
			return l
		}
	}
	return nil
}

func (t *Thread) Process(schema *DAGSchema, node ipld.Node) error {
	// file node?
	if schema.FileSchema != nil {
		// get schema link
		fs := linkByName(node.Links(), "schema")
		if fs == nil {
			return ErrFileNodeMissingSchema
		}
		fsHash := fs.Cid.Hash().B58String()
		if fsHash != schema.FileSchema.B58String() {
			return ErrFileSchemaValidationFailed
		}
		// get data link
		data := linkByName(node.Links(), "data")
		if data == nil {
			return ErrFileNodeMissingData
		}
		dataHash := data.Cid.Hash().B58String()

		// handle file pin
		if schema.Pin {
			if err := ipfs.PinNode(t.node(), node); err != nil {
				return err
			}
			if err := ipfs.PinPath(t.node(), fsHash, false); err != nil {
				return err
			}
			if err := ipfs.PinPath(t.node(), dataHash, false); err != nil {
				return err
			}
		}

		// remote pin the schema hash
		t.cafeOutbox.Add(fsHash, repo.CafeStoreRequest)

		// if not mobile, remote pin the actual file data
		t.cafeOutbox.Add(dataHash, repo.CafeStoreRequest)

	} else {
		// dir node here
		for name, ds := range schema.Nodes {
			// ensure link is present
			link := linkByName(node.Links(), string(name))
			if link == nil {
				return ErrDAGSchemaValidationFailed
			}
			// get next node
			nd, err := ipfs.LinkNode(t.node(), link)
			if err != nil {
				return err
			}
			// keep going
			if err := t.Process(ds, nd); err != nil {
				return err
			}
		}

		// handle dir pin
		if schema.Pin {
			if err := ipfs.PinNode(t.node(), node); err != nil {
				return err
			}
		}
	}

	// remote pin all nodes
	t.cafeOutbox.Add(node.Cid().Hash().B58String(), repo.CafeStoreRequest)
	go t.cafeOutbox.Flush()

	return nil
}

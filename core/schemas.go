package core

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/jsonschema"
	"github.com/textileio/textile-go/ipfs"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"io"
	"time"
)

type Schemas map[string]mh.Multihash

var fileSchemas = []FileSchema{&Blob{}, &Image{}, &ImageExif{}}

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
	// only one built-in so far
	{
		ID:  "TextilePhoto",
		Pin: true,
		Nodes: map[DAGLink]*DAGSchema{
			"raw": {
				SchemaID: "Blob",
			},
			"exif": {
				SchemaID: "ImageExif",
			},
			"large": {
				SchemaID: "Image",
				Props: DAGProps{
					"width":   1600,
					"quality": 75,
				},
			},
			"medium": {
				SchemaID: "Image",
				Props: DAGProps{
					"width":   800,
					"quality": 75,
				},
			},
			"small": {
				SchemaID: "Image",
				Props: DAGProps{
					"width":   320,
					"quality": 75,
				},
			},
			"thumb": {
				SchemaID: "Image",
				Props: DAGProps{
					"width":   100,
					"quality": 75,
				},
				Pin: true,
			},
		},
	},
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
	ID       string                 `json:"$id,omitempty"`
	SchemaID string                 `json:"schema_id,omitempty"`
	Schema   mh.Multihash           `json:"schema,omitempty"`
	Props    DAGProps               `json:"props,omitempty"`
	Pin      bool                   `json:"pin"`
	Nodes    map[DAGLink]*DAGSchema `json:"nodes,omitempty"`
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
	node.Schema = files[node.SchemaID]
	for _, n := range node.Nodes {
		loadDAGSchema(n, files)
	}
}

package mill

import (
	"encoding/json"

	"github.com/textileio/textile-go/schema"
	"github.com/xeipuuv/gojsonschema"
)

type Schema struct{}

func (m *Schema) ID() string {
	return "/schema"
}

func (m *Schema) Encrypt() bool {
	return false
}

func (m *Schema) Pin() bool {
	return true
}

func (m *Schema) AcceptMedia(media string) error {
	return accepts([]string{"application/json"}, media)
}

func (m *Schema) Options() (string, error) {
	return "", nil
}

func (m *Schema) Mill(input []byte, name string) (*Result, error) {
	var node schema.Node
	if err := json.Unmarshal(input, &node); err != nil {
		return nil, err
	}

	if node.Mill == "" {
		if len(node.Links) == 0 {
			return nil, schema.ErrEmptySchema
		}

		for _, link := range node.Links {
			if !schema.ValidateMill(link.Mill) {
				return nil, schema.ErrSchemaInvalidMill
			}

			// extra check for json
			if link.Mill == "/json" {
				if link.JsonSchema == nil {
					return nil, schema.ErrMissingJsonSchema
				}
				if err := validateJsonSchema(link.JsonSchema); err != nil {
					return nil, err
				}
			}
		}

		// ensure link steps are solvable
		if _, err := schema.Steps(node.Links); err != nil {
			return nil, err
		}

	} else {
		if !schema.ValidateMill(node.Mill) {
			return nil, schema.ErrSchemaInvalidMill
		}

		// extra check for json
		if node.Mill == "/json" {
			if node.JsonSchema == nil {
				return nil, schema.ErrMissingJsonSchema
			}
			if err := validateJsonSchema(node.JsonSchema); err != nil {
				return nil, err
			}
		}
	}

	data, err := json.Marshal(&node)
	if err != nil {
		return nil, err
	}

	return &Result{File: data}, nil
}

func validateJsonSchema(jschema map[string]interface{}) error {
	data, err := json.Marshal(&jschema)
	if err != nil {
		return err
	}

	loader := gojsonschema.NewStringLoader(string(data))

	if _, err := gojsonschema.NewSchema(loader); err != nil {
		return schema.ErrBadJsonSchema
	}

	return nil
}

package mobile

import (
	"encoding/json"

	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
)

func (m *Mobile) AddSchema(schemaJSON string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}
	added, err := m.addSchema(schemaJSON)
	if err != nil {
		return "", err
	}

	return toJSON(added)
}

func (m *Mobile) addSchema(schemaJSON string) (*repo.File, error) {
	var node schema.Node
	if err := json.Unmarshal([]byte(schemaJSON), &node); err != nil {
		return nil, err
	}
	data, err := json.Marshal(&node)
	if err != nil {
		return nil, err
	}

	conf := core.AddFileConfig{
		Input: data,
		Media: "application/json",
	}

	return m.node.AddFile(&mill.Schema{}, conf)
}

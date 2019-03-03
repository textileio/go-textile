package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/mill"
)

// AddSchema adds a new schema via schema mill
func (m *Mobile) AddSchema(jsonstr string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	added, err := m.node.AddFileIndex(&mill.Schema{}, core.AddFileConfig{
		Input: []byte(jsonstr),
		Media: "application/json",
	})
	if err != nil {
		return nil, err
	}

	return proto.Marshal(added)
}

package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
)

// Summary calls core Summary
func (m *Mobile) Summary() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.Summary())
}

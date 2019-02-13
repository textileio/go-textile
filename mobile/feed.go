package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
)

// Feed calls core Feed
func (m *Mobile) Feed(offset string, limit int, threadId string, annotated bool) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	items, err := m.node.Feed(offset, limit, threadId, annotated)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(items)
}

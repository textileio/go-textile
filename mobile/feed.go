package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
)

// Feed calls core Feed
func (m *Mobile) Feed(offset string, limit int, threadId string, mode int32) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	items, err := m.node.Feed(offset, limit, threadId, pb.FeedMode(mode))
	if err != nil {
		return nil, err
	}

	return proto.Marshal(items)
}

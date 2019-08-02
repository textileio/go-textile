package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
)

// Feed calls core Feed
func (m *Mobile) Feed(req []byte) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	mreq := new(pb.FeedRequest)
	if err := proto.Unmarshal(req, mreq); err != nil {
		return nil, err
	}

	items, err := m.node.Feed(mreq)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(items)
}

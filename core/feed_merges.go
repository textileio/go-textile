package core

import (
	"github.com/textileio/textile-go/pb"
)

func (t *Textile) merge(block *pb.Block, opts feedItemOpts) (*pb.Merge, error) {
	if block.Type != pb.Block_MERGE {
		return nil, ErrBlockWrongType
	}

	var targets []*pb.FeedItem
	for _, p := range block.Parents {
		parent, err := t.feedItem(t.datastore.Blocks().Get(p), feedItemOpts{})
		if err != nil {
			return nil, err
		}
		targets = append(targets, parent)
	}

	return &pb.Merge{
		Block:   block.Id,
		Date:    block.Date,
		User:    t.User(block.Author),
		Targets: targets,
	}, nil
}

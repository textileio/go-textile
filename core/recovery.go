package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
)

// SearchBackups searches the network for backups
func (t *Textile) SearchBackups(query *pb.BackupQuery, options *pb.QueryOptions) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, nil, nil, err
	}

	resCh, errCh, cancel := t.search(&pb.Query{
		Type:    pb.QueryType_BACKUPS,
		Options: options,
		Payload: &any.Any{
			TypeUrl: "/BackupQuery",
			Value:   payload,
		},
	})
	return resCh, errCh, cancel, nil
}

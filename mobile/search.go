package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
)

// handleSearch handles the response channels from a search
func handleSearch(resultCh <-chan *pb.QueryResult, errCh <-chan error, cancel *broadcast.Broadcaster, cb Callback) (func(), error) {
	doneFn := func() {
		cb.Call(proto.Marshal(&pb.QueryEvent{
			Type: pb.QueryEvent_DONE,
		}))
	}
	cancelFn := func() {
		cancel.Close()
		doneFn()
	}

	go func() {
		select {
		case err := <-errCh:
			cb.Call(nil, err)
			return

		case res, ok := <-resultCh:
			if !ok {
				doneFn()
				return
			}
			cb.Call(proto.Marshal(&pb.QueryEvent{
				Type: pb.QueryEvent_DATA,
				Data: res,
			}))
		}
	}()

	return cancelFn, nil
}

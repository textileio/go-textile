package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
)

type CancelFn struct {
	cancel *broadcast.Broadcaster
	done   func()
}

func (c *CancelFn) Call() {
	c.cancel.Close()
	c.done()
}

// handleSearchStream handles the response channels from a search
func handleSearchStream(resultCh <-chan *pb.QueryResult, errCh <-chan error, cancel *broadcast.Broadcaster, cb Callback) (*CancelFn, error) {
	doneFn := func() {
		cb.Call(proto.Marshal(&pb.QueryEvent{
			Type: pb.QueryEvent_DONE,
		}))
	}
	cancelFn := &CancelFn{cancel: cancel, done: doneFn}

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

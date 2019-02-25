package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
)

// CancelFn is used to cancel an async request
type CancelFn struct {
	cancel *broadcast.Broadcaster
	done   func()
}

// Call is used to invoke the cancel
func (c *CancelFn) Call() {
	c.cancel.Close()
	c.done()
}

// handleSearchStream handles the response channels from a search
func handleSearchStream(resultCh <-chan *pb.QueryResult, errCh <-chan error, cancel *broadcast.Broadcaster, cb Callback) (*CancelFn, error) {
	var done bool
	doneFn := func() {
		if done {
			return
		}
		done = true
		cb.Call(proto.Marshal(&pb.QueryEvent{
			Type: pb.QueryEvent_DONE,
		}))
	}
	cancelFn := &CancelFn{cancel: cancel, done: doneFn}

	go func() {
		for {
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
		}
	}()

	return cancelFn, nil
}

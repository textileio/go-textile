package mobile

import (
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/pb"
)

// SearchHandle is used to cancel an async search request
type SearchHandle struct {
	Id     string
	cancel *broadcast.Broadcaster
	done   func()
}

// Cancel is used to cancel the request
func (h *SearchHandle) Cancel() {
	h.cancel.Close()
	h.done()
}

// handleSearchStream handles the response channels from a search
func (m *Mobile) handleSearchStream(resultCh <-chan *pb.QueryResult, errCh <-chan error, cancel *broadcast.Broadcaster) (*SearchHandle, error) {
	id := ksuid.New().String()

	var done bool
	doneFn := func() {
		if done {
			return
		}
		done = true
		m.notify(pb.MobileEventType_QUERY_RESPONSE, &pb.MobileQueryEvent{
			Id:   id,
			Type: pb.MobileQueryEvent_DONE,
		})
	}
	handle := &SearchHandle{
		Id:     id,
		cancel: cancel,
		done:   doneFn,
	}

	go func() {
		for {
			select {
			case err := <-errCh:
				m.notify(pb.MobileEventType_QUERY_RESPONSE, &pb.MobileQueryEvent{
					Id:   id,
					Type: pb.MobileQueryEvent_ERROR,
					Error: &pb.Error{
						Code:    500,
						Message: err.Error(),
					},
				})
				return

			case res, ok := <-resultCh:
				if !ok {
					doneFn()
					return
				}
				m.notify(pb.MobileEventType_QUERY_RESPONSE, &pb.MobileQueryEvent{
					Id:   id,
					Type: pb.MobileQueryEvent_DATA,
					Data: res,
				})
			}
		}
	}()

	return handle, nil
}

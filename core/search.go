package core

import (
	"sync"

	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
)

// queryResultSet holds a unique set of search results
type queryResultSet struct {
	filter pb.QueryOptions_FilterType
	items  map[string]*pb.QueryResult
	mux    sync.Mutex
}

// newQueryResultSet returns a new queryResultSet
func newQueryResultSet(filter pb.QueryOptions_FilterType) *queryResultSet {
	return &queryResultSet{
		filter: filter,
		items:  make(map[string]*pb.QueryResult, 0),
	}
}

// Add only adds a result to the set if it's newer than last
func (s *queryResultSet) Add(items ...*pb.QueryResult) []*pb.QueryResult {
	s.mux.Lock()
	defer s.mux.Unlock()

	var added []*pb.QueryResult
	for _, i := range items {
		last := s.items[i.Id]
		switch s.filter {
		case pb.QueryOptions_NO_FILTER:
			break
		case pb.QueryOptions_HIDE_OLDER:
			if last != nil && protoTimeToNano(i.Date) < protoTimeToNano(last.Date) {
				continue
			}
		}
		s.items[i.Id] = i
		added = append(added, i)
	}

	return added
}

// List returns the items as a slice
func (s *queryResultSet) List() []*pb.QueryResult {
	s.mux.Lock()
	defer s.mux.Unlock()

	var list []*pb.QueryResult
	for _, i := range s.items {
		list = append(list, i)
	}

	return list
}

// Search searches the network based on the given query
func (t *Textile) search(query *pb.Query) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster) {
	query = applyQueryDefaults(query)

	var searchChs []chan *pb.QueryResult

	// local results channel
	localCh := make(chan *pb.QueryResult)
	searchChs = append(searchChs, localCh)

	// remote results channel(s)
	var cafeChs []chan *pb.QueryResult
	clientCh := make(chan *pb.QueryResult)
	sessions := t.datastore.CafeSessions().List()
	if len(sessions) > 0 {
		for range sessions {
			cafeCh := make(chan *pb.QueryResult)
			cafeChs = append(cafeChs, cafeCh)
			searchChs = append(searchChs, cafeCh)
		}
	} else {
		searchChs = append(searchChs, clientCh)
	}

	resultCh := mergeQueryResults(searchChs)
	errCh := make(chan error)
	cancel := broadcast.NewBroadcaster(0)

	go func() {
		defer func() {
			for _, ch := range searchChs {
				close(ch)
			}
		}()

		// search local
		results, err := t.cafe.searchLocal(query.Type, query.Options.Filter, query.Payload, true)
		if err != nil {
			errCh <- err
			return
		}
		for _, res := range results.items {
			localCh <- res
		}
		if query.Options.Local || len(results.items) >= int(query.Options.Limit) {
			return
		}

		// search the network
		if len(sessions) == 0 {

			// search via pubsub directly
			canceler := cancel.Listen()
			if err := t.cafe.searchPubSub(query, func(res *pb.QueryResults) bool {
				for _, n := range results.Add(res.Items...) {
					clientCh <- n
				}
				return len(results.items) >= int(query.Options.Limit)
			}, canceler.Ch, false); err != nil {
				errCh <- err
				return
			}

		} else {

			// search via cafes
			wg := sync.WaitGroup{}
			for i, session := range sessions {
				cafe, err := peer.IDB58Decode(session.Id)
				if err != nil {
					errCh <- err
					return
				}
				canceler := cancel.Listen()

				wg.Add(1)
				go func(i int, cafe peer.ID, canceler *broadcast.Listener) {
					defer wg.Done()

					// token must be attached per cafe session, use a new query
					q := &pb.Query{}
					*q = *query
					if err := t.cafe.Search(q, cafe, func(res *pb.QueryResult) {
						for _, n := range results.Add(res) {
							cafeChs[i] <- n
						}
					}, canceler.Ch); err != nil {
						errCh <- err
						return
					}
				}(i, cafe, canceler)
			}

			wg.Wait()
		}
	}()

	return resultCh, errCh, cancel
}

// mergeQueryResults merges results from mulitple queries
func mergeQueryResults(cs []chan *pb.QueryResult) chan *pb.QueryResult {
	out := make(chan *pb.QueryResult)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c chan *pb.QueryResult) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

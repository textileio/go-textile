package cafe

import "github.com/textileio/go-textile/pb"

/*
list: return a batch of requests by date until "next" is empty
complete: mark a single request as complete (does not delete until all in the group are deleted
so that we can know the total data size being handled)
stat a group: returns some info:
  - num total reqs
  - num complete
  - percent complete
  - will return nil if all done since they will have been deleted

Payloads for different requests (expose something like get payload?):
- store: PUT /store/:cid, body => raw object data
- unstore: DELETE /store/:cid, body => none
- store thread: PUT /thread/:id, body => encrypted thread object (snapshot)
- unstore thread: DELETE /thread/:id, body => none
- deliver message: POST /inbox/:pid, body => encrypted message

Other, related methods...
- stat a thread: are there any pending cafe requests (per cafe) / block messages for this thread?
- stat a files block: what percentage of total bytes for the group are complete, per cafe?
- when you add a new cafe, a bunch of new request should be created, which results in a not-as-synced state

Thoughts...

There are two types of things that need to be "stored" and count to users' data usage:
1. thread blocks (requests grouped by thread id)
2. files dags (requests grouped by target id)
(thread snapshots do not count toward data usage, nor do inboxed messages, ?)

*/

type RequestHandler interface {
	Flush()
	Store(cids *pb.StringList, cafeId string) (*pb.StringList, error)
	Unstore(cids *pb.StringList, cafeId string) (*pb.StringList, error)
	StoreThread(thrd *pb.Thread, cafeId string) error
	UnstoreThread(id string, cafeId string) error
	DeliverMessage(msgId string, peerId string, cafe *pb.Cafe) error
}

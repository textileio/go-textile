package mobile

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/golang/protobuf/proto"
	icid "github.com/ipfs/go-cid"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

// RegisterCafe is the async flavor of registerCafe
func (m *Mobile) RegisterCafe(id string, token string, cb Callback) {
	m.node.Lock()
	go func() {
		defer m.node.Unlock()

		cb.Call(m.registerCafe(id, token))
	}()
}

// registerCafe calls core RegisterCafe
func (m *Mobile) registerCafe(id string, token string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	_, err := m.node.RegisterCafe(id, token)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// DeegisterCafe is the async flavor of deregisterCafe
func (m *Mobile) DeregisterCafe(id string, cb Callback) {
	m.node.Lock()
	go func() {
		defer m.node.Unlock()

		cb.Call(m.deregisterCafe(id))
	}()
}

// deregisterCafe calls core DeregisterCafe
func (m *Mobile) deregisterCafe(id string) error {
	if !m.node.Online() {
		return core.ErrOffline
	}

	err := m.node.DeregisterCafe(id)
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// RefreshCafeSession is the async flavor of refreshCafeSession
func (m *Mobile) RefreshCafeSession(id string, cb ProtoCallback) {
	m.node.Lock()
	go func() {
		defer m.node.Unlock()

		cb.Call(m.refreshCafeSession(id))
	}()
}

// refreshCafeSession calls core RefreshCafeSession
func (m *Mobile) refreshCafeSession(id string) ([]byte, error) {
	if !m.node.Online() {
		return nil, core.ErrOffline
	}

	session, err := m.node.RefreshCafeSession(id)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(session)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// CheckCafeMessages calls core CheckCafeMessages
func (m *Mobile) CheckCafeMessages() error {
	m.node.Lock()
	go func() {
		defer m.node.Unlock()

		if !m.node.Online() {
			log.Warning("check messages called offline")
			return
		}

		err := m.node.CheckCafeMessages()
		if err != nil {
			log.Errorf("error checking cafe inbox: %s", err)
		}

		m.node.FlushCafes()
	}()

	return nil
}

// CafeSession calls core CafeSession
func (m *Mobile) CafeSession(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	session, err := m.node.CafeSession(id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	bytes, err := proto.Marshal(session)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// CafeSessions calls core CafeSessions
func (m *Mobile) CafeSessions() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	bytes, err := proto.Marshal(m.node.CafeSessions())
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// CafeRequests paginates new requests
func (m *Mobile) CafeRequests(limit int) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	groups := m.node.Datastore().CafeRequests().ListGroups("", limit)
	return proto.Marshal(&pb.Strings{Values: groups})
}

// CafeRequestPending marks a request as pending
func (m *Mobile) CafeRequestPending(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.Datastore().CafeRequests().UpdateGroupStatus(id, pb.CafeRequest_PENDING)
	if err != nil {
		return err
	}

	log.Debugf("request %s pending", id)

	m.notify(pb.MobileEventType_CAFE_SYNC_GROUP_UPDATE, m.cafeSyncGroupStatus(id))
	return nil
}

// CafeRequestNotPending marks a request as not pending (new)
func (m *Mobile) CafeRequestNotPending(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	log.Debugf("request %s not pending", id)

	return m.node.Datastore().CafeRequests().UpdateGroupStatus(id, pb.CafeRequest_NEW)
}

// CompleteCafeRequest marks a request as complete
func (m *Mobile) CompleteCafeRequest(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.Datastore().CafeRequests().UpdateGroupStatus(id, pb.CafeRequest_COMPLETE)
	if err != nil {
		return err
	}

	log.Debugf("request %s completed", id)

	return m.handleCafeRequestDone(id, m.cafeSyncGroupStatus(id))
}

// FailCafeRequest deletes a request
func (m *Mobile) FailCafeRequest(id string, reason string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	status := m.cafeSyncGroupStatus(id)
	status.Error = reason
	status.ErrorId = id

	log.Warningf("request %s failed: %s", id, reason)

	return m.handleCafeRequestDone(id, status)
}

// UpdateCafeRequestProgress updates the request with progress info
func (m *Mobile) UpdateCafeRequestProgress(id string, transferred int64, total int64) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.Datastore().CafeRequests().UpdateGroupProgress(id, transferred, total)
	if err != nil {
		return err
	}

	log.Debugf("request progress: %d / %d transferred", transferred, total)

	m.notify(pb.MobileEventType_CAFE_SYNC_GROUP_UPDATE, m.cafeSyncGroupStatus(id))
	return nil
}

// WriteCafeRequest is the async version of writeCafeRequest
func (m *Mobile) WriteCafeRequest(group string, cb ProtoCallback) {
	m.node.Lock()
	go func() {
		defer m.node.Unlock()

		cb.Call(m.writeCafeRequest(group))
	}()
}

// writeCafeRequest returns an HTTP request object for the given group, writing payload to disk
// - store: PUT /store, body => multipart, one file per req
// - unstore: DELETE /store/:cid, body => noop
// - store thread: PUT /threads/:id, body => encrypted thread object (snapshot)
// - unstore thread: DELETE /threads/:id, body => noop
// - deliver message: POST /inbox/:pid, body => encrypted message
func (m *Mobile) writeCafeRequest(group string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	// ensure tmp exists
	err := util.Mkdirp(m.RepoPath + "/tmp")
	if err != nil {
		return nil, err
	}

	reqs := m.node.Datastore().CafeRequests().GetGroup(group)
	if len(reqs.Items) == 0 {
		return nil, fmt.Errorf("request group not found")
	}
	var hreq *pb.CafeHTTPRequest

	fail := func(reason string) ([]byte, error) {
		err = m.FailCafeRequest(group, reason)
		if err != nil {
			return nil, err
		}
		return proto.Marshal(hreq)
	}

	// group by cafe
	creqs := make(map[string][]*pb.CafeRequest)
	for _, req := range reqs.Items {
		creqs[req.Cafe.Peer] = append(creqs[req.Cafe.Peer], req)
	}
	for cafe, reqs := range creqs {

		// group by type
		treqs := make(map[pb.CafeRequest_Type][]*pb.CafeRequest)
		for _, req := range reqs {
			treqs[req.Type] = append(treqs[req.Type], req)
		}
		if len(treqs) > 1 {
			return fail("request group contains multiple types")
		}
		rtype := reqs[0].Type

		// store reqs can be handled by multipart form data
		if rtype == pb.CafeRequest_STORE {
			session := m.node.Datastore().CafeSessions().Get(cafe)
			if session == nil {
				return nil, fmt.Errorf("session for cafe %s not found", cafe)
			}

			hreq = &pb.CafeHTTPRequest{
				Type: pb.CafeHTTPRequest_PUT,
				Url:  session.Cafe.Url + "/api/" + core.CafeApiVersion + "/store",
				Headers: map[string]string{
					"Authorization": "Basic " + session.Access,
				},
				Path: fmt.Sprintf("%s/tmp/%s", m.RepoPath, group),
			}

			// write each req with a multipart writer
			file, err := os.Create(hreq.Path)
			if err != nil {
				return nil, err
			}
			writer := multipart.NewWriter(file)
			for _, req := range reqs {
				part, err := writer.CreateFormFile("file", req.Target)
				if err != nil {
					return nil, err
				}

				data, err := ipfs.DataAtPath(m.node.Ipfs(), req.Target)
				if err != nil {
					if err == iface.ErrIsDir {
						data, err := ipfs.ObjectAtPath(m.node.Ipfs(), req.Target)
						if err != nil {
							return nil, err
						}
						_, err = part.Write(data)
						if err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
				} else {
					_, err = part.Write(data)
					if err != nil {
						return nil, err
					}
				}
			}
			_ = writer.Close()
			_ = file.Close()

			hreq.Headers["Content-Type"] = writer.FormDataContentType()

		} else {
			if len(reqs) > 1 {
				return fail("type does not allow multiple requests per group")
			}
			unpin := make(map[string]struct{})
			for _, req := range reqs {
				hreq = &pb.CafeHTTPRequest{
					Url: req.Cafe.Url + "/api/" + core.CafeApiVersion,
					Headers: map[string]string{
						"Content-Type": "application/octet-stream",
					},
					Path: fmt.Sprintf("%s/tmp/%s", m.RepoPath, group),
				}

				var body []byte
				switch req.Type {
				case pb.CafeRequest_UNSTORE:
					hreq.Type = pb.CafeHTTPRequest_DELETE
					hreq.Url += "/store/" + req.Target
					body = []byte("noop")

				case pb.CafeRequest_STORE_THREAD:
					hreq.Type = pb.CafeHTTPRequest_PUT
					hreq.Url += "/threads/" + req.Target

					thrd := m.node.Datastore().Threads().Get(req.Target)
					if thrd == nil {
						return nil, core.ErrThreadNotFound
					}
					plaintext, err := proto.Marshal(thrd)
					if err != nil {
						return nil, err
					}
					ciphertext, err := m.Encrypt(plaintext)
					if err != nil {
						return nil, err
					}
					body = ciphertext

				case pb.CafeRequest_UNSTORE_THREAD:
					hreq.Type = pb.CafeHTTPRequest_DELETE
					hreq.Url += "/threads/" + req.Target
					body = []byte("noop")

				case pb.CafeRequest_INBOX:
					hreq.Type = pb.CafeHTTPRequest_POST
					hreq.Url += "/inbox/" + m.node.Ipfs().Identity.Pretty() + "/" + req.Peer
					body, err = ipfs.DataAtPath(m.node.Ipfs(), req.Target)
					if err != nil {
						return nil, err
					}
					unpin[req.Target] = struct{}{}
				}

				// include session token for non-inbox requests
				if req.Type != pb.CafeRequest_INBOX {
					session := m.node.Datastore().CafeSessions().Get(cafe)
					if session == nil {
						return nil, fmt.Errorf("session for cafe %s not found", cafe)
					}
					hreq.Headers["Authorization"] = "Basic " + session.Access
				}

				err = util.WriteFileByPath(hreq.Path, body)
				if err != nil {
					return nil, err
				}
			}

			// unpin tmp objects
			for p := range unpin {
				id, err := icid.Decode(p)
				if err != nil {
					return nil, err
				}
				err = ipfs.UnpinCid(m.node.Ipfs(), id, false)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return proto.Marshal(hreq)
}

// cafeSyncGroupStatus returns the status of the given request's sync group
func (m *Mobile) cafeSyncGroupStatus(id string) *pb.CafeSyncGroupStatus {
	return m.node.Datastore().CafeRequests().SyncGroupStatus(id)
}

// deleteCafeRequestBody removes the file associated with the request
func (m *Mobile) deleteCafeRequestBody(id string) error {
	return os.Remove(fmt.Sprintf("%s/tmp/%s", m.RepoPath, id))
}

// handleCafeRequestDone handles clean up after a request is complete/failed
func (m *Mobile) handleCafeRequestDone(id string, status *pb.CafeSyncGroupStatus) error {
	if status.ErrorId != "" {
		m.notify(pb.MobileEventType_CAFE_SYNC_GROUP_FAILED, status)

		// delete queued block
		// @todo: Uncomment this when sync can only be handled by a single cafe session
		//syncGroup := m.node.Datastore().CafeRequests().GetSyncGroup(id)
		//err := m.node.Datastore().Blocks().Delete(syncGroup)
		//if err != nil {
		//	return err
		//}
		err := m.node.Datastore().CafeRequests().DeleteByGroup(id)
		if err != nil {
			return err
		}
	} else if status.NumComplete == status.NumTotal {
		m.notify(pb.MobileEventType_CAFE_SYNC_GROUP_COMPLETE, status)
	} else {
		m.notify(pb.MobileEventType_CAFE_SYNC_GROUP_UPDATE, status)
	}

	// release pending blocks
	m.node.FlushBlocks()

	return m.deleteCafeRequestBody(id)
}

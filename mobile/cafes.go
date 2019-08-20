package mobile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	icid "github.com/ipfs/go-cid"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/util"
)

// RegisterCafe is the async flavor of registerCafe
func (m *Mobile) RegisterCafe(id string, token string, cb Callback) {
	m.node.WaitAdd(1, "Mobile.RegisterCafe")
	go func() {
		defer m.node.WaitDone("Mobile.RegisterCafe")

		cb.Call(m.registerCafe(id, token))
	}()
}

// registerCafe calls gets a new session from the given host
func (m *Mobile) registerCafe(host string, token string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	url := fmt.Sprintf("%s/api/%s/sessions/challenge/?account_addr=%s",
		host, core.CafeApiVersion, m.Address())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	err = errorCheck(res)
	if err != nil {
		return err
	}

	snonce := new(pb.CafeClientNonce)
	err = jsonpb.Unmarshal(res.Body, snonce)
	if err != nil {
		return err
	}

	// complete the challenge
	nonce := ksuid.New().String()
	sig, err := m.Sign([]byte(snonce.Value + nonce))
	if err != nil {
		return err
	}

	pid := m.node.Ipfs().Identity.Pretty()
	url = fmt.Sprintf("%s/api/%s/sessions/%s/?account_addr=%s&challenge=%s&nonce=%s",
		host, core.CafeApiVersion, pid, m.Address(), snonce.Value, nonce)
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(sig))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+token)
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	err = errorCheck(res)
	if err != nil {
		return err
	}

	session := new(pb.CafeSession)
	err = jsonpb.Unmarshal(res.Body, session)
	if err != nil {
		return err
	}

	// return existing session
	if x := m.node.Datastore().CafeSessions().Get(session.Id); x != nil {
		return nil
	}

	err = m.node.Datastore().CafeSessions().AddOrUpdate(session)
	if err != nil {
		return err
	}

	err = m.node.UpdatePeerInboxes()
	if err != nil {
		return err
	}

	// sync all blocks and files target
	err = m.node.CafeRequestThreadsContent(session.Id)
	if err != nil {
		return err
	}

	for _, thrd := range m.node.Threads() {
		_, err = thrd.Annouce(nil)
		if err != nil {
			return err
		}
	}

	err = m.node.PublishPeer()
	if err != nil {
		return err
	}

	err = m.node.SnapshotThreads()
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// RefreshCafeSession is the async flavor of refreshCafeSession
func (m *Mobile) RefreshCafeSession(id string, cb ProtoCallback) {
	m.node.WaitAdd(1, "Mobile.RefreshCafeSession")
	go func() {
		defer m.node.WaitDone("Mobile.RefreshCafeSession")

		cb.Call(m.refreshCafeSession(id))
	}()
}

// refreshCafeSession refreshes the session with the given id
func (m *Mobile) refreshCafeSession(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	session := m.node.Datastore().CafeSessions().Get(id)
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	pid := m.node.Ipfs().Identity.Pretty()
	url := fmt.Sprintf("%s/api/%s/sessions/%s/refresh", session.Cafe.Url, core.CafeApiVersion, pid)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(session.Access)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+session.Refresh)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = errorCheck(res)
	if err != nil {
		return nil, err
	}

	session = new(pb.CafeSession)
	err = jsonpb.Unmarshal(res.Body, session)
	if err != nil {
		return nil, err
	}

	err = m.node.Datastore().CafeSessions().AddOrUpdate(session)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(session)
}

// DeegisterCafe is the async flavor of deregisterCafe
func (m *Mobile) DeregisterCafe(id string, cb Callback) {
	m.node.WaitAdd(1, "Mobile.DeregisterCafe")
	go func() {
		defer m.node.WaitDone("Mobile.DeregisterCafe")

		cb.Call(m.deregisterCafe(id))
	}()
}

// deregisterCafe deletes the session with the given id
func (m *Mobile) deregisterCafe(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	session := m.node.Datastore().CafeSessions().Get(id)
	if session == nil {
		return nil
	}

	pid := m.node.Ipfs().Identity.Pretty()
	url := fmt.Sprintf("%s/api/%s/sessions/%s", session.Cafe.Url, core.CafeApiVersion, pid)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+session.Access)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	err = errorCheck(res)
	if err != nil {
		return err
	}

	// cleanup
	err = m.node.Datastore().CafeRequests().DeleteByCafe(session.Id)
	if err != nil {
		return err
	}
	err = m.node.Datastore().CafeSessions().Delete(session.Id)
	if err != nil {
		return err
	}

	err = m.node.UpdatePeerInboxes()
	if err != nil {
		return err
	}

	for _, thrd := range m.node.Threads() {
		_, err := thrd.Annouce(nil)
		if err != nil {
			return err
		}
	}

	err = m.node.PublishPeer()
	if err != nil {
		return err
	}

	m.node.FlushCafes()

	return nil
}

// CheckCafeMessages is the async flavor of checkCafeMessages
func (m *Mobile) CheckCafeMessages(cb Callback) {
	m.node.WaitAdd(1, "Mobile.CheckCafeMessages")
	go func() {
		defer m.node.WaitDone("Mobile.CheckCafeMessages")

		cb.Call(m.checkCafeMessages())
	}()
}

// checkCafeMessages queries all sessions for new messages
func (m *Mobile) checkCafeMessages() error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	// get active cafe sessions
	sessions := m.node.Datastore().CafeSessions().List().Items
	if len(sessions) == 0 {
		return nil
	}

	var err error
	for _, s := range sessions {
		err = m.checkSessionMessages(s)
		if err != nil {
			return err
		}
	}

	m.node.FlushCafes()

	return nil
}

// checkSessionMessages checks a session for new messages
func (m *Mobile) checkSessionMessages(session *pb.CafeSession) error {
	pid := m.node.Ipfs().Identity.Pretty()
	url := fmt.Sprintf("%s/api/%s/inbox/%s", session.Cafe.Url, core.CafeApiVersion, pid)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+session.Access)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	err = errorCheck(res)
	if err != nil {
		return err
	}

	msgs := new(pb.CafeMessages)
	err = jsonpb.Unmarshal(res.Body, msgs)
	if err != nil {
		return err
	}

	// save messages to inbox
	for _, msg := range msgs.Messages {
		err = m.node.Inbox().Add(msg)
		if err != nil {
			if !db.ConflictError(err) {
				return err
			}
		}
	}

	m.node.Inbox().Flush()

	// delete them from the remote so that more can be fetched
	if len(msgs.Messages) > 0 {
		return m.deleteCafeMessages(session.Id)
	}
	return nil
}

// deleteCafeMessages deletes a page of cafe messages
func (m *Mobile) deleteCafeMessages(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	session := m.node.Datastore().CafeSessions().Get(id)
	if session == nil {
		return nil
	}

	pid := m.node.Ipfs().Identity.Pretty()
	url := fmt.Sprintf("%s/api/%s/inbox/%s", session.Cafe.Url, core.CafeApiVersion, pid)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+session.Access)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	err = errorCheck(res)
	if err != nil {
		return err
	}

	ack := new(pb.CafeDeleteMessagesAck)
	err = jsonpb.Unmarshal(res.Body, ack)
	if err != nil {
		return err
	}

	if !ack.More {
		return nil
	}

	// apparently there are more new messages waiting...
	return m.checkSessionMessages(session)
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

	res, err := proto.Marshal(session)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CafeSessions calls core CafeSessions
func (m *Mobile) CafeSessions() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	res, err := proto.Marshal(m.node.CafeSessions())
	if err != nil {
		return nil, err
	}
	return res, nil
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
	go func() {
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

// errorCheck returns an error encoded in the response
func errorCheck(res *http.Response) error {
	if res.StatusCode >= http.StatusBadRequest {
		decoder := json.NewDecoder(res.Body)
		e := new(core.CafeError)
		err := decoder.Decode(e)
		if err != nil {
			return err
		}
		return fmt.Errorf(e.Error)
	}
	return nil
}

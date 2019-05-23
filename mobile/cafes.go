package mobile

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/golang/protobuf/proto"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

// RegisterCafe calls core RegisterCafe
func (m *Mobile) RegisterCafe(host string, token string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	if _, err := m.node.RegisterCafe(host, token); err != nil {
		return err
	}
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

// RefreshCafeSession calls core RefreshCafeSession
func (m *Mobile) RefreshCafeSession(id string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
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

// DeegisterCafe calls core DeregisterCafe
func (m *Mobile) DeregisterCafe(id string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	return m.node.DeregisterCafe(id)
}

// CheckCafeMessages calls core CheckCafeMessages
func (m *Mobile) CheckCafeMessages() error {
	if !m.node.Started() {
		return core.ErrOffline
	}

	return m.node.CheckCafeMessages()
}

// CafeRequests paginates new request groups
func (m *Mobile) CafeRequests(limit int) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	groups := m.node.Datastore().CafeRequests().ListGroups("", limit)
	return proto.Marshal(&pb.Strings{Values: groups})
}

// SetCafeRequestComplete marks a request group as complete
func (m *Mobile) SetCafeRequestComplete(group string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.Datastore().CafeRequests().UpdateGroupStatus(group, pb.CafeRequest_COMPLETE)
	if err != nil {
		return err
	}
	return m.node.Datastore().CafeRequests().DeleteCompleteSyncGroups()
}

// SetCafeRequestFailed deletes a cafe request group
func (m *Mobile) SetCafeRequestFailed(group string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	err := m.node.Datastore().CafeRequests().DeleteByGroup(group)
	if err != nil {
		return err
	}
	return m.node.Datastore().CafeRequests().DeleteCompleteSyncGroups()
}

// WriteCafeHTTPRequests a list of request objects for the given group, writing bodies to disk
// Note: This also marks the group as pending
// - store: PUT /store/:cid, body => raw object data
// - unstore: DELETE /store/:cid, body => none
// - store thread: PUT /threads/:id, body => encrypted thread object (snapshot)
// - unstore thread: DELETE /threads/:id, body => none
// - deliver message: POST /inbox/:pid, body => encrypted message
func (m *Mobile) WriteCafeHTTPRequests(group string) ([]byte, error) {
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
	var pending []string

	// group by cafe
	creqs := make(map[string][]*pb.CafeRequest)
	for _, req := range reqs.Items {
		creqs[req.Cafe.Peer] = append(creqs[req.Cafe.Peer], req)
	}

	hreqs := &pb.CafeHTTPRequestList{Items: make([]*pb.CafeHTTPRequest, 0)}
	for cafe, reqs := range creqs {

		// load the session for this cafe
		session := m.node.Datastore().CafeSessions().Get(cafe)
		if session == nil {
			return nil, fmt.Errorf("session for cafe %s not found", cafe)
		}

		// group by type
		treqs := make(map[pb.CafeRequest_Type][]*pb.CafeRequest)
		for _, req := range reqs {
			treqs[req.Type] = append(treqs[req.Type], req)
		}
		for rtype, reqs := range treqs {

			// store reqs can be handled by multipart form data
			if rtype == pb.CafeRequest_STORE {
				hreq := &pb.CafeHTTPRequest{
					Type: pb.CafeHTTPRequest_PUT,
					Url:  session.Cafe.Url + "/api/v1/store",
					Headers: map[string]string{
						"Authorization":  "Basic " + session.Access,
						"X-Textile-Peer": m.node.Ipfs().Identity.Pretty(),
					},
					Path: fmt.Sprintf("%s/tmp/%s", m.RepoPath, ksuid.New().String()),
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
							data, err := ipfs.GetObjectAtPath(m.node.Ipfs(), req.Target)
							if err != nil {
								return nil, err
							}
							hreq.Headers["X-Textile-Store-Type"] = "object"
							_, err = part.Write(data)
							if err != nil {
								return nil, err
							}
						} else {
							return nil, err
						}
					} else {
						hreq.Headers["X-Textile-Store-Type"] = "data"
						_, err = part.Write(data)
						if err != nil {
							return nil, err
						}
					}

					pending = append(pending, req.Id)
				}
				_ = writer.Close()
				_ = file.Close()

				hreq.Headers["Content-Type"] = writer.FormDataContentType()

				hreqs.Items = append(hreqs.Items, hreq)

			} else {
				for _, req := range reqs {
					hreq := &pb.CafeHTTPRequest{
						Url: session.Cafe.Url + "/api/v1",
						Headers: map[string]string{
							"Authorization":  "Basic " + session.Access,
							"X-Textile-Peer": m.node.Ipfs().Identity.Pretty(),
							"Content-Type":   "application/octet-stream",
						},
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
						hreq.Url += "/inbox/" + req.Peer
						body = []byte(req.Target)
					}

					hreq.Headers["Content-Type"] = "application/octet-stream"
					hreq.Path = fmt.Sprintf("%s/tmp/%s", m.RepoPath, req.Id)
					err := util.WriteFileByPath(hreq.Path, body)
					if err != nil {
						return nil, err
					}

					hreqs.Items = append(hreqs.Items, hreq)
					pending = append(pending, req.Id)
				}
			}
		}
	}

	// mark as pending
	for _, id := range pending {
		err = m.node.Datastore().CafeRequests().UpdateStatus(id, pb.CafeRequest_PENDING)
		if err != nil {
			return nil, err
		}
	}

	return proto.Marshal(hreqs)
}

// CafeRequestGroupStatus returns the status of the request group
func (m *Mobile) CafeRequestGroupStatus(group string) ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.Datastore().CafeRequests().SyncGroupStatus(group))
}

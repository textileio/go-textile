package mobile

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	ipld "github.com/ipfs/go-ipld-format"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

// PeerId returns the ipfs peer id
func (m *Mobile) PeerId() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	pid, err := m.node.PeerId()
	if err != nil {
		return "", err
	}
	return pid.Pretty(), nil
}

// SwarmConnect opens a new direct connection to a peer using an IPFS multiaddr
func (m *Mobile) SwarmConnect(address string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	results, err := ipfs.SwarmConnect(m.node.Ipfs(), []string{address})
	if err != nil {
		return "", err
	}

	return results[0], nil
}

// DataAtPath is the async version of dataAtPath
func (m *Mobile) DataAtPath(pth string, cb DataCallback) {
	m.node.WaitAdd(1, "Mobile.DataAtPath")
	go func() {
		defer m.node.WaitDone("Mobile.DataAtPath")
		cb.Call(m.dataAtPath(pth))
	}()
}

// dataAtPath calls core DataAtPath
func (m *Mobile) dataAtPath(pth string) ([]byte, string, error) {
	if !m.node.Started() {
		return nil, "", core.ErrStopped
	}

	data, err := m.node.DataAtPath(pth)
	if err != nil {
		if err == ipld.ErrNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	media, err := m.node.GetMedia(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	return data, media, nil
}

// IpfsPubsubPub publishes a message to a given pubsub topic
func (m *Mobile) IpfsPubsubPub(topic string, data string) error {
	if !m.node.Started() {
		return core.ErrStopped
	}

	payload := []byte(data)
	err := ipfs.Publish(m.node.Ipfs(), topic, payload)
	if err != nil {
		return err
	}

	return nil
}

// CancelIpfsPubsubSub is used to cancel the request
func (m *Mobile) CancelIpfsPubsubSub(queryId string) {
	index := -1
	for i, handle := range ipfsPubsubSubHandles {
		if (queryId == handle.Id) {
			handle.cancel.Close()
			handle.done()
			index = i
			break
		}
	}
	if index != -1 {
		ipfsPubsubSubHandles = append(ipfsPubsubSubHandles[:index], ipfsPubsubSubHandles[index + 1:]...)
	}
}

var ipfsPubsubSubHandles = []*SearchHandle{}

// IpfsPubsubSub Subscribes to messages on a given topic
func (m *Mobile) IpfsPubsubSub(topic string) (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	msgs := make(chan iface.PubSubMessage, 10)
	ctx := m.node.Ipfs().Context()
	id := ksuid.New().String()
	go func() {
		if err := ipfs.Subscribe(m.node.Ipfs(), ctx, topic, true, msgs); err != nil {
			close(msgs)
			m.notify(pb.MobileEventType_QUERY_RESPONSE, &pb.MobileQueryEvent{
				Id:   id,
				Type: pb.MobileQueryEvent_ERROR,
				Error: &pb.Error{
					Code:    500,
					Message: err.Error(),
				},
			})
			log.Errorf("ipfs pubsub sub stopped with error: %s", err.Error())
			return
		}
	}()
	log.Infof("ipfs pubsub sub started for %s", topic)

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
	cancel := broadcast.NewBroadcaster(0)
	ipfsPubsubSubHandles = append(ipfsPubsubSubHandles, &SearchHandle{
		Id:     id,
		cancel: cancel,
		done:   doneFn,
	})
	cancelCh := cancel.Listen().Ch

	go func() {
		for {
			select {
			case <-cancelCh:
				log.Infof("ipfs pubsub sub shutdown for %s", topic)
				return
			case msg, ok := <-msgs:
				if !ok {
					index := -1
					for i, handle := range ipfsPubsubSubHandles {
						if (id == handle.Id) {
							index = i
							break
						}
					}
					if index != -1 {
						ipfsPubsubSubHandles = append(ipfsPubsubSubHandles[:index], ipfsPubsubSubHandles[index + 1:]...)
					}

					doneFn()
					log.Infof("ipfs pubsub sub shutdown for %s", topic)
					return
				}

				mPeer := msg.From()
				if mPeer.Pretty() == m.node.Ipfs().Identity.Pretty() {
					break
				}

				value, err := proto.Marshal(&pb.Strings{
					Values: []string{string(msg.Data())},
				})
				if err != nil {
					m.notify(pb.MobileEventType_QUERY_RESPONSE, &pb.MobileQueryEvent{
						Id:   id,
						Type: pb.MobileQueryEvent_ERROR,
						Error: &pb.Error{
							Code:    500,
							Message: err.Error(),
						},
					})
					break
				}

				res := &pb.QueryResult{
					Id:    fmt.Sprintf("%x", msg.Seq()),
					Value: &any.Any{
						TypeUrl: "/Strings",
						Value:   value,
					},
				}
				m.notify(pb.MobileEventType_QUERY_RESPONSE, &pb.MobileQueryEvent{
					Id:   id,
					Type: pb.MobileQueryEvent_DATA,
					Data: res,
				})
			}
		}
	}()

	return id, nil
}

package service

import (
	"context"
	"fmt"
	"github.com/textileio/textile-go/pb"
	inet "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	ggio "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/io"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"sync"
	"time"
)

// sender is a reusable peer stream
type sender struct {
	protocol   protocol.ID
	stream     inet.Stream
	writer     ggio.WriteCloser
	reader     ggio.ReadCloser
	mux        sync.Mutex
	pid        peer.ID
	service    *Service
	singleMsg  int
	invalid    bool
	requests   map[int32]chan *pb.Envelope
	requestMux sync.Mutex
}

// ReadMessageTimeout specifies timeout for reading a response
var ReadMessageTimeout = time.Minute * 1

// ErrReadTimeout read response has timed out
var ErrReadTimeout = fmt.Errorf("timed out reading response")

// messageSenderForPeer creates a sender for the given peer
func (s *Service) messageSenderForPeer(pid peer.ID, proto protocol.ID) (*sender, error) {
	s.senderMux.Lock()
	ms, ok := s.sender[pid]
	if ok {
		s.senderMux.Unlock()
		return ms, nil
	}
	ms = &sender{
		protocol: proto,
		pid:      pid,
		service:  s,
		requests: make(map[int32]chan *pb.Envelope, 2),
	}
	s.sender[pid] = ms
	s.senderMux.Unlock()

	if err := ms.prepOrInvalidate(); err != nil {
		s.senderMux.Lock()
		defer s.senderMux.Unlock()
		if msCur, ok := s.sender[pid]; ok {
			// Changed. Use the new one, old one is invalid and
			// not in the map so we can just throw it away.
			if ms != msCur {
				return msCur, nil
			}
			// Not changed, remove the now invalid stream from the
			// map.
			delete(s.sender, pid)
		}
		// Invalid but not in map. Must have been removed by a disconnect.
		return nil, err
	}
	// All ready to go.
	return ms, nil
}

// invalidate is called before this sender is removed from the strmap.
// It prevents the sender from being reused/reinitialized and then
// forgotten (leaving the stream open).
func (ms *sender) invalidate() {
	ms.invalid = true
	if ms.stream != nil {
		ms.stream.Reset()
		ms.stream = nil
	}
}

// prepOrInvalidate invalidates a sender if prep fails
func (ms *sender) prepOrInvalidate() error {
	ms.mux.Lock()
	defer ms.mux.Unlock()
	if err := ms.prep(); err != nil {
		ms.invalidate()
		return err
	}
	return nil
}

// prep creates a new stream, reader, and writer for a sender
func (ms *sender) prep() error {
	if ms.invalid {
		return fmt.Errorf("message sender has been invalidated")
	}
	if ms.stream != nil {
		return nil
	}
	nstr, err := ms.service.Node.PeerHost.NewStream(ms.service.Node.Context(), ms.pid, ms.protocol)
	if err != nil {
		return err
	}
	ms.reader = ggio.NewDelimitedReader(nstr, inet.MessageSizeMax)
	ms.writer = ggio.NewDelimitedWriter(nstr)
	ms.stream = nstr
	return nil
}

// streamReuseTries is the number of times we will try to reuse a stream to a
// given peer before giving up and reverting to the old one-message-per-stream
// behaviour.
const streamReuseTries = 3

// SendMessage sends a message to a peer (a response is no expected)
func (ms *sender) SendMessage(ctx context.Context, pmes *pb.Envelope) error {
	ms.mux.Lock()
	defer ms.mux.Unlock()
	retry := false
	for {
		if err := ms.prep(); err != nil {
			return err
		}
		if err := ms.writer.WriteMsg(pmes); err != nil {
			ms.stream.Reset()
			ms.stream = nil
			if retry {
				return err
			} else {
				retry = true
				continue
			}
		}
		if ms.singleMsg > streamReuseTries {
			ms.stream.Close()
			ms.stream = nil
		} else if retry {
			ms.singleMsg++
		}
		return nil
	}
}

// SendRequest sends a message to a peer and expects a response
func (ms *sender) SendRequest(ctx context.Context, pmes *pb.Envelope) (*pb.Envelope, error) {
	returnChan := make(chan *pb.Envelope)
	ms.requestMux.Lock()
	ms.requests[pmes.Message.RequestId] = returnChan
	ms.requestMux.Unlock()

	ms.mux.Lock()
	defer ms.mux.Unlock()
	retry := false
	for {
		if err := ms.prep(); err != nil {
			return nil, err
		}
		if err := ms.writer.WriteMsg(pmes); err != nil {
			ms.stream.Reset()
			ms.stream = nil
			if retry {
				return nil, err
			} else {
				retry = true
				continue
			}
		}
		mes, err := ms.ctxReadMsg(ctx, returnChan)
		if err != nil {
			ms.stream.Reset()
			ms.stream = nil
			return nil, err
		}
		if ms.singleMsg > streamReuseTries {
			ms.stream.Close()
			ms.stream = nil
		} else if retry {
			ms.singleMsg++
		}
		return mes, nil
	}
}

// closeRequest stop listening for responses
func (ms *sender) closeRequest(id int32) {
	ms.requestMux.Lock()
	ch, ok := ms.requests[id]
	if ok {
		close(ch)
		delete(ms.requests, id)
	}
	ms.requestMux.Unlock()
}

// ctxReadMsg reads a response
func (ms *sender) ctxReadMsg(ctx context.Context, returnChan chan *pb.Envelope) (*pb.Envelope, error) {
	t := time.NewTimer(ReadMessageTimeout)
	defer t.Stop()
	select {
	case mes := <-returnChan:
		return mes, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-t.C:
		return nil, ErrReadTimeout
	}
}

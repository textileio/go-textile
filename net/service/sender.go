package service

import (
	"context"
	"fmt"
	"github.com/textileio/textile-go/pb"
	inet "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	ggio "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/io"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"math/rand"
	"sync"
	"time"
)

type sender struct {
	s         inet.Stream
	w         ggio.WriteCloser
	r         ggio.ReadCloser
	lk        sync.Mutex
	p         peer.ID
	service   *TextileService
	singleMes int
	invalid   bool
	requests  map[int32]chan *pb.Envelope
	requestlk sync.Mutex
}

var ReadMessageTimeout = time.Minute * 5
var ErrReadTimeout = fmt.Errorf("timed out reading response")

func (s *TextileService) messageSenderForPeer(pid peer.ID) (*sender, error) {
	defer func() {
		if recover() != nil {
			log.Error("recovered from messageSenderForPeer")
		}
	}()
	s.senderlk.Lock()
	ms, ok := s.sender[pid]
	if ok {
		s.senderlk.Unlock()
		return ms, nil
	}
	ms = &sender{p: pid, service: s, requests: make(map[int32]chan *pb.Envelope, 2)}
	s.sender[pid] = ms
	s.senderlk.Unlock()

	if err := ms.prepOrInvalidate(); err != nil {
		s.senderlk.Lock()
		defer s.senderlk.Unlock()

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

func (s *TextileService) newMessageSender(pid peer.ID) *sender {
	return &sender{
		p:        pid,
		service:  s,
		requests: make(map[int32]chan *pb.Envelope, 2), // low initial capacity
	}
}

// invalidate is called before this sender is removed from the strmap.
// It prevents the sender from being reused/reinitialized and then
// forgotten (leaving the stream open).
func (ms *sender) invalidate() {
	ms.invalid = true
	if ms.s != nil {
		ms.s.Reset()
		ms.s = nil
	}
}

func (ms *sender) prepOrInvalidate() error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	if err := ms.prep(); err != nil {
		ms.invalidate()
		return err
	}
	return nil
}

func (ms *sender) prep() error {
	if ms.invalid {
		return fmt.Errorf("message sender has been invalidated")
	}
	if ms.s != nil {
		return nil
	}

	nstr, err := ms.service.host.NewStream(ms.service.ctx, ms.p, ProtocolTextile)
	if err != nil {
		return err
	}

	ms.r = ggio.NewDelimitedReader(nstr, inet.MessageSizeMax)
	ms.w = ggio.NewDelimitedWriter(nstr)
	ms.s = nstr

	return nil
}

// streamReuseTries is the number of times we will try to reuse a stream to a
// given peer before giving up and reverting to the old one-message-per-stream
// behaviour.
const streamReuseTries = 3

func (ms *sender) SendMessage(ctx context.Context, pmes *pb.Envelope) error {
	defer func() {
		if recover() != nil {
			log.Error("recovered from sender.SendMessage")
		}
	}()
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if err := ms.prep(); err != nil {
			return err
		}

		if err := ms.w.WriteMsg(pmes); err != nil {
			ms.s.Reset()
			ms.s = nil

			if retry {
				return err
			} else {
				retry = true
				continue
			}
		}

		if ms.singleMes > streamReuseTries {
			ms.s.Close()
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return nil
	}
}

func (ms *sender) SendRequest(ctx context.Context, pmes *pb.Envelope) (*pb.Envelope, error) {
	defer func() {
		if recover() != nil {
			log.Error("recovered from sender.SendRequest")
		}
	}()
	pmes.Message.RequestId = rand.Int31()
	returnChan := make(chan *pb.Envelope)
	ms.requestlk.Lock()
	ms.requests[pmes.Message.RequestId] = returnChan
	ms.requestlk.Unlock()

	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if err := ms.prep(); err != nil {
			return nil, err
		}

		if err := ms.w.WriteMsg(pmes); err != nil {
			ms.s.Reset()
			ms.s = nil

			if retry {
				return nil, err
			} else {
				retry = true
				continue
			}
		}

		mes, err := ms.ctxReadMsg(ctx, returnChan)
		if err != nil {
			ms.s.Reset()
			ms.s = nil
			return nil, err
		}

		if ms.singleMes > streamReuseTries {
			ms.s.Close()
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return mes, nil
	}
}

// stop listening for responses
func (ms *sender) closeRequest(id int32) {
	ms.requestlk.Lock()
	ch, ok := ms.requests[id]
	if ok {
		close(ch)
		delete(ms.requests, id)
	}
	ms.requestlk.Unlock()
}

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

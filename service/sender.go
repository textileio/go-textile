package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/helpers"
	inet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-msgio"
	"github.com/textileio/go-textile/pb"
)

func (srv *Service) updateFromMessage(ctx context.Context, p peer.ID) error {
	// Make sure that this node is actually a DHT server, not just a client.
	protos, err := srv.Node().Peerstore.SupportsProtocols(p, string(srv.handler.Protocol()))
	if err == nil && len(protos) > 0 {
		srv.Node().DHT.Update(ctx, p)
	}
	return nil
}

func (srv *Service) messageSenderForPeer(ctx context.Context, p peer.ID) (*messageSender, error) {
	srv.smlk.Lock()
	ms, ok := srv.strmap[p]
	if ok {
		srv.smlk.Unlock()
		return ms, nil
	}
	ms = &messageSender{p: p, srv: srv}
	srv.strmap[p] = ms
	srv.smlk.Unlock()

	if err := ms.prepOrInvalidate(ctx); err != nil {
		srv.smlk.Lock()
		defer srv.smlk.Unlock()

		if msCur, ok := srv.strmap[p]; ok {
			// Changed. Use the new one, old one is invalid and
			// not in the map so we can just throw it away.
			if ms != msCur {
				return msCur, nil
			}
			// Not changed, remove the now invalid stream from the
			// map.
			delete(srv.strmap, p)
		}
		// Invalid but not in map. Must have been removed by a disconnect.
		return nil, err
	}
	// All ready to go.
	return ms, nil
}

type messageSender struct {
	s  inet.Stream
	r  msgio.ReadCloser
	lk sync.Mutex
	p  peer.ID

	srv *Service

	invalid   bool
	singleMes int
}

// invalidate is called before this messageSender is removed from the strmap.
// It prevents the messageSender from being reused/reinitialized and then
// forgotten (leaving the stream open).
func (ms *messageSender) invalidate() {
	ms.invalid = true
	if ms.s != nil {
		_ = ms.s.Reset()
		ms.s = nil
	}
}

func (ms *messageSender) prepOrInvalidate(ctx context.Context) error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	if err := ms.prep(ctx); err != nil {
		ms.invalidate()
		return err
	}
	return nil
}

func (ms *messageSender) prep(ctx context.Context) error {
	if ms.invalid {
		return fmt.Errorf("message sender has been invalidated")
	}
	if ms.s != nil {
		return nil
	}

	nstr, err := ms.srv.Node().PeerHost.NewStream(ctx, ms.p, ms.srv.handler.Protocol())
	if err != nil {
		return err
	}

	ms.r = msgio.NewVarintReaderSize(nstr, inet.MessageSizeMax)
	ms.s = nstr

	return nil
}

// streamReuseTries is the number of times we will try to reuse a stream to a
// given peer before giving up and reverting to the old one-message-per-stream
// behaviour.
const streamReuseTries = 3

func (ms *messageSender) SendMessage(ctx context.Context, pmes *pb.Envelope) error {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if err := ms.prep(ctx); err != nil {
			return err
		}

		if err := ms.writeMsg(pmes); err != nil {
			_ = ms.s.Reset()
			ms.s = nil

			if retry {
				log.Info("error writing message, bailing: ", err)
				return err
			}
			log.Info("error writing message, trying again: ", err)
			retry = true
			continue
		}

		if ms.singleMes > streamReuseTries {
			go helpers.FullClose(ms.s)
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return nil
	}
}

func (ms *messageSender) SendRequest(ctx context.Context, pmes *pb.Envelope) (*pb.Envelope, error) {
	ms.lk.Lock()
	defer ms.lk.Unlock()
	retry := false
	for {
		if err := ms.prep(ctx); err != nil {
			return nil, err
		}

		if err := ms.writeMsg(pmes); err != nil {
			_ = ms.s.Reset()
			ms.s = nil

			if retry {
				log.Info("error writing message, bailing: ", err)
				return nil, err
			}
			log.Info("error writing message, trying again: ", err)
			retry = true
			continue
		}

		mes := new(pb.Envelope)
		if err := ms.ctxReadMsg(ctx, mes); err != nil {
			_ = ms.s.Reset()
			ms.s = nil

			if retry {
				log.Info("error reading message, bailing: ", err)
				return nil, err
			}
			log.Info("error reading message, trying again: ", err)
			retry = true
			continue
		}

		if ms.singleMes > streamReuseTries {
			go helpers.FullClose(ms.s)
			ms.s = nil
		} else if retry {
			ms.singleMes++
		}

		return mes, nil
	}
}

func (ms *messageSender) writeMsg(pmes *pb.Envelope) error {
	return writeMsg(ms.s, pmes)
}

func (ms *messageSender) ctxReadMsg(ctx context.Context, mes *pb.Envelope) error {
	errc := make(chan error, 1)
	go func(r msgio.ReadCloser) {
		bytes, err := r.ReadMsg()
		defer r.ReleaseMsg(bytes)
		if err != nil {
			errc <- err
			return
		}
		errc <- proto.Unmarshal(bytes, mes)
	}(ms.r)

	t := time.NewTimer(dhtReadMessageTimeout)
	defer t.Stop()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return ErrReadTimeout
	}
}

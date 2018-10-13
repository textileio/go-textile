package net

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/cafe"
	"github.com/textileio/textile-go/cafe/client"
	ipfsutil "github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface"
	"io"
	"sync"
	"time"
)

const kFrequency = time.Minute * 10
const groupSize = 5

var ErrTokenExpired = errors.New("token expired")
var ErrPinRequestEmpty = errors.New("pin request empty response")
var ErrPinRequestMismatch = errors.New("pin request content id mismatch")

type PinnerConfig struct {
	Datastore repo.Datastore
	node      *core.IpfsNode
}

type StoreQueue struct {
	datastore repo.Datastore
	node      *core.IpfsNode
	mux       sync.Mutex
}

func NewStoreQueue(config *PinnerConfig) *StoreQueue {
	return &StoreQueue{
		datastore: config.Datastore,
		node:      config.node,
	}
}

func (p *StoreQueue) Run() {
	tick := time.NewTicker(kFrequency)
	defer tick.Stop()
	go p.Pin()
	for {
		select {
		case <-tick.C:
			go p.Pin()
		}
	}
}

func (p *StoreQueue) Pin() {
	p.mux.Lock()
	defer p.mux.Unlock()

	// get tokens
	tokens, err := p.getTokens(false)
	if err != nil {
		log.Errorf("pinner get tokens error: %s", err)
		return
	}
	if tokens == nil {
		return
	}

	// start at no offset
	if err := p.handle(p.datastore.PinRequests().List("", groupSize), tokens); err != nil {
		log.Errorf("pin error: %s", err)
		return
	}
}

func (p *StoreQueue) Put(id string) error {
	pr := &repo.StoreRequest{Id: id, Date: time.Now()}
	if err := p.datastore.PinRequests().Put(pr); err != nil {
		return err
	}

	// run it now
	go p.Pin()

	return nil
}

func (p *StoreQueue) handle(reqs []repo.StoreRequest) error {
	if len(reqs) == 0 {
		return nil
	}

	// process them
	var expired bool
	var toDelete []string
	wg := sync.WaitGroup{}
	for _, r := range reqs {
		wg.Add(1)
		go func(pr repo.StoreRequest) {
			if err := p.send(pr); err != nil {
				if err == ErrTokenExpired {
					expired = true
				} else {
					log.Errorf("pin request %s failed: %s", pr.Id, err)
				}
			} else {
				toDelete = append(toDelete, pr.Id)
			}
			wg.Done()
		}(r)
	}
	wg.Wait()

	// check expired
	if expired {
		if _, err := p.getTokens(true); err != nil {
			return err
		}
		p.Pin()
		return nil
	}

	log.Debugf("handled %d pin requests, deleting...", len(toDelete))

	// next batch
	offset := reqs[len(reqs)-1].Id
	next := p.datastore.PinRequests().List(offset, pinGroupSize)

	// clean up
	for _, id := range toDelete {
		if err := p.datastore.PinRequests().Delete(id); err != nil {
			log.Errorf("failed to delete pin request %s: %s", id, err)
		}
	}

	// keep going
	return p.handlePins(next, tokens)
}

func (p *StoreQueue) send(pr repo.StoreRequest, tokens *repo.CafeSession) error {
	return Pin(p.node(), pr.Id, tokens, p.url)
}

func Pin(ipfs *core.IpfsNode, id string, tokens *repo.CafeSession, url string) error {
	if tokens == nil {
		return errors.New("pin attempted without tokens")
	}

	// load local content
	cType := "application/octet-stream"
	var reader io.Reader
	data, err := ipfsutil.GetDataAtPath(ipfs, id)
	if err != nil {
		if err == iface.ErrIsDir {
			reader, err = ipfsutil.GetArchiveAtPath(ipfs, id)
			if err != nil {
				return err
			}
			cType = "application/gzip"
		} else {
			return err
		}
	} else {
		reader = bytes.NewReader(data)
	}

	// pin to cafe
	res, err := client.Pin(tokens.Access, reader, url, cType)
	if err != nil {
		return err
	}
	if res.Error != nil {
		if *res.Error == cafe.ErrUnauthorized {
			return ErrTokenExpired
		}
		return errors.New(*res.Error)
	}
	if res.Id == nil {
		return ErrPinRequestEmpty
	}
	if *res.Id != id {
		return ErrPinRequestMismatch
	}
	return nil
}

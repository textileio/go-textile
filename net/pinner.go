package net

import (
	"bytes"
	"github.com/pkg/errors"
	cafe "github.com/textileio/textile-go/core/cafe"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core/coreapi/interface"
	"io"
	"sync"
	"time"
)

const kPinFrequency = time.Minute * 10
const pinGroupSize = 5

var errPinRequestEmpty = errors.New("pin request empty response")
var errPinRequestMismatch = errors.New("pin request content id mismatch")

type PinnerConfig struct {
	Datastore repo.Datastore
	Ipfs      func() *core.IpfsNode
	Url       string
	Tokens    *repo.CafeTokens
}

type Pinner struct {
	datastore repo.Datastore
	ipfs      func() *core.IpfsNode
	url       string
	Tokens    *repo.CafeTokens
	mux       sync.Mutex
}

func NewPinner(config *PinnerConfig) *Pinner {
	return &Pinner{
		datastore: config.Datastore,
		ipfs:      config.Ipfs,
		url:       config.Url,
		Tokens:    config.Tokens,
	}
}

func (p *Pinner) Url() string {
	return p.url
}

func (p *Pinner) Run() {
	tick := time.NewTicker(kPinFrequency)
	defer tick.Stop()
	go p.Pin()
	for {
		select {
		case <-tick.C:
			go p.Pin()
		}
	}
}

func (p *Pinner) Pin() {
	p.mux.Lock()
	defer p.mux.Unlock()

	// check tokens
	if p.Tokens == nil {
		return
	}

	if err := p.handlePin(""); err != nil {
		return
	}
}

func (p *Pinner) Put(id string) error {
	pr := &repo.PinRequest{Id: id, Date: time.Now()}
	if err := p.datastore.PinRequests().Put(pr); err != nil {
		return err
	}

	// run it now
	go p.Pin()

	return nil
}

func (p *Pinner) handlePin(offset string) error {
	// get pending pin list
	prs := p.datastore.PinRequests().List(offset, pinGroupSize)
	if len(prs) == 0 {
		return nil
	}
	log.Debugf("handling %d pin requests...", len(prs))

	// process them
	var toDelete []string
	wg := sync.WaitGroup{}
	for _, r := range prs {
		wg.Add(1)
		go func(pr repo.PinRequest) {
			if err := p.send(pr); err != nil {
				log.Errorf("pin request %s failed: %s", pr.Id, err)
			} else {
				toDelete = append(toDelete, pr.Id)
			}
			wg.Done()
		}(r)
	}
	wg.Wait()
	log.Debugf("handled %d pin requests, deleting...", len(toDelete))

	// clean up
	for _, id := range toDelete {
		if err := p.datastore.PinRequests().Delete(id); err != nil {
			log.Errorf("failed to delete pin request %s: %s", id, err)
		}
	}

	// keep going
	return p.handlePin(prs[len(prs)-1].Id)
}

func (p *Pinner) send(pr repo.PinRequest) error {
	return Pin(p.ipfs(), pr.Id, p.Tokens, p.url)
}

func Pin(ipfs *core.IpfsNode, id string, tokens *repo.CafeTokens, url string) error {
	// load local content
	cType := "application/octet-stream"
	var reader io.Reader
	data, err := util.GetDataAtPath(ipfs, id)
	if err != nil {
		if err == iface.ErrIsDir {
			reader, err = util.GetArchiveAtPath(ipfs, id)
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
	res, err := cafe.Pin(tokens.Access, reader, url, cType)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New(*res.Error)
	}
	if res.Id == nil {
		return errPinRequestEmpty
	}
	if *res.Id != id {
		return errPinRequestMismatch
	}
	return nil
}

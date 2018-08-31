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
	"net/http"
	"sync"
	"time"
)

const kPinFrequency = time.Minute * 10
const pinGroupSize = 5

var ErrTokenExpired = errors.New("token expired")
var ErrPinRequestEmpty = errors.New("pin request empty response")
var ErrPinRequestMismatch = errors.New("pin request content id mismatch")

type PinnerConfig struct {
	Datastore repo.Datastore
	Ipfs      func() *core.IpfsNode
	Url       string
	GetTokens func(bool) (*repo.CafeTokens, error)
}

type Pinner struct {
	datastore repo.Datastore
	ipfs      func() *core.IpfsNode
	url       string
	getTokens func(bool) (*repo.CafeTokens, error)
	mux       sync.Mutex
}

func NewPinner(config *PinnerConfig) *Pinner {
	return &Pinner{
		datastore: config.Datastore,
		ipfs:      config.Ipfs,
		url:       config.Url,
		getTokens: config.GetTokens,
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
	if err := p.handlePins(p.datastore.PinRequests().List("", pinGroupSize), tokens); err != nil {
		log.Errorf("pin error: %s", err)
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

func (p *Pinner) handlePins(pins []repo.PinRequest, tokens *repo.CafeTokens) error {
	if len(pins) == 0 {
		return nil
	}

	// process them
	var expired bool
	var toDelete []string
	wg := sync.WaitGroup{}
	for _, r := range pins {
		wg.Add(1)
		go func(pr repo.PinRequest) {
			if err := p.send(pr, tokens); err != nil {
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
	offset := pins[len(pins)-1].Id
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

func (p *Pinner) send(pr repo.PinRequest, tokens *repo.CafeTokens) error {
	return Pin(p.ipfs(), pr.Id, tokens, p.url)
}

func Pin(ipfs *core.IpfsNode, id string, tokens *repo.CafeTokens, url string) error {
	if tokens == nil {
		return errors.New("pin attempted without tokens")
	}

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
	if res.Status == http.StatusUnauthorized {
		return ErrTokenExpired
	}
	if res.Error != nil {
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

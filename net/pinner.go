package net

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/util"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"net/http"
	"sync"
	"time"
)

const kPinFrequency = time.Minute * 10
const pinGroupSize = 5

var errPinRequestFailed = errors.New("pin request failed")

type PinnerConfig struct {
	Datastore repo.Datastore
	Ipfs      func() *core.IpfsNode
	Api       string
}

type Pinner struct {
	datastore repo.Datastore
	ipfs      func() *core.IpfsNode
	api       string
	mux       sync.Mutex
}

func NewPinner(config PinnerConfig) *Pinner {
	return &Pinner{datastore: config.Datastore, ipfs: config.Ipfs, api: config.Api}
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
	if err := p.handlePin(""); err != nil {
		return
	}
}

func (p *Pinner) Put(id string) error {
	pr := &repo.PinRequest{Id: id, Date: time.Now()}
	err := p.datastore.PinRequests().Put(pr)
	if err != nil {
		return err
	}
	log.Debugf("put pin request for %s", id)

	// run it now
	go p.Pin()

	return nil
}

func (p *Pinner) handlePin(offset string) error {
	prs := p.datastore.PinRequests().List(offset, pinGroupSize)
	if len(prs) == 0 {
		return nil
	}
	log.Debugf("handling %d pin requests...", len(prs))

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
	log.Debugf("successfully handled %d pin requests, deleting...", len(toDelete))

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
	// get token
	token, _, err := p.datastore.Profile().GetTokens()
	if err != nil {
		return err
	}

	// load local content
	data, err := util.GetDataAtPath(p.ipfs(), pr.Id)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)

	// make the request
	req, err := http.NewRequest("POST", p.api, reader)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check status
	if resp.StatusCode != 201 {
		return errPinRequestFailed
	}
	return nil
}

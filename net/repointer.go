package net

import (
	"context"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"time"
)

const kRepointFrequency = time.Hour * 12
const kPointerExpiration = time.Hour * 24 * 30

type PointerRepublisher struct {
	ipfs        *core.IpfsNode
	datastore   repo.Datastore
	isModerator func() bool
}

func NewPointerRepublisher(node *core.IpfsNode, datastore repo.Datastore) *PointerRepublisher {
	return &PointerRepublisher{
		ipfs:      node,
		datastore: datastore,
		isModerator: func() bool {
			return false
		},
	}
}

func (r *PointerRepublisher) Run() {
	tick := time.NewTicker(kRepointFrequency)
	defer tick.Stop()
	go r.Republish()
	for range tick.C {
		go r.Republish()
	}
}

func (r *PointerRepublisher) Republish() {
	log.Debug("republishing pointers...")

	// get all pointers
	pointers, err := r.datastore.Pointers().GetAll()
	if err != nil {
		log.Errorf("error republishing: %s", err)
		return
	}
	log.Debugf("found %d pointers to republish", len(pointers))

	// republish or delete each pointer
	ctx := context.Background()
	for _, pointer := range pointers {
		switch pointer.Purpose {
		case repo.MESSAGE:
			if time.Now().Sub(pointer.Date) > kPointerExpiration {
				r.datastore.Pointers().Delete(pointer.Value.ID)
				log.Debugf("deleted pointer %s", pointer.Value.ID.Pretty())
			} else {
				go func(p repo.Pointer) {
					repo.PublishPointer(r.ipfs, ctx, p)
					log.Debugf("published pointer %s", p.Value.ID.Pretty())
				}(pointer)
			}
		default:
			continue
		}
	}
}

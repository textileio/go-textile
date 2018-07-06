package net

import (
	"context"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core"
	"time"
)

const kRepointFrequency = time.Hour * 12
const kPointerExpiration = time.Hour * 24 * 30

type PointerRepublisher struct {
	ipfs        *core.IpfsNode
	db          repo.Datastore
	isModerator func() bool
}

func NewPointerRepublisher(node *core.IpfsNode, database repo.Datastore, isModerator func() bool) *PointerRepublisher {
	return &PointerRepublisher{
		ipfs:        node,
		db:          database,
		isModerator: isModerator,
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
	republishModerator := r.isModerator()
	pointers, err := r.db.Pointers().GetAll()
	if err != nil {
		log.Errorf("error republishing: %s", err)
		return
	}
	ctx := context.Background()

	for _, pointer := range pointers {
		switch pointer.Purpose {
		case repo.MESSAGE:
			if time.Now().Sub(pointer.Date) > kPointerExpiration {
				r.db.Pointers().Delete(pointer.Value.ID)
			} else {
				go repo.PublishPointer(r.ipfs, ctx, pointer)
			}
		case repo.MODERATOR:
			if republishModerator {
				go repo.PublishPointer(r.ipfs, ctx, pointer)
			} else {
				r.db.Pointers().Delete(pointer.Value.ID)
			}
		default:
			continue
		}
	}
}

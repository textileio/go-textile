package core

import (
	"crypto/rand"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
)

func (a *api) addThreads(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing thread name")
		return
	}
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	config := AddThreadConfig{
		Name:      args[0],
		Join:      true,
		Initiator: a.node.account.Address(),
	}

	if opts["key"] != "" {
		config.Key = opts["key"]
	} else {
		config.Key = ksuid.New().String()
	}

	if opts["schema"] != "" {
		config.Schema, err = mh.FromB58String(opts["schema"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
		}
	}

	if opts["type"] != "" {
		var err error
		config.Type, err = repo.ThreadTypeFromString(opts["type"])
		if err != nil {
			g.String(http.StatusBadRequest, "invalid thread type")
			return
		}
	} else {
		config.Type = repo.OpenThread
	}

	// make a new secret
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		a.abort500(g, err)
		return
	}

	thrd, err := a.node.AddThread(sk, config)
	if err != nil {
		a.abort500(g, err)
		return
	}
	info, err := thrd.Info()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusCreated, info)
}

func (a *api) lsThreads(g *gin.Context) {
	infos := make([]*ThreadInfo, 0)
	for _, thrd := range a.node.Threads() {
		info, err := thrd.Info()
		if err != nil {
			a.abort500(g, err)
			return
		}
		infos = append(infos, info)
	}
	g.JSON(http.StatusOK, infos)
}

func (a *api) getThreads(g *gin.Context) {
	id := g.Param("id")
	if id == "default" {
		id = a.node.config.Threads.Defaults.ID
	}
	thrd := a.node.Thread(id)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}
	info, err := thrd.Info()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.JSON(http.StatusOK, info)
}

func (a *api) rmThreads(g *gin.Context) {
	id := g.Param("id")
	thrd := a.node.Thread(id)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}
	if _, err := a.node.RemoveThread(id); err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, "ok")
}

func (a *api) streamThreads(g *gin.Context) {
	id := g.Param("id")
	thrd := a.node.Thread(id)
	if thrd == nil {
		g.String(http.StatusNotFound, "thread not found")
		return
	}
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	listener := a.node.GetThreadUpdateListener()
	g.Stream(func(w io.Writer) bool {
		select {
		case update, ok := <-listener.Ch:
			if !ok {
				return false
			}
			if data, ok := update.(ThreadUpdate); ok {
				info, _ := addBlockInfo(a, data)
				if opts["events"] == "true" {
					g.SSEvent("threadUpdate", info)
				} else {
					g.JSON(http.StatusOK, info)
				}
			}
		default:
		}
		return true
	})

	listener.Close()
}

func addBlockInfo(a *api, update ThreadUpdate) (ThreadUpdate, error) {
	switch update.Block.Type {
	case "FILES":
		info, err := a.node.ThreadFile(update.Block.Id)
		if err != nil {
			log.Errorf("error getting thread file: %s", err)
		}
		return ThreadUpdate{
			Block:      update.Block,
			ThreadId:   update.ThreadId,
			ThreadName: update.ThreadName,
			Info:       info,
		}, nil
	default: // For everything else... we've already go block info
		return update, nil
	}
}

package core

import (
	"crypto/rand"
	"net/http"
	"strings"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"

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

	if opts["sharing"] != "" {
		var err error
		config.Sharing, err = repo.ThreadSharingFromString(opts["sharing"])
		if err != nil {
			g.String(http.StatusBadRequest, "invalid thread sharing")
			return
		}
	} else {
		config.Sharing = repo.NotSharedThread
	}

	if opts["members"] != "" {
		mlist := make([]string, 0)
		for _, m := range strings.Split(opts["members"], ",") {
			if m != "" {
				mlist = append(mlist, m)
			}
		}
		config.Members = mlist
	}

	// make a new secret
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		a.abort500(g, err)
		return
	}

	thrd, err := a.node.AddThread(sk, config)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
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

func (a *api) peersThreads(g *gin.Context) {
	id := g.Param("id")
	if id == "default" {
		id = a.node.config.Threads.Defaults.ID
	}

	thrd := a.node.Thread(id)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}

	contacts := make([]ContactInfo, 0)
	for _, p := range thrd.Peers() {
		contact := a.node.Contact(p.Id)
		if contact != nil {
			contacts = append(contacts, *contact)
		}
	}

	g.JSON(http.StatusOK, contacts)
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

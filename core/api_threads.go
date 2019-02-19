package core

import (
	"crypto/rand"
	"net/http"
	"strings"

	"github.com/textileio/textile-go/pb"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
)

// addThreads godoc
// @Summary Adds and joins a new thread
// @Description Adds a new Thread with given name, type, and sharing and members options, returning
// @Description a ThreadInfo object
// @Tags threads
// @Produce application/json
// @Param X-Textile-Args header string true "name")
// @Param X-Textile-Opts header string false "key: A locally unique key used by an app to identify this thread on recovery, schema: Existing Thread Schema IPFS CID, type: Set the thread type to one of 'private', 'read_only', 'public', or 'open', sharing: Set the thread sharing style to one of 'not_shared','invite_only', or 'shared', members: An array of contact addresses. When supplied, the thread will not allow additional peers beyond those in array, useful for 1-1 chat/file sharing" default(type=private,sharing=not_shared,members=)
// @Success 201 {object} core.ThreadInfo "thread"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads [post]
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

func (a *api) addOrUpdateThreads(g *gin.Context) {
	var thrd pb.Thread
	if err := pbUnmarshaler.Unmarshal(g.Request.Body, &thrd); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if thrd.Id == "" || len(thrd.Sk) == 0 {
		g.String(http.StatusBadRequest, "invalid thread")
		return
	}

	if err := a.node.AddOrUpdateThread(&thrd); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.String(http.StatusOK, "ok")
}

// lsThreads godoc
// @Summary Lists info on all threads
// @Description Lists all local threads, returning an array of ThreadInfo objects
// @Tags threads
// @Produce application/json
// @Success 200 {array} core.ThreadInfo "threads"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads [get]
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

// getThreads godoc
// @Summary Gets a thread
// @Description Gets and displays info about a thread
// @Tags threads
// @Produce application/json
// @Param id path string true "thread id")
// @Success 200 {object} core.ThreadInfo "thread"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id} [get]
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

// peersThreads godoc
// @Summary List all thread peers
// @Description Lists all peers in a thread, optionally listing peers in the default thread
// @Tags threads
// @Produce application/json
// @Param id path string true "thread id")
// @Success 200 {array} core.ContactInfo "contacts"
// @Failure 404 {string} string "Not Found"
// @Router /threads/{id}/peers [get]
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

// rmThreads godoc
// @Summary Leave and remove a thread
// @Description Leaves and removes a thread
// @Tags threads
// @Produce application/json
// @Param id path string true "thread id")
// @Success 200 {string} string "ok"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id} [del]
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

package core

import (
	"crypto/rand"
	"net/http"

	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/util"
)

// addThreads godoc
// @Summary Adds and joins a new thread
// @Description Adds a new Thread with given name, type, and sharing and members options, returning
// @Description a Thread object
// @Tags threads
// @Produce application/json
// @Param X-Textile-Args header string true "name")
// @Param X-Textile-Opts header string false "key: A locally unique key used by an app to identify this thread on recovery, schema: Existing Thread Schema IPFS CID, type: Set the thread type to one of 'private', 'read_only', 'public', or 'open', sharing: Set the thread sharing style to one of 'not_shared','invite_only', or 'shared', members: An array of contact addresses. When supplied, the thread will not allow additional peers beyond those in array, useful for 1-1 chat/file sharing" default(type=private,sharing=not_shared,members=)
// @Success 201 {object} pb.Thread "thread"
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

	config := pb.AddThreadConfig{
		Name: args[0],
	}

	if opts["key"] != "" {
		config.Key = opts["key"]
	} else {
		config.Key = ksuid.New().String()
	}

	if opts["schema"] != "" {
		config.Schema.Id = opts["schema"]
	}

	config.Type = pb.Thread_Type(pbValForEnumString(pb.Thread_Type_value, opts["type"]))
	config.Sharing = pb.Thread_Sharing(pbValForEnumString(pb.Thread_Sharing_value, opts["sharing"]))
	config.Members = util.SplitString(opts["members"], ",")

	// make a new secret
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		a.abort500(g, err)
		return
	}

	thrd, err := a.node.AddThread(config, sk, a.node.account.Address(), true)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	view, err := a.node.ThreadView(thrd.Id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	pbJSON(g, http.StatusCreated, view)
}

// addOrUpdateThreads godoc
// @Summary Add or update a thread directly
// @Description Adds or updates a thread directly, usually from a backup
// @Tags threads
// @Produce application/json
// @Param thread body pb.Thread true "thread")
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Router /threads/{id} [put]
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

	if thrd.Id != g.Param("id") {
		g.String(http.StatusBadRequest, "thread id mismatch")
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
// @Description Lists all local threads, returning a ThreadList object
// @Tags threads
// @Produce application/json
// @Success 200 {object} pb.ThreadList "threads"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads [get]
func (a *api) lsThreads(g *gin.Context) {
	views := &pb.ThreadList{
		Items: make([]*pb.Thread, 0),
	}
	for _, thrd := range a.node.Threads() {
		view, err := a.node.ThreadView(thrd.Id)
		if err != nil {
			a.abort500(g, err)
			return
		}
		views.Items = append(views.Items, view)
	}

	pbJSON(g, http.StatusOK, views)
}

// getThreads godoc
// @Summary Gets a thread
// @Description Gets and displays info about a thread
// @Tags threads
// @Produce application/json
// @Param id path string true "thread id")
// @Success 200 {object} pb.Thread "thread"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id} [get]
func (a *api) getThreads(g *gin.Context) {
	id := g.Param("id")
	if id == "default" {
		id = a.node.config.Threads.Defaults.ID
	}

	view, err := a.node.ThreadView(id)
	if err != nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}

	pbJSON(g, http.StatusOK, view)
}

// peersThreads godoc
// @Summary List all thread peers
// @Description Lists all peers in a thread, optionally listing peers in the default thread
// @Tags threads
// @Produce application/json
// @Param id path string true "thread id")
// @Success 200 {object} pb.ContactList "contacts"
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

	contacts := &pb.ContactList{Items: make([]*pb.Contact, 0)}
	for _, p := range thrd.Peers() {
		contact := a.node.Contact(p.Id)
		if contact != nil {
			contacts.Items = append(contacts.Items, contact)
		}
	}

	pbJSON(g, http.StatusOK, contacts)
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
// @Router /threads/{id} [delete]
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

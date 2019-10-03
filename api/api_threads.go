package api

import (
	"crypto/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

// addThreads godoc
// @Summary Adds and joins a new thread
// @Description Adds a new Thread with given name, type, and sharing and whitelist options, returning
// @Description a Thread object
// @Tags threads
// @Produce application/json
// @Param X-Textile-Args header string true "name"
// @Param X-Textile-Opts header string false "key: A locally unique key used by an app to identify this thread on recovery, schema: Existing Thread Schema IPFS CID, type: Set the thread type to one of 'private', 'read_only', 'public', or 'open', sharing: Set the thread sharing style to one of 'not_shared','invite_only', or 'shared', whitelist: An array of contact addresses. When supplied, the thread will not allow additional peers beyond those in array, useful for 1-1 chat/file sharing" default(type=private,sharing=not_shared,whitelist=)
// @Success 201 {object} pb.Thread "thread"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads [post]
func (a *Api) addThreads(g *gin.Context) {
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
		config.Schema = &pb.AddThreadConfig_Schema{
			Id: opts["schema"],
		}
	}

	config.Type = pb.Thread_Type(pbValForEnumString(pb.Thread_Type_value, opts["type"]))
	config.Sharing = pb.Thread_Sharing(pbValForEnumString(pb.Thread_Sharing_value, opts["sharing"]))
	config.Whitelist = util.SplitString(opts["whitelist"], ",")

	// make a new secret
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		a.abort500(g, err)
		return
	}

	thrd, err := a.Node.AddThread(config, sk, a.Node.Account().Address(), true, true)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	view, err := a.Node.ThreadView(thrd.Id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	a.Node.FlushCafes()

	pbJSON(g, http.StatusCreated, view)
}

// addOrUpdateThreads godoc
// @Summary Add or update a thread directly
// @Description Adds or updates a thread directly, usually from a backup
// @Tags threads
// @Param id path string true "id"
// @Param thread body pb.Thread true "thread"
// @Success 204 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Router /threads/{id} [put]
func (a *Api) addOrUpdateThreads(g *gin.Context) {
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

	if err := a.Node.AddOrUpdateThread(&thrd); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	g.Status(http.StatusNoContent)
}

// renameThreads godoc
// @Summary Rename a thread
// @Description Renames a thread. Only initiators can rename a thread.
// @Tags threads
// @Param id path string true "id"
// @Param X-Textile-Args header string true "name"
// @Success 204 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Router /threads/{id}/name [put]
func (a *Api) renameThreads(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing thread name")
		return
	}

	if err := a.Node.RenameThread(g.Param("id"), args[0]); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	g.Status(http.StatusNoContent)
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
func (a *Api) lsThreads(g *gin.Context) {
	views := &pb.ThreadList{
		Items: make([]*pb.Thread, 0),
	}
	for _, thrd := range a.Node.Threads() {
		view, err := a.Node.ThreadView(thrd.Id)
		if err == nil {
			views.Items = append(views.Items, view)
		} else {
			log.Errorf("error getting thread view %s: %s", thrd.Id, err)
		}
	}

	pbJSON(g, http.StatusOK, views)
}

// getThreads godoc
// @Summary Gets a thread
// @Description Gets and displays info about a thread
// @Tags threads
// @Produce application/json
// @Param id path string true "thread id"
// @Success 200 {object} pb.Thread "thread"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id} [get]
func (a *Api) getThreads(g *gin.Context) {
	id := g.Param("id")

	view, err := a.Node.ThreadView(id)
	if err != nil {
		g.String(http.StatusNotFound, core.ErrThreadNotFound.Error())
		return
	}

	pbJSON(g, http.StatusOK, view)
}

// peersThreads godoc
// @Summary List all thread peers
// @Description Lists all peers in a thread
// @Tags threads
// @Produce application/json
// @Param id path string true "thread id"
// @Success 200 {object} pb.ContactList "contacts"
// @Failure 404 {string} string "Not Found"
// @Router /threads/{id}/peers [get]
func (a *Api) peersThreads(g *gin.Context) {
	id := g.Param("id")

	peers, err := a.Node.ThreadPeers(id)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, peers)
}

// rmThreads godoc
// @Summary Abandons a thread.
// @Description Abandons a thread, and if no one else is participating, then the thread dissipates.
// @Tags threads
// @Param id path string true "thread id"
// @Success 204 {string} string "ok"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id} [delete]
func (a *Api) rmThreads(g *gin.Context) {
	id := g.Param("id")

	thrd := a.Node.Thread(id)
	if thrd == nil {
		g.String(http.StatusNotFound, core.ErrThreadNotFound.Error())
		return
	}

	if _, err := a.Node.RemoveThread(id); err != nil {
		a.abort500(g, err)
		return
	}

	a.Node.FlushCafes()

	g.Status(http.StatusNoContent)
}

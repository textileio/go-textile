package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/core"
)

// addThreadMessages godoc
// @Summary Add a message
// @Description Adds a message to a thread
// @Tags threads
// @Produce application/json
// @Param X-Textile-Args header string true "urlescaped message body"
// @Success 200 {object} pb.Text "message"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id}/messages [post]
func (a *Api) addThreadMessages(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing message body")
		return
	}

	threadId := g.Param("id")
	thrd := a.Node.Thread(threadId)
	if thrd == nil {
		g.String(http.StatusNotFound, core.ErrThreadNotFound.Error())
		return
	}

	// @todo Allow the setting of the target in 0.5.0, which is the new way to comment
	hash, err := thrd.AddMessage("", args[0])
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	msg, err := a.Node.Message(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	pbJSON(g, http.StatusCreated, msg)
}

// lsThreadMessages godoc
// @Summary Paginates thread messages
// @Description Paginates thread messages
// @Tags messages
// @Produce application/json
// @Param X-Textile-Opts header string false "thread: Thread ID (can also use 'default', omit for all), offset: Offset ID to start listing from (omit for latest), limit: List page size (default: 5)" default(thread=,offset=,limit=10)
// @Success 200 {object} pb.TextList "messages"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /messages [get]
func (a *Api) lsThreadMessages(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	threadId := opts["thread"]
	if threadId != "" {
		thrd := a.Node.Thread(threadId)
		if thrd == nil {
			g.String(http.StatusNotFound, core.ErrThreadNotFound.Error())
			return
		}
	}

	limit := 10
	if opts["limit"] != "" {
		limit, err = strconv.Atoi(opts["limit"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	list, err := a.Node.Messages(opts["offset"], limit, threadId)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, list)
}

// getThreadMessages godoc
// @Summary Get thread message
// @Description Gets a thread message by block ID
// @Tags messages
// @Produce application/json
// @Param block path string true "block id"
// @Success 200 {object} pb.Text "message"
// @Failure 400 {string} string "Bad Request"
// @Router /messages/{block} [get]
func (a *Api) getThreadMessages(g *gin.Context) {
	info, err := a.Node.Message(g.Param("block"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}

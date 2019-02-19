package core

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// getThreadsSub godoc
// @Summary Subscribe to thread updates
// @Description Subscribes to updates in a thread or all threads. An update is generated
// @Description when a new block is added to a thread. There are several update types:
// @Description JOIN, ANNOUNCE, LEAVE, MESSAGE, FILES, COMMENT, LIKE, MERGE, IGNORE, FLAG
// @Tags sub
// @Produce application/json
// @Param id path string false "thread id, omit to stream all events"
// @Param X-Textile-Opts header string false "type: Or'd list of event types (e.g., FILES|COMMENTS|LIKES) or empty to include all types, events: Whether to emit Server-Sent Events (SSEvent) or plain JSON" default(type=,events="false")
// @Success 200 {object} core.ThreadUpdate "stream of updates"
// @Failure 500 {string} string "Internal Server Error"
// @Router /sub/{id} [get]
func (a *api) getThreadsSub(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	// Expects or'd list of event types (e.g., FILES|COMMENTS|LIKES).
	types := strings.Split(strings.TrimSpace(strings.ToUpper(opts["type"])), "|")
	threadId := g.Param("id")
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}

	listener := a.node.ThreadUpdateListener()
	g.Stream(func(w io.Writer) bool {
		select {
		case <-g.Request.Context().Done():
			return false

		case update, ok := <-listener.Ch:
			if !ok {
				return false
			}
			if data, ok := update.(ThreadUpdate); ok {
				if threadId != "" && data.ThreadId != threadId {
					break
				}
				for _, t := range types {
					if t == "" || data.Block.Type == t {
						info, err := addBlockInfo(a, data)
						if err != nil {
							log.Error(err)
						}
						if opts["events"] == "true" {
							// Support events option to emit Server-Sent Events (SSEvent),
							// otherwise, emit JSON responses. SSEvents enable browsers/clients
							// to consume the stream using EventSource.
							g.SSEvent("update", info)
						} else {
							g.JSON(http.StatusOK, info)
							g.Writer.Write([]byte("\n"))
						}
					}
				}
			}
		}
		return true
	})

	listener.Close()
}

func addBlockInfo(a *api, update ThreadUpdate) (ThreadUpdate, error) {
	switch update.Block.Type {
	case "FILES":
		info, err := a.node.File(update.Block.Id)
		if err != nil {
			return update, errors.New("error getting thread file: " + err.Error())
		}
		return ThreadUpdate{
			Block:      update.Block,
			ThreadId:   update.ThreadId,
			ThreadName: update.ThreadName,
			Info:       info,
		}, nil
	default:
		return update, nil
	}
}

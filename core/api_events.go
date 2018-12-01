package core

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *api) getThreadsEvents(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	// Expects or'd list of event types (e.g., FILES|COMMENTS|LIKES).
	types := strings.Split(strings.TrimSpace(opts["type"]), "|")
	threadId := g.Param("id")
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	} // If id wasn't supplied, it will be an empty string

	listener := a.node.GetThreadUpdateListener()
	g.Stream(func(w io.Writer) bool {
		select {
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
		info, err := a.node.ThreadFile(update.Block.Id)
		if err != nil {
			return update, errors.New("error getting thread file: " + err.Error())
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

package core

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/pb"
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
// @Success 200 {object} pb.FeedItem "stream of updates"
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

		case value, ok := <-listener.Ch:
			if !ok {
				return false
			}
			if update, ok := value.(*pb.FeedItem); ok {
				if threadId != "" && update.Thread != threadId {
					break
				}

				btype, err := FeedItemType(update)
				if err != nil {
					log.Error(err)
					break
				}

				for _, t := range types {
					if t == "" || btype.String() == t {

						str, err := pbMarshaler.MarshalToString(update)
						if err != nil {
							g.String(http.StatusBadRequest, err.Error())
							break
						}

						if opts["events"] == "true" {
							g.SSEvent("update", str)
						} else {
							g.Data(http.StatusOK, "application/json", []byte(str))
							g.Writer.Write([]byte("\n"))
						}

						break
					}
				}
			}
		}
		return true
	})

	listener.Close()
}

package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/core"
	pb "github.com/textileio/go-textile/pb"
)

// getThreadsObserve godoc
// @Summary Observe thread updates
// @Description Observes updates in a thread or all threads. An update is generated
// @Description when a new block is added to a thread. There are several update types:
// @Description MERGE, IGNORE, FLAG, JOIN, ANNOUNCE, LEAVE, TEXT, FILES, COMMENT, LIKE
// @Tags observe
// @Produce application/json
// @Param thread path string false "thread id, omit to stream all events"
// @Param X-Textile-Opts header string false "type: Or'd list of event types (e.g., FILES|COMMENTS|LIKES) or empty to include all types, events: Whether to emit Server-Sent Events (SSEvent) or plain JSON" default(type=,events="false")
// @Success 200 {object} pb.FeedItem "stream of updates"
// @Failure 500 {string} string "Internal Server Error"
// @Router /observe/{id} [get]
func (a *Api) getThreadsObserve(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	// Expects or'd list of event types (e.g., FILES|COMMENTS|LIKES).
	types := strings.Split(strings.TrimSpace(strings.ToUpper(opts["type"])), "|")
	threadId := g.Param("id")

	listener := a.Node.ThreadUpdateListener()
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

				btype, err := core.FeedItemType(update)
				if err != nil {
					log.Error(err.Error())
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

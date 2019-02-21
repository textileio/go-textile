package core

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/pb"
)

// lsThreadFeed godoc
// @Summary Paginates post and annotation block types
// @Description Paginates post (join|leave|files|message) and annotation (comment|like) block types
// @Description The mode option dictates how the feed is displayed:
// @Description "chrono": All feed block types are shown. Annotations always nest their target post,
// @Description i.e., the post a comment is about.
// @Description "annotated": Annotations are nested under post targets, but are not shown in the
// @Description top-level feed.
// @Description "stacks": Related blocks are chronologically grouped into "stacks". A new stack is
// @Description started if an unrelated block breaks continuity. This mode is used by Textile
// @Description Photos. Stacks may include:
// @Description * The initial post with some nested annotations. Newer annotations may have already
// @Description been listed.
// @Description * One or more annotations about a post. The newest annotation assumes the "top"
// @Description position in the stack. Additional annotations are nested under the target.
// @Description Newer annotations may have already been listed in the case as well.
// @Tags feed
// @Produce application/json
// @Param X-Textile-Opts header string false "thread: Thread ID (can also use 'default'), offset: Offset ID to start listing from (omit for latest), limit: List page size (default: 5), mode: Feed mode (one of 'chrono', 'annotated', or 'stacks')" default(thread=,offset=,limit=5,mode="chrono")
// @Success 200 {object} pb.FeedItemList "feed"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /feed [get]
func (a *api) lsThreadFeed(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	threadId := opts["thread"]
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	if threadId != "" {
		thrd := a.node.Thread(threadId)
		if thrd == nil {
			g.String(http.StatusNotFound, ErrThreadNotFound.Error())
			return
		}
	}

	limit := 5
	if opts["limit"] != "" {
		limit, err = strconv.Atoi(opts["limit"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	feedMode := pb.FeedMode_value[strings.ToUpper(opts["mode"])]
	list, err := a.node.Feed(opts["offset"], limit, threadId, pb.FeedMode(feedMode))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, list)
}

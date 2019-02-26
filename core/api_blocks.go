package core

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// lsBlocks godoc
// @Summary Paginates blocks in a thread
// @Description Paginates blocks in a thread. Blocks are the raw components in a thread.
// @Description Think of them as an append-only log of thread updates where each update is
// @Description hash-linked to its parent(s). New / recovering peers can sync history by simply
// @Description traversing the hash tree.
// @Tags blocks
// @Produce application/json
// @Param X-Textile-Opts header string false "thread: Thread ID (can also use 'default'), offset: Offset ID to start listing from (omit for latest), limit: List page size (default: 5)" default(thread=,offset=,limit=5)
// @Success 200 {object} pb.BlockList "blocks"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks [get]
func (a *api) lsBlocks(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	threadId := opts["thread"]
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	if threadId == "" {
		g.String(http.StatusBadRequest, "missing thread id")
		return
	}

	thrd := a.node.Thread(threadId)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}

	limit := 5
	if opts["limit"] != "" {
		limit, err = strconv.Atoi(opts["limit"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	query := fmt.Sprintf("threadId='%s'", thrd.Id)
	blocks := a.node.datastore.Blocks().List(opts["offset"], limit, query)
	for _, block := range blocks.Items {
		block.User = a.node.User(block.Author)
	}

	pbJSON(g, http.StatusOK, blocks)
}

// getBlocks godoc
// @Summary Gets thread block
// @Description Gets a thread block by ID
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.Block "block"
// @Failure 404 {string} string "Not Found"
// @Router /blocks/{id} [get]
func (a *api) getBlocks(g *gin.Context) {
	id := g.Param("id")

	block, err := a.node.BlockView(id)
	if err != nil {
		g.String(http.StatusNotFound, "block not found")
		return
	}

	pbJSON(g, http.StatusOK, block)
}

// rmBlocks godoc
// @Summary Remove thread block
// @Description Removes a thread block by ID
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 201 {object} pb.Block "block"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks/{id} [delete]
func (a *api) rmBlocks(g *gin.Context) {
	id := g.Param("id")

	thrd := a.getBlockThread(g, id)
	if thrd == nil {
		return
	}

	hash, err := thrd.AddIgnore(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	block, err := a.node.BlockView(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, block)
}

func (a *api) getBlockThread(g *gin.Context, id string) *Thread {
	block, err := a.node.Block(id)
	if err != nil {
		g.String(http.StatusNotFound, "block not found")
		return nil
	}
	thrd := a.node.Thread(block.Thread)
	if thrd == nil {
		g.String(http.StatusNotFound, "thread not found")
		return nil
	}
	return thrd
}

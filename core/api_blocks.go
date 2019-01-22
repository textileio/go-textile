package core

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	infos := make([]BlockInfo, 0)
	query := fmt.Sprintf("threadId='%s'", thrd.Id)
	for _, block := range a.node.datastore.Blocks().List(opts["offset"], limit, query) {
		username, avatar := a.node.ContactDisplayInfo(block.AuthorId)

		infos = append(infos, BlockInfo{
			Id:       block.Id,
			ThreadId: block.ThreadId,
			AuthorId: block.AuthorId,
			Username: username,
			Avatar:   avatar,
			Type:     block.Type.Description(),
			Date:     block.Date,
			Parents:  block.Parents,
			Target:   block.Target,
			Body:     block.Body,
		})
	}

	g.JSON(http.StatusOK, infos)
}

func (a *api) getBlocks(g *gin.Context) {
	id := g.Param("id")

	info, err := a.node.BlockInfo(id)
	if err != nil {
		g.String(http.StatusNotFound, "block not found")
		return
	}

	g.JSON(http.StatusOK, info)
}

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

	info, err := a.node.BlockInfo(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, info)
}

func (a *api) getBlockThread(g *gin.Context, id string) *Thread {
	block, err := a.node.Block(id)
	if err != nil {
		g.String(http.StatusNotFound, "block not found")
		return nil
	}
	thrd := a.node.Thread(block.ThreadId)
	if thrd == nil {
		g.String(http.StatusNotFound, "thread not found")
		return nil
	}
	return thrd
}

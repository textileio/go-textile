package core

import (
	"net/http"
	"strconv"

	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
)

func (a *api) addThreadFiles(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	threadId := g.Param("id")
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	thrd := a.node.Thread(threadId)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}

	var node ipld.Node
	var keys Keys

	var dir Directory
	if err := g.BindJSON(&dir); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(dir) == 0 {
		g.String(http.StatusBadRequest, "no files found")
		return
	}

	if dir[schema.SingleFileTag].Hash != "" {
		node, keys, err = a.node.AddNodeFromFiles([]repo.File{dir[schema.SingleFileTag]})
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		node, keys, err = a.node.AddNodeFromDirs([]Directory{dir})
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	if node == nil {
		g.String(http.StatusBadRequest, "no files found")
		return
	}

	hash, err := thrd.AddFiles(node, opts["caption"], keys)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	info, err := a.node.BlockInfo(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, info)
}

func (a *api) lsThreadFiles(g *gin.Context) {
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

	limit := 25
	if opts["limit"] != "" {
		limit, err = strconv.Atoi(opts["limit"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	list, err := a.node.ThreadFiles(opts["offset"], limit, threadId)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, list)
}

func (a *api) getThreadFiles(g *gin.Context) {
	info, err := a.node.ThreadFile(g.Param("block"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, info)
}

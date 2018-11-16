package core

import (
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"net/http"
	"strconv"

	"github.com/textileio/textile-go/repo"

	"github.com/gin-gonic/gin"
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

	// parse file or directory
	var dir Directory
	if err := g.BindJSON(&dir); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(dir) > 0 {
		node, keys, err = a.node.AddNodeFromDirs([]Directory{dir})
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		var file repo.File
		if err := g.BindJSON(&file); err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		node, keys, err = a.node.AddNodeFromFiles([]repo.File{file})
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

	threadId := g.Param("id")
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	thrd := a.node.Thread(threadId)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}

	limit := 25
	if opts["limit"] != "" {
		limit, err = strconv.Atoi(opts["limit"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	list, err := a.node.Files(thrd.Id, opts["offset"], limit)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, list)
}

func (a *api) getThreadFiles(g *gin.Context) {
	threadId := g.Param("id")
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	thrd := a.node.Thread(threadId)
	if thrd == nil {
		g.String(http.StatusNotFound, ErrThreadNotFound.Error())
		return
	}

	info, err := a.node.File(threadId, g.Param("block"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, info)
}

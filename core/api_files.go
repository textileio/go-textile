package core

import (
	"net/http"
	"strconv"

	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/ipfs"
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

	var dirs []Directory
	if err := g.BindJSON(&dirs); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(dirs) == 0 {
		g.String(http.StatusBadRequest, "no files found")
		return
	}

	if dirs[0][schema.SingleFileTag].Hash != "" {
		var files []repo.File
		for _, dir := range dirs {
			if len(dir) > 0 && dir[schema.SingleFileTag].Hash != "" {
				files = append(files, dir[schema.SingleFileTag])
			}
		}
		node, keys, err = a.node.AddNodeFromFiles(files)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		node, keys, err = a.node.AddNodeFromDirs(dirs)
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

	limit := 5
	if opts["limit"] != "" {
		limit, err = strconv.Atoi(opts["limit"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	list, err := a.node.Files(opts["offset"], limit, threadId)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, list)
}

func (a *api) getThreadFiles(g *gin.Context) {
	info, err := a.node.FeedFile(g.Param("block"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, info)
}

func (a *api) lsThreadFileTargetKeys(g *gin.Context) {
	target := g.Param("target")

	node, err := ipfs.NodeAtPath(a.node.Ipfs(), target)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	keys, err := a.node.TargetNodeKeys(node)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, keys)
}

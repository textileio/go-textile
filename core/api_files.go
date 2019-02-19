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

// addThreadFiles godoc
// @Summary Adds a file or directory of files to a thread
// @Description Adds a file or directory of files to a thread. Files not supported by the thread
// @Description schema are ignored. Nested directories are included. An existing file hash may
// @Description also be used as input.
// @Tags threads
// @Accept application/json
// @Produce application/json
// @Param dir body core.Directory true "milled dir (output from mill endpoint)"
// @Param X-Textile-Opts header string false "caption: Caption to add to file(s)" default(caption=)
// @Success 201 {object} core.BlockInfo "block"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id}/files [post]
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

// lsThreadFiles godoc
// @Summary Paginates thread files
// @Description Paginates thread files. If thread id not provided, paginate all files. Specify
// @Description "default" to use the default thread (if set).
// @Tags files
// @Produce application/json
// @Param X-Textile-Opts header string false "thread: Thread ID. Omit for all, offset: Offset ID to start listing from. Omit for latest, limit: List page size. (default: 5)" default(thread=,offset=,limit=5)
// @Success 200 {object} pb.FilesList "files"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /files [get]
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

	pbJSON(g, http.StatusOK, list)
}

// getThreadFiles godoc
// @Summary Get thread file
// @Description Gets a thread file by block ID
// @Tags files
// @Produce application/json
// @Param block path string true "block id"
// @Success 200 {object} pb.Files "file"
// @Failure 400 {string} string "Bad Request"
// @Router /files/{block} [get]
func (a *api) getThreadFiles(g *gin.Context) {
	info, err := a.node.File(g.Param("block"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}

// lsThreadFileTargetKeys godoc
// @Summary Show file keys
// @Description Shows file keys under the given target from an add
// @Tags files
// @Produce application/json
// @Param blotargetck path string true "target id"
// @Success 200 {object} core.Keys "keys"
// @Failure 400 {string} string "Bad Request"
// @Router /keys/{target} [get]
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

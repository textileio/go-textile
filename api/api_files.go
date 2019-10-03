package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema"
)

// addThreadFiles godoc
// @Summary Adds a file or directory of files to a thread
// @Description Adds a file or directory of files to a thread. Files not supported by the thread
// @Description schema are ignored. Nested directories are included. An existing file hash may
// @Description also be used as input.
// @Tags threads
// @Accept application/json
// @Produce application/json
// @Param dir body pb.DirectoryList true "list of milled dirs (output from mill endpoint)"
// @Param X-Textile-Opts header string false "caption: Caption to add to file(s)" default(caption=)
// @Success 201 {object} pb.Files "file"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /threads/{id}/files [post]
func (a *Api) addThreadFiles(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	threadId := g.Param("id")
	thrd := a.Node.Thread(threadId)
	if thrd == nil {
		g.String(http.StatusNotFound, core.ErrThreadNotFound.Error())
		return
	}

	var node ipld.Node
	var keys *pb.Keys

	dirs := new(pb.DirectoryList)
	if err := pbUnmarshaler.Unmarshal(g.Request.Body, dirs); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(dirs.Items) == 0 {
		g.String(http.StatusBadRequest, "no files found")
		return
	}

	if dirs.Items[0].Files[schema.SingleFileTag] != nil {
		var files []*pb.FileIndex
		for _, dir := range dirs.Items {
			if len(dir.Files) > 0 && dir.Files[schema.SingleFileTag].Hash != "" {
				files = append(files, dir.Files[schema.SingleFileTag])
			}
		}
		node, keys, err = a.Node.AddNodeFromFiles(files)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		node, keys, err = a.Node.AddNodeFromDirs(dirs)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	if node == nil {
		g.String(http.StatusBadRequest, "no files found")
		return
	}

	// @todo Allow the setting of the target in 0.5.0
	hash, err := thrd.AddFiles(node, "", opts["caption"], keys.Files)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	files, err := a.Node.File(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	pbJSON(g, http.StatusCreated, files)
}

// lsThreadFiles godoc
// @Summary Paginates thread files
// @Description Paginates thread files. If thread id not provided, paginate all files.
// @Tags files
// @Produce application/json
// @Param X-Textile-Opts header string false "thread: Thread ID. Omit for all, offset: Offset ID to start listing from. Omit for latest, limit: List page size. (default: 5)" default(thread=,offset=,limit=5)
// @Success 200 {object} pb.FilesList "files"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /files [get]
func (a *Api) lsThreadFiles(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	threadId := opts["thread"]
	if threadId != "" {
		thrd := a.Node.Thread(threadId)
		if thrd == nil {
			g.String(http.StatusNotFound, core.ErrThreadNotFound.Error())
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

	list, err := a.Node.Files(opts["offset"], limit, threadId)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, list)
}

// lsThreadFileTargetKeys godoc
// @Summary Show file keys
// @Description Shows file keys under the given target from an add
// @Tags files
// @Produce application/json
// @Param target path string true "target id"
// @Success 200 {object} pb.Keys "keys"
// @Failure 400 {string} string "Bad Request"
// @Router /keys/{target} [get]
func (a *Api) lsThreadFileTargetKeys(g *gin.Context) {
	target := g.Param("target")

	node, err := ipfs.NodeAtPath(a.Node.Ipfs(), target, ipfs.CatTimeout)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	keys, err := a.Node.TargetNodeKeys(node)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, keys)
}

// getFileMeta godoc
// @Summary File metadata at hash
// @Description Returns the metadata for file
// @Tags files
// @Produce application/json
// @Param hash path string true "file hash"
// @Success 200 {object} pb.FileIndex "file"
// @Failure 404 {string} string "Not Found"
// @Router /file/{target}/meta [get]
func (a *Api) getFileMeta(g *gin.Context) {
	file, err := a.Node.FileMeta(g.Param("hash"))
	if err != nil {
		g.String(http.StatusNotFound, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, file)
}

// getFileContent godoc
// @Summary File content at hash
// @Description Returns decrypted raw content for file
// @Tags files
// @Produce application/octet-stream
// @Param hash path string true "file hash"
// @Success 200 {string} byte
// @Failure 404 {string} string "Not Found"
// @Router /file/{hash}/content [get]
func (a *Api) getFileContent(g *gin.Context) {
	reader, file, err := a.Node.FileContent(g.Param("hash"))
	if err != nil {
		g.String(http.StatusNotFound, err.Error())
		return
	}

	g.DataFromReader(http.StatusOK, file.Size, file.Media, reader, map[string]string{})
}

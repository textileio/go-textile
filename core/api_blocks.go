package core

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
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
		block.User = a.node.PeerUser(block.Author)
	}

	var dots bool
	if opts["dots"] != "" {
		dots, err = strconv.ParseBool(opts["dots"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	if !dots {
		pbJSON(g, http.StatusOK, blocks)
		return
	}

	var nextOffset string
	if len(blocks.Items) > 0 {
		nextOffset = blocks.Items[len(blocks.Items)-1].Id

		// see if there's actually more
		if len(a.node.datastore.Blocks().List(nextOffset, 1, query).Items) == 0 {
			nextOffset = ""
		}
	}

	dotsf, err := a.toDots(blocks)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	viz := &pb.BlockViz{
		Dots:  dotsf,
		Count: int32(len(blocks.Items)),
		Next:  nextOffset,
	}

	pbJSON(g, http.StatusOK, viz)
}

// getBlocks godoc
// @Summary Backwards compatible redirect to /blocks/{id}/meta
// @Router /blocks/{id} [get]
func (a *api) getBlocks(g *gin.Context) {
	id := g.Param("id")
	g.Redirect(http.StatusPermanentRedirect, "/blocks/" + id + "/meta")
}

// getBLockMeta godoc
// @Summary Gets thread block
// @Description Gets a thread block by ID
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.Block "block"
// @Failure 404 {string} string "Not Found"
// @Router /blocks/{id}/meta [get]
func (a *api) getBlockMeta(g *gin.Context) {
	id := g.Param("id")

	block, err := a.node.BlockView(id)
	if err != nil {
		g.String(http.StatusNotFound, "block not found")
		return
	}

	pbJSON(g, http.StatusOK, block)
}

// getBlockFiles godoc
// @Summary Get thread file
// @Description Gets a thread file by block ID
// @Tags files
// @Produce application/json
// @Param block path string true "block id"
// @Success 200 {object} pb.Files "file"
// @Failure 400 {string} string "Bad Request"
// @Router /files/{block} [get]
func (a *api) getBlockFiles(g *gin.Context) {
	files, err := a.node.File(g.Param("id"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, files)
}

// Helper method
func (a *api) getBlockFile(id string, indexStr string, path string) (*pb.FileIndex, error) {
	var file *pb.FileIndex

	files, err := a.node.File(id)
	if err != nil {
		return file, err
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return file, err
	}

	file = files.Files[index].Links[path]
	return file, nil
}

// getBlockFileMeta godoc
// @todo
func (a *api) getBlockFileMeta(g *gin.Context) {
	id := g.Param("id")
	index := g.Param("index")
	path := g.Param("path")
	file, err := a.getBlockFile(id, index, path)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	pbJSON(g, http.StatusOK, file)
}


// getBlockFileContent godoc
// @todo
func (a *api) getBlockFileContent(g *gin.Context) {
	id := g.Param("id")
	index := g.Param("index")
	path := g.Param("path")
	file, err := a.getBlockFile(id, index, path)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	reader, err := a.node.FileIndexContent(file)
	if err != nil {
		g.String(http.StatusNotFound, err.Error())
		return
	}
	g.DataFromReader(http.StatusOK, file.Size, file.Media, reader, map[string]string{})
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

func (a *api) toDots(blocks *pb.BlockList) (string, error) {
	dots := `digraph {
    rankdir="BT";`

	for _, b := range blocks.Items {
		dot := toDot(b)

		for _, p := range b.Parents {
			if strings.TrimSpace(p) == "" {
				continue
			}
			pp, err := a.node.Block(p)
			if err != nil {
				log.Warningf("block %s: %s", p, err)
				dots += "\n    " + dot + " -> MISSING_" + pre(p) + ";"
				continue
			}
			dots += "\n    " + dot + " -> " + toDot(pp) + ";"
		}
	}

	return dots + "\n}", nil
}

func toDot(block *pb.Block) string {
	t := block.Type.String()
	var a string
	if block.Type != pb.Block_MERGE {
		a = "_" + ipfs.ShortenID(block.Author)
	}
	return t + a + "_" + pre(block.Id)
}

func pre(hash string) string {
	if len(hash) < 7 {
		return hash
	}
	return hash[:7]
}

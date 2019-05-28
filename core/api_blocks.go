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

func getBlock(node *Textile, id string) (*pb.Block, error, int) {
	block, err := node.BlockView(id)
	if err != nil {
		return nil, fmt.Errorf("block not found [id=%s]", id), http.StatusNotFound
	}
	return block, nil, http.StatusOK
}

func getThread(node *Textile, id string) (*Thread, error, int) {
	thread := node.Thread(id)
	if thread == nil {
		return nil, fmt.Errorf("thread not found [id=%s]", id), http.StatusNotFound
	}
	return thread, nil, http.StatusOK
}

func getBlockThread(node *Textile, id string) (*Thread, error, int) {
	block, err, code := getBlock(node, id)
	if err != nil {
		return nil, err, code
	}
	return getThread(node, block.Thread)
}

func getFiles(node *Textile, id string) (*pb.Files, error, int) {
	files, err := node.File(id) // despite naming, this is files
	if err != nil {
		return nil, err, http.StatusNotFound
	}
	return files, nil, http.StatusOK
}

func getFile(files *pb.Files, indexStr string, path string) (*pb.FileIndex, error, int) {
	// file
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}
	file, ok := files.Files[index].Links[path]
	if !ok {
		return nil, err, http.StatusNotFound
	}
	return file, nil, http.StatusOK
}

func getFilesFile(node *Textile, id string, indexStr string, path string) (*pb.FileIndex, error, int) {
	files, err, code := getFiles(node, id)
	if err != nil {
		return nil, err, code
	}
	return getFile(files, indexStr, path)
}

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

	thread := a.node.Thread(threadId)
	if thread == nil {
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

	query := fmt.Sprintf("threadId='%s'", thread.Id)
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

// getBlockMeta godoc
// @Summary Gets the metadata for a block
// @Tags blocks
// @Produce application/json
// @Router /blocks/{id}/meta [get]
// @Param id path string true "block id"
// @Success 200 {object} pb.Block "block"
// @Failure 404 {string} string "Not Found"
func (a *api) getBlockMeta(g *gin.Context) {
	block, err, code := getBlock(a.node, g.Param("id"))
	if err != nil {
		sendError(g, err, code)
		return
	}
	pbJSON(g, http.StatusOK, block)
}

// getBlockFiles godoc
// @Summary Gets the metadata for a files block
// @Tags files
// @Produce application/json
// @Router /blocks/{id}/files [get]
// @Param id path string true "block id"
// @Success 200 {object} pb.Files "files"
// @Failure 404 {string} string "Not Found"
func (a *api) getBlockFiles(g *gin.Context) {
	files, err, code := getFiles(a.node, g.Param("id"))
	if err != nil {
		sendError(g, err, code)
		return
	}
	pbJSON(g, http.StatusOK, files)
}

// getBlockFileMeta godoc
// @Summary Gets the metadata of a file within a files block
// @Tags files
// @Produce application/json
// @Router /blocks/{id}/files/{index}/{path}/meta [get]
// @Param id path string true "block id"
// @Param index path string true "file index"
// @Param path path string true "file path"
// @Success 200 {object} pb.FileIndex "file"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
func (a *api) getBlockFileMeta(g *gin.Context) {
	file, err, code := getFilesFile(a.node, g.Param("id"), g.Param("index"), g.Param("path"))
	if err != nil {
		sendError(g, err, code)
		return
	}
	pbJSON(g, http.StatusOK, file)
}

// getBlockFileContent godoc
// @Summary Gets the decrypted file content of a file within a files block
// @Tags files
// @Produce application/json
// @Router /blocks/{id}/files/{index}/{path}/content [get]
// @Param id path string true "block id"
// @Param index path string true "file index"
// @Param path path string true "file path"
// @Success 200 {string} byte
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
func (a *api) getBlockFileContent(g *gin.Context) {
	file, err, code := getFilesFile(a.node, g.Param("id"), g.Param("index"), g.Param("path"))
	if err != nil {
		sendError(g, err, code)
		return
	}
	reader, err := a.node.FileIndexContent(file)
	if err != nil {
		sendError(g, err, http.StatusNotFound)
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
	blockID := g.Param("id")

	thread, err, code := getBlockThread(a.node, blockID)
	if err != nil {
		sendError(g, err, code)
		return
	}

	hash, err := thread.AddIgnore(blockID)
	if err != nil {
		a.abort500(g, err)
		return
	}

	block, err, code := getBlock(a.node, hash.B58String())
	if err != nil {
		sendError(g, err, code)
		return
	}

	go a.node.cafeOutbox.Flush()

	pbJSON(g, http.StatusCreated, block)
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

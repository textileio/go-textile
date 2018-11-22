package core

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/schema"
)

func (a *api) schemaMill(g *gin.Context) {
	var node schema.Node
	if err := g.BindJSON(&node); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	mill := &m.Schema{}

	data, err := json.Marshal(&node)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	conf := AddFileConfig{
		Input: data,
		Media: "application/json",
	}

	added, err := a.node.AddFile(mill, conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

func (a *api) blobMill(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	mill := &m.Blob{}

	plaintext := opts["plaintext"] == "true"

	conf, err := a.getFileConfig(g, mill, opts["use"], plaintext)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	added, err := a.node.AddFile(mill, *conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

func (a *api) imageResizeMill(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	mill := &m.ImageResize{
		Opts: m.ImageResizeOpts{
			Quality: "75",
		},
	}

	// width is required
	if opts["width"] == "" {
		g.String(http.StatusBadRequest, "missing width")
		return
	}
	mill.Opts.Width = opts["width"]

	// quality defaults to 75
	if opts["quality"] != "" {
		mill.Opts.Quality = opts["quality"]
	}

	plaintext := opts["plaintext"] == "true"

	conf, err := a.getFileConfig(g, mill, opts["use"], plaintext)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	added, err := a.node.AddFile(mill, *conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

func (a *api) imageExifMill(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	mill := &m.ImageExif{}

	plaintext := opts["plaintext"] == "true"

	conf, err := a.getFileConfig(g, mill, opts["use"], plaintext)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	conf.Media = "application/json"

	added, err := a.node.AddFile(mill, *conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	m "github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/schema"
	"io/ioutil"
	"net/http"
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

	added, err := a.node.AddFile(data, "", "application/json", mill)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

func (a *api) blobMill(g *gin.Context) {
	mill := &m.Blob{}

	file, name, err := a.openFile(g)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	media, err := a.node.FileMedia(file, mill)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	file.Seek(0, 0)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	added, err := a.node.AddFile(data, name, media, mill)
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

	file, name, err := a.openFile(g)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	media, err := a.node.FileMedia(file, mill)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	file.Seek(0, 0)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	added, err := a.node.AddFile(data, name, media, mill)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

func (a *api) imageExifMill(g *gin.Context) {
	mill := &m.ImageExif{}

	file, name, err := a.openFile(g)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	if _, err := a.node.FileMedia(file, mill); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	file.Seek(0, 0)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	added, err := a.node.AddFile(data, name, "application/json", mill)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

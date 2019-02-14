package core

import (
	"io/ioutil"
	"net/http"

	"github.com/textileio/textile-go/repo"

	"github.com/gin-gonic/gin"
	m "github.com/textileio/textile-go/mill"
)

func (a *api) schemaMill(g *gin.Context) {
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer g.Request.Body.Close()

	mill := &m.Schema{}

	conf := AddFileConfig{
		Input: body,
		Media: "application/json",
	}

	added, err := a.node.AddFileIndex(mill, conf)
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

	added, err := a.node.AddFileIndex(mill, *conf)
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

	added, err := a.node.AddFileIndex(mill, *conf)
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

	added, err := a.node.AddFileIndex(mill, *conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

func (a *api) jsonMill(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	mill := &m.Json{}

	conf := AddFileConfig{
		Media:     "application/json",
		Plaintext: opts["plaintext"] == "true",
	}

	if opts["use"] == "" {
		body, err := ioutil.ReadAll(g.Request.Body)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		defer g.Request.Body.Close()

		if body == nil {
			g.String(http.StatusBadRequest, "missing doc")
			return
		}
		conf.Input = body

	} else {
		var file *repo.File
		reader, file, err := a.node.FileData(opts["use"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		conf.Use = file.Checksum

		conf.Input, err = ioutil.ReadAll(reader)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	added, err := a.node.AddFileIndex(mill, conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

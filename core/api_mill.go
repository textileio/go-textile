package core

import (
	"github.com/gin-gonic/gin"
	m "github.com/textileio/textile-go/mill"
	"image/jpeg"
	"net/http"
	"strconv"
)

func (a *api) blobMill(g *gin.Context) {
	mill := &m.Blob{}

	file, name, err := a.openFile(g)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	added, err := a.node.AddFile(file, name, mill)
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
			Quality: jpeg.DefaultQuality,
		},
	}

	// width is required
	if opts["width"] == "" {
		g.String(http.StatusBadRequest, "missing width")
		return
	}
	mill.Opts.Width, err = strconv.Atoi(opts["width"])
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// quality defaults to 75
	if opts["quality"] != "" {
		mill.Opts.Quality, err = strconv.Atoi(opts["quality"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	file, name, err := a.openFile(g)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	added, err := a.node.AddFile(file, name, mill)
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

	added, err := a.node.AddFile(file, name, mill)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusCreated, added)
}

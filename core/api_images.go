package core

import (
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/images"
	"github.com/textileio/textile-go/repo"
	"image/jpeg"
	"net/http"
	"strconv"
)

func (a *api) encodeImages(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if opts["width"] == "" {
		g.String(http.StatusBadRequest, "missing width")
		return
	}
	width, err := strconv.Atoi(opts["width"])
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	quality := jpeg.DefaultQuality
	if opts["quality"] != "" {
		quality, err = strconv.Atoi(opts["quality"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	// handle file
	form, err := g.MultipartForm()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	fileHeaders := form.File["file"]
	var adds []*repo.File
	for _, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}

		// decode image
		reader, format, err := images.DecodeImage(file)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
		}

		// encode with opts
		buff, err := images.EncodeImage(reader, format, width, quality)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
		}

		res, err := a.node.AddFile(file)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		file.Close()
		adds = append(adds, res)
	}
	g.JSON(http.StatusCreated, adds)
}

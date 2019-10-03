package api

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/core"
	m "github.com/textileio/go-textile/mill"
)

// schemaMill godoc
// @Summary Validate, add, and pin a new Schema
// @Description Takes a JSON-based Schema, validates it, adds it to IPFS, and returns a file object
// @Tags mills
// @Accept application/json
// @Produce application/json
// @Param schema body pb.Node true "schema"
// @Success 201 {object} pb.FileIndex "file"
// @Failure 400 {string} string "Bad Request"
// @Router /mills/schema [post]
func (a *Api) schemaMill(g *gin.Context) {
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer g.Request.Body.Close()

	mill := &m.Schema{}

	conf := core.AddFileConfig{
		Input: body,
		Media: "application/json",
	}

	added, err := a.Node.AddFileIndex(mill, conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, added)
}

// blobMill godoc
// @Summary Process raw data blobs
// @Description Takes a binary data blob, and optionally encrypts it, before adding to IPFS,
// @Description and returns a file object
// @Tags mills
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file false "multipart/form-data file"
// @Param X-Textile-Opts header string false "plaintext: whether to leave unencrypted), use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS" default(plaintext=false,use="")
// @Success 201 {object} pb.FileIndex "file"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /mills/blob [post]
func (a *Api) blobMill(g *gin.Context) {
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

	added, err := a.Node.AddFileIndex(mill, *conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, added)
}

// imageResizeMill godoc
// @Summary Resize an image
// @Description Takes an input image, and resizes/resamples it (optionally encrypting output),
// @Description before adding to IPFS, and returns a file object
// @Tags mills
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file false "multipart/form-data file"
// @Param X-Textile-Opts header string true "plaintext: whether to leave unencrypted, use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS, width: the requested image width (required), quality: the requested JPEG image quality" default(plaintext=false,use="",quality=75,width=100)
// @Success 201 {object} pb.FileIndex "file"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /mills/image/resize [post]
func (a *Api) imageResizeMill(g *gin.Context) {
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

	added, err := a.Node.AddFileIndex(mill, *conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, added)
}

// imageExifMill godoc
// @Summary Extract EXIF data from image
// @Description Takes an input image, and extracts its EXIF data (optionally encrypting output),
// @Description before adding to IPFS, and returns a file object
// @Tags mills
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file false "multipart/form-data file"
// @Param X-Textile-Opts header string false "plaintext: whether to leave unencrypted, use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS" default(plaintext=false,use="")
// @Success 201 {object} pb.FileIndex "file"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /mills/image/exif [post]
func (a *Api) imageExifMill(g *gin.Context) {
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

	added, err := a.Node.AddFileIndex(mill, *conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, added)
}

// jsonMill godoc
// @Summary Process input JSON data
// @Description Takes an input JSON document, validates it according to its json-schema.org definition,
// @Description optionally encrypts the output before adding to IPFS, and returns a file object
// @Tags mills
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file false "multipart/form-data file"
// @Param X-Textile-Opts header string false "plaintext: whether to leave unencrypted, use: if empty, assumes body contains multipart form file data, otherwise, will attempt to fetch given CID from IPFS" default(plaintext="false",use="")
// @Success 201 {object} pb.FileIndex "file"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /mills/json [post]
func (a *Api) jsonMill(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	mill := &m.Json{}

	conf := core.AddFileConfig{
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
		reader, file, err := a.Node.FileContent(opts["use"])
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

	added, err := a.Node.AddFileIndex(mill, conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, added)
}

package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *api) addImages(g *gin.Context) {
	form, err := g.MultipartForm()
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	fileHeaders := form.File["file"]
	var adds []*AddDataResult
	for _, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		added, err := a.node.AddImage(file, header.Filename)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		file.Close()
		adds = append(adds, added)
	}
	g.JSON(http.StatusCreated, adds)
}

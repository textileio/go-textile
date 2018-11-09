package core

import (
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/repo"
	"net/http"
)

func (a *api) addImages(g *gin.Context) {
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
		original, err := a.node.AddFile(file)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}

		file.Close()
		adds = append(adds, original)
	}
	g.JSON(http.StatusCreated, adds)
}

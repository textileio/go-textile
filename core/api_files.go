package core

import (
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/repo"
	"net/http"
)

func (a *api) addFiles(g *gin.Context) {
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

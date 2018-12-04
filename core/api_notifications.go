package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) lsNotifications(g *gin.Context) {
	list := a.node.Notifications("", -1)
	g.JSON(http.StatusOK, list)
}

func (a *api) readNotifications(g *gin.Context) {
	id := g.Param("id")

	var err error
	if id == "all" {
		err = a.node.ReadAllNotifications()
	} else {
		err = a.node.ReadNotification(id)
	}
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, "ok")
}

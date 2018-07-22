package cafe

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Cafe) pin(g *gin.Context) {

	// ship it
	g.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
	})
}

package cafe

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Cafe) health(g *gin.Context) {
	g.Writer.WriteHeader(http.StatusNoContent)
}

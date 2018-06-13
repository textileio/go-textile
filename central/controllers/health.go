package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNoContent)
}

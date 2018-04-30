package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"verion": os.Getenv("VERSION"),
	})
}

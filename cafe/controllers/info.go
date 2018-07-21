package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": os.Getenv("VERSION"),
	})
}

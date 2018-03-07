package main

import (
	"time"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	. "gitlab.com/textileio/mill-go/api/config"
	. "gitlab.com/textileio/mill-go/api/dao"
	. "gitlab.com/textileio/mill-go/api/models"
)

var config = Config{}
var dao = DAO{}

func registerUser(c *gin.Context) {
	var reg Registration
	if err := c.ShouldBindJSON(&reg); err == nil {
		now := time.Now()
		user := User{
			ID: bson.NewObjectId(),
			Created: now,
			LastSeen: now,
			Identities: []Identity{reg.Identity}}
		if err := dao.Insert(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "resourceId": user.ID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// Parse the config and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", registerUser)
	}
	router.Run(":8000")
}

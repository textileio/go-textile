package main

import (
	"os"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/controllers"
	"github.com/textileio/textile-go/central/auth"
)

// Establish a connection to DB
func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// initialize a dao
	dao.Dao = &dao.DAO{
		Hostname: os.Getenv("HOSTNAME"),
		DatabaseName: os.Getenv("DATABASENAME"),
	}
	dao.Dao.Connect()
}

// Define HTTP request routes
func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.Use(auth.Auth(os.Getenv("TOKEN_SECRET")))
	{
		v1.PUT("/users", controllers.SignUp)
		v1.POST("/users", controllers.SignIn)
	}
	router.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

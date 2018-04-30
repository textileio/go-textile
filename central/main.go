package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/textileio/textile-go/central/controllers"
	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/middleware"
)

// Establish a connection to DB
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// parse db tls setting
	var tls bool
	tls, err = strconv.ParseBool(os.Getenv("DB_TLS"))
	if err != nil {
		tls = false
	}

	// initialize a dao
	dao.Dao = &dao.DAO{
		Hosts:    os.Getenv("DB_HOSTS"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		TLS:      tls,
	}
	dao.Dao.Connect()
	dao.Dao.Index()
}

// Define HTTP request routes
func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.Use(middleware.Auth(os.Getenv("TOKEN_SECRET")))
	{
		v1.PUT("/users", controllers.SignUp)
		v1.POST("/users", controllers.SignIn)
	}
	router.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

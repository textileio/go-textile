package main

import (
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/central/controllers"
	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/middleware"
	"os"
)

func init() {
	// establish a connection to DB
	dao.Dao = &dao.DAO{
		Hosts:    os.Getenv("DB_HOSTS"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		TLS:      os.Getenv("DB_TLS") == "yes",
	}
	dao.Dao.Connect()

	// ensure we're indexed
	dao.Dao.Index()
}

func main() {
	// build http router
	router := gin.Default()
	router.GET("/", controllers.Info)
	router.GET("/health", controllers.Health)

	// api routes
	v1 := router.Group("/api/v1")
	v1.Use(middleware.Auth(os.Getenv("TOKEN_SECRET")))
	{
		v1.PUT("/users", controllers.SignUp)
		v1.POST("/users", controllers.SignIn)
		v1.POST("/referrals", controllers.CreateReferral)
		v1.GET("/referrals", controllers.ListReferrals)
	}
	router.Run(os.Getenv("BIND"))
}

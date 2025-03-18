package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/controllers"
	"github.com/matheusabido/kfofo-api/middleware"
	"github.com/matheusabido/kfofo-api/validator"
)

func SetupRoutes() {
	router := gin.Default()
	validator.SetupValidator()

	router.Use(createCors())
	protected := router.Group("/").Use(middleware.AuthMiddleware)

	protected.GET("/user/:id", controllers.GetUser)
	router.POST("/user", controllers.PostUser)

	router.POST("/session", controllers.PostLogin)

	router.Run()
}

func createCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	})
}

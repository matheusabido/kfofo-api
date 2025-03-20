package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/controllers"
	"github.com/matheusabido/kfofo-api/middleware"
	"github.com/matheusabido/kfofo-api/validator"
)

func SetupRoutes() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.MaxMultipartMemory = 2 * 1024 * 1024
	validator.SetupValidator()

	router.Use(createCors())
	protected := router.Group("/").Use(middleware.AuthMiddleware)

	protected.GET("/user/:id", controllers.GetUser)
	router.POST("/user", controllers.PostUser)
	protected.PUT("/user/:id", controllers.PutUser)
	protected.DELETE("/user/:id", controllers.DeleteUser)

	router.GET("/homes", controllers.GetHomes)
	router.GET("/home/:id", controllers.GetHome)
	protected.POST("/home", controllers.PostHome)
	protected.PUT("/home/:id", controllers.PutHome)
	protected.DELETE("/home/:id", controllers.DeleteHome)

	router.POST("/session", controllers.PostLogin)

	router.GET("/home/picture", controllers.GetHomePicture)
	protected.POST("/home/picture", controllers.PostHomePicture)

	router.GET("/utensils", controllers.GetUtensils)
	protected.PUT("/utensils", controllers.UpdateUtensils)

	protected.GET("/bookings", controllers.GetBookings)
	protected.POST("/booking", controllers.PostBooking)
	protected.DELETE("/booking/:id", controllers.DeleteBooking)

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

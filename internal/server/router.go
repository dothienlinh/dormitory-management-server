package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{

			auth := v1.Group("/auth")
			{
				auth.POST("/login", handlers.Login)
				auth.POST("/register", handlers.Register)
				auth.POST("/refresh-token", handlers.RefreshToken)
			}

			protected := v1.Group("/")
			protected.Use(middleware.AuthMiddleware)
			{
				auth := protected.Group("/auth")
				{
					auth.GET("/me", handlers.Me)
					auth.POST("/logout", handlers.Logout)
				}
			}
		}
	}

	log.Println("Server is running on port 8080")

	router.Run()
}

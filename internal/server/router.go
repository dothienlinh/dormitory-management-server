package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(middleware *middleware.Middleware, handler *handlers.Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{

			auth := v1.Group("/auth")
			{
				auth.POST("/login", handler.Login())
				auth.POST("/register", handler.Register())
				auth.POST("/refresh-token", handler.RefreshToken())
			}

			protected := v1.Group("/")
			protected.Use(middleware.AuthMiddleware())
			{
				auth := protected.Group("/auth")
				{
					auth.GET("/me", handler.Me())
					auth.POST("/logout", handler.Logout())
				}

				// API for admin
				admin := protected.Group("/admin")
				admin.Use(middleware.RoleAdminMiddleware())
				{
					rooms := admin.Group("/rooms")
					{
						rooms.POST("/", handlers.CreateRoom)
						rooms.PUT("/:id")
						rooms.DELETE("/:id")
					}
				}

			}
		}
	}

	return router
}

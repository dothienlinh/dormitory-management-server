package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAmenityRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	amenities := router.Group("/amenities")
	{
		amenities.GET("", handlers.Amenity.List())

		amenitiesAdmin := amenities.Group("", mw.AuthMiddleware(), mw.ManagerMiddleware())
		{
			amenitiesAdmin.POST("", handlers.Amenity.Create())
			amenitiesAdmin.PUT("/:id", handlers.Amenity.Update())
			amenitiesAdmin.DELETE("/:id", handlers.Amenity.Delete())
		}
	}
}

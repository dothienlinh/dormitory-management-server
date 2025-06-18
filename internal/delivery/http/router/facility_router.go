package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupFacilityRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	facilities := router.Group("/facilities")
	{
		facilities.GET("", handlers.Facility.List())

		facilitiesAdmin := facilities.Group("", mw.AuthMiddleware(), mw.AdminMiddleware())
		{
			facilitiesAdmin.POST("", handlers.Facility.Create())
			facilitiesAdmin.PUT("/:id", handlers.Facility.Update())
			facilitiesAdmin.DELETE("/:id", handlers.Facility.Delete())
		}
	}
}

package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupMaintenanceHistoryRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	maintenanceHistory := router.Group("/maintenance-histories")

	maintenanceHistoryAdmin := maintenanceHistory.Group("", mw.AuthMiddleware(), mw.AdminMiddleware())
	{
		maintenanceHistoryAdmin.POST("", handlers.MaintenanceHistory.Create())
		maintenanceHistoryAdmin.GET("/:id", handlers.MaintenanceHistory.Detail())
		maintenanceHistoryAdmin.PUT("/:id", handlers.MaintenanceHistory.Update())
		maintenanceHistoryAdmin.DELETE("/:id", handlers.MaintenanceHistory.Delete())
	}
}

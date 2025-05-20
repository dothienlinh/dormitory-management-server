package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDashboardRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	dashboard := router.Group("/dashboard", mw.AuthMiddleware(), mw.AdminMiddleware())
	{
		dashboard.GET("/stats", handlers.Dashboard.GetStats())
	}
}

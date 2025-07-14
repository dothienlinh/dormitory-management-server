package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDashboardRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	// Admin dashboard routes
	adminDashboard := router.Group("/dashboard", mw.AuthMiddleware(), mw.AdminMiddleware())
	{
		adminDashboard.GET("/stats", handlers.Dashboard.GetStats())
	}

	// Student dashboard routes
	studentDashboard := router.Group("/dashboard", mw.AuthMiddleware())
	{
		studentDashboard.GET("/student/overview", handlers.Dashboard.StudentOverview())
		studentDashboard.GET("/notifications", handlers.Dashboard.GetNotifications())
		studentDashboard.PATCH("/notifications/:id/read", handlers.Dashboard.MarkNotificationAsRead())
		studentDashboard.PATCH("/notifications/mark-all-read", handlers.Dashboard.MarkAllNotificationsAsRead())
		studentDashboard.GET("/events", handlers.Dashboard.GetEvents())
		studentDashboard.GET("/service-requests", handlers.Dashboard.GetServiceRequests())
		studentDashboard.GET("/quick-stats", handlers.Dashboard.GetQuickStats())
	}
}

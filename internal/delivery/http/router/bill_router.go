package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupBillRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	bill := router.Group("/bills", mw.AuthMiddleware())
	{
		bill.GET("", handlers.Bill.MyListBills())
	}
}

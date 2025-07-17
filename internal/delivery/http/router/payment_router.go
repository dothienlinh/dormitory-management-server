package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	payments := router.Group("/payments", mw.AuthMiddleware())

	{
		payments.POST("/vietqr", handlers.PaymentHandler.CreateLinkPaymentVietQR())
	}

	router.POST("/payments/vietqr/receive-hook", handlers.PaymentHandler.ReceiveHookVietQR())
}

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

		// Student payment routes for consistency with frontend
		payments.GET("/my-payments", handlers.PaymentHistory.GetMyPaymentHistory())
		payments.POST("/:id/pay", handlers.PaymentHistory.MakePayment())
		payments.GET("/:id/receipt", handlers.PaymentHistory.DownloadReceipt())
	}

	router.POST("/payments/vietqr/receive-hook", handlers.PaymentHandler.ReceiveHookVietQR())
}

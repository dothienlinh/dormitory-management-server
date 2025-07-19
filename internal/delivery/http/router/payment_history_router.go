package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPaymentHistoryRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	paymentHistory := router.Group("/payment-history", mw.AuthMiddleware())
	{
		paymentHistory.GET("", handlers.PaymentHistory.GetPaymentHistoryList())
		paymentHistory.GET("/:id", handlers.PaymentHistory.GetPaymentHistoryByID())
		paymentHistory.GET("/contract/:contract_id", handlers.PaymentHistory.GetPaymentHistoryByContractID())

		// Student payment history routes
		paymentHistory.GET("/my-payments", handlers.PaymentHistory.GetMyPaymentHistory())
		paymentHistory.POST("/:id/pay", handlers.PaymentHistory.MakePayment())
		paymentHistory.GET("/:id/receipt", handlers.PaymentHistory.DownloadReceipt())

		paymentHistoryAdmin := paymentHistory.Group("", mw.ManagerMiddleware())
		{
			paymentHistoryAdmin.POST("", handlers.PaymentHistory.CreatePaymentHistory())
			paymentHistoryAdmin.PUT("/:id", handlers.PaymentHistory.UpdatePaymentHistory())
			paymentHistoryAdmin.DELETE("/:id", handlers.PaymentHistory.DeletePaymentHistory())
		}
	}
}

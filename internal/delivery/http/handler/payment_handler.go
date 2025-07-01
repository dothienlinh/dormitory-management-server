package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewPaymentHandler(useCases usecase.UseCases, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *PaymentHandler) CreateLinkPaymentVietQR() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("CreateLinkPaymentVietQR")

		var payload entity.CreateLinkPaymentVietQR
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Payment().CreateLinkPaymentVietQR(ctx, &payload)
		ctx.JSON(resp.Status, resp.Response)
	}
}

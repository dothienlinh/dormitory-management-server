package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
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

func (h *PaymentHandler) ReceiveHookVietQR() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		webhookDataReq := payos.WebhookType{}

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			h.logger.Error("Failed to read body", zap.Error(err))
			ctx.Error(err)
			return
		}

		if err := json.Unmarshal(body, &webhookDataReq); err != nil {
			h.logger.Error("Failed to unmarshal body", zap.Error(err))
			ctx.Error(err)
			return
		}
		webhookData, err := payos.VerifyPaymentWebhookData(webhookDataReq)
		if err != nil {
			h.logger.Error("Failed to verify payment webhook data", zap.Error(err))
			ctx.Error(err)
			return
		}

		if err := h.useCases.Payment().ReceiveHookVietQR(ctx, webhookData); err != nil {
			h.logger.Error("Failed to receive hook vietqr", zap.Error(err))
			ctx.Error(err)
			return
		}
	}
}

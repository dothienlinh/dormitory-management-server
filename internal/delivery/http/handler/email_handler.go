package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmailHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewEmailHandler(useCases usecase.UseCases, logger logger.Logger) *EmailHandler {
	return &EmailHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *EmailHandler) SendEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		var payload entity.SendCodeEmail
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		_, err := h.useCases.Email().SendOTP(ctx, payload)
		if err != nil {
			h.logger.Error("Failed to send email", zap.Error(err))
			resp = response.InternalServerError("Failed to send email")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = response.Success(nil, 0)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *EmailHandler) VerifyCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		var payload entity.VerifyCodeEmail
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Email().VerifyCodeEmail(ctx, payload)

		ctx.JSON(resp.Status, resp.Response)
	}
}

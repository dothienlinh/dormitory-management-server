package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"net/http"

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
		var payload entity.SendCodeEmail
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		_, err := h.useCases.Email().SendOTP(ctx, payload)
		if err != nil {
			h.logger.Error("Failed to send email", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to send email",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Email sent successfully",
		})
	}
}

func (h *EmailHandler) VerifyCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload entity.VerifyCodeEmail
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Email().VerifyCodeEmail(ctx, payload)

		ctx.JSON(resp.Status, resp.Response)
	}
}

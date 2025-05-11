package handlers

import (
	"dormitory_management/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) CreateContract() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := models.CreateContract{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("failed to bind json", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.service.Contract.CreateContract(ctx, payload)
	}
}

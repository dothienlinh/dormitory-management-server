package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) CreateContract() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := models.CreateContract{}
		if err := ctx.ShouldBindJSON(payload); err != nil {
			h.logger.Error("failed to bind json", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			h.logger.Error("failed to validate struct", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contract, err := h.service.Contract.CreateContract(payload)
		if err != nil {
			h.logger.Error("failed to create contract", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		Success(ctx, contract, http.StatusCreated)
	}
}

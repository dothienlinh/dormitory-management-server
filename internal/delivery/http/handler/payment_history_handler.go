package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentHistoryHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewPaymentHistoryHandler(useCases usecase.UseCases, logger logger.Logger) *PaymentHistoryHandler {
	return &PaymentHistoryHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *PaymentHistoryHandler) CreatePaymentHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("CreatePaymentHistory")

		var createPaymentHistory entity.CreatePaymentHistory
		if err := c.ShouldBindJSON(&createPaymentHistory); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.PaymentHistory().CreatePaymentHistory(c, &createPaymentHistory)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) GetPaymentHistoryByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetPaymentHistoryByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse payment history ID", zap.Error(err))
			resp = response.BadRequest("Invalid payment history ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.PaymentHistory().GetPaymentHistoryByID(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) GetPaymentHistoryByContractID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetPaymentHistoryByContractID")

		contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var filter entity.PaymentHistoryFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		filter.ContractID = uint64(contractID)
		resp = h.useCases.PaymentHistory().GetPaymentHistoryList(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) GetMyPaymentHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetMyPaymentHistory")

		userID, exists := c.Get("user_id")
		if !exists {
			h.logger.Error("User ID not found in context")
			resp = response.Unauthorized("User not authenticated")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var filter entity.PaymentHistoryFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.PaymentHistory().GetMyPaymentHistory(c, userID.(uint64), &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) GetPaymentHistoryList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetPaymentHistoryList")

		var filter entity.PaymentHistoryFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.PaymentHistory().GetPaymentHistoryList(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) UpdatePaymentHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("UpdatePaymentHistory")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse payment history ID", zap.Error(err))
			resp = response.BadRequest("Invalid payment history ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var updateData entity.UpdatePaymentHistory
		if err := c.ShouldBindJSON(&updateData); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.PaymentHistory().UpdatePaymentHistory(c, uint(id), &updateData)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) MakePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("MakePayment")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse payment history ID", zap.Error(err))
			resp = response.BadRequest("Invalid payment history ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			h.logger.Error("User ID not found in context")
			resp = response.Unauthorized("User not authenticated")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var makePaymentReq entity.MakePaymentRequest
		if err := c.ShouldBindJSON(&makePaymentReq); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.PaymentHistory().MakePayment(c, uint(id), userID.(uint64), &makePaymentReq)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) DeletePaymentHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DeletePaymentHistory")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse payment history ID", zap.Error(err))
			resp = response.BadRequest("Invalid payment history ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.PaymentHistory().DeletePaymentHistory(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *PaymentHistoryHandler) DownloadReceipt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DownloadReceipt")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse payment history ID", zap.Error(err))
			resp = response.BadRequest("Invalid payment history ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			h.logger.Error("User ID not found in context")
			resp = response.Unauthorized("User not authenticated")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.PaymentHistory().DownloadReceipt(c, uint(id), userID.(uint64))
		c.JSON(resp.Status, resp.Response)
	}
}

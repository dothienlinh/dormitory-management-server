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

type ContractHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewContractHandler(useCases usecase.UseCases, logger logger.Logger) *ContractHandler {
	return &ContractHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *ContractHandler) CreateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("CreateContract")

		var createContract entity.CreateContract
		if err := c.ShouldBindJSON(&createContract); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Contract().CreateContract(c, &createContract)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) GetContractByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetContractByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().GetContractByID(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) GetContractByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetContractByUserID")

		userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse user ID", zap.Error(err))
			resp = response.BadRequest("Invalid user ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().GetContractByUserID(c, uint(userID))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) GetListContracts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetListContracts")

		var filter entity.ContractFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Contract().GetListContracts(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) UpdateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("UpdateContract")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var updateData entity.UpdateContract
		if err := c.ShouldBindJSON(&updateData); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Contract().UpdateContract(c, uint(id), &updateData)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) DeleteContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DeleteContract")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().DeleteContract(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) GetMyContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetMyContract")

		userID, exists := c.Get("user_id")
		if !exists {
			h.logger.Error("User ID not found in context")
			resp = response.Unauthorized("User not authenticated")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().GetMyContract(c, userID.(uint64))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) DownloadContractPDF() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DownloadContractPDF")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
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

		resp = h.useCases.Contract().DownloadContractPDF(c, uint(id), userID.(uint64))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *ContractHandler) GetContractPaymentHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetContractPaymentHistory")

		contractID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
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

		resp = h.useCases.Contract().GetContractPaymentHistory(c, uint(contractID), userID.(uint64))
		c.JSON(resp.Status, resp.Response)
	}
}
